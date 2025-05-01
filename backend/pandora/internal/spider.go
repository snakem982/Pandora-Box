package internal

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/metacubex/mihomo/adapter"
	"github.com/metacubex/mihomo/common/convert"
	mu "github.com/metacubex/mihomo/common/utils"
	"github.com/metacubex/mihomo/config"
	"github.com/metacubex/mihomo/log"
	"github.com/snakem982/pandora-box/api/models"
	"github.com/snakem982/pandora-box/pkg/cache"
	"github.com/snakem982/pandora-box/pkg/constant"
	"github.com/snakem982/pandora-box/pkg/utils"
	"gopkg.in/yaml.v3"
	"regexp"
	"strings"
	"sync"
	"time"
)

// 分享链接
var shareLinkReg = regexp.MustCompile("(vless|vmess|trojan|ss|ssr|tuic|hysteria|hysteria2|hy2)://[-A-Za-z0-9\u4e00-\u9ea5+&@#/%?=~_!:,.;]+[-A-Za-z0-9\u4e00-\u9ea5+&@#/%=~_]")

// 订阅地址
var subReg = regexp.MustCompile("(https|http)://[-A-Za-z0-9\u4e00-\u9ea5+&@#/%?=~_!:,.;]+[-A-Za-z0-9\u4e00-\u9ea5+&@#/%=~_]")

// ScanShareLinks 扫描分享链接
func ScanShareLinks(content string) []string {
	return shareLinkReg.FindAllString(content, -1)
}

// ScanSubs 扫描订阅地址
func ScanSubs(content string) []string {
	return subReg.FindAllString(content, -1)
}

// Deduplicate 节点去重
func Deduplicate(proxies []map[string]any) []map[string]any {
	seen := make(map[string]bool) // 存储已经出现过的 key 值
	var result []map[string]any

	for _, proxy := range proxies {

		proxyType, existType := proxy["type"].(string)
		if !existType {
			continue
		}

		var (
			proxyId string
			err     error
		)
		server := proxy["server"]
		port := proxy["port"]
		password := proxy["password"]
		uuid := proxy["uuid"]
		switch proxyType {
		case "ss":
			proxyId = fmt.Sprintf("%s|%v|%v|%v", "ss", server, port, password)
		case "ssr":
			proxyId = fmt.Sprintf("%s|%v|%v|%v", "ssr", server, port, password)
		case "vmess":

			proxyId = fmt.Sprintf("%s|%v|%v|%v", "vmess", server, port, uuid)
		case "vless":
			flow, existFlow := proxy["flow"].(string)
			if existFlow && flow != "" && flow != "xtls-rprx-vision" {
				continue
			}
			proxyId = fmt.Sprintf("%s|%v|%v|%v", "vless", server, port, uuid)
		case "trojan":
			proxyId = fmt.Sprintf("%s|%v|%v|%v", "trojan", server, port, password)
		case "hysteria":
			authStr, exist := proxy["auth_str"]
			if !exist {
				authStr = proxy["auth-str"]
			}
			proxyId = fmt.Sprintf("%s|%v|%v|%v", "hysteria", server, port, authStr)
		case "hysteria2":
			proxyId = fmt.Sprintf("%s|%v|%v|%v", "hysteria2", server, port, password)
		case "wireguard":
			authStr := proxy["private-key"]
			proxyId = fmt.Sprintf("%s|%v|%v|%v", "wireguard", server, port, authStr)
		case "tuic":
			proxyId = fmt.Sprintf("%s|%v|%v|%v|%v", "tuic", server, port, uuid, password)
		case "socks5":
			username := proxy["username"]
			proxyId = fmt.Sprintf("%s|%v|%v|%v|%v", "socks5", server, port, username, password)
		case "mieru":
			username := proxy["username"]
			proxyId = fmt.Sprintf("%s|%v|%v|%v|%v", "mieru", server, port, username, password)
		case "http":
			username := proxy["username"]
			proxyId = fmt.Sprintf("%s|%v|%v|%v|%v", "http", server, port, username, password)
		case "anytls":
			proxyId = fmt.Sprintf("%s|%v|%v|%v", "anytls", server, port, password)
		default:
			err = fmt.Errorf("unsupport proxy type: %s", proxyType)
		}

		if err != nil {
			continue
		}

		// 如果 key 的值尚未出现，加入结果集
		if !seen[proxyId] {
			seen[proxyId] = true
			result = append(result, proxy)
		}
	}

	return result
}

var lock sync.Mutex
var nullValue models.Void

