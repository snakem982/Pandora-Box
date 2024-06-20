package premium

import (
	C "github.com/metacubex/mihomo/constant"
	"gopkg.in/yaml.v3"
	"net"
	"os"
	"sync"
	"time"
)

type CfIps struct {
	UpdateTime time.Time `json:"update-time" yaml:"update-time"`
	HttpIps    []string  `json:"http-ips" yaml:"http-ips"`
	HttpsIps   []string  `json:"https-ips" yaml:"https-ips"`
}

var preLock = sync.Mutex{}

func GetExcellentIps(cdnType CdnType) (Ips CfIps) {
	if preLock.TryLock() {
		defer preLock.Unlock()

		cfIps := GetIpsFromFile(cdnType)
		if cfIps != nil {
			LoadIPRanges(cdnType)
			Ips = *cfIps
			return
		}

		// 置随机数种子
		InitRandSeed()
		// 开始延迟测速 + 过滤延迟/丢包
		TcpPort = 443
		httpsIps := NewPing(cdnType).Run().FilterDelay().FilterLossRate()

		Ips.UpdateTime = time.Now()

		Ips.HttpsIps = make([]string, httpsIps.Len())
		for index, value := range httpsIps {
			Ips.HttpsIps[index] = value.IP.String()
		}

		out, err := yaml.Marshal(Ips)
		if err != nil {
			println("marshal cloudflare ip error: " + err.Error())
			return
		}

		if err := os.WriteFile(C.Path.HomeDir()+cdnType.String(), out, 0666); err != nil {
			println("save cloudflare ip error: " + err.Error())
		}
	}

	return
}

func GetIpsFromFile(cdnType CdnType) *CfIps {
	var cf CfIps
	fi, err := os.Open(C.Path.HomeDir() + cdnType.String())
	if err != nil {
		return nil
	}
	defer func(fi *os.File) {
		err := fi.Close()
		if err != nil {

		}
	}(fi)

	decoder := yaml.NewDecoder(fi)
	err = decoder.Decode(&cf)
	if err != nil {
		return nil
	}

	if time.Since(cf.UpdateTime) > 168*time.Hour {
		return nil
	}

	return &cf
}

func IsCdnIp(CIDR []net.IPNet, addr string) bool {
	ip := net.ParseIP(addr)
	for _, ipNet := range CIDR {
		if ipNet.Contains(ip) {
			return true
		}
	}

	return false
}
