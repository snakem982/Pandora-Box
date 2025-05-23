package internal

import (
	"github.com/metacubex/mihomo/hub/executor"
	"github.com/metacubex/mihomo/tunnel"
	"github.com/snakem982/pandora-box/pkg/constant"
	sysProxy "github.com/snakem982/pandora-box/pkg/sys/proxy"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/metacubex/mihomo/config"
	C "github.com/metacubex/mihomo/constant"
	"github.com/metacubex/mihomo/log"
	plog "github.com/sirupsen/logrus"
	"github.com/snakem982/pandora-box/api/models"
	"github.com/snakem982/pandora-box/pkg/cache"
	"github.com/snakem982/pandora-box/pkg/utils"
)

// Init meta 启动前的初始化
func Init() {
	// 设置工作目录
	C.SetHomeDir(utils.GetUserHomeDir())

	// 设置日志输出目录
	logName := "px-server.log"
	logFilePath := utils.GetUserHomeDir("logs", logName)
	f, err := utils.CreateFileForAppend(logFilePath)
	if err != nil {
		return
	}

	// 组合一下即可，os.Stdout代表标准输出流
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
	cache.GetMetaDB()

	// 输出日志
	log.Infoln("[CacheDB] initialized")
	log.Infoln("[HomePath] is %s", utils.GetUserHomeDir())

	// 修改权限
	pathTemp := utils.GetUserHomeDir("logs", "px-client.log")
	_ = utils.SetPermissions(pathTemp)
	pathTemp = utils.GetUserHomeDir("px-electron.db", "config.json")
	_ = utils.SetPermissions(pathTemp)
	log.Infoln("[Permission] is ok")

	// 释放资源文件
	_, _ = utils.SaveFile(utils.GetUserHomeDir("geoip.metadb"), GeoIp)
	_, _ = utils.SaveFile(utils.GetUserHomeDir("GeoSite.dat"), GeoSite)
	_, _ = utils.SaveFile(utils.GetUserHomeDir("ASN.mmdb"), ASN)

	GeoIp = nil
	GeoSite = nil
	ASN = nil
}

var NowConfig *config.Config
var havaStartCore bool
var StartLock = sync.Mutex{}

// startCore 函数用于启动核心功能
func startCore(profile models.Profile, reload bool) {

	// 获取规则分组
	useTemplate, templateBuf := getTemplate(profile)

	// 获取配置文件
	providerBuf, err := os.ReadFile(filepath.Join(C.Path.HomeDir(), profile.Path))
	if err != nil {
		log.Warnln("Read config error: %s", err.Error())
		return
	}

	// 解析配置文件1
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
	rawCfg.ExternalUI = ""
	rawCfg.ExternalUIURL = ""
	rawCfg.Tun.DNSHijack = []string{"any:53"}
	rawCfg.Tun.AutoRoute = true
	rawCfg.Tun.AutoDetectInterface = true
	rawCfg.Tun.Device = "Pandora"
	rawCfg.UnifiedDelay = true

	// 从数据库中获取 mihomo 配置,进行 rawCfg 赋值
	var mi models.Mihomo
	_ = cache.Get(constant.Mihomo, &mi)
	if mi.BindAddress == "" {
		mi = models.Mihomo{
			Mode:        "rule",
			Proxy:       false,
			Tun:         false,
			Port:        9697,
			BindAddress: "127.0.0.1",
			Stack:       "Mixed",
			Dns:         false,
			Ipv6:        false,
		}
	}
	rawCfg.Mode = tunnel.ModeMapping[mi.Mode]
	rawCfg.AllowLan = true
	rawCfg.MixedPort = mi.Port
	rawCfg.BindAddress = mi.BindAddress
	rawCfg.Tun.Stack = C.StackTypeMapping[strings.ToLower(mi.Stack)]
	rawCfg.IPv6 = mi.Ipv6

	// 保存规则数
	_ = cache.Put("Rule_No", len(rawCfg.Rule))

	// 解析配置文件2
	NowConfig, err = config.ParseRawConfig(rawCfg)
	if err != nil {
		log.Errorln("ParseRawConfig error: %v", err)
		return
	}

	// 覆盖dns
	if mi.Dns {
		var dns models.Dns
		_ = cache.Get(constant.Dns, &dns)

		if dns.Content == "" {
			dns.Content = DefaultDNS
		}

		cfg, _ := executor.ParseWithBytes([]byte(dns.Content))
		NowConfig.DNS = cfg.DNS
	}

	// 应用配置
	if reload {
		NowConfig.General.Tun.Enable = mi.Tun
	} else {
		// 检测端口占用
		err = utils.IsPortAvailable(mi.BindAddress, mi.Port)
		if err != nil {
			log.Errorln("IsPortAvailable error: %v", err)
			mi.Port, _ = utils.GetRandomPort(mi.BindAddress)
			NowConfig.General.MixedPort = mi.Port
		}

		// 初次加载不能开启tun,不然在windows上会崩
		NowConfig.General.Tun.Enable = false
	}

	// 激活配置
	go executor.ApplyConfig(NowConfig, !reload)

	// 代理开启
	if mi.Proxy {
		_ = sysProxy.EnableProxy(mi.BindAddress, mi.Port)
	}
	// 存储配置
	_ = cache.Put(constant.Mihomo, mi)
	// 更新启动标志
	havaStartCore = true
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

// SwitchProfile 切换配置
func SwitchProfile(reload bool) {
	StartLock.Lock()
	defer StartLock.Unlock()

	// 应用配置
	var profile models.Profile

	// 获取切换配置
	var profiles []models.Profile
	_ = cache.GetList(constant.PrefixProfile, &profiles)

	if len(profiles) == 0 {
		return
	}

	haveSelected := false
	for _, p := range profiles {
		if p.Selected {
			profile = p
			haveSelected = true
		}
	}

	if !haveSelected {
		profile = profiles[0]
		profile.Selected = true
		_ = cache.Put(profile.Id, profile)
	}

	if !havaStartCore {
		reload = false
	}

	startCore(profile, reload)
}
