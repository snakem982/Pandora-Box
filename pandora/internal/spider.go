package internal

import (
	_ "embed"
	"fmt"
	"github.com/metacubex/mihomo/common/convert"
	"github.com/sagernet/sing/common/json"
	"github.com/snakem982/pandora-box/pandora/api/models"
	"gopkg.in/yaml.v3"
	"regexp"
	"strings"
)

//go:embed flags.json
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

// 分享链接
var shareLinkReg = regexp.MustCompile("(vless|vmess|trojan|ss|ssr|tuic|hysteria|hysteria2|hy2)://([A-Za-z0-9+/_&?=@:%.-])+")

// 订阅地址
var subReg = regexp.MustCompile("(https|http)://[-A-Za-z0-9\u4e00-\u9ea5+&@#/%?=~_!:,.;]+[-A-Za-z0-9\u4e00-\u9ea5+&@#/%=~_]")

// ScanShareLink 扫描分享链接
func ScanShareLink(content string) []string {
	return shareLinkReg.FindAllString(content, -1)
}

// ScanSub 扫描订阅地址
func ScanSub(content string) []string {
	return subReg.FindAllString(content, -1)
}

// IsYAML 判断是否为 yaml
func IsYAML(data string) bool {
	var yml map[string]interface{}
	return yaml.Unmarshal([]byte(data), &yml) == nil
}

// Parse 解析节点
func Parse(data string) (proxies []map[string]any) {
	// 返回节点
	proxies = make([]map[string]any, 0)

	// 清除首尾空格
	tempStr := strings.TrimSpace(data)

	// 进行 sing 解析
	if strings.HasPrefix(tempStr, "{") {
		sing, err := convert.ConvertsSingBox([]byte(tempStr))
		if err == nil && sing != nil {
			proxies = sing
			return
		}
	}

	// 进行 yaml 解析
	if IsYAML(data) {
		yml := struct {
			proxies []map[string]any
		}{
			proxies: proxies,
		}
		err := yaml.Unmarshal([]byte(data), &yml)
		if err == nil {
			proxies = yml.proxies
			return
		}
	}

	// 进行 base64 解析
	v2ray, err := convert.ConvertsV2Ray([]byte(tempStr))
	if err == nil && v2ray != nil {
		proxies = v2ray
	}

	return proxies
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

// CrawlProxy 进行节点爬取
func CrawlProxy(getter models.Getter) []map[string]any {

	return nil
}
