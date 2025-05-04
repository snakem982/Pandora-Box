package job

import (
	"github.com/metacubex/mihomo/log"
	"github.com/snakem982/pandora-box/api/models"
	"github.com/snakem982/pandora-box/internal"
	"github.com/snakem982/pandora-box/pkg/cache"
	"github.com/snakem982/pandora-box/pkg/constant"
	"github.com/snakem982/pandora-box/pkg/cron"
	"github.com/snakem982/pandora-box/pkg/utils"
	"strconv"
	"sync"
	"time"
)

var refreshLock sync.Mutex

func RefreshJob() {
	cron.AddTask(30*time.Minute, doRefresh)
}

func doRefresh() {
	if refreshLock.TryLock() {
		defer refreshLock.Unlock()
	} else {
		return
	}

	log.Infoln("refresh job start")

	// 获取需要更新的订阅
	var profiles []models.Profile
	_ = cache.GetList(constant.PrefixProfile, &profiles)
	if profiles == nil || len(profiles) == 0 {
		return
	}

	// 过滤远程订阅
	var filteredProfiles []models.Profile
	for _, profile := range profiles {
		if profile.Type == 1 && profile.Interval != "" {
			filteredProfiles = append(filteredProfiles, profile)
		}
	}

	log.Infoln("refresh job find %d profiles", len(filteredProfiles))

	// 进行更新逻辑
	for _, fp := range filteredProfiles {
		// 获取指针
		profile := &fp

		// 获取更新间隔
		interval, err := strconv.Atoi(profile.Interval)
		if err != nil {
			continue
		}
		// 计算上次更新时间到现在时间的间隔
		duration := time.Now().Sub(profile.GetUpdateTime())
		if int(duration.Hours()) < interval {
			continue
		}

		log.Infoln("refresh job profile %v need fresh", profile.Title)

		// 进行更新
		go func() {
			// 标题
			title := profile.Title

			// 发送请求
			sub := profile.Content
			headers := map[string]string{}
			res, err := utils.FastGet(sub, headers, internal.GetProxyUrl())
			if err != nil {
				log.Errorln("[Refresh] Sub=%s, URL = %s, Request Error:%v", title, sub, err)
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

				log.Infoln("refresh job profile %v fresh success", profile.Title)
			} else {
				log.Errorln("[Refresh] Sub=%s, URL = %s, Resolve Error:%v", title, sub, err)
			}
		}()
	}
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
