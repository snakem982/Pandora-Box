package models

type Getter struct {
	Id        string            `json:"id" yaml:"id"`
	Content   string            `json:"content" yaml:"content"` // 可以为任意内容 url base64 json yaml 等
	Headers   map[string]string `json:"headers,omitempty" yaml:"headers,omitempty"`
	Cache     int               `json:"cache,omitempty" yaml:"cache,omitempty"`
	Crawl     int               `json:"crawl,omitempty" yaml:"crawl,omitempty"`
	Available int               `json:"available,omitempty" yaml:"available,omitempty"`
	Interval  string            `json:"interval,omitempty" yaml:"interval,omitempty"`
	Update    string            `json:"update,omitempty" yaml:"update,omitempty"`
}
