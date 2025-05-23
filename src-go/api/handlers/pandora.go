package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/metacubex/mihomo/hub/executor"
	"github.com/metacubex/mihomo/log"
	"github.com/metacubex/mihomo/tunnel/statistic"
	"github.com/snakem982/pandora-box/api"
	"github.com/snakem982/pandora-box/api/job"
	"github.com/snakem982/pandora-box/api/models"
	"github.com/snakem982/pandora-box/pkg/cache"
	"github.com/snakem982/pandora-box/pkg/constant"
	sys "github.com/snakem982/pandora-box/pkg/sys/proxy"
	"github.com/snakem982/pandora-box/pkg/utils"
	"net/http"
	"os"
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

	// 配置目录
	r.Get("/configDir", configDir)

	// 退出px
	r.Get("/exit", exitPx)

	return r
}

func getPandoraVersion(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, render.M{"version": api.Version})
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
	if !executor.GetGeneral().Tun.Enable {
		statistic.DefaultManager.Range(func(c statistic.Tracker) bool {
			_ = c.Close()
			return true
		})
	}
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
	var mc models.Mihomo
	_ = cache.Get(constant.Mihomo, &mc)
	if mc.BindAddress == mi.BindAddress && mc.Port == mi.MixedPort {
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

func configDir(w http.ResponseWriter, r *http.Request) {
	render.PlainText(w, r, utils.GetUserHomeDir())
}

func exitPx(w http.ResponseWriter, r *http.Request) {
	job.Exit(false)
	render.PlainText(w, r, "ok")
	os.Exit(0)
}
