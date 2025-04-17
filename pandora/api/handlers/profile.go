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
	"net/http"
	"path/filepath"
	"strings"
)

func Profile(r chi.Router) {
	r.Mount("/profile", profileRouter())
}

func profileRouter() http.Handler {
	r := chi.NewRouter()
	// 增加
	r.Post("/", addFromWeb)
	r.Post("/file", addFromFile)
	// 删除
	r.Delete("/", deleteProfile)
	// 修改
	r.Put("/", putProfile)
	// 查找
	r.Get("/", getProfile)
	// 更新订阅
	r.Put("/refresh", refreshProfile)
	// 切换订阅
	r.Patch("/", switchProfile)

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

func addFromFile(w http.ResponseWriter, r *http.Request) {
	// 获取数据
	profile := &models.Profile{}
	if err := render.DecodeJSON(r.Body, profile); err != nil {
		ErrorResponse(w, r, err)
		return
	}

	// 解析存盘
	err := internal.Resolve(profile.Content, profile, false)
	if err != nil {
		log.Errorln("[addFromFile] Resolve Error:%v", err)
		ErrorResponse(w, r, err)
		return
	}

	// 更新数据库
	UpdateDb(profile, 2)

	render.NoContent(w, r)
}

func addFromWeb(w http.ResponseWriter, r *http.Request) {
	// 获取数据
	profile := &models.Profile{}
	if err := render.DecodeJSON(r.Body, profile); err != nil {
		ErrorResponse(w, r, err)
		return
	}

	// 返回页面错误
	var tempErr error

	// 解析存盘
	err := internal.Resolve(profile.Content, profile, false)
	if err == nil {
		profile.Title = "Local-" + utils.RandString(5)
		UpdateDb(profile, 2)
		render.NoContent(w, r)
		return
	} else {
		tempErr = err
		log.Errorln("[addFromWeb] Resolve Error:%v", err)
	}

	// 扫描订阅
	subs := internal.ScanSubs(profile.Content)
	ok := false
	for _, sub := range subs {
		headers := map[string]string{}
		res, err := utils.FastGet(sub, headers, internal.GetProxyUrl())
		if err != nil {
			tempErr = err
			log.Errorln("[addFromWeb] URL = %s, Request Error:%v", sub, err)
			continue
		}

		// 解析存盘
		subProfile := &models.Profile{
			Content: sub,
		}
		err = internal.Resolve(res.Body, subProfile, false)
		if err == nil {
			// 进行请求头解析
			internal.ParseHeaders(res.Headers, sub, subProfile)
			UpdateDb(subProfile, 1)
			ok = true
		} else {
			tempErr = err
			log.Errorln("[addFromWeb] URL = %s, Resolve Error:%v", sub, err)
		}
	}
	if !ok {
		ErrorResponse(w, r, tempErr)
		return
	}

	render.NoContent(w, r)
}

func refreshProfile(w http.ResponseWriter, r *http.Request) {
	// 获取数据
	profile := &models.Profile{}
	if err := render.DecodeJSON(r.Body, profile); err != nil {
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
	err = internal.Resolve(res.Body, profile, true)
	if err == nil {
		// 进行请求头解析
		internal.ParseHeaders(res.Headers, sub, profile)
		UpdateDb(profile, 1)
	} else {
		ErrorResponse(w, r, err)
		log.Errorln("[refreshProfile] URL = %s, Resolve Error:%v", sub, err)
		return
	}

	render.NoContent(w, r)
}

func putProfile(w http.ResponseWriter, r *http.Request) {
	profile := models.Profile{}
	if err := render.DecodeJSON(r.Body, &profile); err != nil {
		ErrorResponse(w, r, err)
		return
	}

	_ = cache.Put(profile.Id, profile)

	render.NoContent(w, r)
}

// 删除配置
func deleteProfile(w http.ResponseWriter, r *http.Request) {
	profile := models.Profile{}
	if err := render.DecodeJSON(r.Body, &profile); err != nil {
		ErrorResponse(w, r, err)
		return
	}

	path := utils.GetUserHomeDir(profile.Path)
	dir := filepath.Dir(path)
	if strings.HasSuffix(dir, "profiles") {
		_ = utils.DeletePath(path)
	} else {
		_ = utils.DeletePath(dir)
	}
	_ = cache.Delete(profile.Id)

	render.NoContent(w, r)
}

// 切换配置
func switchProfile(w http.ResponseWriter, r *http.Request) {
	var profile models.Profile
	if err := render.DecodeJSON(r.Body, &profile); err != nil {
		ErrorResponse(w, r, err)
		return
	}

	var profiles []models.Profile
	_ = cache.GetList(constant.PrefixProfile, &profiles)
	for _, p := range profiles {
		if p.Selected {
			p.Selected = false
			_ = cache.Put(p.Id, p)
			break
		} else {
			continue
		}
	}
	profile.Selected = true
	_ = cache.Put(profile.Id, profile)

	internal.StartCore(profile, true)

	render.NoContent(w, r)
}
