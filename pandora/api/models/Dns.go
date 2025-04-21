package models

type Dns struct {
	Enable  bool   `json:"enable" yaml:"enable"`
	Content string `json:"content" yaml:"content"`
}
