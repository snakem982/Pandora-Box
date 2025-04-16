package models

type Profile struct {
	Id       string `json:"id,omitempty"`
	Type     int    `json:"type"` // 1: 远程订阅 2：本地配置 3：爬取合并
	Title    string `json:"title,omitempty"`
	Order    int64  `json:"order"`
	Selected bool   `json:"selected,omitempty"`
	Path     string `json:"path"`
	Url      string `json:"url,omitempty"`
	Upload   string `json:"upload,omitempty"`
	Download string `json:"download,omitempty"`
	Total    string `json:"total,omitempty"`
	Expire   string `json:"expire,omitempty"`
	Interval string `json:"interval,omitempty"`
	Home     string `json:"home,omitempty"`
	Update   string `json:"update,omitempty"`
}
