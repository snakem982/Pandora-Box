package models

import (
	"github.com/snakem982/pandora-box/pkg/utils"
	"time"
)

type Getter struct {
	Id        string            `json:"id" yaml:"id"`
	Order     int64             `json:"order" yaml:"order"`
	Content   string            `json:"content" yaml:"content"` // 可以为任意内容 url base64 json yaml 等
	TestUrl   string            `json:"testUrl" yaml:"testUrl"`
	Headers   map[string]string `json:"headers,omitempty" yaml:"headers,omitempty"`
	Cache     int               `json:"cache,omitempty" yaml:"cache,omitempty"`
	Crawl     int               `json:"crawl,omitempty" yaml:"crawl,omitempty"`
	Available int               `json:"available,omitempty" yaml:"available,omitempty"`
	Interval  string            `json:"interval,omitempty" yaml:"interval,omitempty"`
	Update    string            `json:"update,omitempty" yaml:"update,omitempty"`
}

type Yml struct {
	Proxies []map[string]any `json:"proxies,omitempty" yaml:"proxies,omitempty"`
}

type Void struct{}

type RealIp struct {
	Key         string `json:"key" yaml:"key"`
	CountryCode string `json:"country_code" yaml:"country_code"`
}

func (g *Getter) GetUpdateTime() time.Time {
	dateTime, _ := utils.ParseDateTime(g.Update)
	return dateTime
}

func (g *Getter) SetUpdateTime() {
	g.Update = utils.GetDateTime()
}
