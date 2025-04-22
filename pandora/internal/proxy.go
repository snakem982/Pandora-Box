package internal

import (
	"fmt"
	"github.com/snakem982/pandora-box/pandora/api/models"
	"github.com/snakem982/pandora-box/pandora/pkg/cache"
	"github.com/snakem982/pandora-box/pandora/pkg/constant"
	sys "github.com/snakem982/pandora-box/pandora/pkg/sys/proxy"
)

// GetProxyUrl 获取代理
func GetProxyUrl() string {
	// 从系统获取
	addr, err := sys.GetHttp()
	if err == nil && addr != nil {
		return fmt.Sprintf("http://%s:%d", addr.Host, addr.Port)
	}

	// 从数据库中获取
	var mi models.Mihomo
	_ = cache.Get(constant.Mihomo, &mi)
	if mi.BindAddress != "" {
		return fmt.Sprintf("http://%s:%d", mi.BindAddress, mi.Port)
	}

	// 都获取不到返回空
	return ""
}
