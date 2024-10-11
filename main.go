package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"github.com/metacubex/mihomo/hub/executor"
	"github.com/metacubex/mihomo/hub/route"
	"github.com/metacubex/mihomo/log"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"os"
	"os/exec"
	"pandora-box/backend/api"
	"pandora-box/backend/cache"
	"pandora-box/backend/constant"
	"pandora-box/backend/meta"
	IsAdmin "pandora-box/backend/system/admin"
	"pandora-box/backend/system/proxy"
	"pandora-box/backend/tools"
	"runtime"
	"strings"
	"time"
)

var devFlag = flag.Bool("dev", false, "布尔类型参数")

func init() {
	flag.Parse()
}

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/540x540.png
var icon []byte

func main() {

	if !*devFlag && runtime.GOOS == "darwin" && !IsAdmin.Check() {
		status, pwd := GetAcStatus()
		if status == "3" {
			startMacInAdmin(pwd)
			return
		}
	}

	meta.Init()

	log.Infoln("Pandora-Box %s %s %s with %s",
		constant.PandoraVersion, runtime.GOOS, runtime.GOARCH, runtime.Version())

	route.Register(api.Hello)
	route.Register(api.Version)
	route.Register(api.Profile)
	route.Register(api.Getter)
	route.Register(api.System)
	route.Register(api.Ignore)
	route.Register(api.MyRules)

	addr := startHttpApi()

	meta.SwitchProfile(false)

	app := NewApp(addr)

	option := &options.App{
		Title:     "Pandora-Box",
		Width:     1200,
		Height:    780,
		MinWidth:  925,
		MinHeight: 675,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.startup,
		OnShutdown: func(ctx context.Context) {
			executor.Shutdown()
			proxy.RemoveProxy()
		},
		Bind: []interface{}{
			app,
		},
	}

	if runtime.GOOS == "darwin" {
		AppMenu := menu.NewMenu()
		AppMenu.Append(menu.AppMenu())
		AppMenu.Append(menu.EditMenu())
		option.Menu = AppMenu
		option.HideWindowOnClose = true
		option.Mac = &mac.Options{
			TitleBar: mac.TitleBarHidden(),
			About: &mac.AboutInfo{
				Title:   constant.PandoraVersion,
				Message: "Copyright © 2024 snakem982",
				Icon:    icon,
			},
		}
		option.CSSDragProperty = "widows"
		option.CSSDragValue = "1"
	}

	err := wails.Run(option)
	if err != nil {
		log.Errorln("wails.Run Error:", err)
	}
}

func startHttpApi() (addr string) {
	var secret string
	value := cache.Get(constant.SecretKey)
	if value != nil {
		secret = string(value)
	} else {
		secret = tools.String(32)
		_ = cache.Put(constant.SecretKey, []byte(secret))
	}
	addr = route.StartByPandora(secret)
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", secret),
	}
	timeOut := 500 * time.Millisecond
	for i := 0; i < 3; i++ {
		okUrl := fmt.Sprintf("http://%s/ok", addr)
		body, _, err := tools.HttpGetWithTimeout(okUrl, timeOut, false, headers)
		if err == nil && string(body) == "ok" {
			log.Infoln("Start Http Serve Success.Addr is %s", addr)
			break
		} else {
			log.Errorln("Start Http Serve Error: %v.Addr is %s", err, addr)
		}

		time.Sleep(timeOut)
	}

	return
}

func startMacInAdmin(pwd string) {
	exePath, err := os.Executable()
	if err != nil {
		log.Errorln("get exe path error：%s", err.Error())
		return
	}

	cmd := exec.Command("sudo", "-b", exePath, ">", "/dev/null")
	cmd.Stdin = strings.NewReader(pwd)
	err = cmd.Run()
	if err != nil {
		log.Errorln("cmd.Run() error：%s", err.Error())
	}
}
