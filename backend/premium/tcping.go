package premium

import (
	"fmt"
	"net"
	"pandora-box/backend/mypool"
	"sort"
	"sync"
	"time"
)

const (
	tcpConnectTimeout = time.Millisecond * 400
	defaultRoutines   = 256
	pingTimes         = 4
)

type CdnType string

const (
	CdnTypeCloudflare CdnType = "Cloudflare"
	CdnTypeFastly     CdnType = "Fastly"
	CdnTypeGcore      CdnType = "Gcore"
)

func (cdnType CdnType) String() string {
	switch cdnType {
	case CdnTypeCloudflare:
		return "/Cloudflare.yaml"
	case CdnTypeFastly:
		return "/Fastly.yaml"
	case CdnTypeGcore:
		return "/Gcore.yaml"
	default:
		return "/Cloudflare.yaml"
	}
}

var TcpPort = 443

type Ping struct {
	m   *sync.Mutex
	ips []*net.IPAddr
	csv PingDelaySet
}

func NewPing(cdnType CdnType) *Ping {
	ips := LoadIPRanges(cdnType)
	return &Ping{
		m:   &sync.Mutex{},
		ips: ips,
		csv: make(PingDelaySet, 0),
	}
}

func (p *Ping) Run() PingDelaySet {
	if len(p.ips) == 0 {
		return p.csv
	}
	fmt.Printf("开始延迟测速（模式：TCP，端口：%d，平均延迟上限：%v ms，丢包几率上限：0 )\n", TcpPort, InputMaxDelay.Milliseconds())
	pool := mypool.NewTimeoutPool(defaultRoutines)
	pool.WaitCount(len(p.ips))
	for _, ipp := range p.ips {
		ip := ipp
		pool.SubmitWithTimeout(func(done chan struct{}) {
			defer func() {
				if e := recover(); e != nil {
				}
				done <- struct{}{}
			}()
			p.tcpPingHandler(ip)
		}, 500*time.Millisecond)
	}
	pool.StartAndWait()
	fmt.Printf("结束延迟测速（模式：TCP，端口：%d，平均延迟上限：%v ms，丢包几率上限：0 )\n", TcpPort, InputMaxDelay.Milliseconds())
	sort.Sort(p.csv)
	return p.csv
}

// tcpPing bool connectionSucceed float32 time
func tcpPing(ip *net.IPAddr) (bool, time.Duration) {
	startTime := time.Now()
	var fullAddress string
	if isIPv4(ip.String()) {
		fullAddress = fmt.Sprintf("%s:%d", ip.String(), TcpPort)
	} else {
		fullAddress = fmt.Sprintf("[%s]:%d", ip.String(), TcpPort)
	}
	conn, err := net.DialTimeout("tcp", fullAddress, tcpConnectTimeout)
	if err != nil {
		return false, 0
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)
	duration := time.Since(startTime)
	return true, duration
}

// CheckConnection pingReceived pingTotalTime
func CheckConnection(ip *net.IPAddr) (recv int, totalDelay time.Duration) {
	for i := 0; i < pingTimes; i++ {
		if ok, delay := tcpPing(ip); ok {
			recv++
			totalDelay += delay
		}
	}
	return
}

func (p *Ping) appendIPData(data *PingData) {
	p.m.Lock()
	defer p.m.Unlock()
	p.csv = append(p.csv, CloudflareIPData{
		PingData: data,
	})
}

// tcpPingHandler handle tcp ping
func (p *Ping) tcpPingHandler(ip *net.IPAddr) {
	rev, totalDlay := CheckConnection(ip)
	nowAble := len(p.csv)
	if rev != 0 {
		nowAble++
	}
	if rev == 0 {
		return
	}
	data := &PingData{
		IP:       ip,
		Sent:     pingTimes,
		Received: rev,
		Delay:    totalDlay / time.Duration(rev),
	}
	p.appendIPData(data)
}
