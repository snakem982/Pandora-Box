package sys

import (
	"fmt"
	"github.com/snakem982/pandora-box/pandora/pkg/utils"
	"strconv"
)

func SetProxy(port string) error {
	temp, _ := strconv.Atoi(port)
	if !utils.IsPortAvailable(temp) {
		return fmt.Errorf("port %s is not available", port)
	}

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

	return nil
}

func RemoveProxy() {
	_ = OffAll()
}
