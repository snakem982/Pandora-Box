package spider

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/metacubex/mihomo/adapter"
	"github.com/metacubex/mihomo/common/utils"
	"github.com/metacubex/mihomo/component/mmdb"
	"github.com/metacubex/mihomo/config"
	C "github.com/metacubex/mihomo/constant"
	"github.com/metacubex/mihomo/dns"
	"github.com/metacubex/mihomo/log"
	"gopkg.in/yaml.v3"
	"math/rand"
	"net"
	"os"
	"pandora-box/backend/cache"
	"pandora-box/backend/constant"
	"pandora-box/backend/meta"
	"pandora-box/backend/mypool"
	"pandora-box/backend/premium"
	"pandora-box/backend/tools"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

var emojiMap = make(map[string]string)

//go:embed flags.json
var fsEmoji []byte

// 优选ip
var _rand = rand.New(rand.NewSource(time.Now().Unix()))
var CloudflareCIDR = premium.LoadCIDR(premium.CdnTypeCloudflare)

func init() {
	type countryEmoji struct {
		Code  string `json:"code"`
		Emoji string `json:"emoji"`
	}
	var countryEmojiList = make([]countryEmoji, 0)
	_ = json.Unmarshal(fsEmoji, &countryEmojiList)
	for _, i := range countryEmojiList {
		emojiMap[i.Code] = i.Emoji
	}
}

func Crawl() bool {
	// 加载默认配置中的节点
	defaultBuf, defaultErr := os.ReadFile(filepath.Join(C.Path.HomeDir(), "uploads/"+constant.DefaultProfile+".yaml"))
	proxies := make([]map[string]any, 0)
	if defaultErr == nil && len(defaultBuf) > 0 {
		rawCfg, err := config.UnmarshalRawConfig(defaultBuf)
		if err == nil && len(rawCfg.Proxy) > 0 {
			proxies = rawCfg.Proxy
			log.Infoln("load default config proxies success %d", len(rawCfg.Proxy))
		}
	}

	// 获取getters
	getters := make([]Getter, 0)
	values := cache.GetList(constant.PrefixGetter)
	if len(values) > 0 {
		for _, value := range values {
			getter := Getter{}
			_ = json.Unmarshal(value, &getter)
			getters = append(getters, getter)
		}
	}

	// 进行抓取
	if len(getters) > 0 {
		wg := &sync.WaitGroup{}
		var pc = make(chan []map[string]any)
		for _, g := range getters {
			collect, err := NewCollect(g.Type, g)
			if err != nil {
				continue
			}
			wg.Add(1)
			go collect.Get2ChanWG(pc, wg)
		}
		go func() {
			wg.Wait()
			close(pc)
		}()
		for p := range pc {
			if p != nil {
				proxies = append(proxies, p...)
			}
		}
	}

	// 优选ip
	isNeedPremium := false
	for _, proxy := range proxies {
		server := proxy["server"].(string)
		if premium.IsCdnIp(CloudflareCIDR, server) {
			isNeedPremium = true
			break
		}
	}
	if isNeedPremium {
		httpsIps := premium.GetExcellentIps(premium.CdnTypeCloudflare).HttpsIps
		if len(httpsIps) > 0 {
			_rand.Seed(time.Now().UnixNano())
			httpsIpsLen := len(httpsIps) - 1
			for _, proxy := range proxies {
				proxyCopy := proxy
				server := proxy["server"].(string)
				if premium.IsCdnIp(CloudflareCIDR, server) {
					c1 := make(map[string]any)
					c2 := make(map[string]any)
					marshal, _ := json.Marshal(proxyCopy)
					_ = json.Unmarshal(marshal, &c1)
					_ = json.Unmarshal(marshal, &c2)
					x := _rand.Intn(httpsIpsLen)
					c1["server"] = httpsIps[x]
					x = _rand.Intn(httpsIpsLen)
					c2["server"] = httpsIps[x]
					proxies = append(proxies, c1)
					proxies = append(proxies, c2)
				}
			}
		}
	}

	// 去重
	maps := Unique(proxies, true)
	if len(maps) == 0 {
		return false
	}

	// 转换
	nodes := map2proxies(maps)
	if len(nodes) == 0 {
		return false
	}

	// url测速
	keys := urlTest(nodes)
	if len(keys) == 0 {
		return false
	}

	// 国家代码查询
	proxies = GetCountryName(keys, maps)

	// 排序添加emoji
	SortAddEmoji(proxies)

	if len(proxies) > 255 {
		proxies = proxies[0:256]
	}

	// 存盘
	data := make(map[string]any)
	data["proxies"] = proxies
	all, _ := yaml.Marshal(data)
	filePath := C.Path.HomeDir() + "/uploads/" + constant.PrefixProfile + "0.yaml"
	_ = os.Remove(filePath)
	_ = os.WriteFile(filePath, all, 0666)

	return true
}

