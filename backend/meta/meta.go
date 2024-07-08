package meta

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/metacubex/mihomo/component/profile/cachefile"
	"github.com/metacubex/mihomo/config"
	C "github.com/metacubex/mihomo/constant"
	"github.com/metacubex/mihomo/hub/executor"
	"github.com/metacubex/mihomo/log"
	"github.com/metacubex/mihomo/tunnel"
	"github.com/sagernet/bbolt"
	plog "github.com/sirupsen/logrus"
	"io"
	"os"
	"pandora-box/backend/cache"
	"pandora-box/backend/constant"
	"pandora-box/backend/resolve"
	isadmin "pandora-box/backend/system/admin"
	"pandora-box/backend/tools"
	"path/filepath"
	"runtime"
	"strings"
)

//go:embed geoip.metadb
var metadata []byte

//go:embed Profile_0.yaml
var DefProfile []byte

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		currentDir, _ := os.Getwd()
		homeDir = filepath.Join(currentDir, "Pandora-Box")
	} else {
		homeDir = filepath.Join(homeDir, "Pandora-Box")
	}
	C.SetHomeDir(homeDir)

	path := C.Path.HomeDir() + "/uploads"
	_ = os.MkdirAll(path, 0777)
	_ = os.WriteFile(C.Path.HomeDir()+"/geoip.metadb", metadata, 0777)
}

// Init
//
//	@Description: 初始化
func Init() {

	defer func() {
		if e := recover(); e != nil {
			log.Errorln("meta.Init field:", e)
		}
	}()

	// 设置输出目录
	logFilePath := filepath.Join(C.Path.HomeDir(), "log.log")
	f, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
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

	plog.AddHook(&LogHook{Path: logFilePath})
	logCheck(logFilePath)

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

	value := cache.Get(constant.DefaultProfile)
	if value == nil || string(value) == "" {
		profile := resolve.Profile{
			Id:       constant.DefaultProfile,
			Type:     1,
			Order:    0,
			Path:     "uploads/" + constant.DefaultProfile + ".yaml",
			Selected: true,
		}
		bytes, _ := json.Marshal(profile)
		_ = cache.Put(profile.Id, bytes)
		_ = saveProfile2Local(constant.DefaultProfile, "yaml", DefProfile)
	}
}

// 将数据保存到本地文件
// 参数：
//   - name: 文件名
//   - suffix: 文件后缀
//   - all: 需要保存的数据
//
// 返回值：错误信息
func saveProfile2Local(name, suffix string, all []byte) error {
	return os.WriteFile(C.Path.HomeDir()+"/uploads/"+name+"."+suffix, all, 0777)
}

var NowConfig *config.Config

// StartCore 函数用于启动核心功能，接收两个参数：profile和reload，分别为配置文件和是否自动reload的标志位
func StartCore(profile resolve.Profile, reload bool) {
	on := cache.Get(constant.DefaultTemplate)
	useTemplate := false
	var providerBuf []byte
	var templateBuf []byte
	needRuleCheck := false
	path := profile.Path

	if string(on) == "on" {
		template, err := os.ReadFile(filepath.Join(C.Path.HomeDir(), constant.DefaultTemplate))
		if err == nil && len(template) > 0 {
			useTemplate = true
			templateBuf = template
		}
	}

	if useTemplate {
		switch profile.Type {
		case 31, 41:
			var err error
			providerBuf, err = os.ReadFile(filepath.Join(C.Path.HomeDir(), path))
			if err != nil {
				log.Warnln("Read config error: %s", err.Error())
				return
			}
			if !strings.Contains(string(providerBuf), "proxy-providers") {
				replace := strings.Replace(string(templateBuf),
					resolve.PandoraDefaultPlace,
					path,
					1)
				providerBuf = []byte(replace)
			}
		default:
			replace := strings.Replace(string(templateBuf),
				resolve.PandoraDefaultPlace,
				path,
				1)
			providerBuf = []byte(replace)
		}
	} else {
		switch profile.Type {
		case 31, 41:
			var err error
			providerBuf, err = os.ReadFile(filepath.Join(C.Path.HomeDir(), path))
			if err != nil {
				log.Warnln("Read config error: %s", err.Error())
				return
			}
			if !strings.Contains(string(providerBuf), "proxy-groups") {
				replace := strings.Replace(string(resolve.PandoraDefaultConfig),
					resolve.PandoraDefaultPlace,
					path,
					1)
				providerBuf = []byte(replace)
			} else {
				needRuleCheck = true
			}
		default:
			replace := strings.Replace(string(resolve.PandoraDefaultConfig),
				resolve.PandoraDefaultPlace,
				path,
				1)
			providerBuf = []byte(replace)
		}
	}

	rawCfg, err := config.UnmarshalRawConfig(providerBuf)
	if err != nil {
		log.Warnln("Unmarshal config error: %s", err.Error())
		return
	}
	if !useTemplate && needRuleCheck && len(rawCfg.Rule) > 10240 {
		replace := strings.Replace(string(resolve.PandoraDefaultConfig),
			resolve.PandoraDefaultPlace,
			path,
			1)
		providerBuf = []byte(replace)
		rawCfg, _ = config.UnmarshalRawConfig(providerBuf)
	}

	rawCfg.Port = 0
	rawCfg.SocksPort = 0
	rawCfg.TProxyPort = 0
	rawCfg.RedirPort = 0
	if reload {
		rawCfg.MixedPort = NowConfig.General.MixedPort
		rawCfg.AllowLan = NowConfig.General.AllowLan
		rawCfg.IPv6 = NowConfig.General.IPv6
		rawCfg.Tun.Enable = NowConfig.General.Tun.Enable
		rawCfg.UnifiedDelay = NowConfig.General.UnifiedDelay
	}
	if rawCfg.AllowLan {
		rawCfg.BindAddress = "*"
	}

	rawCfg.ExternalController = ""
	rawCfg.Mode = tunnel.Rule
	rawCfg.GeodataMode = false

	if isadmin.Check() {
		log.Infoln("Is Admin Check: %s", "true")
		rawCfg.Tun.DNSHijack = []string{"any:53"}
		rawCfg.Tun.AutoRoute = true
		rawCfg.Tun.AutoDetectInterface = true
		rawCfg.Tun.Device = "Pandora"
	} else {
		log.Infoln("Is Admin Check: %s", "false")
		rawCfg.Tun.Enable = false
	}

	cfg, err := config.ParseRawConfig(rawCfg)
	if err != nil {
		log.Warnln("Parse config error: %s", err.Error())
		return
	}

	C.UA = "Clash.Meta/1.18.6"
	if !reload {
		freePort, err := tools.GetFreeWithPort(10000)
		if err != nil {
			freePort, _ = tools.GetFreePort()
		}
		cfg.General.MixedPort = freePort
	}

	NowConfig = cfg
	executor.ApplyConfig(cfg, true)
}
