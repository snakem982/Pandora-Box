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
	Url string
}

func (c *Clash) Get() []map[string]any {
	return ComputeFuzzy(GetBytes(c.Url))
}

func (c *Clash) Get2ChanWG(pc chan []map[string]any, wg *sync.WaitGroup) {
	defer wg.Done()
	nodes := c.Get()
	log.Infoln("STATISTIC: Clash count=%d url=%s", len(nodes), c.Url)
	if len(nodes) > 0 {
		pc <- nodes
	}
}

func NewClashCollect(getter Getter) Collect {
	return &Clash{Url: getter.Url}
}
