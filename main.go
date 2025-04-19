package main

import (
	"embed"
	_ "embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

func main() {

	app := application.New(application.Options{
		Name:        "Pandora-Box",
		Description: "A Simple Mihomo Gui",
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		Icon: icon,
	})

	app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title: "Pandora-Box",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 80,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		URL:       "/?host=127.0.0.1&port=25834&secret=Y8IUaPeFLTRvsrdf2mUJkLMBuphVZRE5",
		Width:     1100,
		Height:    760,
		MinWidth:  960,
		MinHeight: 660,
	})

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
