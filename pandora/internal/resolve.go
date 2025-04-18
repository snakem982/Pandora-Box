package internal

import (
	"encoding/json"
	"fmt"
	"github.com/metacubex/mihomo/adapter"
	"github.com/metacubex/mihomo/common/convert"
	"github.com/metacubex/mihomo/config"
	"github.com/metacubex/mihomo/log"
	"github.com/snakem982/pandora-box/pandora/api/models"
	"github.com/snakem982/pandora-box/pandora/pkg/constant"
	"github.com/snakem982/pandora-box/pandora/pkg/utils"
	"gopkg.in/yaml.v3"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 保存文件
func saveProfile(proxies []map[string]any, profile *models.Profile) {
	yml := models.Yml{Proxies: proxies}
	out, _ := yaml.Marshal(yml)
	savePath := utils.GetUserHomeDir(profile.Path)
	_, _ = utils.SaveFile(savePath, out)
}

// MapsToProxies 将任意数量的 map[string]any 切片转换为任意数量的 map[string]any 切片，
// 仅包含通过 adapter.ParseProxy 解析成功的元素。
func MapsToProxies(ray []map[string]any) []map[string]any {
	pool := utils.NewTimeoutPoolWithDefaults()
	pool.WaitCount(len(ray))
	mutex := sync.Mutex{}

	proxies := make([]map[string]any, 0)
	for _, m := range ray {
		proxy := m
		pool.SubmitWithTimeout(func(done chan struct{}) {
			defer func() {
				if e := recover(); e != nil {
					log.Errorln("[MapsToProxies] Error:%v", e)
				}
				done <- struct{}{}
			}()
			proxy["skip-cert-verify"] = true
			_, err := adapter.ParseProxy(proxy)
			if err == nil {
				mutex.Lock()
				proxies = append(proxies, proxy)
				mutex.Unlock()
			} else {
				marshal, err2 := json.Marshal(proxy)
				if err2 == nil {
					log.Warnln("[MapsToProxies] proxy: %s ,err: %s", string(marshal), err.Error())
				}
			}
		}, 2*time.Second)
	}
	pool.StartAndWait()

	return proxies
}

// Resolve 解析内容，保存成 profile
func Resolve(content string, profile *models.Profile, refresh bool) error {
	// 解析内容预处理
	tempStr := strings.TrimSpace(content)
	tempBytes := []byte(tempStr)

	// 如果不是刷新则创建 id
	if !refresh {
		snowflakeId := utils.SnowflakeId()
		profile.Id = fmt.Sprintf("%s%d", constant.PrefixProfile, snowflakeId)
		profile.Order = snowflakeId
		profile.Path = "./profiles/" + profile.Id + ".yaml"
	}

	// Sing解析
	if utils.IsJSON(tempStr) {
		sing, err := convert.ConvertsSingBox(tempBytes)
		if err == nil {
			saveProfile(sing, profile)
			return nil
		}

		return err
	}

	// Base64解析
	if utils.IsBase64(tempStr) {
		v2ray, err := convert.ConvertsV2Ray(tempBytes)
		if err == nil {
			saveProfile(v2ray, profile)
			return nil
		}

		return err
	}

	// 分享链接解析
	shareLinks := ScanShareLinks(tempStr)
	var builder strings.Builder
	for _, link := range shareLinks {
		builder.WriteString(link + "\n")
	}
	if builder.Len() > 0 {
		share, err := convert.ConvertsV2Ray([]byte(builder.String()))
		if err == nil {
			saveProfile(share, profile)
			return nil
		}

		return err
	}

	// Yaml解析
	rawCfg, err := config.UnmarshalRawConfig(tempBytes)
	if err == nil {
		_, yamlError := config.ParseRawConfig(rawCfg)
		if yamlError != nil {
			// 配置校验失败，尝试提取可用节点
			rails := MapsToProxies(rawCfg.Proxy)
			if len(rails) == 0 {
				return yamlError
			} else {
				saveProfile(rails, profile)
				return nil
			}
		}

		// 保存yaml
		if len(rawCfg.ProxyProvider) > 0 || len(rawCfg.Proxy) > 0 {

			// 对 provider 进行路径替换
			findProvider := changeProvidersPath(profile.Order, rawCfg)
			var yml []byte
			if findProvider {
				yml, _ = yaml.Marshal(rawCfg)
				profile.Path = fmt.Sprintf("./profiles/%d/%s.yaml", profile.Order, profile.Id)
			} else {
				yml = tempBytes
			}

			// 保存操作
			savePath := utils.GetUserHomeDir(profile.Path)
			_, _ = utils.SaveFile(savePath, yml)
			return nil
		} else {
			return fmt.Errorf("proxy or provider is 0")
		}

	}

	return err
}

func changeProvidersPath(snowflakeId int64, config *config.RawConfig) (findProvider bool) {
	findProvider = false

	dir := fmt.Sprintf("./profiles/%d/", snowflakeId)
	proxyProviders := config.ProxyProvider
	for _, provider := range proxyProviders {

		if path, findPath := provider["path"]; findPath {
			provider["path"] = dir + ReplaceTwoPoint(path.(string))
		} else {
			if u, findUrl := provider["url"]; findUrl {
				provider["path"] = dir + utils.MD5(u.(string))
			}
		}

		findProvider = true
	}

	ruleProviders := config.RuleProvider
	for _, ruleProvider := range ruleProviders {

		if path, findPath := ruleProvider["path"]; findPath {
			ruleProvider["path"] = dir + ReplaceTwoPoint(path.(string))
		} else {
			if u, findUrl := ruleProvider["url"]; findUrl {
				ruleProvider["path"] = dir + utils.MD5(u.(string))
			}
		}

		findProvider = true
	}

	return
}

func ReplaceTwoPoint(path string) string {
	path = filepath.Join(path)
	return strings.Replace(path, "../", "", 1)
}

func parseFields(input string) map[string]int64 {
	// 分割字段
	pairs := strings.Split(input, ";")
	result := make(map[string]int64)

	// 处理每个键值对
	for _, pair := range pairs {
		// 去掉可能的空格
		pair = strings.TrimSpace(pair)
		parts := strings.Split(pair, "=")
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			number, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				result[key] = number
			}
		}
	}

	return result
}

