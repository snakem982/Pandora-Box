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
	Getter
}

func (s *SingBox) Get() []map[string]any {
	proxies := make([]map[string]any, 0)

	all := GetBytes(s.Url, s.Headers)
	if all != nil {
		sing, err := convert.ConvertsSingBox(all)
		if err == nil && sing != nil {
			proxies = sing
		}
	}

	return proxies
}

func (s *SingBox) Get2ChanWG(pc chan []map[string]any, wg *sync.WaitGroup) {
	defer wg.Done()
	nodes := s.Get()
	log.Infoln("STATISTIC: SingBox count=%d url=%s", len(nodes), s.Url)
	AddIdAndUpdateGetter(pc, nodes, s.Getter)
}

func NewSingBoxCollect(g Getter) Collect {
	return &SingBox{g}
}
