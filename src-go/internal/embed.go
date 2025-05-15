package internal

import _ "embed"

//go:embed em/Template_0.yaml
var Template_0 []byte

//go:embed em/Template_1.yaml
var Template_1 []byte

//go:embed em/Template_2.yaml
var Template_2 []byte

//go:embed em/config_download.yaml
var PandoraDefaultDownloadConfig []byte

//go:embed em/geoip.metadb
var GeoIp []byte

//go:embed em/GeoSite.dat
var GeoSite []byte

//go:embed em/GeoLite2-ASN.mmdb
var ASN []byte

//go:embed em/webtest.json
var DefaultWebTest []byte

//go:embed em/dns.yaml
var DefaultDNS string
