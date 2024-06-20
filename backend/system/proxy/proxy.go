package proxy

import (
	"strconv"
)

func SetProxy(port string) {
	temp, _ := strconv.Atoi(port)
	_ = OnHttp(Addr{
		Host: "127.0.0.1",
		Port: temp,
	})
	_ = OnHttps(Addr{
		Host: "127.0.0.1",
		Port: temp,
	})
	_ = OnSocks(Addr{
		Host: "127.0.0.1",
		Port: temp,
	})
}

func RemoveProxy() {
	_ = OffAll()
}
