package main

import (
	"embed"
	_ "embed"
	"github.com/snakem982/pandora-box/pandora"
	"github.com/snakem982/pandora-box/systray"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
	"log"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {

	pandora.StartCore()

	app := application.New(application.Options{
		Name:        "Pandora-Box",
		Description: "A Simple Mihomo Gui",
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ActivationPolicy: application.ActivationPolicyRegular,
			ApplicationShouldTerminateAfterLastWindowClosed: false,
		},
		Icon: systray.Icon,
	})

	systemTray := app.NewSystemTray()

	window := app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title: "Pandora-Box",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 80,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		URL:       "/?host=127.0.0.1&port=9686&secret=Y8IUaPeFLTRvsrdf2mUJkLMBuphVZRE5",
		Width:     1100,
		Height:    760,
		MinWidth:  960,
		MinHeight: 660,
		Windows: application.WindowsWindow{
			HiddenOnTaskbar: true,
		},
	})

	// 处理窗口显示
	window.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
		window.Hide()
		e.Cancel()
	})
	app.OnApplicationEvent(events.Mac.ApplicationShouldHandleReopen, func(event *application.ApplicationEvent) {
		window.Show()
	})

	systray.Run(app, systemTray, window)

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}

}
