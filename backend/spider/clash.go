package spider

import (
	"github.com/metacubex/mihomo/log"
	"pandora-box/backend/constant"
	"sync"
)

func init() {
	Register(constant.CollectClash, NewClashCollect)
}

type Clash struct {
	Url     string
	Headers map[string]string
}

func (c *Clash) Get() []map[string]any {
	content := GetBytes(c.Url, c.Headers)
	return ComputeFuzzy(content, c.Headers)
}

func (c *Clash) Get2ChanWG(pc chan []map[string]any, wg *sync.WaitGroup) {
	defer wg.Done()
	nodes := c.Get()
	log.Infoln("STATISTIC: Clash count=%d url=%s", len(nodes), c.Url)
	if len(nodes) > 0 {
		pc <- nodes
	}
}

func NewClashCollect(g Getter) Collect {
	return &Clash{Url: g.Url, Headers: g.Headers}
}
