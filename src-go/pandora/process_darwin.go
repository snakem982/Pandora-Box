//go:build darwin

package pandora

import (
	"os"
	"os/signal"
	"syscall"
)

func InitProcessGroup() {
	// 仅在 macOS 提权（root）且当前不是进程组组长时才脱离原进程组与会话，
	// 避免受 authtrampoline 空闲回收时波及同进程组进程
	if os.Geteuid() == 0 && os.Getpid() != syscall.Getpgrp() {
		_, _ = syscall.Setsid()
	}

	// 明确忽略 SIGHUP，防止授权会话或终端断开误杀后台进程
	signal.Ignore(syscall.SIGHUP)
}