// ScanProxies 获取节点
func ScanProxies(content string, headers map[string]string, deep int) (proxies []map[string]any) {
	proxies = make([]map[string]any, 0)
	if deep > 2 {
		return
	}

	tempStr := strings.TrimSpace(content)

	// Sing解析
	if utils.IsJSON(tempStr) {
		sing, err := convert.ConvertsSingBox([]byte(tempStr))
		if err == nil && sing != nil {
			return sing
		}
	}

	// Base64解析
	if utils.IsBase64(tempStr) {
		v2ray, err := convert.ConvertsV2Ray([]byte(tempStr))
		if err == nil && v2ray != nil {
			return v2ray
		}
	}

	// 初始化urls
	var urls = make(map[string]models.Void)
	// 处理 ruleProvider
	var ruleProviderUrl = map[string]bool{
		"http://www.gstatic.com/generate_204":  true,
		"https://www.gstatic.com/generate_204": true,
		"https://www.google.com/blank.html":    true,
	}
	rawCfg, err := config.UnmarshalRawConfig([]byte(tempStr))
	if err == nil {
		for _, m := range rawCfg.ProxyProvider {
			if url, find := m["url"].(string); find && strings.HasPrefix(url, "http") {
				urls[url] = nullValue
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

		if len(urls) == 0 && len(rawCfg.Proxy) > 0 {
			return rawCfg.Proxy
		}
	}

	// 扫描分享链接
	shareLinks := ScanShareLinks(tempStr)
	var builder strings.Builder
	for _, link := range shareLinks {
		builder.WriteString(link + "\n")
	}
	if builder.Len() > 0 {
		v2ray, err := convert.ConvertsV2Ray([]byte(builder.String()))
		if err == nil && v2ray != nil {
			proxies = append(proxies, v2ray...)
		}
	}

	// 扫描URL
	if len(urls) == 0 {
		subs := ScanSubs(tempStr)
		for _, sub := range subs {
			if ruleProviderUrl[sub] {
				continue
			}
			urls[sub] = nullValue
		}
	}

	// 无订阅内容
	i := len(urls)
	if i == 0 {
		return
	}

	// 只有一个 url 直接请求
	if i == 1 {
		for url := range urls {
			Worker(url, &proxies, headers, deep+1)
		}
		return
	}

	// 进行订阅请求
	pool := utils.NewTimeoutPoolWithDefaults()
	pool.WaitCount(i)
	for url := range urls {
		pool.Submit(func(done chan struct{}) {
			defer func() {
				if err := recover(); err != nil {
					log.Errorln("爬取失败 URL= %s, 错误: %v", url, err)
				}
				done <- struct{}{}
			}()
			Worker(url, &proxies, headers, deep+1)
		})
	}
	pool.StartAndWait()

	return
}

// Worker 发起请求
func Worker(url string, proxies *[]map[string]any, headers map[string]string, deep int) {
	// 函数执行计时
	start := time.Now()
	var err error
	var resp *utils.ResponseResult
	defer func() {
		elapsed := time.Since(start)
		if err != nil {
			log.Errorln("请求失败 URL= %s, 耗时: %v，错误: %v", url, elapsed, err)
		}
	}()

	// 并发Get请求
	resp, err = utils.FastGet(url, headers, GetProxyUrl())
	if err != nil {
		return
	}
	if resp == nil {
		err = fmt.Errorf("响应为空")
		return
	}

	// 对返回内容扫描
	scanProxies := ScanProxies(resp.Body, headers, deep)
	if len(scanProxies) > 0 {
		lock.Lock()
		*proxies = append(*proxies, scanProxies...)
		lock.Unlock()
	}
}

// UrlTest 延迟测试
func UrlTest(proxies []map[string]any, testUrl string) []map[string]any {
	pool := utils.NewTimeoutPoolWithDefaults()
	result := make([]map[string]any, 0)
	m := sync.Mutex{}

	expectedStatus, _ := mu.NewUnsignedRanges[uint16]("200/204/301/302/304")
	url := testUrl
	if url == "" {
		url = "https://www.google.com/blank.html"
	}

	pool.WaitCount(len(proxies))
	for _, p := range proxies {
		proxy := p
		pool.SubmitWithTimeout(func(done chan struct{}) {
			defer func() {
				if e := recover(); e != nil {
					log.Errorln("延迟测试失败, 错误: %v", e)
				}
				done <- struct{}{}
			}()

			switch proxy["type"] {
			case "wireguard":
				return
			case "ss":
				delete(proxy, "dialer-proxy")
			default:
				delete(proxy, "dialer-proxy")
				proxy["skip-cert-verify"] = true
			}

			proxyAdapter, err := adapter.ParseProxy(proxy)
			if err != nil {
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*4500)
			defer cancel()
			pass := proxyAdapter.URLTestByPandora(ctx, url, expectedStatus)
			if pass {
				m.Lock()
				result = append(result, proxy)
				m.Unlock()
			}
		}, 5*time.Second)
	}
	pool.StartAndWait()

	return result
}

// CrawlProxy 进行节点爬取
func CrawlProxy(getter models.Getter) (proxies []map[string]any) {
	proxies = make([]map[string]any, 0)

	// 加载缓存的节点
	cachePath := utils.GetUserHomeDir(constant.DefaultCrawlDir, getter.Id+".yaml")
	content, err := utils.ReadFile(cachePath)
	if err != nil {
		yml := models.Yml{
			Proxies: proxies,
		}
		_ = yaml.Unmarshal([]byte(content), &yml)
		proxies = yml.Proxies
		getter.Cache = getter.Available
	}

	// 爬取新节点
	scanProxies := ScanProxies(getter.Content, getter.Headers, 0)
	getter.Crawl = len(scanProxies)
	proxies = append(proxies, scanProxies...)

	// 节点去重
	proxies = Deduplicate(proxies)

	// 连通性测试
	proxies = UrlTest(proxies, getter.TestUrl)

	// 国家代码查询
	proxies = GetCountryName(proxies, true)

	// 排序添加Emoji
	SortAddEmoji(proxies)

	// 存盘
	yml := models.Yml{
		Proxies: proxies,
	}
	out, _ := yaml.Marshal(yml)
	savePath := utils.GetUserHomeDir(constant.DefaultCrawlDir, getter.Id+".yaml")
	_, _ = utils.SaveFile(savePath, out)

	// 更新db
	getter.Available = len(proxies)
	getter.SetUpdateTime()
	_ = cache.Put(getter.Id, getter)

	// 清理realIp缓存
	go CleanRealIpCache()

	return
}
