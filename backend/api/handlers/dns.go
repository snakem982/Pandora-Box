package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/metacubex/mihomo/hub/executor"
	"github.com/metacubex/mihomo/log"
	"github.com/snakem982/pandora-box/api/models"
	"github.com/snakem982/pandora-box/internal"
	"github.com/snakem982/pandora-box/pkg/cache"
	"github.com/snakem982/pandora-box/pkg/constant"
	"net/http"
)

func DNS(r chi.Router) {
	r.Mount("/pDns", DNSRouter())
}

func DNSRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/", getDNS)
	r.Put("/", updateDNS)
	r.Post("/switch", switchDNS)

	return r
}

func getDNS(w http.ResponseWriter, r *http.Request) {
	var dns models.Dns
	_ = cache.Get(constant.Dns, &dns)

	if dns.Content == "" {
		dns.Content = internal.DefaultDNS
		_ = cache.Put(constant.Dns, dns)
	}

	render.JSON(w, r, dns)
}

func updateDNS(w http.ResponseWriter, r *http.Request) {
	// 读取请求体
	req := struct {
		Data string `json:"data"`
	}{}
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		ErrorResponse(w, r, err)
		return
	}

	_, err := executor.ParseWithBytes([]byte(req.Data))
	if err != nil {
		log.Warnln("[testDNS] error: %v", err)
		ErrorResponse(w, r, err)
		return
	}

	// 保存dns 内容
	dns := models.Dns{
		Content: req.Data,
	}
	_ = cache.Put(constant.Dns, dns)

	// 进行 配置重载
	var mi models.Mihomo
	_ = cache.Get(constant.Mihomo, &mi)
	if mi.Dns {
		internal.SwitchProfile(true)
	}

	render.NoContent(w, r)
}

func switchDNS(w http.ResponseWriter, r *http.Request) {
	// 读取请求体
	req := struct {
		Enable bool `json:"enable"`
	}{}
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		ErrorResponse(w, r, err)
		return
	}

	// 进行 配置重载
	var mi models.Mihomo
	_ = cache.Get(constant.Mihomo, &mi)
	mi.Dns = req.Enable
	_ = cache.Put(constant.Mihomo, mi)
	internal.SwitchProfile(true)

	render.NoContent(w, r)
}
