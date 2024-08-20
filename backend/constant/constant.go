package constant

import _ "embed"

const (
	DefaultProfile  = "Profile_0"
	DefaultTemplate = "Template_0.yaml"
	PrefixProfile   = "Profile_"
	PrefixGetter    = "Getter_"
	RealIpHeader    = "RealIp_"
)

const (
	CollectLocal     = "local"
	CollectClash     = "clash"
	CollectV2ray     = "v2ray"
	CollectSharelink = "share"
	CollectFuzzy     = "fuzzy"
)

//go:embed version.txt
var PandoraVersion string

const PandoraVersionUrl = "https://raw.githubusercontent.com/snakem982/Pandora-Box/main/backend/constant/version.txt"
const PandoraDownloadUrl = "https://github.com/snakem982/Pandora-Box/releases/download/%s/%s-%s.zip"
