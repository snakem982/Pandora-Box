package main

import (
	"flag"
	"github.com/metacubex/mihomo/hub/executor"
	"github.com/metacubex/mihomo/log"
	"github.com/snakem982/pandora-box/pandora"
	sys "github.com/snakem982/pandora-box/pkg/sys/proxy"
	"github.com/snakem982/pandora-box/pkg/utils"
	"go.uber.org/automaxprocs/maxprocs"
	"net/url"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	// 优化线程资源配置
	_, _ = maxprocs.Set(maxprocs.Logger(func(string, ...any) {}))

	// 回调地址
	addr := flag.String("addr", "", "callback address")
	home := flag.String("home", "", "home directory")
	proxyUID := flag.String("uid", "", "desktop user ID for system proxy settings")

	// 解析命令行参数
	flag.Parse()
	if *proxyUID != "" {
		if err := os.Setenv("PX_PROXY_UID", *proxyUID); err != nil {
			panic(err)
		}
	}

	if addr == nil || *addr == "" {
		panic("callback address is required")
	}

	if home == nil || *home == "" {
		panic("home directory is required")
	}

	homeDir, err := url.QueryUnescape(*home)
	if err != nil {
		panic(err)
	}

	// 设置工作目录
	utils.InitHomeDir(homeDir)

	// 保持单例
	if utils.NotSingleton("px-server.pid") {
		os.Exit(1)
	}

	// 初始化工作目录
	pandora.Init()

	// 开启后端api
	pandora.StartCore(*addr)

	termSign := make(chan os.Signal, 1)
	signal.Notify(termSign, syscall.SIGINT, syscall.SIGTERM)
	<-termSign
	log.Warnln("received termination signal")
	pandora.Release()
	utils.UnlockSingleton()
	executor.Shutdown()
	sys.DisableProxy()

}