func Unique(mappings []map[string]any, needTls bool) (maps map[string]map[string]any) {

	maps = make(map[string]map[string]any)

	for _, mapping := range mappings {
		proxyType, existType := mapping["type"].(string)
		if !existType {
			continue
		}

		var (
			proxyId string
			err     error
		)
		server := mapping["server"]
		port := mapping["port"]
		password := mapping["password"]
		uuid := mapping["uuid"]
		switch proxyType {
		case "ss":
			proxyId = fmt.Sprintf("%s|%v|%v|%v", "ss", server, port, password)
		case "ssr":
			proxyId = fmt.Sprintf("%s|%v|%v|%v", "ssr", server, port, password)
		case "vmess":
			if needTls {
				tls, existTls := mapping["tls"].(bool)
				if !existTls || !tls {
					continue
				}
			}
			proxyId = fmt.Sprintf("%s|%v|%v|%v", "vmess", server, port, uuid)
		case "vless":
			if needTls {
				tls, existTls := mapping["tls"].(bool)
				if !existTls || !tls {
					continue
				}
			}
			flow, existFlow := mapping["flow"].(string)
			if existFlow && flow != "" && flow != "xtls-rprx-vision" {
				continue
			}
			proxyId = fmt.Sprintf("%s|%v|%v|%v", "vless", server, port, uuid)
		case "trojan":
			if needTls {
				_, existSni := mapping["sni"].(string)
				if !existSni {
					continue
				}
			}
			proxyId = fmt.Sprintf("%s|%v|%v|%v", "trojan", server, port, password)
		case "hysteria":
			authStr, exist := mapping["auth_str"]
			if !exist {
				authStr = mapping["auth-str"]
			}
			proxyId = fmt.Sprintf("%s|%v|%v|%v", "hysteria", server, port, authStr)
		case "hysteria2":
			proxyId = fmt.Sprintf("%s|%v|%v|%v", "hysteria2", server, port, password)
		case "wireguard":
			authStr := mapping["private-key"]
			proxyId = fmt.Sprintf("%s|%v|%v|%v", "wireguard", server, port, authStr)
		case "tuic":
			proxyId = fmt.Sprintf("%s|%v|%v|%v|%v", "tuic", server, port, uuid, password)
		default:
			err = fmt.Errorf("unsupport proxy type: %s", proxyType)
		}

		if err != nil {
			continue
		}
		temp := mapping
		temp["name"] = proxyId
		maps[proxyId] = temp
	}

	return
}

func map2proxies(maps map[string]map[string]any) (proxies []C.Proxy) {
	pool := mypool.NewTimeoutPoolWithDefaults()
	pool.WaitCount(len(maps))
	mutex := sync.Mutex{}

	proxies = make([]C.Proxy, 0)
	for _, m := range maps {
		proxy := m
		pool.SubmitWithTimeout(func(done chan struct{}) {
			defer func() {
				if e := recover(); e != nil {
					log.Errorln("===map2proxies===%s", e)
				}
				done <- struct{}{}
			}()
			proxyT, err := adapter.ParseProxy(proxy)
			if err == nil {
				mutex.Lock()
				proxies = append(proxies, proxyT)
				mutex.Unlock()
			}
		}, 2*time.Second)
	}
	pool.StartAndWait()

	return
}

