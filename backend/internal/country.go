package internal

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/metacubex/mihomo/adapter"
	"github.com/metacubex/mihomo/component/mmdb"
	C "github.com/metacubex/mihomo/constant"
	"github.com/metacubex/mihomo/dns"
	"github.com/metacubex/mihomo/log"
	"github.com/snakem982/pandora-box/api/models"
	"github.com/snakem982/pandora-box/pkg/cache"
	"github.com/snakem982/pandora-box/pkg/constant"
	"github.com/snakem982/pandora-box/pkg/utils"
	"io"
	"net"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

//go:embed em/flags.json
var fsEmoji []byte

var emojiMap = make(map[string]string)

// 初始化国旗
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

func CleanRealIpCache() {
	var list []models.RealIp
	_ = cache.GetList(constant.RealIpHeader, &list)
	vLen := len(list)
	if vLen < 20480 {
		log.Infoln("real ip cache len is %d", vLen)
		return
	}

	log.Infoln("real ip cache len is %d,clean start...", vLen)
	m := make(map[string]any)
	for i := 0; i < vLen/2; i++ {
		key := list[i].Key
		m[key] = 1
	}
	_ = cache.DeleteList(m)
	log.Infoln("real ip cache len is %d,clean end.", vLen)
}

func getRealIpCountryCode(ctx context.Context, proxy map[string]any) (string, error) {
	var value models.RealIp
	realIpKey := fmt.Sprintf("%s%v:%v", constant.RealIpHeader, proxy["server"], proxy["port"])
	_ = cache.Get(realIpKey, &value)
	if value.CountryCode != "" {
		return value.CountryCode, nil
	}

	req, err := http.NewRequest(http.MethodGet, "http://ip-api.com/json/", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")
	req = req.WithContext(ctx)

	transport := &http.Transport{
		// from http.DefaultTransport
		DisableKeepAlives: true,
		DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
			proxy, _ := adapter.ParseProxy(proxy)
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
		Ip          string `json:"query"`
		CountryCode string `json:"countryCode"`
	}{}
	err = json.Unmarshal(body, &ipInfo)
	if err != nil {
		return "", err
	}

	_ = cache.Put(realIpKey, models.RealIp{
		Key:         realIpKey,
		CountryCode: ipInfo.CountryCode,
	})
	return ipInfo.CountryCode, nil
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

func GetCountryName(proxies []map[string]any, need bool) []map[string]any {
	mainServer, _ := dns.ParseNameServer([]string{
		"8.8.4.4",
		"180.76.76.76",
	})

	fallbackServer, _ := dns.ParseNameServer([]string{
		"1.1.1.1",
		"202.175.3.3",
	})

	cfg := dns.Config{
		Main:     mainServer,
		Fallback: fallbackServer,
	}

	if NowConfig != nil {
		c := NowConfig.DNS
		cfg = dns.Config{
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
	}

	r := dns.NewResolver(cfg)

	result := make([]map[string]any, 0)
	ipLock := sync.Mutex{}

	pool := utils.NewTimeoutPoolWithDefaults()
	pool.WaitCount(len(proxies))

	for _, proxy := range proxies {
		p := proxy
		p["name"] = "ZZ"

		pool.SubmitWithTimeout(func(done chan struct{}) {
			defer func() {
				if e := recover(); e != nil {
					log.Errorln("国家代码获取失败，错误：%v", e)
				}
				done <- struct{}{}
			}()

			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*3600)
			defer cancel()

			ipOrDomain := p["server"].(string)
			if utils.CheckStringAlphabet(ipOrDomain) {
				ips, err := r.LookupIP(ctx, ipOrDomain)
				if err == nil && ips != nil && len(ips) > 0 {
					p["name"] = getCountryCode(ips[0].String())
				}
			} else {
				p["name"] = getCountryCode(ipOrDomain)
			}

			if need {
				// 获取落地ip国家代码
				realIpCC, err := getRealIpCountryCode(ctx, p)
				if err == nil {
					p["name"] = realIpCC
				}
			}

			ipLock.Lock()
			result = append(result, p)
			ipLock.Unlock()

		}, 4*time.Second)
	}

	pool.StartAndWait()
	return result
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

		return utils.Reverse(proxies[i]["server"].(string)) < utils.Reverse(proxies[j]["server"].(string))
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
