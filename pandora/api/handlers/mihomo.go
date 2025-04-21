package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/snakem982/pandora-box/pandora/api/models"
	"github.com/snakem982/pandora-box/pandora/pkg/cache"
	"github.com/snakem982/pandora-box/pandora/pkg/constant"
	"net/http"
)

func Mihomo(r chi.Router) {
	r.Mount("/mihomo", MihomoRouter())
}

func MihomoRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/", getMihomo)
	r.Put("/", updateMihomo)

	return r
}

func getMihomo(w http.ResponseWriter, r *http.Request) {
	var mi models.Mihomo
	_ = cache.Get(constant.Mihomo, &mi)

	if mi.BindAddress == "" {
		mi = models.Mihomo{
			Mode:        "rule",
			Proxy:       false,
			Tun:         false,
			Port:        9697,
			BindAddress: "127.0.0.1",
			Stack:       "Mixed",
			Dns:         false,
			Ipv6:        false,
		}
		_ = cache.Put(constant.Mihomo, mi)
	}

	render.JSON(w, r, mi)
}

func updateMihomo(w http.ResponseWriter, r *http.Request) {
	// 读取请求体
	var mi models.Mihomo
	if err := render.DecodeJSON(r.Body, &mi); err != nil {
		ErrorResponse(w, r, err)
		return
	}

	_ = cache.Put(constant.Mihomo, mi)

	render.NoContent(w, r)
}
