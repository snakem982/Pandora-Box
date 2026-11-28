package main

import (
	"flag"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/metacubex/mihomo/hub/executor"
	"github.com/metacubex/mihomo/log"
	"github.com/snakem982/pandora-box/pandora"
	sys "github.com/snakem982/pandora-box/pkg/sys/proxy"
	"github.com/snakem982/pandora-box/pkg/utils"
	"go.uber.org/automaxprocs/maxprocs"
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

	// 明确忽略 SIGHUP，防止授权会话或终端断开时误杀后台进程
	signal.Ignore(syscall.SIGHUP)

	termSign := make(chan os.Signal, 2)
	signal.Notify(termSign, syscall.SIGINT, syscall.SIGTERM)
	sig := <-termSign
	log.Warnln("received termination signal: %s (%d)", sig.String(), sig)
	sys.DisableProxy()
	executor.Shutdown()
	pandora.Release()
	utils.UnlockSingleton()
}
