package main

import (
	"flag"
	"fmt"
	"github.com/metacubex/mihomo/hub/executor"
	"github.com/snakem982/pandora-box"
	"github.com/snakem982/pandora-box/pkg/utils"
	webview "github.com/webview/webview_go"
	"os"
	"os/signal"
	"pandora-box/traymenu"
	"pandora-box/window"
	"syscall"
)

func main() {

	// 是否调试模式 -debug，默认是 false
	debug := flag.Bool("debug", false, "enable debug mode")
	// 是否后台运行 -back，默认是 false
	background := flag.Bool("back", false, "enable background")
	// 地址
	addr := flag.String("addr", "", "enable address")

	// 解析命令行参数
	flag.Parse()

	// 加载后端，成功发送数据
	if *background {
		// 保持单例
		if utils.NotSingleton("pandora.pid") {
			os.Exit(1)
		}

		// 开启后端api
		pandora.StartCore(*addr, false)

		termSign := make(chan os.Signal, 1)
		signal.Notify(termSign, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-termSign:
			executor.Shutdown()
			utils.UnlockSingleton()
			pandora.Release()
		}
		return
	}

	// 保持单例
	if utils.NotSingleton("pandora-box.pid") {
		os.Exit(1)
	}

	// 获取网页地址
	var url string
	if *debug {
		port, secret := pandora.StartCore("", true)
		url = fmt.Sprintf("http://localhost:1420/?port=%d&secret=%s", port, secret)
	} else {
		// 初始化工作目录
		pandora.Init(true)

		// 启动api
		url = window.TryAdmin()
	}

	// webview初始化
	w := webview.New(*debug)
	defer w.Destroy()
	// 开启托盘
	traymenu.Start(w)
	// 绑定窗口函数
	window.Init(w)
	// storage
	window.Storage(w)
	// 进行地址加载
	w.Navigate(url)

	// 捕获 Ctrl+C 和 kill信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-c
		os.Exit(0)
	}()

	w.Run()

}
