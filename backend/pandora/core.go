package pandora

import (
	"fmt"
	"github.com/metacubex/mihomo/hub/route"
	"github.com/metacubex/mihomo/log"
	"github.com/snakem982/pandora-box/api/handlers"
	"github.com/snakem982/pandora-box/internal"
	"github.com/snakem982/pandora-box/pkg/cache"
	"github.com/snakem982/pandora-box/pkg/constant"
	"github.com/snakem982/pandora-box/pkg/utils"
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
	Init(isClient)

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
	internal.SwitchProfile()

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

	return port, secret
}
