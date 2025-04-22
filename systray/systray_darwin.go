package systray

import (
	_ "embed"
	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed icon/icon-128.png
var Icon []byte

func Run(app *application.App, systemTray *application.SystemTray, window *application.WebviewWindow) {
	systemTray.SetIcon(Icon)
	myMenu := app.NewMenu()

	item := myMenu.Add("显示窗口")
	item.OnClick(func(ctx *application.Context) {
		window.Show()
	})

	myMenu.AddSeparator()

	// Callbacks can be shared. This is useful for radio groups
	radioCallback := func(ctx *application.Context) {
		menuItem := ctx.ClickedMenuItem()
		menuItem.SetLabel(menuItem.Label() + "!")
	}

	// Radio groups are created implicitly by placing radio items next to each other in a menu
	myMenu.AddRadio("规则模式", true).OnClick(radioCallback)
	myMenu.AddRadio("全局模式", false).OnClick(radioCallback)
	myMenu.AddRadio("直连模式", false).OnClick(radioCallback)

	myMenu.AddSeparator()

	subMenu := myMenu.AddSubmenu("订阅")
	subMenu.Add("订阅1").OnClick(func(ctx *application.Context) {
		ctx.ClickedMenuItem().SetLabel("Clicked!")
	})
	subMenu.Add("订阅2").OnClick(func(ctx *application.Context) {
		ctx.ClickedMenuItem().SetLabel("Clicked!")
	})
	subMenu.Add("订阅3").OnClick(func(ctx *application.Context) {
		ctx.ClickedMenuItem().SetLabel("Clicked!")
	})
	myMenu.AddSeparator()

	item2 := myMenu.Add("系统代理")
	item2.OnClick(func(ctx *application.Context) {
		item2.SetChecked(!ctx.ClickedMenuItem().Checked())
	})

	item3 := myMenu.Add("TUN 模式")
	item3.OnClick(func(ctx *application.Context) {
		item3.SetChecked(!ctx.ClickedMenuItem().Checked())
	})
	myMenu.AddSeparator()

	myMenu.Add("退出").OnClick(func(ctx *application.Context) {
		app.Quit()
	})

	systemTray.SetMenu(myMenu)
	systemTray.WindowOffset(2)
}
