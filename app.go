//go:build darwin

package main

import (
	"bytes"
	"context"
	"github.com/keybase/go-keychain"
	C "github.com/metacubex/mihomo/constant"
	"github.com/metacubex/mihomo/log"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
	"os/exec"
	"pandora-box/backend/cache"
	"pandora-box/backend/constant"
	"pandora-box/backend/meta"
	isadmin "pandora-box/backend/system/admin"
	"pandora-box/backend/system/open"
	"pandora-box/backend/system/update"
	"pandora-box/backend/tools"
	"path/filepath"
	"strings"
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
	return "true"
}

func (a *App) GetMacAcStatus() string {
	status, _ := GetAcStatus()
	return status
}

var KeyChainId = "Pandora-Box"

func GetAcStatus() (string, string) {

	if *devFlag {
		return "1", ""
	}

	query := keychain.NewItem()
	query.SetSecClass(keychain.SecClassGenericPassword)
	query.SetService(KeyChainId)
	query.SetAccount(KeyChainId)
	query.SetAccessGroup(KeyChainId)
	query.SetMatchLimit(keychain.MatchLimitOne)
	query.SetReturnData(true)
	results, err := keychain.QueryItem(query)
	// 密码不存在 1:需要授权
	if err != nil {
		return "1", ""
	} else if len(results) != 1 {
		return "1", ""
	}

	password := string(results[0].Data)

	// 密码存在错误 2：重新授权
	// 校验密码
	cmd := exec.Command("sudo", "-S", "echo", "ok")
	cmd.Stdin = strings.NewReader(password)
	err = cmd.Run()
	if err != nil {
		return "2", ""
	}

	// 3：授权校验通过
	return "3", password
}

func (a *App) GetClipboard() string {
	// 定义 osascript 命令
	cmd := exec.Command("osascript", "-e", "return (the clipboard)")

	// 捕获输出
	var out bytes.Buffer
	cmd.Stdout = &out

	// 执行命令并检查错误
	err := cmd.Run()
	if err != nil {
		return ""
	}

	return out.String()
}

func (a *App) SetMacAc(pwd string) string {
	// 校验密码 1:密码错误
	cmd := exec.Command("sudo", "-S", "echo", "ok")
	cmd.Stdin = strings.NewReader(pwd)
	err := cmd.Run()
	if err != nil {
		return "1"
	}

	// 存储密码 2:存储密码错误
	item := keychain.NewItem()
	item.SetSecClass(keychain.SecClassGenericPassword)
	item.SetService(KeyChainId)
	item.SetAccount(KeyChainId)
	_ = keychain.DeleteItem(item)
	item = keychain.NewItem()
	item.SetSecClass(keychain.SecClassGenericPassword)
	item.SetService(KeyChainId)
	item.SetAccount(KeyChainId)
	item.SetAccessGroup(KeyChainId)
	item.SetLabel("Pandora-Box.app")
	item.SetData([]byte(pwd))
	item.SetSynchronizable(keychain.SynchronizableNo)
	item.SetAccessible(keychain.AccessibleWhenUnlocked)
	err = keychain.AddItem(item)

	if err != nil {
		return "2"
	}

	return "3"
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
	desktopPath := filepath.Join(homeDir, "Downloads")
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
	runtime.Quit(a.ctx)
}

func (a *App) ExportCrawl() string {

	homeDir := filepath.Dir(C.Path.HomeDir())
	desktopPath := filepath.Join(homeDir, "Downloads")
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
