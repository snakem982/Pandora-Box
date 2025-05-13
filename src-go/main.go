package main

import (
	"flag"
	"github.com/metacubex/mihomo/hub/executor"
	"github.com/metacubex/mihomo/log"
	"github.com/snakem982/pandora-box/pandora"
	sys "github.com/snakem982/pandora-box/pkg/sys/proxy"
	"github.com/snakem982/pandora-box/pkg/utils"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	// 回调地址
	addr := flag.String("addr", "", "callback address")

	// 解析命令行参数
	flag.Parse()

	if addr == nil || *addr == "" {
		os.Exit(1)
	}

	// 保持单例
	if utils.NotSingleton("px-server.pid") {
		os.Exit(1)
	}

	// 初始化工作目录
	pandora.Init(false)

	// 开启后端api
	pandora.StartCore(*addr, false)

	termSign := make(chan os.Signal, 1)
	signal.Notify(termSign, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-termSign:
		log.Warnln("received termination signal")
		pandora.Release()
		utils.UnlockSingleton()
		executor.Shutdown()
		sys.DisableProxy()
	}

}
