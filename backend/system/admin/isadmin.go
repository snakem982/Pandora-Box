//go:build !windows

package isadmin

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Check if the program has administrative privileges.
func Check() bool {
	return os.Getuid() == 0
}

// KillProcessesByName 杀死所有名字为指定名称的进程
func KillProcessesByName(name string) error {
	// 使用 ps 和 grep 命令来查找进程
	cmd := exec.Command("bash", "-c", fmt.Sprintf("ps aux | grep %s | grep -v grep | awk '{print $2}'", name))
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	pids := strings.Fields(string(output))
	for _, pid := range pids {
		// 使用 kill 命令来杀死进程
		killCmd := exec.Command("kill", "-9", pid)
		err := killCmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}
