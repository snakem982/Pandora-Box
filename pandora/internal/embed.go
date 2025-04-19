package internal

import _ "embed"

//go:embed em/config.yaml
var PandoraDefaultConfig []byte

//go:embed em/config_download.yaml
var PandoraDefaultDownloadConfig []byte

//go:embed em/geoip.metadb
var GeoIp []byte

//go:embed em/GeoSite.dat
var GeoSite []byte

//go:embed em/GeoLite2-ASN.mmdb
var ASN []byte

//go:embed em/version.txt
var PandoraVersion string

//go:embed em/webtest.json
var DefaultWebTest []byte
