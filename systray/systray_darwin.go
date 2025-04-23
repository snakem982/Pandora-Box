package systray

import (
	_ "embed"
	"encoding/json"

	"github.com/snakem982/pandora-box/pandora/api/models"
	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed icon/icon-128.png
var Icon []byte

var t = GetI18nInstance().Translate

// 保存下拉列表元素
var i18nMenuItem = make(map[string]*application.MenuItem)

func Run(app *application.App, systemTray *application.SystemTray, window *application.WebviewWindow) {
	systemTray.SetIcon(Icon)
	systemTray.SetTooltip("Pandora-Box")

	myMenu := app.NewMenu()

	show := myMenu.Add("显示窗口")
	show.OnClick(func(ctx *application.Context) {
		window.Show()
	})
	i18nMenuItem["show"] = show
	myMenu.AddSeparator()

	rule := myMenu.AddRadio("规则模式", true)
	rule.OnClick(func(ctx *application.Context) {
		menuItem := ctx.ClickedMenuItem()
		if menuItem.Checked() {
			return
		}
		app.EmitEvent("switchMode", "rule")
	})
	i18nMenuItem["rule"] = rule

	global := myMenu.AddRadio("全局模式", false)
	global.OnClick(func(ctx *application.Context) {
		menuItem := ctx.ClickedMenuItem()
		if menuItem.Checked() {
			return
		}
		app.EmitEvent("switchMode", "global")
	})
	i18nMenuItem["global"] = global

	direct := myMenu.AddRadio("直连模式", false)
	direct.OnClick(func(ctx *application.Context) {
		menuItem := ctx.ClickedMenuItem()
		if menuItem.Checked() {
			return
		}
		app.EmitEvent("switchMode", "direct")
	})
	i18nMenuItem["direct"] = direct
	myMenu.AddSeparator()

	profiles := myMenu.AddSubmenu("订阅")
	i18nMenuItem["profiles"] = myMenu.FindByLabel("订阅")
	myMenu.AddSeparator()

	proxy := myMenu.Add("系统代理")
	proxy.OnClick(func(ctx *application.Context) {
		proxy.SetChecked(!ctx.ClickedMenuItem().Checked())
	})
	i18nMenuItem["proxy"] = proxy

	tun := myMenu.Add("TUN 模式")
	tun.OnClick(func(ctx *application.Context) {
		tun.SetChecked(!ctx.ClickedMenuItem().Checked())
	})
	i18nMenuItem["tun"] = tun
	myMenu.AddSeparator()

	quit := myMenu.Add("退出")
	quit.OnClick(func(ctx *application.Context) {
		app.Quit()
	})
	i18nMenuItem["quit"] = quit

	systemTray.SetMenu(myMenu)
	systemTray.WindowOffset(2)

	listenMode(app)
	listenTranslate(app)
	listenProfiles(app, myMenu, profiles)
}

// 监听模式切换
func listenMode(app *application.App) {
	// Custom event handling
	app.OnEvent("mode", func(e *application.CustomEvent) {
		key := e.Data.(string)
		i18nMenuItem["rule"].SetChecked(false)
		i18nMenuItem["global"].SetChecked(false)
		i18nMenuItem["direct"].SetChecked(false)
		i18nMenuItem[key].SetChecked(true)
	})
}

// 监听语言切换
func listenTranslate(app *application.App) {
	// Custom event handling
	app.OnEvent("translate", func(e *application.CustomEvent) {
		lang := e.Data.(string)
		for key, value := range i18nMenuItem {
			value.SetLabel(t(lang, key))
		}
	})
}

// 监听订阅
func listenProfiles(app *application.App, myMenu, profiles *application.Menu) {
	// Custom event handling
	app.OnEvent("profiles", func(e *application.CustomEvent) {
		var p []models.Profile
		bytes, _ := json.Marshal(e.Data)
		_ = json.Unmarshal(bytes, &p)

		profiles.Clear()

		for _, profile := range p {
			direct1 := profiles.AddRadio(profile.Title, profile.Selected)
			direct1.OnClick(func(ctx *application.Context) {
				if profile.Selected {
					return
				}
				app.EmitEvent("switchProfiles", profile)
			})
		}

		profiles.Update()
		myMenu.Update()
	})

}
