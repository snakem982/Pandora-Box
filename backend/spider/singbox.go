package spider

import (
	"github.com/metacubex/mihomo/common/convert"
	"github.com/metacubex/mihomo/log"
	"pandora-box/backend/constant"
	"sync"
)

func init() {
	Register(constant.CollectSingBox, NewSingBoxCollect)
}

type SingBox struct {
	Url     string
	Headers map[string]string
}

func (c *SingBox) Get() []map[string]any {
	proxies := make([]map[string]any, 0)

	all := GetBytes(c.Url, c.Headers)
	if all != nil {
		sing, err := convert.ConvertsSingBox(all)
		if err == nil && sing != nil {
			proxies = sing
		}
	}

	return proxies
}

func (c *SingBox) Get2ChanWG(pc chan []map[string]any, wg *sync.WaitGroup) {
	defer wg.Done()
	nodes := c.Get()
	log.Infoln("STATISTIC: SingBox count=%d url=%s", len(nodes), c.Url)
	if len(nodes) > 0 {
		pc <- nodes
	}
}

func NewSingBoxCollect(g Getter) Collect {
	return &SingBox{Url: g.Url, Headers: g.Headers}
}
