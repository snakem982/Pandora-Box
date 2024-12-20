package spider

import (
	"github.com/metacubex/mihomo/common/convert"
	"github.com/metacubex/mihomo/config"
	"github.com/metacubex/mihomo/log"
	"pandora-box/backend/constant"
	"pandora-box/backend/mypool"
	"regexp"
	"strings"
	"sync"
	"time"
)

func init() {
	Register(constant.CollectFuzzy, NewFuzzyCollect)
	Register(constant.CollectAuto, NewFuzzyCollect)
}

type Fuzzy struct {
	Url     string
	Headers map[string]string
}

func (c *Fuzzy) Get() []map[string]any {
	content := GetBytes(c.Url, c.Headers)
	return ComputeFuzzy(content, c.Headers)
}

func (c *Fuzzy) Get2ChanWG(pc chan []map[string]any, wg *sync.WaitGroup) {
	defer wg.Done()
	nodes := c.Get()
	log.Infoln("STATISTIC: Auto|Fuzzy count=%d url=%s", len(nodes), c.Url)
	if len(nodes) > 0 {
		pc <- nodes
	}
}

func NewFuzzyCollect(g Getter) Collect {
	return &Fuzzy{Url: g.Url, Headers: g.Headers}
}

type void struct{}

var nullValue void

var re = regexp.MustCompile(`proxies|api|clash|Clash|v2ray|token|raw|subscribe|txt|yaml|yml|sub|uuid`)
var not = regexp.MustCompile(`svg|png|mp4|mp3|jpg|jpeg|m3u8|flv|gif|icon|ktv|mov|webcam`)
var urlRe = regexp.MustCompile("(https|http)://[-A-Za-z0-9\u4e00-\u9ea5+&@#/%?=~_!:,.;]+[-A-Za-z0-9\u4e00-\u9ea5+&@#/%=~_]")

func grepFuzzy(all []byte, proxyProviderUrl map[string]void, ruleProviderUrl map[string]bool) map[string]void {
	set := proxyProviderUrl

	subUrls := urlRe.FindAllString(string(all), -1)
	for _, url := range subUrls {
		if !re.MatchString(url) || not.MatchString(url) || ruleProviderUrl[url] {
			continue
		}
		set[url] = nullValue
	}

	return set
}

func ComputeFuzzy(content []byte, headers map[string]string) []map[string]any {

	proxies := make([]map[string]any, 0)
	if content == nil {
		return proxies
	}

	// 处理 proxyProviderUrl
	var proxyProviderUrl = make(map[string]void)
	// 处理 ruleProvider
	var ruleProviderUrl = make(map[string]bool)

	// 尝试clash解析 成功返回
	rawCfg, err := config.UnmarshalRawConfig(content)
	if err == nil {
		for _, m := range rawCfg.ProxyProvider {
			if _, find := m["url"]; !find {
				continue
			}
			s := m["url"].(string)
			if strings.HasPrefix(s, "http") {
				proxyProviderUrl[s] = nullValue
			}
		}

		for _, m := range rawCfg.RuleProvider {
			if _, find := m["url"]; !find {
				continue
			}
			s := m["url"].(string)
			if strings.HasPrefix(s, "http") {
				ruleProviderUrl[s] = true
			}
		}

		proxies = rawCfg.Proxy
	}

	if len(proxyProviderUrl) == 0 && len(proxies) == 0 {
		// 尝试v2ray解析 成功返回
		v2ray, err := convert.ConvertsV2Ray(content)
		if err == nil && v2ray != nil {
			return v2ray
		}

		// 尝试sing解析 成功返回
		sing, err := convert.ConvertsSingBox(content)
		if err == nil && sing != nil {
			return sing
		}
	}

	// 进行订阅爬取
	fuzzy := grepFuzzy(content, proxyProviderUrl, ruleProviderUrl)
	pool := mypool.NewTimeoutPoolWithDefaults()
	pool.WaitCount(len(fuzzy))

	var cFlag = regexp.MustCompile(`proxies|provider|clash|yaml|yml`)
	lock := sync.Mutex{}
	for temp := range fuzzy {
		url := temp
		pool.SubmitWithTimeout(func(done chan struct{}) {
			defer func() {
				// 每个协程内部使用recover捕获可能在调用逻辑中发生的panic
				if e := recover(); e != nil {
					log.Errorln("====爬取错误====%s", e)
				}
				done <- struct{}{}
			}()

			getter := Getter{Url: url, Headers: headers}
			var ok []map[string]any
			if cFlag.MatchString(url) {
				all := GetBytes(url, headers)
				if all != nil {
					rawCfgInner, err := config.UnmarshalRawConfig(all)
					if err == nil && rawCfgInner.Proxy != nil {
						ok = rawCfgInner.Proxy
					}
				}
			} else if strings.Contains(url, "README.md") {
				collect, _ := NewCollect(constant.CollectSharelink, getter)
				ok = collect.Get()
			} else {
				all := GetBytes(url, headers)
				if all == nil || len(all) < 32 {
					return
				}
				isGo := true
				rawCfgInner, err := config.UnmarshalRawConfig(all)
				if err == nil && len(rawCfgInner.Proxy) > 0 {
					ok = rawCfgInner.Proxy
					isGo = false
				}
				if isGo {
					v2ray, err := convert.ConvertsV2Ray(all)
					if err == nil && v2ray != nil {
						ok = v2ray
						isGo = false
					}
				}
				if isGo {
					sing, err := convert.ConvertsSingBox(all)
					if err == nil && sing != nil {
						ok = sing
						isGo = false
					}
				}
				if isGo {
					builder := strings.Builder{}
					for _, link := range grepShareLink(all) {
						builder.WriteString(link + "\n")
					}
					if builder.Len() > 0 {
						all = []byte(builder.String())
						v2ray, err := convert.ConvertsV2Ray(all)
						if err == nil && v2ray != nil {
							ok = v2ray
						}
					}
				}
			}

			if ok != nil && len(ok) > 0 {
				lock.Lock()
				proxies = append(proxies, ok...)
				lock.Unlock()
			}
		}, time.Minute)
	}
	pool.StartAndWait()

	// 进行分享链接爬取
	builder := strings.Builder{}
	for _, link := range grepShareLink(content) {
		builder.WriteString(link + "\n")
	}
	if builder.Len() > 0 {
		content = []byte(builder.String())
		v2ray, err := convert.ConvertsV2Ray(content)
		if err != nil && v2ray != nil {
			lock.Lock()
			proxies = append(proxies, v2ray...)
			lock.Unlock()
		}
	}

	return proxies
}
