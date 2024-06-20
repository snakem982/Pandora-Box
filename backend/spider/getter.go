package spider

import (
	"errors"
	"golang.org/x/net/html"
	"pandora-box/backend/tools"
	"strings"
	"sync"
)

type Getter struct {
	Id   string `json:"id,omitempty" yaml:"id,omitempty"`
	Type string `json:"type" yaml:"type"`
	Url  string `json:"url" yaml:"url"`
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

var ErrorCreateNotSupported = errors.New("type not supported")

func NewCollect(sourceType string, getter Getter) (Collect, error) {
	if c, ok := collectorMap[sourceType]; ok {
		return c(getter), nil
	}

	return nil, ErrorCreateNotSupported
}

func GetBytes(url string) []byte {
	all, _ := tools.ConcurrentHttpGet(url)
	if all != nil {
		temp := html.UnescapeString(string(all))
		temp = strings.Replace(temp, "\"HOST\"", "\"Host\"", -1)
		all = []byte(temp)
	}

	return all
}
