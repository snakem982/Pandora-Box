//go:build darwin

package open

import (
	C "github.com/metacubex/mihomo/constant"
	"pandora-box/backend/system/proxy"
)

func OpenConfigDirectory() (string, error) {
	return proxy.Command("open", C.Path.HomeDir())
}
