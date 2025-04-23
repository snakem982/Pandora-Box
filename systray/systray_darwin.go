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
	i18nMenuItem["rule"] = rule

	global := myMenu.AddRadio("全局模式", false)
	i18nMenuItem["global"] = global

	direct := myMenu.AddRadio("直连模式", false)
	i18nMenuItem["direct"] = direct
	myMenu.AddSeparator()

	profiles := myMenu.AddSubmenu("订阅")
	i18nMenuItem["profiles"] = myMenu.FindByLabel("订阅")
	myMenu.AddSeparator()

	proxy := myMenu.AddCheckbox("系统代理", false)
	i18nMenuItem["proxy"] = proxy

	tun := myMenu.AddCheckbox("TUN 模式", false)
	i18nMenuItem["tun"] = tun
	myMenu.AddSeparator()

	quit := myMenu.Add("退出")
	quit.OnClick(func(ctx *application.Context) {
		app.Quit()
	})
	i18nMenuItem["quit"] = quit

	systemTray.SetMenu(myMenu)
	systemTray.WindowOffset(2)

	listenTranslate(app)
	listenMode(app, myMenu)
	listenProfiles(app, myMenu, profiles)
	listenProxy(app, proxy)
	listenTun(app, myMenu, tun)
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

// 监听模式切换
func listenMode(app *application.App, myMenu *application.Menu) {

	modes := []string{"rule", "global", "direct"}
	now := "rule"
	for _, mode := range modes {
		i18nMenuItem[mode].OnClick(func(ctx *application.Context) {
			if now == mode {
				return
			}
			app.EmitEvent("switchMode", mode)
		})
	}

	// Custom event handling
	app.OnEvent("mode", func(e *application.CustomEvent) {
		key := e.Data.(string)
		if now == key {
			return
		}
		for _, mode := range modes {
			i18nMenuItem[mode].SetChecked(false)
		}
		i18nMenuItem[key].SetChecked(true)
		now = key
		myMenu.Update()
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

// 监听系统代理
func listenProxy(app *application.App, proxy *application.MenuItem) {
	// Custom event handling
	app.OnEvent("proxy", func(e *application.CustomEvent) {
		proxy.SetChecked(e.Data.(bool))
	})
	proxy.OnClick(func(ctx *application.Context) {
		app.EmitEvent("switchProxy")
	})
}

// 监听 TUN 模式
func listenTun(app *application.App, myMenu *application.Menu, tun *application.MenuItem) {
	// Custom event handling
	app.OnEvent("tun", func(e *application.CustomEvent) {
		tun.SetChecked(e.Data.(bool))
		myMenu.Update()
	})
	tun.OnClick(func(ctx *application.Context) {
		app.EmitEvent("switchTun")
	})
}
