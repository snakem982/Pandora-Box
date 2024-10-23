package api

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/metacubex/mihomo/config"
	C "github.com/metacubex/mihomo/constant"
	"github.com/metacubex/mihomo/hub/route"
	"github.com/metacubex/mihomo/log"
	"gopkg.in/yaml.v3"
	"net/http"
	"os"
	"pandora-box/backend/constant"
	"pandora-box/backend/resolve"
	"pandora-box/backend/spider"
	"pandora-box/backend/tools"
	"path/filepath"
	"sort"
	"strings"
)

func Filter(r chi.Router) {

	r.Get("/nodeHave", func(w http.ResponseWriter, r *http.Request) {
		defaultBuf, defaultErr := spider.GetNodesCache()
		if defaultErr == nil && len(defaultBuf) > 0 {
			rawCfg, err := config.UnmarshalRawConfig(defaultBuf)
			if err == nil && len(rawCfg.Proxy) > 0 {
				render.PlainText(w, r, "true")
			} else {
				render.PlainText(w, r, "false")
			}
		} else {
			render.PlainText(w, r, "false")
		}
	})

	r.Get("/nodeCache", func(w http.ResponseWriter, r *http.Request) {
		// 加载默认配置中的节点
		defaultBuf, defaultErr := spider.GetNodesCache()
		if defaultErr == nil && len(defaultBuf) > 0 {
			rawCfg, err := config.UnmarshalRawConfig(defaultBuf)
			if err == nil && len(rawCfg.Proxy) > 0 {
				log.Infoln("load default config proxies success %d", len(rawCfg.Proxy))
				render.JSON(w, r, getStatus(rawCfg.Proxy))
			} else {
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, render.M{"message": "nodes is 0"})
			}
		} else {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, render.M{"message": "nodes is 0"})
		}
	})

	r.Post("/nodeFilter", func(w http.ResponseWriter, r *http.Request) {
		req := struct {
			Protocol []string `json:"protocol"`
			Country  []string `json:"country"`
			Count    int      `json:"count"`
			Option   int      `json:"option"`
		}{}
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, route.ErrBadRequest)
			return
		}
		allProtocol := len(req.Protocol) == 0
		allCountry := len(req.Country) == 0

		// 加载默认配置中的节点
		defaultBuf, defaultErr := spider.GetNodesCache()
		if defaultErr == nil && len(defaultBuf) > 0 {
			rawCfg, err := config.UnmarshalRawConfig(defaultBuf)
			if err == nil && len(rawCfg.Proxy) > 0 {
				log.Infoln("load default config proxies success %d", len(rawCfg.Proxy))
				proxies := make([]map[string]any, 0)
				for _, proxy := range rawCfg.Proxy {
					if allProtocol || judge(req.Protocol, proxy["type"].(string)) {
						name := proxy["name"].(string)
						name = strings.Split(name, "_")[0]
						if allCountry || judge(req.Country, name) {
							proxies = append(proxies, proxy)
						}
					}
				}
				spider.SortProxies(proxies)
				if len(proxies) >= req.Count {
					proxies = proxies[0:req.Count]
				}
				for i, _ := range proxies {
					name := proxies[i]["name"].(string)
					name = strings.Split(name, "_")[0]
					name = fmt.Sprintf("%s_%+02v", name, i+1)
					proxies[i]["name"] = name
				}

				// 1:筛选 2：覆盖默认 3：生成新配置 4：导出使用
				switch req.Option {
				case 1:
					render.JSON(w, r, getStatus(proxies))
					return
				case 2:
					spider.Save2Local(proxies, "0.yaml")
					render.PlainText(w, r, "true")
					return
				case 3:
					nodes := make(map[string]any)
					nodes["proxies"] = proxies
					content, _ := yaml.Marshal(nodes)
					_ = ResolveConfig(false, false, "", "", "Filter_"+tools.Dec(8), 41, content)
					render.PlainText(w, r, "true")
					return
				case 4:
					nodes := make(map[string]any)
					nodes["proxies"] = proxies
					content, _ := yaml.Marshal(nodes)
					replace := strings.Replace(string(resolve.PandoraDefaultDownloadConfig),
						resolve.PandoraDefaultPlace,
						string(content),
						1)
					_ = os.WriteFile(filepath.Join(C.Path.HomeDir(), constant.DefaultDownload), []byte(replace), 0777)

					render.PlainText(w, r, "true")
					return
				}

			} else {
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, render.M{"message": "nodes is 0"})
				return
			}
		} else {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, render.M{"message": "nodes is 0"})
			return
		}
	})

	r.Get("/Pandora-Box-Download", func(w http.ResponseWriter, r *http.Request) {
		file, err := os.ReadFile(filepath.Join(C.Path.HomeDir(), constant.DefaultDownload))
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, "server error: "+err.Error())
			return
		}
		w.Header().Set("Content-Disposition", "attachment; filename=\"Pandora-Box-Share.yaml\"")
		render.Data(w, r, file)
	})
}

type Node struct {
	Name  string `json:"value"`
	Label string `json:"label"`
	Value int    `json:"count"`
}

func sortMap(input map[string]int) []Node {
	output := make([]Node, 0)

	for k, v := range input {
		output = append(output, Node{k, k, v})
	}

	sort.Slice(output, func(i, j int) bool {

		if output[i].Value == output[j].Value {
			return output[i].Name < output[j].Name
		}

		return output[i].Value > output[j].Value
	})

	return output
}

func getStatus(proxies []map[string]any) map[string]any {
	// 遍历
	protocol := make(map[string]int)
	country := make(map[string]int)
	for _, proxy := range proxies {
		// 各种协议节点数量统计
		proxyType := proxy["type"].(string)
		if _, ok := protocol[proxyType]; ok {
			protocol[proxyType] += 1
		} else {
			protocol[proxyType] = 1
		}
		// 各个国家节点数量统计
		name := proxy["name"].(string)
		name = strings.Split(name, "_")[0]
		if _, ok := country[name]; ok {
			country[name] += 1
		} else {
			country[name] = 1
		}
	}

	// node cache 总结
	status := make(map[string]any)
	status["count"] = len(proxies)
	status["protocol"] = sortMap(protocol)
	status["country"] = sortMap(country)

	return status
}

func judge(array []string, key string) bool {
	for _, value := range array {
		if value == key {
			return true
		}
	}

	return false
}

func saveFile() {

}
