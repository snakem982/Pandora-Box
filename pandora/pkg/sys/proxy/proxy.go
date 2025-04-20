package sys

import (
	"fmt"
	"github.com/snakem982/pandora-box/pandora/pkg/utils"
)

// EnableProxy 开启系统代理
func EnableProxy(host string, port int) error {
	// 检测端口，不可用进行报错
	if !utils.IsPortAvailable(port) {
		return fmt.Errorf("port %d is not available", port)
	}

	_ = OnHttp(Addr{
		Host: host,
		Port: port,
	})
	_ = OnHttps(Addr{
		Host: host,
		Port: port,
	})
	_ = OnSocks(Addr{
		Host: host,
		Port: port,
	})

	return nil
}

// DisableProxy 关闭代理
func DisableProxy() {
	_ = OffAll()
}
