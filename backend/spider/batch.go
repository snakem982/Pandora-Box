package spider

import (
	"github.com/metacubex/mihomo/log"
	"pandora-box/backend/constant"
	"sync"
)

func init() {
	Register(constant.CollectLocal, NewBatchCollect)
	Register(constant.CollectBatch, NewBatchCollect)
}

type Batch struct {
	Getter
}

func (b *Batch) Get() []map[string]any {
	return ComputeFuzzy([]byte(b.Url), b.Headers)
}

func (b *Batch) Get2ChanWG(pc chan []map[string]any, wg *sync.WaitGroup) {
	defer wg.Done()
	nodes := b.Get()

	str := b.Url
	runes := []rune(str)
	if len(runes) > 128 {
		str = string(runes[:128]) + "..."
	}
	log.Infoln("STATISTIC: Batch|Local count=%d content=%s", len(nodes), str)
	AddIdAndUpdateGetter(pc, nodes, b.Getter)
}

func NewBatchCollect(g Getter) Collect {
	return &Batch{g}
}
