package utils

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"syscall"
)

// 检查进程是否存在
func processExists(pid int) bool {
	if pid <= 0 {
		return false
	}
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	// Windows: os.FindProcess 总是返回一个对象，必须 Signal 测试
	// Linux/macOS: Signal(0) 不会发送信号，只检查存在性
	err = process.Signal(syscall.Signal(0))
	if err == nil {
		return true
	}
	// 在某些系统上，错误是 "operation not permitted"，表示进程存在但无权限访问
	if errors.Is(err, syscall.EPERM) {
		return true
	}
	return false
}

func NotSingleton(name string) bool {
	// 读取pid
	pidFile := GetUserHomeDir("pid", name)
	data, err := ReadFile(pidFile)
	if err == nil {
		if pid, err := strconv.Atoi(data); err == nil {
			if processExists(pid) {
				fmt.Printf("\nAnother instance is running (PID %d).", pid)
				return true
			}
		}
	}

	// 写入当前 PID
	currentPid := os.Getpid()
	pid := []byte(strconv.Itoa(currentPid))
	if ok, err := SaveFile(pidFile, pid); !ok {
		fmt.Printf("\nFailed to write PID file:%v", err)
		return true
	}

	return false
}

func CleanPid(name string) {
	pidFile := GetUserHomeDir("pid", name)
	_ = DeletePath(pidFile)
}
