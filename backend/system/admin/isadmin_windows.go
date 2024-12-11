//go:build windows

package isadmin

import (
	"golang.org/x/sys/windows"
	"pandora-box/backend/system/proxy"
	"strings"
)

// Check if the program has administrative privileges.
func Check() bool {
	var sid *windows.SID

	_ = windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		2,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
		&sid)

	admin, _ := windows.Token(0).IsMember(sid)

	return admin
}

// KillProcessesByName 杀死所有名字为指定名称的进程
func KillProcessesByName(name string) error {
	// 使用 tasklist 命令查找进程
	output, err := proxy.Command("tasklist")
	if err != nil {
		return err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.ToLower(line)
		name = strings.ToLower(name)
		if strings.Contains(line, name) {
			fields := strings.Fields(line)
			if len(fields) > 1 {
				pid := fields[1]

				// 使用 taskkill 命令来杀死进程
				_, err = proxy.Command("taskkill", "/F", "/PID", pid)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
