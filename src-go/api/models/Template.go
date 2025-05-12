package models

type Template struct {
	Id       string `json:"id" yaml:"id"`
	Order    int64  `json:"order" yaml:"order"`
	Title    string `json:"title" yaml:"title"`
	Path     string `json:"path" yaml:"path"`
	Selected bool   `json:"selected" yaml:"selected"`
}
