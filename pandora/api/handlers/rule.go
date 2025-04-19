package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/metacubex/mihomo/hub/executor"
	"github.com/metacubex/mihomo/log"
	"github.com/snakem982/pandora-box/pandora/api/models"
	"github.com/snakem982/pandora-box/pandora/pkg/cache"
	"github.com/snakem982/pandora-box/pandora/pkg/constant"
	"github.com/snakem982/pandora-box/pandora/pkg/utils"
	"net/http"
)

func Rule(r chi.Router) {
	r.Mount("/rule", ruleRouter())
}

func ruleRouter() chi.Router {
	r := chi.NewRouter()
	// bypass
	r.Get("/bypass", getBypass)
	r.Post("/bypass", deleteBypass)

	// 统一规则分组
	r.Get("/list", getTemplateList)
	r.Get("/template/{id}", getTemplate)
	r.Post("/template", addTemplate)
	r.Delete("/template/{id}", deleteTemplate)
	r.Put("/template", updateTemplate)

	// 校验测试
	r.Post("/test", testTemplate)
	r.Post("/switch", switchTemplate)

	return r
}

func getBypass(w http.ResponseWriter, r *http.Request) {

	render.NoContent(w, r)
}

func deleteBypass(w http.ResponseWriter, r *http.Request) {

	render.NoContent(w, r)
}

func getTemplateList(w http.ResponseWriter, r *http.Request) {
	var list []models.Template
	_ = cache.GetList(constant.PrefixTemplate, &list)

	render.JSON(w, r, list)
}

func getTemplate(w http.ResponseWriter, r *http.Request) {
	// 获取路径参数中的ID
	id := chi.URLParam(r, "id")
	var template models.Template
	_ = cache.Get(id, &template)

	// 解析出模板的内容返回页面
	path := utils.GetUserHomeDir(template.Path)
	body, err := utils.ReadFile(path)
	if err != nil {
		ErrorResponse(w, r, err)
		return
	}

	render.PlainText(w, r, body)
}

func addTemplate(w http.ResponseWriter, r *http.Request) {
	// 读取请求体
	req := struct {
		Data  []byte `json:"data"`
		Title string `json:"title"`
	}{}
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		ErrorResponse(w, r, err)
		return
	}

	// 获取id
	order := utils.SnowflakeId()
	id := fmt.Sprintf("%s%d", constant.PrefixTemplate, order)
	path := fmt.Sprintf("/%s/%s.yaml", constant.DefaultTemplateDir, id)

	// 存盘
	_, err := utils.SaveFile(utils.GetUserHomeDir(path), req.Data)
	if err != nil {
		ErrorResponse(w, r, err)
		return
	}

	// 存数据库
	template := models.Template{
		Id:    id,
		Path:  path,
		Order: order,
		Title: req.Title,
	}
	_ = cache.Put(id, template)

	render.NoContent(w, r)
}

func deleteTemplate(w http.ResponseWriter, r *http.Request) {
	// 获取路径参数中的ID
	id := chi.URLParam(r, "id")

	_ = cache.Delete(id)

	render.NoContent(w, r)
}

func updateTemplate(w http.ResponseWriter, r *http.Request) {
	// 读取请求体
	req := struct {
		Data     []byte          `json:"data"`
		Template models.Template `json:"template"`
	}{}
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		ErrorResponse(w, r, err)
		return
	}

	// 存盘
	_, err := utils.SaveFile(utils.GetUserHomeDir(req.Template.Path), req.Data)
	if err != nil {
		ErrorResponse(w, r, err)
		return
	}

	// 存数据库
	_ = cache.Put(req.Template.Id, req.Template)

	render.NoContent(w, r)
}

func testTemplate(w http.ResponseWriter, r *http.Request) {
	// 读取请求体
	req := struct {
		Data []byte `json:"data"`
	}{}
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		ErrorResponse(w, r, err)
		return
	}

	_, err := executor.ParseWithBytes(req.Data)
	if err != nil {
		log.Warnln("[testTemplate] error: %v", err)
		ErrorResponse(w, r, err)
		return
	}

	render.NoContent(w, r)
}

func switchTemplate(w http.ResponseWriter, r *http.Request) {
	// 读取请求体
	var template models.Template
	if err := render.DecodeJSON(r.Body, &template); err != nil {
		ErrorResponse(w, r, err)
		return
	}

	if template.Selected {
		var list []models.Template
		_ = cache.GetList(constant.PrefixTemplate, &list)
		for _, m := range list {
			if m.Selected {
				m.Selected = false
				_ = cache.Put(m.Id, m)
				break
			}
		}
		// todo 切换
	} else {
		// todo 切换
	}
	_ = cache.Put(template.Id, template)

	render.NoContent(w, r)
}
