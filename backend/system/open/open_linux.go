//go:build linux

package open

import (
	C "github.com/metacubex/mihomo/constant"
)

func OpenConfigDirectory() (string, error) {
	return proxy.Command("nautilus", C.Path.HomeDir())
}
