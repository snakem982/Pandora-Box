package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"pandora-box/backend/constant"
)

// Hello 返回hello pandora-box
func Hello(r chi.Router) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, render.M{"hello": "pandora-box"})
	})

	r.Get("/ok", func(w http.ResponseWriter, r *http.Request) {
		render.PlainText(w, r, "ok")
	})
}

// Version 返回Version
func Version(r chi.Router) {
	r.Get("/version", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, render.M{"version": constant.PandoraVersion})
	})
}
