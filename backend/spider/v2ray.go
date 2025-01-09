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
	Getter
}

func (v *V2ray) Get() []map[string]any {
	content := GetBytes(v.Url, v.Headers)
	return ComputeFuzzy(content, v.Headers)
}

func (v *V2ray) Get2ChanWG(pc chan []map[string]any, wg *sync.WaitGroup) {
	defer wg.Done()
	nodes := v.Get()
	log.Infoln("STATISTIC: V2ray count=%d url=%s", len(nodes), v.Url)
	AddIdAndUpdateGetter(pc, nodes, v.Getter)
}

func NewV2rayCollect(g Getter) Collect {
	return &V2ray{g}
}
