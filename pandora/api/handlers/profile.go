package handlers

import (
	"net/http"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/metacubex/mihomo/hub/route"
	"github.com/metacubex/mihomo/log"
	"github.com/snakem982/pandora-box/pandora/api/models"
	"github.com/snakem982/pandora-box/pandora/internal"
	"github.com/snakem982/pandora-box/pandora/pkg/cache"
	"github.com/snakem982/pandora-box/pandora/pkg/constant"
	"github.com/snakem982/pandora-box/pandora/pkg/utils"
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
	r.Post("/delete", deleteProfile)
	// 修改
	r.Put("/", putProfile)
	// 查找
	r.Get("/", getProfile)
	// 更新订阅
	r.Put("/refresh", refreshProfile)
	// 切换订阅
	r.Patch("/", switchProfile)
	// 存储排序
	r.Get("/order", saveProfileOrder)

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
	_ = cache.Put(profile.Id, *profile)
}

func getProfile(w http.ResponseWriter, r *http.Request) {
	var res []models.Profile
	_ = cache.GetList(constant.PrefixProfile, &res)

	var order []models.Profile
	_ = cache.Get(constant.ProfileOrder, &order)

	// 创建一个 map 用于快速查找 order 中的元素
	orderMap := make(map[string]int)
	for index, item := range order {
		orderMap[item.Id] = index
	}

	// 对 res 进行排序
	sort.SliceStable(res, func(i, j int) bool {
		// 如果 res[i] 和 res[j] 都在 order 中，按 order 中的顺序排序
		indexI, existsI := orderMap[res[i].Id]
		indexJ, existsJ := orderMap[res[j].Id]
		if existsI && existsJ {
			return indexI < indexJ
		}
		// 如果只有一个在 order 中，优先排序在 order 中的
		if existsI {
			return true
		}
		if existsJ {
			return false
		}
		// 如果都不在 order 中，按 Order 字段排序
		return res[i].Order < res[j].Order
	})

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

	// 返回页面list
	ps := make([]*models.Profile, 0)

	// 返回页面错误
	var tempErr error

	// 解析存盘
	err := internal.Resolve(profile.Content, profile, false)
	if err == nil {
		if profile.Title == "" {
			profile.Title = "Local-" + utils.GetDateTime()
		}
		UpdateDb(profile, 2)
		ps = append(ps, profile)
		render.JSON(w, r, ps)
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
			ps = append(ps, subProfile)
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

	render.JSON(w, r, ps)
}

func refreshProfile(w http.ResponseWriter, r *http.Request) {
	// 获取数据
	profile := &models.Profile{}
	if err := render.DecodeJSON(r.Body, profile); err != nil {
		ErrorResponse(w, r, err)
		return
	}
	title := profile.Title

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
		if title != "" {
			profile.Title = title
		}
		UpdateDb(profile, 1)

		// 如果配置正在使用中  进行配置更新
		if profile.Selected {
			internal.StartCore(*profile)
		}
	} else {
		ErrorResponse(w, r, err)
		log.Errorln("[refreshProfile] URL = %s, Resolve Error:%v", sub, err)
		return
	}

	render.JSON(w, r, profile)
}

func putProfile(w http.ResponseWriter, r *http.Request) {
	profile := &models.Profile{}
	if err := render.DecodeJSON(r.Body, profile); err != nil {
		ErrorResponse(w, r, err)
		return
	}

	_ = cache.Put(profile.Id, profile)

	render.NoContent(w, r)
}

// 删除配置
func deleteProfile(w http.ResponseWriter, r *http.Request) {
	profile := &models.Profile{}
	if err := render.DecodeJSON(r.Body, profile); err != nil {
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

	internal.StartCore(profile)

	render.NoContent(w, r)
}
