package api

import (
	"os"
	"time"

	"github.com/metacubex/mihomo/hub/executor"
	"github.com/snakem982/pandora-box/pkg/cache"
	sys "github.com/snakem982/pandora-box/pkg/sys/proxy"
	"github.com/snakem982/pandora-box/pkg/utils"
)

var Version string

var ControllerPort int

// Exit 退出px
func Exit() {
	sys.DisableProxy()
	cache.Close()
	executor.Shutdown()
	utils.UnlockSingleton()
	// 给异步日志刷盘留出缓冲时间
	time.Sleep(200 * time.Millisecond)
	os.Exit(0)
}
