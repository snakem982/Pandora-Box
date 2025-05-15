package job

import (
	"github.com/metacubex/mihomo/log"
	"github.com/snakem982/pandora-box/api/models"
	"github.com/snakem982/pandora-box/internal"
	"github.com/snakem982/pandora-box/pkg/cache"
	"github.com/snakem982/pandora-box/pkg/constant"
	"github.com/snakem982/pandora-box/pkg/cron"
	"github.com/snakem982/pandora-box/pkg/proxy"
	"github.com/snakem982/pandora-box/pkg/utils"
	"strconv"
	"strings"
	"sync"
	"time"
)

var refreshLock sync.Mutex

func RefreshJob() {
	cron.AddTask("Refresh", 30*time.Minute, DoRefresh)
}

func DoRefresh() {
	if refreshLock.TryLock() {
		defer refreshLock.Unlock()
	} else {
		return
	}

	log.Infoln("[Refresh] job start")

	// 获取需要更新的订阅
	var profiles []models.Profile
	_ = cache.GetList(constant.PrefixProfile, &profiles)
	if profiles == nil || len(profiles) == 0 {
		return
	}

	// 过滤远程订阅
	var filteredProfiles []models.Profile
	for _, profile := range profiles {
		if profile.Type != 1 || profile.Interval == "" {
			continue
		}

		// 获取更新间隔
		interval, err := strconv.Atoi(profile.Interval)
		if err != nil {
			continue
		}

		// 计算上次更新时间到现在时间的间隔
		diff := utils.GetHourDiff(profile.GetUpdateTime())
		if diff < interval {
			continue
		}

		filteredProfiles = append(filteredProfiles, profile)
	}

	log.Infoln("[Refresh] job find %d profile need fresh", len(filteredProfiles))

	// 进行更新逻辑
	for _, fp := range filteredProfiles {
		// 获取指针
		profile := &fp

		log.Infoln("[Refresh] job profile %v fresh start", profile.Title)

		// 进行更新
		title := profile.Title

		// 发送请求
		sub := profile.Content
		headers := map[string]string{}
		res, err := utils.FastGet(sub, headers, proxy.GetProxyUrl())
		if err != nil {
			log.Errorln("[Refresh] Sub=%s, URL = %s, Request Error:%v", title, sub, err)
			continue
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

			log.Infoln("[Refresh] job profile %v fresh success", profile.Title)
		} else {
			log.Errorln("[Refresh] Sub=%s, URL = %s, Resolve Error:%v", title, sub, err)
		}
	}
}

// UpdateDb 更新数据库
func UpdateDb(profile *models.Profile, kind int) {
	profile.Type = kind
	profile.SetUpdateTime()
	if kind == 2 {
		profile.Content = ""
	}
	checkTitle(profile)
	_ = cache.Put(profile.Id, *profile)
}

func checkTitle(profile *models.Profile) {
	profile.Title = strings.TrimSpace(profile.Title)
	if profile.Title != "" {
		return
	}
	if profile.Type == 1 {
		profile.Title = "Sub-" + utils.GetDateTime()
	} else if profile.Type == 2 {
		profile.Title = "Local-" + utils.GetDateTime()
	} else {

	}
}
