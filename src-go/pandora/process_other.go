//go:build !darwin

package pandora

func InitProcessGroup() {
	// 非 macOS 平台不执行 darwin 专属的进程组/会话脱离逻辑
}
