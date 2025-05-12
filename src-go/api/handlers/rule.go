package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/metacubex/mihomo/hub/executor"
	"github.com/metacubex/mihomo/log"
	"github.com/snakem982/pandora-box/api/models"
	"github.com/snakem982/pandora-box/internal"
	"github.com/snakem982/pandora-box/pkg/cache"
	"github.com/snakem982/pandora-box/pkg/constant"
	sys "github.com/snakem982/pandora-box/pkg/sys/proxy"
	"github.com/snakem982/pandora-box/pkg/utils"
)

func Rule(r chi.Router) {
	r.Mount("/rule", ruleRouter())
}

func ruleRouter() chi.Router {
	r := chi.NewRouter()
	// 忽略的域名
	r.Get("/ignore", getIgnore)
	r.Put("/ignore", updateIgnore)

	// 统一规则分组
	r.Get("/list", getTemplateList)
	r.Get("/template/{id}", getTemplate)
	r.Post("/template", addTemplate)
	r.Delete("/template/{id}", deleteTemplate)
	r.Put("/template", updateTemplate)

	// 校验测试
	r.Post("/test", testTemplate)
	r.Post("/switch", switchTemplate)

	// 忽略的域名
	r.Get("/num", getNum)

	return r
}

func getIgnore(w http.ResponseWriter, r *http.Request) {
	ignore, _ := sys.GetIgnore()
	render.JSON(w, r, ignore)
}

func updateIgnore(w http.ResponseWriter, r *http.Request) {
	var body []string
	if err := render.DecodeJSON(r.Body, &body); err != nil {
		ErrorResponse(w, r, err)
		return
	}

	log.Infoln("[updateIgnore] %v", body)

	if err := sys.SetIgnore(body); err != nil {
		ErrorResponse(w, r, err)
		return
	}

	render.NoContent(w, r)
}

func getTemplateList(w http.ResponseWriter, r *http.Request) {
	var list []models.Template
	_ = cache.GetList(constant.PrefixTemplate, &list)

	if len(list) == 0 {
		// 如果没有数据，使用默认的模板
		list2 := [3][]byte{internal.Template_0, internal.Template_1, internal.Template_2}
		titles := [3]string{"m1", "m2", "m3"}
		for i := 0; i < 3; i++ {
			template := models.Template{
				Id:       fmt.Sprintf("%s%d", constant.PrefixTemplate, i),
				Order:    int64(i),
				Title:    titles[i],
				Path:     fmt.Sprintf("/%s/%s.yaml", constant.DefaultTemplateDir, fmt.Sprintf("%s%d", constant.PrefixTemplate, i)),
				Selected: false,
			}

			// 存盘
			_, _ = utils.SaveFile(utils.GetUserHomeDir(template.Path), list2[i])

			// 存数据库
			_ = cache.Put(template.Id, template)

			// 返回页面
			list = append(list, template)
		}
	}

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
		Data     string          `json:"data"`
		Template models.Template `json:"template"`
	}{}
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		ErrorResponse(w, r, err)
		return
	}

	// 存盘
	_, err := utils.SaveFile(utils.GetUserHomeDir(req.Template.Path), []byte(req.Data))
	if err != nil {
		ErrorResponse(w, r, err)
		return
	}

	render.NoContent(w, r)
}

func testTemplate(w http.ResponseWriter, r *http.Request) {
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

	// 状态改变
	var list []models.Template
	_ = cache.GetList(constant.PrefixTemplate, &list)
	for _, m := range list {
		if m.Selected {
			m.Selected = false
			_ = cache.Put(m.Id, m)
			break
		}
	}
	if template.Selected {
		_ = cache.Put(template.Id, template)
	}

	// 进行配置切换
	internal.SwitchProfile(true)

	render.NoContent(w, r)
}

func getNum(w http.ResponseWriter, r *http.Request) {
	var num int
	_ = cache.Get("Rule_No", &num)
	res := struct {
		Data int `json:"data"`
	}{num}

	render.JSON(w, r, res)
}
