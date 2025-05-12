package sys

import (
	"fmt"
	"strings"

	sys "github.com/snakem982/pandora-box/pkg/sys/cmd"
	"golang.org/x/sys/windows"
)

// IsAdmin 检查程序是否具有管理员权限
func IsAdmin() bool {
	var sid *windows.SID

	// 分配并初始化管理员组的 SID
	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		2,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
		&sid)
	if err != nil {
		// 如果 SID 初始化失败，没有管理员权限
		return false
	}
	defer windows.FreeSid(sid) // 确保释放 SID 资源

	// 检查当前进程是否属于管理员组
	isMember, err := windows.Token(0).IsMember(sid)
	if err != nil {
		// 如果检查失败，没有管理员权限
		return false
	}

	return isMember
}

// KillProcessesByName 杀死所有名字为指定名称的进程
func KillProcessesByName(name string) error {
	// 将目标进程名称转换为小写以便精确匹配
	name = strings.ToLower(name)

	// 使用 tasklist 命令查找所有进程
	output, err := sys.Command("tasklist")
	if err != nil {
		return fmt.Errorf("failed to execute tasklist: %w", err)
	}

	// 按行解析 tasklist 输出
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		// 将每行转换为小写以便匹配
		lineLower := strings.ToLower(line)

		// 检查是否包含目标进程名称
		if strings.HasPrefix(lineLower, name+" ") || strings.Contains(lineLower, " "+name+" ") {
			fields := strings.Fields(line)
			if len(fields) > 1 {
				pid := fields[1] // 获取进程 ID

				// 使用 taskkill 命令杀死进程
				_, err = sys.Command("taskkill", "/F", "/PID", pid)
				if err != nil {
					return fmt.Errorf("failed to kill process %s (PID: %s): %w", name, pid, err)
				}
			}
		}
	}

	return nil
}
