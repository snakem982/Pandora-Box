package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/metacubex/mihomo/hub/route"
	"github.com/metacubex/mihomo/log"
	"net/http"
	"pandora-box/backend/system/proxy"
)

func System(r chi.Router) {
	r.Put("/system/{port}", func(w http.ResponseWriter, r *http.Request) {
		port := chi.URLParam(r, "port")
		proxy.SetProxy(port)

		render.NoContent(w, r)
	})

	r.Delete("/system", func(w http.ResponseWriter, r *http.Request) {
		proxy.RemoveProxy()

		render.NoContent(w, r)
	})
}

func Ignore(r chi.Router) {
	r.Get("/ignore", func(w http.ResponseWriter, r *http.Request) {
		ignore, _ := proxy.GetIgnore()
		render.JSON(w, r, ignore)
	})

	r.Put("/ignore", func(w http.ResponseWriter, r *http.Request) {
		body := make([]string, 0)
		if err := render.DecodeJSON(r.Body, &body); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, route.ErrBadRequest)
			return
		}

		log.Infoln("ignore: %v", body)

		if err := proxy.SetIgnore(body); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, route.ErrBadRequest)
			return
		}

		render.NoContent(w, r)
	})
}
