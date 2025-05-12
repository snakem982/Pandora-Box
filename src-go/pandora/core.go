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
	"net/url"
	"time"
)

func Init(isClient bool) {
	internal.Init(isClient)
}

func Release() {
	cache.Close()
}

func StartCore(server string, isClient bool) (port int, secret string) {

	route.Register(handlers.Profile)
	route.Register(handlers.WebTest)
	route.Register(handlers.Rule)
	route.Register(handlers.DNS)
	route.Register(handlers.Mihomo)
	route.Register(handlers.Pandora)

	// 设置地址
	host := "127.0.0.1"

	// 获取端口
	if utils.IsPortAvailable(host, 9686) == nil {
		port = 9686
	} else {
		port, _ = utils.GetRandomPort(host)
	}

	// 获取密钥
	_ = cache.Get(constant.SecretKey, &secret)
	if secret == "" {
		secret = utils.RandString(16)
		_ = cache.Put(constant.SecretKey, secret)
	}

	cors := route.Cors{AllowOrigins: []string{"*"}, AllowPrivateNetwork: true}
	route.StartByPandoraBox(host, port, secret, cors)
	log.Infoln("Routing startup completed")

	// 开启mihomo
	internal.SwitchProfile(false)

	// 进行回调
	if server != "" {
		encodedHome := url.QueryEscape(utils.GetUserHomeDir())
		callbackUrl := fmt.Sprintf("http://%s/pxStore?port=%v&secret=%v&home=%v", server, port, secret, encodedHome)
		for {
			log.Infoln("向地址发送数据：%s", callbackUrl)
			body, _, err := utils.SendGet(callbackUrl, map[string]string{}, "")
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
		waitUrl := fmt.Sprintf("http://%s:%d/wait", host, port)
		headers := map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", secret),
		}
		_, _, _ = utils.SendGet(waitUrl, headers, "")

		job.LogJob("px-server.log")
		job.RefreshJob()
		job.AliveJob("alive", server)
	}
	// 开启定时任务
	go cron.Start()

	return port, secret
}
