package pandora

import (
	"fmt"
	"github.com/metacubex/mihomo/hub/route"
	"github.com/metacubex/mihomo/log"
	"github.com/snakem982/pandora-box/api/handlers"
	"github.com/snakem982/pandora-box/api/job"
	"github.com/snakem982/pandora-box/internal"
	"github.com/snakem982/pandora-box/pkg/cache"
	"github.com/snakem982/pandora-box/pkg/constant"
	"github.com/snakem982/pandora-box/pkg/cron"
	"github.com/snakem982/pandora-box/pkg/utils"
	"strconv"
	"strings"
	"time"
)

func Init(isClient bool) {
	internal.Init(isClient)
}

func Release() {
	cache.Close()
}

func StartCore(server string, isClient bool) (port int, secret string) {
	// 初始化
	internal.Init(isClient)

	route.Register(handlers.Profile)
	route.Register(handlers.WebTest)
	route.Register(handlers.Rule)
	route.Register(handlers.DNS)
	route.Register(handlers.Mihomo)
	route.Register(handlers.Pandora)

	// 获取密钥
	_ = cache.Get(constant.SecretKey, &secret)
	if secret == "" {
		secret = utils.RandString(16)
		_ = cache.Put(constant.SecretKey, secret)
	}

	// 启动api
	apiAddr := route.StartByPandora(false, secret)
	split := strings.Split(apiAddr, ":")
	host := split[0]
	port, _ = strconv.Atoi(split[1])
	log.Infoln("Routing startup completed")

	// 开启mihomo
	internal.SwitchProfile(false)

	// 进行回调
	if server != "" {
		url := fmt.Sprintf("http://%s/pxStore?port=%v&secret=%v", server, port, secret)
		for {
			log.Infoln("向地址发送数据：%s", url)
			body, _, err := utils.SendGet(url, map[string]string{}, "")
			if err != nil {
				time.Sleep(500 * time.Millisecond)
				continue
			}
			if body == "ok" {
				break
			}
		}
	}

	// 定时更新订阅
	if !isClient {
		url := fmt.Sprintf("http://%s:%d/wait", host, port)
		headers := map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", secret),
		}
		_, _, _ = utils.SendGet(url, headers, "")

		job.LogJob("px-server.log")
		job.RefreshJob()
	}
	// 开启定时任务
	go cron.Start()

	return port, secret
}
