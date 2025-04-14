package models

import _ "embed"

//go:embed config.yaml
var PandoraDefaultConfig []byte

//go:embed config_download.yaml
var PandoraDefaultDownloadConfig []byte

type Profile struct {
	Id         string `json:"id,omitempty"`
	Type       int    `json:"type"` // 1: 默认yaml 2:分享txt 31:订阅yaml 32:订阅txt 41:导入yaml 42:导入txt
	Title      string `json:"title,omitempty"`
	Path       string `json:"path"`
	Url        string `json:"url,omitempty"`
	HomePage   string `json:"homePage,omitempty"`
	Selected   bool   `json:"selected,omitempty"`
	Order      int64  `json:"order"`
	IsWarpPlus bool   `json:"isWarpPlus,omitempty"`
}
