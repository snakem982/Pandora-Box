//go:build !darwin

package main

import (
	"context"
	"github.com/metacubex/mihomo/log"
	"pandora-box/backend/meta"
	isadmin "pandora-box/backend/system/admin"
	"pandora-box/backend/system/open"
	"runtime"
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
