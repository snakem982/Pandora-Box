package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	C "github.com/metacubex/mihomo/constant"
	"github.com/metacubex/mihomo/hub/executor"
	"github.com/metacubex/mihomo/log"
	"net/http"
	"os"
	"pandora-box/backend/cache"
	"pandora-box/backend/constant"
	"pandora-box/backend/meta"
	"pandora-box/backend/resolve"
	"path/filepath"
	"strings"
)

// MyRules 返回hello pandora-box
func MyRules(r chi.Router) {
	r.Get("/myRules/default", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, render.M{"status": "ok", "buf": string(resolve.PandoraDefaultConfig)})
	})

	r.Get("/myRules", func(w http.ResponseWriter, r *http.Request) {
		buf, err := os.ReadFile(filepath.Join(C.Path.HomeDir(), constant.DefaultTemplate))
		if err != nil {
			log.Warnln("Read DefaultTemplate error: %s", err.Error())
			render.JSON(w, r, render.M{"status": "read failed"})
			return
		}

		render.JSON(w, r, render.M{"status": "ok", "buf": string(buf)})
	})

	r.Post("/myRules/test", func(w http.ResponseWriter, r *http.Request) {
		body := struct {
			Data string `json:"data"`
		}{}
		if err := render.DecodeJSON(r.Body, &body); err != nil {
			log.Warnln("json parse error: %s", err.Error())
			render.JSON(w, r, render.M{"status": "json parse error " + err.Error()})
			return
		}
		replace := strings.Replace(body.Data,
			resolve.PandoraDefaultPlace,
			"uploads/"+constant.DefaultProfile+".yaml",
			1)

		_, err := executor.ParseWithBytes([]byte(replace))
		if err != nil {
			log.Warnln("mihomo parse error: %s", err.Error())
			render.JSON(w, r, render.M{"status": "mihomo parse error " + err.Error()})
			return
		}

		render.JSON(w, r, render.M{"status": "ok"})
	})

	r.Post("/myRules/save", func(w http.ResponseWriter, r *http.Request) {
		body := struct {
			Data string `json:"data"`
		}{}
		if err := render.DecodeJSON(r.Body, &body); err != nil {
			log.Warnln("json parse error: %s", err.Error())
			render.JSON(w, r, render.M{"status": "json parse error " + err.Error()})
			return
		}

		err := os.WriteFile(filepath.Join(C.Path.HomeDir(), constant.DefaultTemplate), []byte(body.Data), 0666)
		if err != nil {
			log.Warnln("save file error: %s", err.Error())
			render.JSON(w, r, render.M{"status": "save file error " + err.Error()})
			return
		}

		on := cache.Get(constant.DefaultTemplate)

		if string(on) == "on" {
			meta.SwitchProfile(true)
		}

		render.JSON(w, r, render.M{"status": "ok"})
	})

	r.Get("/myRules/on", func(w http.ResponseWriter, r *http.Request) {
		on := cache.Get(constant.DefaultTemplate)
		render.JSON(w, r, render.M{"status": "ok", "on": string(on)})
	})

	r.Put("/myRules/on", func(w http.ResponseWriter, r *http.Request) {
		body := struct {
			Data string `json:"data"`
		}{}
		if err := render.DecodeJSON(r.Body, &body); err != nil {
			log.Warnln("json parse error: %s", err.Error())
			render.JSON(w, r, render.M{"status": "json parse error " + err.Error()})
			return
		}
		_ = cache.Put(constant.DefaultTemplate, []byte(body.Data))

		meta.SwitchProfile(true)

		render.JSON(w, r, render.M{"status": "ok"})
	})

}
