package models

type WebTest struct {
	Id      string `json:"id" yaml:"id"`
	Order   int64  `json:"order" yaml:"order"`
	Src     string `json:"src" yaml:"src"`
	TestUrl string `json:"testUrl" yaml:"testUrl"`
}
