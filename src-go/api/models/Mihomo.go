package models

type Mihomo struct {
	Mode  string `json:"mode"`
	Proxy bool   `json:"proxy"`
	Tun   bool   `json:"tun"`

	Port        int    `json:"port"`
	BindAddress string `json:"bindAddress"`
	Stack       string `json:"stack"`
	Dns         bool   `json:"dns"`
	Ipv6        bool   `json:"ipv6"`
}
