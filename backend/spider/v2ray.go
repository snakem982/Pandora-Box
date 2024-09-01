package spider

import (
	"github.com/metacubex/mihomo/log"
	"pandora-box/backend/constant"
	"sync"
)

func init() {
	Register(constant.CollectV2ray, NewV2rayCollect)
}

type V2ray struct {
	Url     string
	Headers map[string]string
}

func (c *V2ray) Get() []map[string]any {
	content := GetBytes(c.Url, c.Headers)
	return ComputeFuzzy(content, c.Headers)
}

func (c *V2ray) Get2ChanWG(pc chan []map[string]any, wg *sync.WaitGroup) {
	defer wg.Done()
	nodes := c.Get()
	log.Infoln("STATISTIC: V2ray count=%d url=%s", len(nodes), c.Url)
	if len(nodes) > 0 {
		pc <- nodes
	}
}

func NewV2rayCollect(g Getter) Collect {
	return &V2ray{Url: g.Url, Headers: g.Headers}
}
