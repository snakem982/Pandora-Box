package spider

import (
	"encoding/json"
	"errors"
	"golang.org/x/net/html"
	"pandora-box/backend/cache"
	"pandora-box/backend/tools"
	"sync"
)

type Getter struct {
	Id             string            `json:"id" yaml:"id"`
	Type           string            `json:"type" yaml:"type"`
	Url            string            `json:"url" yaml:"url"`
	Headers        map[string]string `json:"headers,omitempty" yaml:"headers,omitempty"`
	CrawlNodes     int               `json:"crawl_nodes,omitempty" yaml:"crawl_nodes,omitempty"`
	AvailableNodes int               `json:"available_nodes,omitempty" yaml:"available_nodes,omitempty"`
}

type Collect interface {
	Get() []map[string]any
	Get2ChanWG(pc chan []map[string]any, wg *sync.WaitGroup)
}

type collector func(getter Getter) Collect

var collectorMap = make(map[string]collector)

func Register(sourceType string, c collector) {
	collectorMap[sourceType] = c
}

func NewCollect(sourceType string, getter Getter) (Collect, error) {
	if c, ok := collectorMap[sourceType]; ok {
		return c(getter), nil
	}

	return nil, errors.New("type not supported")
}

func GetBytes(url string, headers map[string]string) []byte {
	all, _ := tools.ConcurrentHttpGet(url, headers)
	if all != nil {
		temp := html.UnescapeString(string(all))
		all = []byte(temp)
	}

	return all
}

func AddIdAndUpdateGetter(pc chan []map[string]any, nodes []map[string]any, g Getter) {
	i := len(nodes)

	// 更新getter
	g.CrawlNodes = i
	g.AvailableNodes = 0
	bytes, _ := json.Marshal(g)
	_ = cache.Put(g.Id, bytes)

	// 添加id
	if i > 0 {
		for _, node := range nodes {
			node["gid"] = g.Id
		}
		pc <- nodes
	}
}

func AvailableAndUpdateGetter(proxies []map[string]any) {
	// 遍历
	gs := make(map[string]int)
	for _, proxy := range proxies {
		if _, ok := proxy["gid"]; !ok {
			continue
		}
		gid := proxy["gid"].(string)
		if _, ok := gs[gid]; ok {
			gs[gid] += 1
		} else {
			gs[gid] = 1
		}
	}

	// 更新
	for k, v := range gs {
		value := cache.Get(k)
		if value == nil {
			continue
		}
		g := Getter{}
		_ = json.Unmarshal(value, &g)
		g.AvailableNodes = v
		bytes, _ := json.Marshal(g)
		_ = cache.Put(g.Id, bytes)
	}
}
