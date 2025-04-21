package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/metacubex/mihomo/hub/executor"
	"github.com/metacubex/mihomo/log"
	"github.com/metacubex/mihomo/tunnel/statistic"
	"github.com/snakem982/pandora-box/pandora/internal"
	sys "github.com/snakem982/pandora-box/pandora/pkg/sys/proxy"
	"github.com/snakem982/pandora-box/pandora/pkg/utils"
	"net/http"
)

func Pandora(r chi.Router) {
	r.Get("/version", getPandoraVersion)
	r.Mount("/pandora", PandoraRouter())
}

func PandoraRouter() chi.Router {
	r := chi.NewRouter()
	// 代理相关
	r.Put("/enableProxy", enableProxy)
	r.Get("/disableProxy", disableProxy)

	// 地址相关
	r.Put("/checkAddressPort", checkAddressPort)

	return r
}

func getPandoraVersion(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, render.M{"version": internal.PandoraVersion})
}

func enableProxy(w http.ResponseWriter, r *http.Request) {
	// 读取请求体
	mi := struct {
		BindAddress string `json:"bindAddress"`
		Port        int    `json:"port"`
	}{}
	if err := render.DecodeJSON(r.Body, &mi); err != nil {
		ErrorResponse(w, r, err)
		return
	}

	// 开启
	_ = sys.EnableProxy(mi.BindAddress, mi.Port)

	render.NoContent(w, r)
}

func disableProxy(w http.ResponseWriter, r *http.Request) {
	sys.DisableProxy()
	log.Warnln("System proxy disabled")
	statistic.DefaultManager.Range(func(c statistic.Tracker) bool {
		_ = c.Close()
		return true
	})
	log.Warnln("All connections disconnected")
	render.NoContent(w, r)
}

func checkAddressPort(w http.ResponseWriter, r *http.Request) {
	// 读取请求体
	mi := struct {
		BindAddress string `json:"bindAddress"`
		MixedPort   int    `json:"port"`
	}{}
	if err := render.DecodeJSON(r.Body, &mi); err != nil {
		ErrorResponse(w, r, err)
		return
	}

	// 检测到Px相同地址端口则跳过
	BindAddress := executor.GetGeneral().BindAddress
	MixedPort := executor.GetGeneral().MixedPort
	if BindAddress == mi.BindAddress && MixedPort == mi.MixedPort {
		render.NoContent(w, r)
		return
	}

	// 检测地址端口是否可用
	err := utils.IsPortAvailable(mi.BindAddress, mi.MixedPort)
	if err != nil {
		ErrorResponse(w, r, err)
		return
	}

	render.NoContent(w, r)
}
