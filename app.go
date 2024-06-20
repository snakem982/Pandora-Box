package main

import (
	"context"
	"github.com/keybase/go-keychain"
	"github.com/metacubex/mihomo/log"
	"os/exec"
	isadmin "pandora-box/backend/system/admin"
	"pandora-box/backend/system/open"
	"runtime"
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
	if runtime.GOOS == "darwin" {
		return "true"
	}

	return "false"
}

func (a *App) GetMacAcStatus() string {
	status, _ := GetAcStatus()
	return status
}

var KeyChainId = "Pandora-Box"

func GetAcStatus() (string, string) {

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

func (a *App) OpenConfigDirectory() {
	_, err := open.OpenConfigDirectory()
	if err != nil {
		log.Errorln("OpenConfigDirectory error:", err)
	}
}
