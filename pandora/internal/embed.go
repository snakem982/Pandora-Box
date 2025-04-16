package internal

import _ "embed"

//go:embed embed/config.yaml
var PandoraDefaultConfig []byte

//go:embed embed/config_download.yaml
var PandoraDefaultDownloadConfig []byte

//go:embed embed/geoip.metadb
var GeoIp []byte

//go:embed embed/GeoSite.dat
var GeoSite []byte

//go:embed embed/GeoLite2-ASN.mmdb
var ASN []byte
