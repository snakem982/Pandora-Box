package spider

import (
	"github.com/metacubex/mihomo/common/convert"
	"github.com/metacubex/mihomo/log"
	"pandora-box/backend/constant"
	"regexp"
	"strings"
	"sync"
)

func init() {
	Register(constant.CollectSharelink, NewShareLinkCollect)
}

type ShareLink struct {
	Getter
}

func (l *ShareLink) Get() []map[string]any {
	proxies := make([]map[string]any, 0)

	all := GetBytes(l.Url, l.Headers)
	if all != nil {
		builder := strings.Builder{}
		for _, link := range grepShareLink(all) {
			builder.WriteString(link + "\n")
		}
		if builder.Len() > 0 {
			v2ray, err := convert.ConvertsV2Ray([]byte(builder.String()))
			if err == nil && v2ray != nil {
				proxies = v2ray
			}
		}
	}

	return proxies
}

func (l *ShareLink) Get2ChanWG(pc chan []map[string]any, wg *sync.WaitGroup) {
	defer wg.Done()
	nodes := l.Get()
	log.Infoln("STATISTIC: ShareLink count=%d url=%s", len(nodes), l.Url)
	AddIdAndUpdateGetter(pc, nodes, l.Getter)
}

func NewShareLinkCollect(g Getter) Collect {
	return &ShareLink{g}
}

var shareLinkReg = regexp.MustCompile("(vless|vmess|trojan|ss|ssr|tuic|hysteria|hysteria2|hy2)://([A-Za-z0-9+/_&?=@:%.-])+")

// grepShareLink
//
//	@Description: 爬取分享链接
//	@param all
//	@return []string
func grepShareLink(all []byte) []string {
	return shareLinkReg.FindAllString(string(all), -1)
}
