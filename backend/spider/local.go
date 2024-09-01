package spider

import (
	"github.com/metacubex/mihomo/log"
	"pandora-box/backend/constant"
	"sync"
)

func init() {
	Register(constant.CollectLocal, NewLocalCollect)
}

type Local struct {
	Url     string
	Headers map[string]string
}

func (c *Local) Get() []map[string]any {
	return ComputeFuzzy([]byte(c.Url), c.Headers)
}

func (c *Local) Get2ChanWG(pc chan []map[string]any, wg *sync.WaitGroup) {
	defer wg.Done()
	nodes := c.Get()

	str := c.Url
	runes := []rune(str)
	if len(runes) > 128 {
		str = string(runes[:128]) + "..."
	}
	log.Infoln("STATISTIC: Local count=%d content=%s", len(nodes), str)

	if len(nodes) > 0 {
		pc <- nodes
	}
}

func NewLocalCollect(g Getter) Collect {
	return &Local{Url: g.Url, Headers: g.Headers}
}
