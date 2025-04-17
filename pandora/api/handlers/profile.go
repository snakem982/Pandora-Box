package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/metacubex/mihomo/hub/route"
	"github.com/metacubex/mihomo/log"
	"github.com/snakem982/pandora-box/pandora/api/models"
	"github.com/snakem982/pandora-box/pandora/internal"
	"github.com/snakem982/pandora-box/pandora/pkg/cache"
	"github.com/snakem982/pandora-box/pandora/pkg/constant"
	"github.com/snakem982/pandora-box/pandora/pkg/utils"
	"io"
	"mime/multipart"
	"net/http"
)

func Profile(r chi.Router) {
	r.Mount("/profile", profileRouter())
}

func profileRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", getProfile)
	r.Post("/file", postFile)
	r.Post("/", postProfile)
	r.Put("/refresh", refreshProfile)

	//r.Route("/{id}", func(r chi.Router) {
	//	r.Put("/", putProfile)
	//	r.Delete("/", deleteProfile)
	//	r.Patch("/", patchProfile)
	//})

	return r
}

// ErrorResponse 是一个共通的方法，用于返回错误信息到客户端
func ErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	render.Status(r, http.StatusBadRequest)
	render.JSON(w, r, route.HTTPError{Message: err.Error()})
}

// UpdateDb 更新数据库
func UpdateDb(profile *models.Profile, kind int) {
	profile.Type = kind
	profile.SetUpdateTime()
	if kind == 2 {
		profile.Content = ""
	}
	_ = cache.Put(profile.Id, profile)
}

func getProfile(w http.ResponseWriter, r *http.Request) {
	var res []models.Profile
	_ = cache.GetList(constant.PrefixProfile, &res)

	render.JSON(w, r, res)
}

func postFile(w http.ResponseWriter, r *http.Request) {
	// 获取数据
	_, header, _ := r.FormFile("file")
	open, _ := header.Open()
	defer func(open multipart.File) {
		err := open.Close()
		if err != nil {

		}
	}(open)
	content, _ := io.ReadAll(open)

	// 解析存盘
	profile := &models.Profile{}
	err := internal.Resolve(string(content), profile, false)
	if err != nil {
		log.Errorln("[postFile] Resolve Error:%v", err)
		ErrorResponse(w, r, err)
		return
	}

	// 更新数据库
	profile.Title = header.Filename
	UpdateDb(profile, 2)

	render.NoContent(w, r)
}

func postProfile(w http.ResponseWriter, r *http.Request) {
	// 数据解析
	body := struct {
		Data string `json:"data"`
	}{}
	if err := render.DecodeJSON(r.Body, &body); err != nil {
		ErrorResponse(w, r, route.ErrBadRequest)
		return
	}

	// 解析存盘
	profile := &models.Profile{}
	err := internal.Resolve(body.Data, profile, false)
	if err == nil {
		profile.Title = "Local-" + utils.RandString(5)
		UpdateDb(profile, 2)
		render.NoContent(w, r)
		return
	} else {
		log.Errorln("[postProfile] Resolve Error:%v", err)
	}

	// 扫描订阅
	subs := internal.ScanSubs(body.Data)
	ok := false
	var tempErr error
	for _, sub := range subs {
		headers := map[string]string{}
		res, err := utils.FastGet(sub, headers, internal.GetProxyUrl())
		if err != nil {
			tempErr = err
			log.Errorln("[postProfile] URL = %s, Request Error:%v", sub, err)
			continue
		}

		// 解析存盘
		subProfile := &models.Profile{}
		err = internal.Resolve(res.Body, subProfile, false)
		if err == nil {
			// 进行请求头解析
			internal.ParseHeaders(res.Headers, sub, subProfile)
			UpdateDb(subProfile, 1)
			ok = true
		} else {
			tempErr = err
			log.Errorln("[postProfile] URL = %s, Resolve Error:%v", sub, err)
		}
	}
	if !ok {
		ErrorResponse(w, r, tempErr)
		return
	}

	render.NoContent(w, r)
}

func refreshProfile(w http.ResponseWriter, r *http.Request) {
	profile := models.Profile{}
	if err := render.DecodeJSON(r.Body, &profile); err != nil {
		ErrorResponse(w, r, err)
		return
	}

	// 发送请求
	sub := profile.Content
	headers := map[string]string{}
	res, err := utils.FastGet(sub, headers, internal.GetProxyUrl())
	if err != nil {
		ErrorResponse(w, r, err)
		log.Errorln("[refreshProfile] URL = %s, Request Error:%v", sub, err)
		return
	}

	// 解析存盘
	subProfile := &models.Profile{}
	err = internal.Resolve(res.Body, subProfile, true)
	if err == nil {
		// 进行请求头解析
		internal.ParseHeaders(res.Headers, sub, subProfile)
		UpdateDb(subProfile, 1)
	} else {
		ErrorResponse(w, r, err)
		log.Errorln("[refreshProfile] URL = %s, Resolve Error:%v", sub, err)
	}

	render.NoContent(w, r)
}
