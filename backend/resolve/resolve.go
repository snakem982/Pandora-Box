package resolve

import (
	_ "embed"
	"encoding/json"
	"github.com/metacubex/mihomo/adapter"
	"github.com/metacubex/mihomo/log"
	"pandora-box/backend/mypool"
	"sync"
	"time"
)

//go:embed config.yaml
var PandoraDefaultConfig []byte

//go:embed config_download.yaml
var PandoraDefaultDownloadConfig []byte

var PandoraDefaultPlace = "{{PANDORA-BOX}}"

// MapsToProxies 将任意数量的 map[string]any 切片转换为任意数量的 map[string]any 切片，
// 仅包含通过 adapter.ParseProxy 解析成功的元素。
func MapsToProxies(ray []map[string]any) []map[string]any {
	pool := mypool.NewTimeoutPoolWithDefaults()
	pool.WaitCount(len(ray))
	mutex := sync.Mutex{}

	proxies := make([]map[string]any, 0)
	for _, m := range ray {
		proxy := m
		pool.SubmitWithTimeout(func(done chan struct{}) {
			defer func() {
				if e := recover(); e != nil {
					log.Errorln("===MapsToProxies===%s", e)
				}
				done <- struct{}{}
			}()
			_, err := adapter.ParseProxy(proxy)
			if err == nil {
				mutex.Lock()
				proxies = append(proxies, proxy)
				mutex.Unlock()
			} else {
				marshal, err2 := json.Marshal(proxy)
				if err2 == nil {
					log.Warnln("===MapsToProxies=== proxy: %s ,err: %s", string(marshal), err.Error())
				}
			}
		}, 2*time.Second)
	}
	pool.StartAndWait()

	return proxies
}
