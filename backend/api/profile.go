package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/metacubex/mihomo/common/convert"
	"github.com/metacubex/mihomo/config"
	C "github.com/metacubex/mihomo/constant"
	"github.com/metacubex/mihomo/hub/executor"
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
	all, _ := io.ReadAll(open)
	// 对内容进行html解码
	temp := html.UnescapeString(string(all))
	temp = strings.Replace(temp, "\"HOST\"", "\"Host\"", -1)
	all = []byte(temp)
	ko, yamlError := executor.ParseWithBytes(all)
	suffix := "yaml"
	kind := 41
	if yamlError != nil {
		log.Errorln("postFileProfile parse error-1: %s", yamlError.Error())
		var ray []map[string]any
		var base64Error error
		rawCfg, err := config.UnmarshalRawConfig(all)
		if err == nil && len(rawCfg.Proxy) > 0 {
			ray = rawCfg.Proxy
		} else {
			log.Errorln("postFileProfile parse error-2: %s", err.Error())
			ray, base64Error = convert.ConvertsV2Ray(all)
			if base64Error != nil {
				log.Errorln("postFileProfile parse error-3: %s", base64Error.Error())
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, route.HTTPError{Message: yamlError.Error()})
				return
			}
			suffix = "txt"
			kind = 42
		}
		ray = resolve.MapsToProxies(ray)
		rails := spider.SortAddIndex(ray)
		if len(rails) == 0 {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, route.HTTPError{Message: "节点数为零<br/>Node size is 0."})
			return
		}
		if len(rails) > 512 {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, route.HTTPError{Message: "节点数超过限制512<br/>Node size is more than 512."})
			return
		}
		proxies := make(map[string]any)
		proxies["proxies"] = rails
		all, _ = yaml.Marshal(proxies)
	} else {
		if len(ko.Proxies) < 7 {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, route.HTTPError{Message: "节点数为零<br/>Node size is 0."})
			return
		}
		if len(ko.Proxies) > 512 {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, route.HTTPError{Message: "节点数超过限制512<br/>Node size is more than 512."})
			return
		}
	}

	snowflakeId := tools.SnowflakeId()
	profile := resolve.Profile{}
	profile.Id = fmt.Sprintf("%s%d", constant.PrefixProfile, snowflakeId)
	profile.Type = kind
	profile.Title = header.Filename
	profile.Order = snowflakeId
	profile.Path = "uploads/" + profile.Id + "." + suffix

	fileSaveError := saveProfile2Local(profile.Id, suffix, all)
	if fileSaveError != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, route.HTTPError{Message: fileSaveError.Error()})
		return
	}

	bytes, _ := json.Marshal(profile)
	_ = cache.Put(profile.Id, bytes)

	render.NoContent(w, r)
}

