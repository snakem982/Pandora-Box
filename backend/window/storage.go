package window

import (
	"github.com/snakem982/pandora-box/pkg/cache"
	webview "github.com/webview/webview_go"
)

func Storage(w webview.WebView) {
	_ = w.Bind("pxGetItem", getItem)
	_ = w.Bind("pxSetItem", setItem)
}

func getItem(key string) string {
	var value string
	err := cache.Get(key, &value)
	if err != nil {
		return ""
	}

	return value
}

func setItem(key, val string) {
	_ = cache.Put(key, val)
}
