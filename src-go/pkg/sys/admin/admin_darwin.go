package sys

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// IsAdmin if the program has administrative privileges.
func IsAdmin() bool {
	return os.Getuid() == 0
}

// KillProcessesByName 杀死所有名字为指定名称的进程
func KillProcessesByName(name string) error {
	// 确保目标进程名称不为空
	if name == "" {
		return fmt.Errorf("process name cannot be empty")
	}

	// 使用 ps 和 grep 查找进程，并确保精确匹配进程名称
	cmd := exec.Command("bash", "-c", fmt.Sprintf("ps -eo pid,comm | grep -w %s | grep -v grep | awk '{print $1}'", name))
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to find processes with name %s: %w", name, err)
	}

	// 解析输出，获取所有匹配的进程 ID
	pids := strings.Fields(string(output))
	if len(pids) == 0 {
		return fmt.Errorf("no processes found with name %s", name)
	}

	// 遍历所有进程 ID，逐个杀死
	for _, pid := range pids {
		killCmd := exec.Command("kill", "-9", pid)
		err := killCmd.Run()
		if err != nil {
			return fmt.Errorf("failed to kill process with PID %s: %w", pid, err)
		}
	}

	return nil
}
