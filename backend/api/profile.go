package api

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/metacubex/mihomo/common/convert"
	"github.com/metacubex/mihomo/config"
	C "github.com/metacubex/mihomo/constant"
	"github.com/metacubex/mihomo/hub/route"
	"github.com/metacubex/mihomo/log"
	"golang.org/x/net/html"
	"gopkg.in/yaml.v3"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"pandora-box/backend/cache"
	"pandora-box/backend/constant"
	"pandora-box/backend/meta"
	"pandora-box/backend/resolve"
	"pandora-box/backend/spider"
	"pandora-box/backend/tools"
	"path/filepath"
	"strconv"
	"strings"
)

func Profile(r chi.Router) {
	r.Mount("/profile", profileRouter())
}

func profileRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", getProfile)
	r.Post("/", postProfile)
	r.Post("/file", postFileProfile)
	r.Put("/refresh", refreshProfile)

	r.Route("/{id}", func(r chi.Router) {
		r.Put("/", putProfile)
		r.Delete("/", deleteProfile)
		r.Patch("/", patchProfile)
	})

	return r
}

func getProfile(w http.ResponseWriter, r *http.Request) {
	res := make([]resolve.Profile, 0)

	values := cache.GetList(constant.PrefixProfile)
	if len(values) > 0 {
		for _, value := range values {
			profile := resolve.Profile{}
			_ = json.Unmarshal(value, &profile)
			res = append(res, profile)
		}
	}

	render.JSON(w, r, res)
}

func postFileProfile(w http.ResponseWriter, r *http.Request) {

	if r.ContentLength > 2097152 {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, route.HTTPError{Message: "文件大小超过限制<br/>File size is more than 2MB."})
		return
	}

	_, header, _ := r.FormFile("file")
	open, _ := header.Open()
	defer func(open multipart.File) {
		err := open.Close()
		if err != nil {

		}
	}(open)
	content, _ := io.ReadAll(open)
	err := ResolveConfig(false, false, "", "", header.Filename, 41, content)
	if err != nil {
		log.Errorln("[%s] %v", header.Filename, err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, route.HTTPError{Message: err.Error()})
		return
	}

	render.NoContent(w, r)
}

func postProfile(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength > 2097152 {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, route.HTTPError{Message: "内容大小超过限制<br/>Content size is more than 2MB."})
		return
	}

	body := struct {
		Data string `json:"data"`
	}{}
	if err := render.DecodeJSON(r.Body, &body); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, route.ErrBadRequest)
		return
	}

	// 尝试clash解析 成功返回
	bodyData := []byte(body.Data)
	_, err := config.UnmarshalRawConfig(bodyData)
	if err == nil {
		err = ResolveConfig(false, false, "", "", tools.Dec(15), 41, bodyData)
		if err == nil {
			render.NoContent(w, r)
			return
		}
	}

	builder := strings.Builder{}
	urls := make([]string, 0)

	// 按行读取文件
	reader := bytes.NewReader(bodyData)
	bufReader := bufio.NewReader(reader)
	for {
		line, _, err := bufReader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil || len(line) == 0 {
			continue
		}
		sub := strings.TrimSpace(string(line))
		sub = strings.Split(sub, " ")[0]
		if strings.HasPrefix(sub, "http") {
			urls = append(urls, sub)
		} else if strings.Contains(sub, "://") {
			builder.WriteString(sub + "\n")
		} else {
		}
	}

	for _, url := range urls {
		content, fileName := tools.ConcurrentHttpGet(url, nil)
		err := ResolveConfig(false, false, "", url, fileName, 31, content)
		if err != nil {
			log.Errorln("url[%s] %v", url, err)
			continue
		}
	}

	if builder.Len() > 0 || len(urls) == 0 {
		var all []byte
		if builder.Len() > 0 {
			all = []byte(builder.String())
		} else {
			all = bodyData
		}

		ray, base64Error := convert.ConvertsV2Ray(all)
		if base64Error == nil && len(ray) > 0 {
			ray = resolve.MapsToProxies(ray)
			rails := spider.SortAddIndex(ray)
			if len(rails) == 0 {
				render.NoContent(w, r)
				return
			}
			if len(rails) > 511 {
				rails = rails[0:512]
			}
			proxies := make(map[string]any)
			proxies["proxies"] = rails
			all, _ = yaml.Marshal(proxies)

			suffix := "txt"
			kind := 2
			snowflakeId := tools.SnowflakeId()
			profile := resolve.Profile{}
			profile.Id = fmt.Sprintf("%s%d", constant.PrefixProfile, snowflakeId)
			profile.Type = kind
			profile.Title = fmt.Sprintf("%d", snowflakeId)
			profile.Order = snowflakeId
			profile.Path = "uploads/" + profile.Id + "." + suffix

			fileSaveError := saveProfile2Local(profile.Path, all)
			if fileSaveError == nil {
				marshal, _ := json.Marshal(profile)
				_ = cache.Put(profile.Id, marshal)
			}
		}
	}

	render.NoContent(w, r)
}

