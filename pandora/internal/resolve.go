package internal

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/schema"
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
	"strings"
	"sync"
	"time"
)

// 保存文件
func saveProfile(proxies []map[string]any, profile models.Profile) {
	yml := models.Yml{Proxies: proxies}
	out, _ := yaml.Marshal(yml)
	savePath := utils.GetUserHomeDir("profiles", profile.Id+".yaml")
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
func Resolve(input models.Profile, refresh bool) (models.Profile, error) {
	profile := input
	tempStr := strings.TrimSpace(profile.Content)
	tempBytes := []byte(tempStr)

	// 如果不是刷新则创建 id
	if !refresh {
		snowflakeId := utils.SnowflakeId()
		profile.Id = fmt.Sprintf("%s%d", constant.PrefixProfile, snowflakeId)
		profile.Order = snowflakeId
	}
	profile.SetUpdateTime()

	// Sing解析
	if utils.IsJSON(tempStr) {
		sing, err := convert.ConvertsSingBox(tempBytes)
		if err == nil {
			saveProfile(sing, profile)
			return profile, nil
		}

		return profile, err
	}

	// Base64解析
	if utils.IsBase64(tempStr) {
		v2ray, err := convert.ConvertsV2Ray(tempBytes)
		if err == nil {
			saveProfile(v2ray, profile)
			return profile, nil
		}

		return profile, err
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
			return profile, nil
		}

		return profile, err
	}

	// Yaml解析
	rawCfg, err := config.UnmarshalRawConfig(tempBytes)
	if err == nil {
		_, yamlError := config.ParseRawConfig(rawCfg)
		if yamlError != nil {
			// 配置校验失败，尝试提取可用节点
			rails := MapsToProxies(rawCfg.Proxy)
			if len(rails) == 0 {
				return profile, yamlError
			} else {
				saveProfile(rails, profile)
				return profile, nil
			}
		}

		// 保存yaml
		if len(rawCfg.ProxyProvider) > 0 || len(rawCfg.Proxy) > 0 {

			// 对 provider 进行路径替换
			findProvider := changeProvidersPath(profile.Order, rawCfg)
			var content []byte
			if findProvider {
				content, _ = yaml.Marshal(rawCfg)
				profile.Path = fmt.Sprintf("profiles/%d/%s.yaml", profile.Order, profile.Id)
			} else {
				content = tempBytes
				profile.Path = "profiles/" + profile.Id + ".yaml"
			}

			// 保存操作
			savePath := utils.GetUserHomeDir(profile.Path)
			_, _ = utils.SaveFile(savePath, content)
			return profile, nil
		}

		return profile, fmt.Errorf("proxy or provider is 0")
	}

	return profile, err
}

func changeProvidersPath(snowflakeId int64, config *config.RawConfig) (findProvider bool) {
	findProvider = false

	dir := fmt.Sprintf("./profiles/%d/", snowflakeId)
	proxyProviders := config.ProxyProvider
	for _, provider := range proxyProviders {

		if path, findPath := provider["path"]; findPath {
			provider["path"] = dir + ReplaceTwoPoint(path.(string))
		} else {
			if url, findUrl := provider["url"]; findUrl {
				provider["path"] = dir + utils.MD5(url.(string))
			}
		}

		findProvider = true
	}

	ruleProviders := config.RuleProvider
	for _, ruleProvider := range ruleProviders {

		if path, findPath := ruleProvider["path"]; findPath {
			ruleProvider["path"] = dir + ReplaceTwoPoint(path.(string))
		} else {
			if url, findUrl := ruleProvider["url"]; findUrl {
				ruleProvider["path"] = dir + utils.MD5(url.(string))
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

func parseFields(input string) map[string]string {
	// 分割字段
	pairs := strings.Split(input, ";")
	result := make(map[string]string)

	// 处理每个键值对
	for _, pair := range pairs {
		// 去掉可能的空格
		pair = strings.TrimSpace(pair)
		parts := strings.Split(pair, "=")
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			result[key] = value
		}
	}

	return result
}

func parseFilename(contentDisposition, key string) string {
	if strings.Contains(contentDisposition, key) {
		values := make(map[string][]string)
		if err := schema.NewDecoder().Decode(&values, url.Values{key: {contentDisposition}}); err == nil {
			if filenames, ok := values[key]; ok && len(filenames) > 0 {
				return filenames[0]
			}
		}
	}
	return ""
}

func parseContentDisposition(header http.Header, urlStr string) string {
	contentDisposition := header.Get("Content-Disposition")
	if contentDisposition != "" {
		contentDisposition = strings.Trim(contentDisposition, "\"")
		if filename := parseFilename(contentDisposition, "filename*"); filename != "" {
			return filename
		}
		if filename := parseFilename(contentDisposition, "filename"); filename != "" {
			return filename
		}
	}

	// Fallback: extract the last part of the URL
	if parsedURL, err := url.Parse(urlStr); err == nil {
		segments := strings.Split(parsedURL.Path, "/")
		return segments[len(segments)-1]
	}

	return "Remote File"
}

// ParseHeaders 对请求头进行解析
func ParseHeaders(header http.Header, url string, profile *models.Profile) {
	// 流量
	if value := header.Get("Subscription-Userinfo"); value != "" {
		subInfo := parseFields(value)
		profile.Upload = subInfo["upload"]
		profile.Download = subInfo["download"]
		profile.Total = subInfo["total"]
		profile.Expire = subInfo["expire"]
	}

	// 文件名
	profile.Title = parseContentDisposition(header, url)

	// 更新间隔
	if val := header.Get("Profile-Update-Interval"); val != "" {
		profile.Interval = val
	}

	// 主页
	if val := header.Get("Profile-Web-Page-Url"); val != "" {
		profile.Home = val
	}

}
