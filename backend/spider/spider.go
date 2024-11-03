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
	"io"
	"math/rand"
	"net"
	"net/http"
	"os"
	"pandora-box/backend/cache"
	"pandora-box/backend/constant"
	"pandora-box/backend/meta"
	"pandora-box/backend/mypool"
	"pandora-box/backend/premium"
	"pandora-box/backend/tools"
	"path/filepath"
	"runtime"
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
	defaultBuf, defaultErr := GetNodesCache()
	proxies := make([]map[string]any, 0)
	if defaultErr == nil && len(defaultBuf) > 0 {
		rawCfg, err := config.UnmarshalRawConfig(defaultBuf)
		if err == nil && len(rawCfg.Proxy) > 0 {
			proxies = rawCfg.Proxy
			log.Infoln("load default config proxies success %d", len(rawCfg.Proxy))
		}
	}

	// 低于节点阈值进行爬取
	if len(proxies) < 1024 {
		proxies = append(doCrawl(), proxies...)
	}

	// 去重
	maps := Unique(proxies, false)
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
	proxies = GetCountryName(keys, maps, true)

	// 排序添加emoji
	SortAddEmoji(proxies)

	// 放入缓存
	Save2Local(proxies, "0_cache.yaml")

	if len(proxies) > 255 {
		proxies = proxies[0:256]
	}

	// 存盘
	Save2Local(proxies, "0.yaml")

	// 清理realIp缓存
	go cleanRealIpCache()

	// 垃圾回收
	go func() {
		time.Sleep(2 * time.Minute)
		runtime.GC()
	}()

	return true
}

func doCrawl() []map[string]any {
	proxies := make([]map[string]any, 0)

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

	// 进行爬取
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
		if proxy["server"] == nil {
			continue
		}
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
				if proxy["server"] == nil {
					continue
				}
				proxyCopy := proxy
				server := proxy["server"].(string)
				if premium.IsCdnIp(CloudflareCIDR, server) {
					c1 := make(map[string]any)
					marshal, _ := json.Marshal(proxyCopy)
					_ = json.Unmarshal(marshal, &c1)
					x := _rand.Intn(httpsIpsLen)
					c1["server"] = httpsIps[x]
					proxies = append(proxies, c1)
				}
			}
		}
	}

	return proxies
}

func Save2Local(proxies []map[string]any, fileName string) {
	data := make(map[string]any)
	data["proxies"] = proxies
	all, _ := yaml.Marshal(data)
	filePath := C.Path.HomeDir() + "/uploads/" + constant.PrefixProfile + fileName
	_ = os.Remove(filePath)
	_ = os.WriteFile(filePath, all, 0777)
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
		case "socks5":
			username := mapping["username"]
			proxyId = fmt.Sprintf("%s|%v|%v|%v|%v", "socks5", server, port, username, password)
		case "http":
			username := mapping["username"]
			proxyId = fmt.Sprintf("%s|%v|%v|%v|%v", "http", server, port, username, password)
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

func cleanRealIpCache() {
	values := cache.GetList(constant.RealIpHeader)
	vLen := len(values)
	if vLen < 20480 {
		log.Infoln("real ip cache len is %d", vLen)
		return
	}
	log.Infoln("real ip cache len is %d,clean start...", vLen)
	m := make(map[string]any)
	for i := 0; i < vLen/2; i++ {
		key := string(values[i])
		m[key] = 1
	}
	_ = cache.DeleteList(m)
	log.Infoln("real ip cache len is %d,clean end.", vLen)
}

func getRealIpCountryCode(ctx context.Context, m map[string]any) (string, error) {
	realIpKey := fmt.Sprintf("%s%v:%v", constant.RealIpHeader, m["server"], m["port"])
	value := cache.Get(realIpKey)
	if string(value) != "" {
		return string(value), nil
	}

	req, err := http.NewRequest(http.MethodGet, "http://ip-api.com/json/", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")
	req = req.WithContext(ctx)

	transport := &http.Transport{
		// from http.DefaultTransport
		DisableKeepAlives:     true,
		MaxIdleConns:          32,
		IdleConnTimeout:       16 * time.Second,
		TLSHandshakeTimeout:   8 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
			proxy, _ := adapter.ParseProxy(m)
			addr := C.Metadata{
				Host:    "ip-api.com",
				DstPort: 80,
			}
			return proxy.DialContext(ctx, &addr)
		},
	}

	client := http.Client{Transport: transport}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("get real ip error, code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	ipInfo := struct {
		Ip      string `json:"query"`
		Country string `json:"countryCode"`
	}{}
	err = json.Unmarshal(body, &ipInfo)
	if err != nil {
		return "", err
	}

	_ = cache.Put(realIpKey, []byte(ipInfo.Country))
	return ipInfo.Country, nil
}

func GetCountryName(keys []string, maps map[string]map[string]any, need bool) []map[string]any {
	c := meta.NowConfig.DNS
	cfg := dns.Config{
		Main:         c.NameServer,
		Fallback:     c.Fallback,
		IPv6:         false,
		IPv6Timeout:  c.IPv6Timeout,
		EnhancedMode: c.EnhancedMode,
		Pool:         c.FakeIPRange,
		Hosts:        c.Hosts,
		Default:      c.DefaultNameserver,
		Policy:       c.NameServerPolicy,
		ProxyServer:  c.ProxyServerNameserver,
	}

	r := dns.NewResolver(cfg)

	proxies := make([]map[string]any, 0)
	ipLock := sync.Mutex{}

	pool := mypool.NewTimeoutPoolWithDefaults()
	pool.WaitCount(len(keys))

	for _, key := range keys {
		m := maps[key]
		m["name"] = "ZZ"

		pool.SubmitWithTimeout(func(done chan struct{}) {
			defer func() {
				if e := recover(); e != nil {
					log.Errorln("===GetCountryName===%s", e)
				}
				done <- struct{}{}
			}()

			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*3600)
			defer cancel()

			ipOrDomain := m["server"].(string)
			if tools.CheckStringAlphabet(ipOrDomain) {
				ips, err := r.LookupIP(ctx, ipOrDomain)
				if err == nil && ips != nil && len(ips) > 0 {
					m["name"] = getCountryCode(ips[0].String())
				}
			} else {
				m["name"] = getCountryCode(ipOrDomain)
			}

			if need {
				// 获取落地ip国家代码
				realIpCC, err := getRealIpCountryCode(ctx, m)
				if err == nil {
					m["name"] = realIpCC
				}
			}

			ipLock.Lock()
			proxies = append(proxies, m)
			ipLock.Unlock()

		}, 4*time.Second)
	}

	pool.StartAndWait()
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

func SortProxies(proxies []map[string]any) {
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
}

func SortAddEmoji(proxies []map[string]any) {
	SortProxies(proxies)

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
	proxies = GetCountryName(keys, maps, false)
	// 排序添加emoji
	SortAddEmoji(proxies)

	return proxies
}

func GetNodesCache() ([]byte, error) {
	return os.ReadFile(filepath.Join(C.Path.HomeDir(), "uploads/"+constant.PrefixProfile+"0_cache.yaml"))
}