func putProfile(w http.ResponseWriter, r *http.Request) {
	req := resolve.Profile{}
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, route.ErrBadRequest)
		return
	}
	marshal, _ := json.Marshal(req)
	_ = cache.Put(req.Id, marshal)

	render.NoContent(w, r)
}

func deleteProfile(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	value := cache.Get(id)
	profile := resolve.Profile{}
	err := json.Unmarshal(value, &profile)
	if err != nil {
		render.NoContent(w, r)
		return
	}
	path := C.Path.HomeDir() + "/" + profile.Path
	dir := filepath.Dir(path)
	if strings.HasSuffix(dir, "uploads") {
		_ = os.Remove(path)
	} else {
		_ = os.RemoveAll(dir)
	}
	_ = cache.Delete(id)

	render.NoContent(w, r)
}

func saveProfile2Local(profilePath string, all []byte) error {
	path := C.Path.HomeDir() + "/" + profilePath
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// 必须分成两步：先创建文件夹、再修改权限
		_ = os.Mkdir(dir, 0777)
		_ = os.Chmod(dir, 0777)
	}
	_ = os.Remove(path)
	return os.WriteFile(path, all, 0777)
}

func patchProfile(w http.ResponseWriter, r *http.Request) {
	req := resolve.Profile{}
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, route.ErrBadRequest)
		return
	}
	meta.StartCore(req, true)

	render.NoContent(w, r)
}

func refreshProfile(w http.ResponseWriter, r *http.Request) {
	req := resolve.Profile{}
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, route.ErrBadRequest)
		return
	}

	url := req.Url
	content, _ := tools.ConcurrentHttpGet(url, nil)
	err := ResolveConfig(true, req.Selected, req.Id, url, req.Title, req.Type, content)
	if err != nil {
		log.Errorln("url[%s] %v", url, err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, route.HTTPError{Message: err.Error()})
		return
	}

	render.NoContent(w, r)
}