func parseContentDisposition(header http.Header, urlStr string) string {
	disposition := header.Get("Content-Disposition")
	if disposition != "" {
		disposition, _ = url.QueryUnescape(disposition)
		split := strings.Split(disposition, "=")
		fileName := split[len(split)-1]
		if strings.Contains(fileName, "''") {
			fileName = strings.Split(fileName, "''")[1]
		}
		if fileName != "" {
			return fileName
		}
	}

	// Fallback: extract the last part of the URL
	if parsedURL, err := url.Parse(urlStr); err == nil {
		segments := strings.Split(parsedURL.Path, "/")
		return segments[len(segments)-1]
	}

	return "Remote-" + utils.GetDateTime()
}

// ParseHeaders 对请求头进行解析
func ParseHeaders(header http.Header, url string, profile *models.Profile) {
	// 流量
	if value := header.Get("Subscription-Userinfo"); value != "" {
		subInfo := parseFields(value)
		profile.Total = subInfo["total"]
		profile.Used = subInfo["upload"] + subInfo["download"]
		profile.Available = profile.Total - profile.Used
		if profile.Available < 0 {
			profile.Available = 0
		}
		if subInfo["expire"] != 0 {
			// 转换为时间
			t := time.Unix(subInfo["expire"], 0)
			profile.Expire = t.Local().Format("2006-01-02 15:04")
		}
	}

	// 文件名
	if profile.Title == "" {
		profile.Title = parseContentDisposition(header, url)
	}

	// 更新间隔
	if val := header.Get("Profile-Update-Interval"); val != "" {
		profile.Interval = val
	}

	// 主页
	if val := header.Get("Profile-Web-Page-Url"); val != "" {
		profile.Home = val
	}

}
