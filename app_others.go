//go:build !darwin

package main

import (
	"context"
	C "github.com/metacubex/mihomo/constant"
	"github.com/metacubex/mihomo/log"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
	"pandora-box/backend/cache"
	"pandora-box/backend/constant"
	"pandora-box/backend/meta"
	isadmin "pandora-box/backend/system/admin"
	"pandora-box/backend/system/open"
	"pandora-box/backend/system/update"
	"pandora-box/backend/tools"
	"path/filepath"
)

// App struct
type App struct {
	ctx   context.Context
	addr  string
	apiOk bool
}

func NewApp(addr string) *App {
	return &App{
		addr: addr,
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) IsMac() string {
	return "false"
}

func (a *App) GetMacAcStatus() string {
	status, _ := GetAcStatus()
	return status
}

func GetAcStatus() (string, string) {
	return "1", ""
}

func (a *App) SetMacAc(pwd string) string {
	return "1"
}

func (a *App) IsAdmin() string {
	if isadmin.Check() {
		return "true"
	}

	return "false"
}

func (a *App) GetFreePort() string {
	return a.addr
}

func (a *App) GetSecret() string {
	value := cache.Get(constant.SecretKey)
	return string(value)
}

func (a *App) OpenConfigDirectory() {
	_, err := open.OpenConfigDirectory()
	if err != nil {
		log.Errorln("OpenConfigDirectory error:", err)
	}
}

func (a *App) IsUnifiedDelay() string {
	meta.StartLock.Lock()
	defer meta.StartLock.Unlock()

	if meta.NowConfig.General.UnifiedDelay {
		return "true"
	}

	return "false"
}

func (a *App) IsNeedUpdate() string {
	needUpdate, _ := update.IsNeedUpdate()
	if needUpdate {
		return "true"
	}

	return "false"
}

func (a *App) ImportConfig() string {
	selection, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择文件 Select File",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "文件类型 Type (*.zip)",
				Pattern:     "*.zip",
			},
		},
	})

	if err != nil || selection == "" {
		return "false"
	}

	secret := cache.Get(constant.SecretKey)

	err = meta.Recovery(selection)

	if err != nil {
		return err.Error()
	}

	if secret != nil {
		_ = cache.Put(constant.SecretKey, secret)
	}

	return "true"
}

func (a *App) ExportConfig() string {
	homeDir := filepath.Dir(C.Path.HomeDir())
	desktopPath := filepath.Join(homeDir, "Desktop")
	_, err := os.Stat(desktopPath)
	if !os.IsNotExist(err) {
		homeDir = desktopPath
	}

	selection, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:            "选择导出位置 Select Export Directory",
		DefaultDirectory: homeDir,
		DefaultFilename:  "Pandora-Box-Config.zip",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "文件类型 Type (*.zip)",
				Pattern:     "*.zip",
			},
		},
	})

	if err != nil || selection == "" {
		return "false"
	}

	err = meta.Dump(selection)

	if err != nil {
		return err.Error()
	}

	return "true"
}

func (a *App) SfQuit() {
	_ = cache.Put(constant.QuitSignal, []byte("1"))
	runtime.Quit(a.ctx)
}

func (a *App) ExportCrawl() string {

	homeDir := filepath.Dir(C.Path.HomeDir())
	desktopPath := filepath.Join(homeDir, "Desktop")
	_, err := os.Stat(desktopPath)
	if !os.IsNotExist(err) {
		homeDir = desktopPath
	}

	selection, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:            "选择导出位置 Select Export Directory",
		DefaultDirectory: homeDir,
		DefaultFilename:  "Pandora-Box-Share.yaml",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "文件类型 Type (*.yaml)",
				Pattern:     "*.yaml",
			},
		},
	})

	if err != nil || selection == "" {
		return "false"
	}

	err = tools.CopyFile(filepath.Join(C.Path.HomeDir(), constant.DefaultDownload), selection)
	if err != nil {
		return err.Error()
	}

	return "true"
}