func changeProvidersPath(snowflakeId int64, config *config.RawConfig) (findProvider bool) {
	findProvider = false

	dir := fmt.Sprintf("./uploads/%d/", snowflakeId)
	proxyProviders := config.ProxyProvider
	for _, provider := range proxyProviders {

		if path, findPath := provider["path"]; findPath {
			provider["path"] = dir + ReplaceTwoPoint(path.(string))
		} else {
			if url, findUrl := provider["url"]; findUrl {
				provider["path"] = dir + tools.MD5(url.(string))
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
				ruleProvider["path"] = dir + tools.MD5(url.(string))
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

func ResolveConfig(refresh, selected bool,
	id string, url string, fileName string,
	kind int, content []byte) error {

	if content == nil || len(content) < 32 {
		log.Errorln("ResolveConfig error: %s", "content is nil or length less 32")
		return fmt.Errorf("content is nil or length less 32")
	}

	// 如果不是刷新创建snowflakeId
	snowflakeId := tools.SnowflakeId()
	if refresh {
		if kind == 32 || kind == 42 {
			kind = kind - 1
		}
		id = strings.TrimLeft(id, constant.PrefixProfile)
		i, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return fmt.Errorf("strconv.ParseInt: %v", err)
		} else {
			snowflakeId = i
		}
	}
	// 是否使用Provider
	findProvider := false
	// 对内容进行html解码
	temp := html.UnescapeString(string(content))
	temp = strings.Replace(temp, "\"HOST\"", "\"Host\"", -1)
	content = []byte(temp)
	rawCfg, rawErr := config.UnmarshalRawConfig(content)
	var ray []map[string]any
	suffix := "yaml"
	// yaml解析失败，尝试base64解码
	if rawErr != nil {
		log.Errorln("config.UnmarshalRawConfig error: %s", rawErr.Error())
		var base64Error error
		var jsonError error
		ray, base64Error = convert.ConvertsV2Ray(content)
		if base64Error != nil {
			log.Errorln("convert.ConvertsV2Ray error: %s", base64Error.Error())
			// base64解析失败，尝试json解码
			ray, jsonError = convert.ConvertsSingBox(content)
			if jsonError != nil {
				log.Errorln("convert.ConvertsSingBox error: %s", jsonError.Error())
				return fmt.Errorf("convert.ConvertsSingBox error: %s", jsonError.Error())
			}
		}
		rails := resolve.MapsToProxies(ray)
		if len(rails) == 0 {
			log.Errorln("resolve.MapsToProxies error: %s", "Node is 0")
			return fmt.Errorf("resolve.MapsToProxies error: %s", "Node is 0")
		}
		spider.SortProxies(rails)
		if len(rails) > 511 {
			rails = rails[0:512]
		}
		proxies := make(map[string]any)
		proxies["proxies"] = rails
		content, _ = yaml.Marshal(proxies)
		kind = kind + 1
		suffix = "txt"
	} else {
		// yaml解析成功，进行配置校验
		rawCfg.GeodataMode = false
		ko, yamlError := config.ParseRawConfig(rawCfg)
		if yamlError != nil {
			log.Errorln("config.ParseRawConfig error: %s", yamlError.Error())
			// 配置校验失败，尝试提取可用节点
			rails := resolve.MapsToProxies(rawCfg.Proxy)
			if len(rails) == 0 {
				log.Errorln("resolve.MapsToProxies error: %s", "Node is 0")
				return fmt.Errorf("resolve.MapsToProxies error: %s", "Node is 0")
			}
			spider.SortProxies(rails)
			if len(rails) > 511 {
				rails = rails[0:512]
			}
			proxies := make(map[string]any)
			proxies["proxies"] = rails
			content, _ = yaml.Marshal(proxies)
		} else {
			findProvider = changeProvidersPath(snowflakeId, rawCfg)
			if !findProvider {
				if len(ko.Proxies) < 7 {
					return fmt.Errorf("config.ParseRawConfig error: %s", "Node is 0")
				}
			}

			if len(ko.Proxies) > 512 {
				// 对于超过512的节点进行截取
				log.Infoln("config.ParseRawConfig : %s Try to cut", "Node size is more than 512.")
				ray = resolve.MapsToProxies(rawCfg.Proxy)
				spider.SortProxies(ray)
				rails := ray[0:512]
				proxies := make(map[string]any)
				proxies["proxies"] = rails
				content, _ = yaml.Marshal(proxies)
			}
		}
	}

	profile := resolve.Profile{}
	profile.Id = fmt.Sprintf("%s%d", constant.PrefixProfile, snowflakeId)
	profile.Type = kind
	if fileName == "" {
		fileName = fmt.Sprintf("sub-%d", snowflakeId)
	}

	profile.Title = fileName
	profile.Url = url
	profile.Order = snowflakeId
	profile.Selected = selected

	if findProvider {
		pg := struct {
			ProxyGroup []map[string]any `yaml:"proxy-groups"`
		}{}
		_ = yaml.Unmarshal(content, &pg)
		rawCfg.ProxyGroup = pg.ProxyGroup
		content, _ = yaml.Marshal(rawCfg)
		profile.Path = fmt.Sprintf("uploads/%d/%s%d.%s", snowflakeId, constant.PrefixProfile, snowflakeId, suffix)
	} else {
		profile.Path = "uploads/" + profile.Id + "." + suffix
	}

	fileSaveError := saveProfile2Local(profile.Path, content)
	if fileSaveError != nil {
		return fmt.Errorf("fileSaveError:%v", fileSaveError)
	}
	marshal, _ := json.Marshal(profile)
	_ = cache.Put(profile.Id, marshal)

	return nil
}
