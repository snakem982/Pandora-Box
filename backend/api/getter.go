package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/metacubex/mihomo/hub/route"
	"net/http"
	"pandora-box/backend/cache"
	"pandora-box/backend/constant"
	"pandora-box/backend/meta"
	"pandora-box/backend/resolve"
	"pandora-box/backend/spider"
	"pandora-box/backend/tools"
)

func Getter(r chi.Router) {
	r.Get("/crawl", crawling)
	r.Mount("/getter", profileGetter())
}

func profileGetter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", getGetter)
	r.Post("/", postGetter)

	r.Route("/{id}", func(r chi.Router) {
		r.Put("/", putGetter)
		r.Delete("/", deleteGetter)
	})

	return r
}

func crawling(w http.ResponseWriter, r *http.Request) {
	crawl := spider.Crawl()
	if !crawl {
		render.NoContent(w, r)
		return
	}

	bytes := cache.Get(constant.DefaultProfile)
	profile := resolve.Profile{}
	err := json.Unmarshal(bytes, &profile)
	if err != nil {
		render.NoContent(w, r)
		return
	}
	if profile.Selected {
		meta.StartCore(profile, true)
	}

	render.NoContent(w, r)
}

func getGetter(w http.ResponseWriter, r *http.Request) {
	res := make([]spider.Getter, 0)

	values := cache.GetList(constant.PrefixGetter)
	if len(values) > 0 {
		for _, value := range values {
			getter := spider.Getter{}
			_ = json.Unmarshal(value, &getter)
			res = append(res, getter)
		}
	}

	render.JSON(w, r, res)
}

func postGetter(w http.ResponseWriter, r *http.Request) {
	req := spider.Getter{}
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, route.ErrBadRequest)
		return
	}
	req.Id = fmt.Sprintf("%s%d", constant.PrefixGetter, tools.SnowflakeId())
	bytes, _ := json.Marshal(req)
	_ = cache.Put(req.Id, bytes)

	render.NoContent(w, r)
}

func putGetter(w http.ResponseWriter, r *http.Request) {
	req := spider.Getter{}
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, route.ErrBadRequest)
		return
	}
	bytes, _ := json.Marshal(req)
	_ = cache.Put(req.Id, bytes)

	render.NoContent(w, r)
}

func deleteGetter(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_ = cache.Delete(id)

	render.NoContent(w, r)
}
