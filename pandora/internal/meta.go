package internal

import (
	"github.com/snakem982/pandora-box/pandora/pkg/constant"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/metacubex/mihomo/config"
	C "github.com/metacubex/mihomo/constant"
	"github.com/metacubex/mihomo/hub/executor"
	"github.com/metacubex/mihomo/log"
	plog "github.com/sirupsen/logrus"
	"github.com/snakem982/pandora-box/pandora/api/models"
	"github.com/snakem982/pandora-box/pandora/pkg/cache"
	"github.com/snakem982/pandora-box/pandora/pkg/utils"
)

// Init meta 启动前的初始化
func Init() {
	// 设置工作目录
	C.SetHomeDir(utils.GetUserHomeDir())

	// 设置日志输出目录
	logFilePath := utils.GetUserHomeDir("logs", "px.log")
	f, err := utils.CreateFile(logFilePath)
	if err != nil {
		return
	}
	if runtime.GOOS != "windows" {
		// 组合一下即可，os.Stdout代表标准输出流
		multiWriter := io.MultiWriter(os.Stdout, f)
		plog.SetOutput(multiWriter)
	} else {
		plog.SetOutput(f)
	}

	// 设置cache db
	db := cache.GetDBInstance()
	if db == nil {
		os.Exit(1)
	}

	// 输出日志
	log.Infoln("[CacheDB] initialized")
	log.Infoln("[HomePath] is %s", utils.GetUserHomeDir())

	// 释放资源文件
	_, _ = utils.SaveFile(utils.GetUserHomeDir("geoip.metadb"), GeoIp)
	_, _ = utils.SaveFile(utils.GetUserHomeDir("GeoSite.dat"), GeoSite)
	_, _ = utils.SaveFile(utils.GetUserHomeDir("ASN.mmdb"), ASN)
	GeoIp = nil
	GeoSite = nil
	ASN = nil
}

var NowConfig *config.Config
var StartLock = sync.Mutex{}

// StartCore 函数用于启动核心功能，接收两个参数：profile和reload，分别为配置文件和是否自动reload的标志位
func StartCore(profile models.Profile, reload bool) {
	StartLock.Lock()
	defer StartLock.Unlock()

	// 获取规则分组
	useTemplate, templateBuf := getTemplate(profile)

	// 获取配置文件
	providerBuf, err := os.ReadFile(filepath.Join(C.Path.HomeDir(), profile.Path))
	if err != nil {
		log.Warnln("Read config error: %s", err.Error())
		return
	}

	// 解析配置文件
	rawCfg, err := config.UnmarshalRawConfig(providerBuf)
	if err != nil {
		log.Warnln("Unmarshal config error: %s", err.Error())
		return
	}

	// 统一规则模板
	if useTemplate || len(rawCfg.Rule) == 0 {
		provider := rawCfg.ProxyProvider
		proxy := rawCfg.Proxy
		rawCfg, _ = config.UnmarshalRawConfig(templateBuf)
		rawCfg.ProxyProvider = provider
		rawCfg.Proxy = proxy
	}

	// Pandora-Box 默认配置
	rawCfg.Port = 0
	rawCfg.SocksPort = 0
	rawCfg.TProxyPort = 0
	rawCfg.RedirPort = 0
	rawCfg.ExternalController = ""
	rawCfg.GeodataMode = false
	rawCfg.Tun.Device = "Pandora"
	rawCfg.UnifiedDelay = true

	// todo 从数据库中获取 mihomo 配置

	if reload && NowConfig != nil {
		general := NowConfig.General
		rawCfg.MixedPort = general.MixedPort
		rawCfg.AllowLan = general.AllowLan
		rawCfg.IPv6 = general.IPv6
		rawCfg.Tun.Enable = general.Tun.Enable
	}

	NowConfig, err = config.ParseRawConfig(rawCfg)
	if err != nil {
		log.Warnln("[StartCore] error: %s", err.Error())
		return
	}

	// 应用配置
	executor.ApplyConfig(NowConfig, !reload)
}

// 获取统一规则分组模板
func getTemplate(profile models.Profile) (bool, []byte) {

	// 优先启用个性模板
	var template models.Template
	if profile.Template != "" {
		_ = cache.Get(profile.Template, &template)
	}
	if template.Path != "" {
		body, err := utils.ReadFile(utils.GetUserHomeDir(template.Path))
		if err == nil {
			return true, []byte(body)
		}
	}

	// 其次启用通用模板
	var list []models.Template
	_ = cache.GetList(constant.PrefixTemplate, &list)
	for _, m := range list {
		if m.Selected {
			template = m
			break
		}
	}
	if template.Path != "" {
		body, err := utils.ReadFile(utils.GetUserHomeDir(template.Path))
		if err == nil {
			return true, []byte(body)
		}
	}

	// 最后返回默认模板
	return false, Template_0
}