func postProfile(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Data string `json:"data"`
	}{}
	if err := render.DecodeJSON(r.Body, &body); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, route.ErrBadRequest)
		return
	}

	builder := strings.Builder{}
	urls := make([]string, 0)
	b64 := ""

	subs := strings.Split(body.Data, "\n")
	for _, subTemp := range subs {
		sub := strings.TrimRight(subTemp, " \r")
		if sub == "" {
			continue
		}
		if strings.HasPrefix(sub, "http") {
			urls = append(urls, sub)
		} else if strings.Contains(sub, "://") {
			builder.WriteString(sub + "\n")
		} else {
			b64 = sub
			break
		}
	}

	for _, url := range urls {
		all, fileName := tools.ConcurrentHttpGet(url)
		if all != nil && len(all) > 128 {
			// 对内容进行html解码
			temp := html.UnescapeString(string(all))
			temp = strings.Replace(temp, "\"HOST\"", "\"Host\"", -1)
			all = []byte(temp)
			ko, yamlError := executor.ParseWithBytes(all)
			suffix := "yaml"
			kind := 31
			if yamlError != nil {
				log.Errorln("postProfile parse error-1: %s", yamlError.Error())
				var ray []map[string]any
				var base64Error error
				rawCfg, err := config.UnmarshalRawConfig(all)
				if err == nil && len(rawCfg.Proxy) > 0 {
					ray = rawCfg.Proxy
				} else {
					log.Errorln("postProfile parse error-2: %s", err.Error())
					ray, base64Error = convert.ConvertsV2Ray(all)
					if base64Error != nil {
						log.Errorln("postProfile parse error-3: %s", base64Error.Error())
						continue
					}
					suffix = "txt"
					kind = 32
				}
				ray = resolve.MapsToProxies(ray)
				rails := spider.SortAddIndex(ray)
				if len(rails) == 0 {
					continue
				}
				if len(rails) > 512 {
					continue
				}
				proxies := make(map[string]any)
				proxies["proxies"] = rails
				all, _ = yaml.Marshal(proxies)
			} else {
				if len(ko.Proxies) < 7 {
					continue
				}
				if len(ko.Proxies) > 512 {
					continue
				}
			}

			snowflakeId := tools.SnowflakeId()
			profile := resolve.Profile{}
			profile.Id = fmt.Sprintf("%s%d", constant.PrefixProfile, snowflakeId)
			profile.Type = kind
			if fileName == "" {
				fileName = fmt.Sprintf("sub-%d", snowflakeId)
			}
			profile.Title = fileName
			profile.Url = url
			profile.Order = snowflakeId
			profile.Path = "uploads/" + profile.Id + "." + suffix

			fileSaveError := saveProfile2Local(profile.Id, suffix, all)
			if fileSaveError != nil {
				continue
			}

			bytes, _ := json.Marshal(profile)
			_ = cache.Put(profile.Id, bytes)
		} else {
			log.Errorln("postProfile url get null or length less 128 [%s]", url)
		}
	}

	if builder.Len() > 0 || b64 != "" {
		var all []byte
		if builder.Len() > 0 {
			all = []byte(builder.String())
		} else {
			all = []byte(b64)
		}

		ray, base64Error := convert.ConvertsV2Ray(all)
		if base64Error == nil && len(ray) > 0 {
			ray = resolve.MapsToProxies(ray)
			rails := spider.SortAddIndex(ray)
			if len(rails) == 0 {
				render.NoContent(w, r)
				return
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

			fileSaveError := saveProfile2Local(profile.Id, suffix, all)
			if fileSaveError == nil {
				bytes, _ := json.Marshal(profile)
				_ = cache.Put(profile.Id, bytes)
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
	bytes, _ := json.Marshal(req)
	_ = cache.Put(req.Id, bytes)

	render.NoContent(w, r)
}

func deleteProfile(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	bytes := cache.Get(id)
	profile := resolve.Profile{}
	err := json.Unmarshal(bytes, &profile)
	if err != nil {
		render.NoContent(w, r)
		return
	}
	_ = os.Remove(C.Path.HomeDir() + "/" + profile.Path)
	_ = cache.Delete(id)

	render.NoContent(w, r)
}

func saveProfile2Local(name, suffix string, all []byte) error {
	return os.WriteFile(C.Path.HomeDir()+"/uploads/"+name+"."+suffix, all, 0666)
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
	all, _ := tools.ConcurrentHttpGet(url)

	if all != nil && len(all) > 128 {
		// 对内容进行html解码
		temp := html.UnescapeString(string(all))
		temp = strings.Replace(temp, "\"HOST\"", "\"Host\"", -1)
		all = []byte(temp)
		ko, yamlError := executor.ParseWithBytes(all)
		kind := 31
		if yamlError != nil {
			log.Errorln("refreshProfile parse error-1: %s", yamlError.Error())
			var ray []map[string]any
			var base64Error error
			rawCfg, err := config.UnmarshalRawConfig(all)
			if err == nil && len(rawCfg.Proxy) > 0 {
				ray = rawCfg.Proxy
			} else {
				log.Errorln("refreshProfile parse error-2: %s", err.Error())
				ray, base64Error = convert.ConvertsV2Ray(all)
				if base64Error != nil {
					log.Errorln("refreshProfile parse error-3: %s", base64Error.Error())
					render.Status(r, http.StatusBadRequest)
					render.JSON(w, r, route.HTTPError{Message: yamlError.Error()})
					return
				}
				kind = 32
			}
			ray = resolve.MapsToProxies(ray)
			rails := spider.SortAddIndex(ray)
			if len(rails) == 0 {
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, route.HTTPError{Message: "节点数为零<br/>Node size is 0."})
				return
			}
			if len(rails) > 512 {
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, route.HTTPError{Message: "节点数超过限制512<br/>Node size is more than 512."})
				return
			}
			proxies := make(map[string]any)
			proxies["proxies"] = rails
			all, _ = yaml.Marshal(proxies)
		} else {
			if len(ko.Proxies) < 7 {
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, route.HTTPError{Message: "节点数为零<br/>Node size is 0."})
				return
			}
			if len(ko.Proxies) > 512 {
				render.Status(r, http.StatusBadRequest)
				render.JSON(w, r, route.HTTPError{Message: "节点数超过限制512<br/>Node size is more than 512."})
				return
			}
		}

		req.Type = kind

		filePath := C.Path.HomeDir() + "/" + req.Path
		_ = os.Remove(filePath)
		fileSaveError := os.WriteFile(filePath, all, 0666)
		if fileSaveError != nil {
			render.Status(r, http.StatusAccepted)
			render.JSON(w, r, route.HTTPError{Message: "服务内部错误"})
			return
		}

		bytes, _ := json.Marshal(req)
		_ = cache.Put(req.Id, bytes)

		render.NoContent(w, r)
	} else {
		log.Errorln("refreshProfile url get null or length less 128 [%s]", url)
		render.Status(r, http.StatusAccepted)
		render.JSON(w, r, route.ErrRequestTimeout)
		return
	}
}
