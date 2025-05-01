package traymenu

import (
	_ "embed"
	"fmt"
	"github.com/energye/systray"
	"github.com/sagernet/sing/common/json"
	"github.com/snakem982/pandora-box/api/models"
	sys "github.com/snakem982/pandora-box/pkg/sys/proxy"
	"github.com/snakem982/pandora-box/pkg/utils"
	webview "github.com/webview/webview_go"
	"pandora-box/static"
	"pandora-box/window"
	"sync"
)

//go:embed favicon.ico
var TrayIcon []byte

var t = GetI18nInstance().Translate

func Start(w webview.WebView) {
	// 启动前准备
	BindBeforeStart(w)
	// 开启托盘
	start, _ := systray.RunWithExternalLoop(func() {
		onReady(w)
	}, onExit)
	start()
}

// tray 操作锁
var trayLock sync.Mutex

// 保存下拉列表元素
var i18nMenuItem = make(map[string]*systray.MenuItem)

// 保存订阅
var profilesMenuItem = make(map[string]*systray.MenuItem)
var modes = []string{"rule", "global", "direct"}

func onReady(w webview.WebView) {
	systray.SetIcon(TrayIcon)
	systray.SetTooltip("Pandora-Box")
	systray.SetOnClick(func(menu systray.IMenu) {
		_ = menu.ShowMenu()
	})

	show := systray.AddMenuItem("显示窗口", "")
	show.Click(func() {
		window.ShowWindow()
	})
	i18nMenuItem["show"] = show
	systray.AddSeparator()

	rule := systray.AddMenuItemCheckbox("规则模式", "", true)
	i18nMenuItem["rule"] = rule

	global := systray.AddMenuItemCheckbox("全局模式", "", false)
	i18nMenuItem["global"] = global

	direct := systray.AddMenuItemCheckbox("直连模式", "", false)
	i18nMenuItem["direct"] = direct
	systray.AddSeparator()

	profiles := systray.AddMenuItem("订阅", "")
	i18nMenuItem["profiles"] = profiles
	systray.AddSeparator()

	proxy := systray.AddMenuItemCheckbox("系统代理", "", false)
	i18nMenuItem["proxy"] = proxy

	tun := systray.AddMenuItemCheckbox("TUN 模式", "", false)
	i18nMenuItem["tun"] = tun
	systray.AddSeparator()

	quit := systray.AddMenuItem("退出", "")
	quit.Click(func() {
		systray.Quit()
	})
	i18nMenuItem["quit"] = quit

	listenMode(w)
	listenProxy(w)
	listenTun(w)
}

func onExit() {
	sys.DisableProxy()
	utils.UnlockSingleton()
	port := static.Get("port")
	secret := static.Get("secret")
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", secret),
	}
	url := fmt.Sprintf("http://127.0.0.1:%s/mihomo/exit", port)
	_, _, _ = utils.SendGet(url, headers, "")
}

func BindBeforeStart(w webview.WebView) {
	// 监听语言切换
	_ = w.Bind("px_translate", func(lang string) {
		trayLock.Lock()
		defer trayLock.Unlock()

		for key, value := range i18nMenuItem {
			value.SetTitle(t(lang, key))
		}
		window.RefreshMenu(lang)
	})

	// 监听模式切换
	_ = w.Bind("px_mode", func(key string) {
		trayLock.Lock()
		defer trayLock.Unlock()

		if i18nMenuItem[key].Checked() {
			return
		}
		for _, mode := range modes {
			i18nMenuItem[mode].Uncheck()
		}
		i18nMenuItem[key].Check()
	})

	// 监听订阅
	_ = w.Bind("px_profiles", func(e interface{}) {
		trayLock.Lock()
		defer trayLock.Unlock()

		var p []models.Profile
		bytes, _ := json.Marshal(e)
		_ = json.Unmarshal(bytes, &p)

		for _, item := range profilesMenuItem {
			if item == nil {
				continue
			}
			item.Hide()
			item = nil
		}

		for _, profile := range p {
			direct1 := i18nMenuItem["profiles"].AddSubMenuItemCheckbox(profile.Title, "", profile.Selected)
			direct1.Click(func() {
				if profile.Selected {
					return
				}
				w.Dispatch(func() {
					w.Eval(getJsCode("switchProfiles", profile))
				})
			})
			profilesMenuItem[profile.Id] = direct1
		}
	})

	// 监听系统代理
	_ = w.Bind("px_proxy", func(e bool) {
		trayLock.Lock()
		defer trayLock.Unlock()

		if e {
			i18nMenuItem["proxy"].Check()
		} else {
			i18nMenuItem["proxy"].Uncheck()
		}
	})

	// 监听 TUN 模式
	_ = w.Bind("px_tun", func(e bool) {
		trayLock.Lock()
		defer trayLock.Unlock()

		if e {
			i18nMenuItem["tun"].Check()
		} else {
			i18nMenuItem["tun"].Uncheck()
		}
	})

	// 关闭
	_ = w.Bind("px_close", func(e bool) {
		systray.Quit()
	})
}

func toJsonAsObject(values ...interface{}) string {
	obj := make(map[string]interface{})

	// 单独处理，如果只有一个参数，直接用这个参数，而不是包在数组
	for i, v := range values {
		key := fmt.Sprintf("%d", i)
		obj[key] = v
	}

	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

func getJsCode(name string, values ...interface{}) string {
	return fmt.Sprintf("window[\"px_%s\"](%v);", name, toJsonAsObject(values...))
}

// 监听模式切换
func listenMode(w webview.WebView) {
	for _, mode := range modes {
		i18nMenuItem[mode].Click(func() {
			trayLock.Lock()
			defer trayLock.Unlock()

			if i18nMenuItem[mode].Checked() {
				return
			}
			w.Dispatch(func() {
				w.Eval(getJsCode("switchMode", mode))
			})
		})
	}
}

// 监听系统代理
func listenProxy(w webview.WebView) {
	i18nMenuItem["proxy"].Click(func() {
		w.Dispatch(func() {
			w.Eval(getJsCode("switchProxy"))
		})
	})
}

// 监听 TUN 模式
func listenTun(w webview.WebView) {
	i18nMenuItem["tun"].Click(func() {
		w.Dispatch(func() {
			w.Eval(getJsCode("switchTun"))
		})
	})
}
