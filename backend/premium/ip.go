package premium

import (
	_ "embed"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

func InitRandSeed() {
	rand.Seed(time.Now().UnixNano())
}

func isIPv4(ip string) bool {
	return strings.Contains(ip, ".")
}

func randIPEndWith(num byte) byte {
	if num == 0 { // 对于 /32 这种单独的 IP
		return byte(0)
	}
	return byte(rand.Intn(int(num)))
}

type IPRanges struct {
	ips     []*net.IPAddr
	mask    string
	firstIP net.IP
	ipNet   *net.IPNet
	CIDR    []net.IPNet
}

func newIPRanges() *IPRanges {
	return &IPRanges{
		ips:  make([]*net.IPAddr, 0),
		CIDR: make([]net.IPNet, 0),
	}
}

// 如果是单独 IP 则加上子网掩码，反之则获取子网掩码(r.mask)
func (r *IPRanges) fixIP(ip string) string {
	// 如果不含有 '/' 则代表不是 IP 段，而是一个单独的 IP，因此需要加上 /32 /128 子网掩码
	if i := strings.IndexByte(ip, '/'); i < 0 {
		if isIPv4(ip) {
			r.mask = "/32"
		} else {
			r.mask = "/128"
		}
		ip += r.mask
	} else {
		r.mask = ip[i:]
	}
	return ip
}

// 解析 IP 段，获得 IP、IP 范围、子网掩码
func (r *IPRanges) parseCIDR(ip string) {
	var err error
	if r.firstIP, r.ipNet, err = net.ParseCIDR(r.fixIP(ip)); err != nil {
		log.Fatalln("ParseCIDR err", err)
		return
	}
	r.CIDR = append(r.CIDR, *r.ipNet)
}

func (r *IPRanges) appendIPv4(d byte) {
	r.appendIP(net.IPv4(r.firstIP[12], r.firstIP[13], r.firstIP[14], d))
}

func (r *IPRanges) appendIP(ip net.IP) {
	r.ips = append(r.ips, &net.IPAddr{IP: ip})
}

// 返回第四段 ip 的最小值及可用数目
func (r *IPRanges) getIPRange() (minIP, hosts byte) {
	minIP = r.firstIP[15] & r.ipNet.Mask[3] // IP 第四段最小值

	// 根据子网掩码获取主机数量
	m := net.IPv4Mask(255, 255, 255, 255)
	for i, v := range r.ipNet.Mask {
		m[i] ^= v
	}
	total, _ := strconv.ParseInt(m.String(), 16, 32) // 总可用 IP 数
	if total > 255 {                                 // 矫正 第四段 可用 IP 数
		hosts = 255
		return
	}
	hosts = byte(total)
	return
}

func (r *IPRanges) chooseIPv4() {
	if r.mask == "/32" { // 单个 IP 则无需随机，直接加入自身即可
		r.appendIP(r.firstIP)
	} else {
		minIP, hosts := r.getIPRange()    // 返回第四段 IP 的最小值及可用数目
		for r.ipNet.Contains(r.firstIP) { // 只要该 IP 没有超出 IP 网段范围，就继续循环随机
			r.appendIPv4(minIP + randIPEndWith(hosts))
			r.firstIP[14]++ // 0.0.(X+1).X
			if r.firstIP[14] == 0 {
				r.firstIP[13]++ // 0.(X+1).X.X
				if r.firstIP[13] == 0 {
					r.firstIP[12]++ // (X+1).X.X.X
				}
			}
		}
	}
}

func (r *IPRanges) chooseIPv6() {
	if r.mask == "/128" { // 单个 IP 则无需随机，直接加入自身即可
		r.appendIP(r.firstIP)
	} else {
		var tempIP uint8                  // 临时变量，用于记录前一位的值
		for r.ipNet.Contains(r.firstIP) { // 只要该 IP 没有超出 IP 网段范围，就继续循环随机
			r.firstIP[15] = randIPEndWith(255) // 随机 IP 的最后一段
			r.firstIP[14] = randIPEndWith(255) // 随机 IP 的最后一段

			targetIP := make([]byte, len(r.firstIP))
			copy(targetIP, r.firstIP)
			r.appendIP(targetIP) // 加入 IP 地址池

			for i := 13; i >= 0; i-- { // 从倒数第三位开始往前随机
				tempIP = r.firstIP[i]              // 保存前一位的值
				r.firstIP[i] += randIPEndWith(255) // 随机 0~255，加到当前位上
				if r.firstIP[i] >= tempIP {        // 如果当前位的值大于等于前一位的值，说明随机成功了，可以退出该循环
					break
				}
			}
		}
	}
}

//go:embed cloudflare.txt
var FsCloudflareIps []byte

//go:embed fastly.txt
var FsFastlyIps []byte

//go:embed gcore.txt
var FsGcoreIps []byte

func LoadIPRanges(cdnType CdnType) []*net.IPAddr {
	ranges := newIPRanges()
	var ips []byte
	switch cdnType {
	case CdnTypeFastly:
		ips = FsFastlyIps
	case CdnTypeCloudflare:
		ips = FsCloudflareIps
	case CdnTypeGcore:
		ips = FsGcoreIps
	default:
	}

	direct := strings.Split(string(ips), "\n")
	for _, s := range direct { // 循环遍历文件每一行
		line := strings.TrimSpace(s) // 去除首尾的空白字符（空格、制表符、换行符等）
		if line == "" {              // 跳过空行
			continue
		}
		ranges.parseCIDR(line) // 解析 IP 段，获得 IP、IP 范围、子网掩码
		if isIPv4(line) {      // 生成要测速的所有 IPv4 / IPv6 地址（单个/随机/全部）
			ranges.chooseIPv4()
		} else {
			ranges.chooseIPv6()
		}
	}

	return ranges.ips
}

func LoadCIDR(cdnType CdnType) []net.IPNet {
	ranges := newIPRanges()
	var ips []byte
	switch cdnType {
	case CdnTypeFastly:
		ips = FsFastlyIps
	case CdnTypeCloudflare:
		ips = FsCloudflareIps
	case CdnTypeGcore:
		ips = FsGcoreIps
	default:
	}

	direct := strings.Split(string(ips), "\n")
	for _, s := range direct { // 循环遍历文件每一行
		line := strings.TrimSpace(s) // 去除首尾的空白字符（空格、制表符、换行符等）
		if line == "" {              // 跳过空行
			continue
		}
		ranges.parseCIDR(line) // 解析 IP 段，获得 IP、IP 范围、子网掩码
	}

	return ranges.CIDR
}
