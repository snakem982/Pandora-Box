//go:build windows

package open

import (
	C "github.com/metacubex/mihomo/constant"
	"os/exec"
	"syscall"
)

func OpenConfigDirectory() (string, error) {
	c := exec.Command(`cmd`, `/c`, `explorer`, C.Path.HomeDir())
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return "", c.Start()
}