func urlTest(proxies []C.Proxy) []string {
	pool := mypool.NewTimeoutPoolWithDefaults()
	keys := make([]string, 0)
	m := sync.Mutex{}

	expectedStatus, _ := utils.NewUnsignedRanges[uint16]("200/204/301/302")
	url := "https://www.gstatic.com/generate_204"

	pool.WaitCount(len(proxies))
	for _, p := range proxies {
		proxy := p
		pool.SubmitWithTimeout(func(done chan struct{}) {
			defer func() {
				if e := recover(); e != nil {
					log.Errorln("===urlTest===%s", e)
				}
				done <- struct{}{}
			}()

			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*4500)
			defer cancel()
			_, err := proxy.URLTest(ctx, url, expectedStatus)
			if err == nil {
				m.Lock()
				keys = append(keys, proxy.Name())
				m.Unlock()
			}
		}, 5*time.Second)
	}
	pool.StartAndWait()

	return keys
}

func GetCountryName(keys []string, maps map[string]map[string]any) []map[string]any {
	c := meta.NowConfig.DNS
	cfg := dns.Config{
		Main:         c.NameServer,
		Fallback:     c.Fallback,
		IPv6:         false,
		IPv6Timeout:  c.IPv6Timeout,
		EnhancedMode: c.EnhancedMode,
		Pool:         c.FakeIPRange,
		Hosts:        c.Hosts,
		FallbackFilter: dns.FallbackFilter{
			GeoIP:     c.FallbackFilter.GeoIP,
			GeoIPCode: c.FallbackFilter.GeoIPCode,
			IPCIDR:    c.FallbackFilter.IPCIDR,
			Domain:    c.FallbackFilter.Domain,
			GeoSite:   c.FallbackFilter.GeoSite,
		},
		Default:     c.DefaultNameserver,
		Policy:      c.NameServerPolicy,
		ProxyServer: c.ProxyServerNameserver,
	}

	r := dns.NewResolver(cfg)

	proxies := make([]map[string]any, 0)
	ipLock := sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(len(keys))
	for _, key := range keys {
		m := maps[key]
		m["name"] = "ZZ"
		go func() {
			defer wg.Done()

			ipOrDomain := m["server"].(string)
			if tools.CheckStringAlphabet(ipOrDomain) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
				defer cancel()
				ips, err := r.LookupIP(ctx, ipOrDomain)
				if err == nil && ips != nil {
					m["name"] = getCountryCode(ips[0].String())
				}
			} else {
				m["name"] = getCountryCode(ipOrDomain)
			}
			ipLock.Lock()
			proxies = append(proxies, m)
			ipLock.Unlock()
		}()
	}

	wg.Wait()
	return proxies
}

func getCountryCode(ip string) string {
	netIp := net.ParseIP(ip)
	codes := mmdb.IPInstance().LookupCode(netIp)
	codesLen := len(codes)
	if codesLen == 0 {
		return "ZZ"
	}
	code := codes[codesLen-1]
	if len(code) == 2 {
		return strings.ToUpper(code)
	}

	return code
}

func getIndex(at string) int {
	switch at {
	case "hysteria2":
		return 1
	case "hysteria":
		return 2
	case "tuic":
		return 3
	case "ss":
		return 5
	case "vless":
		return 6
	default:
		return 10
	}
}

func SortAddEmoji(proxies []map[string]any) {
	sort.Slice(proxies, func(i, j int) bool {
		iProtocol := proxies[i]["type"].(string)
		jProtocol := proxies[j]["type"].(string)

		if getIndex(iProtocol) != getIndex(jProtocol) {
			return getIndex(iProtocol) < getIndex(jProtocol)
		}

		if proxies[i]["name"].(string) != proxies[j]["name"].(string) {
			return proxies[i]["name"].(string) < proxies[j]["name"].(string)
		}

		return tools.Reverse(proxies[i]["server"].(string)) < tools.Reverse(proxies[j]["server"].(string))
	})

	for i, _ := range proxies {
		name := proxies[i]["name"].(string)
		name = fmt.Sprintf("%s %s_%+02v", emojiMap[name], name, i+1)
		proxies[i]["name"] = strings.TrimSpace(name)
	}
}

func SortAddIndex(proxies []map[string]any) []map[string]any {
	// 去重
	maps := Unique(proxies, false)
	keys := make([]string, 0)
	for k := range maps {
		keys = append(keys, k)
	}
	// 国家代码查询
	proxies = GetCountryName(keys, maps)
	// 排序添加emoji
	SortAddEmoji(proxies)

	return proxies
}
