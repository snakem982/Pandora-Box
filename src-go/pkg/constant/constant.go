package constant

import _ "embed"

const (
	DefaultWorkDir     = "Pandora-Box-V3"
	DefaultCrawlDir    = "crawl"
	DefaultTemplateDir = "template"
	DefaultServerDB    = "px-server.db"
	DefaultClientDB    = "px-client.db"
	DefaultDownload    = "Download_0.yaml"
	PrefixProfile      = "Profile_"
	ProfileOrder       = "ProfileOrder"
	PrefixWebTest      = "WebTest_"
	WebTestOrder       = "WebTestOrder"
	PrefixGetter       = "Getter_"
	PrefixTemplate     = "Template_"
	TemplateSwitch     = "TemplateSwitch"
	RealIpHeader       = "RealIp_"
	SecretKey          = "SecretKey_pb"
	RecoverTmp         = "RecoverTmp"
	QuitSignal         = "QuitSignal"
	Dns                = "DNS"
	Mihomo             = "Mihomo"
)

const (
	CollectLocal     = "local"
	CollectBatch     = "batch"
	CollectClash     = "clash"
	CollectV2ray     = "v2ray"
	CollectSharelink = "share"
	CollectFuzzy     = "fuzzy"
	CollectAuto      = "auto"
	CollectSingBox   = "sing"
)

const PandoraVersionUrl = "https://raw.githubusercontent.com/snakem982/Pandora-Box/main/backend/constant/version.txt"
const PandoraDownloadUrl = "https://github.com/snakem982/Pandora-Box/releases/download/%s/%s-%s.zip"
