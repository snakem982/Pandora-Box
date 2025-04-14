package internal

import (
	"fmt"
	"github.com/metacubex/bbolt"
	"github.com/metacubex/mihomo/component/profile/cachefile"
	"github.com/metacubex/mihomo/config"
	C "github.com/metacubex/mihomo/constant"
	"github.com/metacubex/mihomo/hub/executor"
	"github.com/metacubex/mihomo/log"
	plog "github.com/sirupsen/logrus"
	"github.com/snakem982/pandora-box/pandora/api/models"
	"github.com/snakem982/pandora-box/pandora/pkg/cache"
	"github.com/snakem982/pandora-box/pandora/pkg/constant"
	"github.com/snakem982/pandora-box/pandora/pkg/utils"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

// Init meta 启动前的初始化
func Init(homeDir string) {
	// 设置工作目录
	C.SetHomeDir(homeDir)

	// 设置日志输出目录
	logFilePath := filepath.Join(homeDir, "/logs/px.log")
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
	cache.BDb = cachefile.Cache().DB
	if cache.BDb == nil {
		os.Exit(1)
	}
	_ = cache.BDb.Batch(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(cache.BName)
		if err != nil {
			log.Warnln("[CacheFile] can't create bucket: %s", err.Error())
			return fmt.Errorf("create bucket: %v", err)
		}
		return nil
	})
}

var NowConfig *config.Config
var StartLock = sync.Mutex{}

// StartCore 函数用于启动核心功能，接收两个参数：profile和reload，分别为配置文件和是否自动reload的标志位
func StartCore(profile models.Profile, reload bool) {
	StartLock.Lock()
	defer StartLock.Unlock()

	templateBuf := models.PandoraDefaultConfig
	useTemplate := false
	path := profile.Path

	template, err := os.ReadFile(filepath.Join(C.Path.HomeDir(), constant.DefaultTemplate))
	if err == nil && len(template) > 0 {
		templateBuf = template
		on := cache.Get(constant.DefaultTemplate)
		if string(on) == "on" {
			useTemplate = true
		}
	}

	providerBuf, err := os.ReadFile(filepath.Join(C.Path.HomeDir(), path))
	if err != nil {
		log.Warnln("Read config error: %s", err.Error())
		return
	}

	rawCfg, err := config.UnmarshalRawConfig(providerBuf)
	if err != nil {
		log.Warnln("Unmarshal config error: %s", err.Error())
		return
	}

	if useTemplate || len(rawCfg.Rule) == 0 {
		provider := rawCfg.ProxyProvider
		proxy := rawCfg.Proxy
		rawCfg, _ = config.UnmarshalRawConfig(templateBuf)
		rawCfg.ProxyProvider = provider
		rawCfg.Proxy = proxy
	}

	rawCfg.Port = 0
	rawCfg.SocksPort = 0
	rawCfg.TProxyPort = 0
	rawCfg.RedirPort = 0
	if reload {
		general := NowConfig.General
		rawCfg.MixedPort = general.MixedPort
		rawCfg.AllowLan = general.AllowLan
		rawCfg.IPv6 = general.IPv6
		rawCfg.Tun.Enable = general.Tun.Enable
		rawCfg.UnifiedDelay = general.UnifiedDelay
	}

	rawCfg.ExternalController = ""
	rawCfg.GeodataMode = false
	rawCfg.Tun.Device = "Pandora"
	rawCfg.UnifiedDelay = true

	NowConfig, err = config.ParseRawConfig(rawCfg)
	if err != nil {
		log.Warnln("Parse config error: %s", err.Error())
		return
	}

	executor.ApplyConfig(NowConfig, !reload)
}
