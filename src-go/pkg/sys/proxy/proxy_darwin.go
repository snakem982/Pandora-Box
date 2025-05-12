package sys

import (
	"bufio"
	"bytes"
	"fmt"
	sys "github.com/snakem982/pandora-box/pkg/sys/cmd"
	"io"
	"net/textproto"
	"regexp"
	"strconv"
	"strings"
)

func OffAll() error {
	if err := OffHttps(); err != nil {
		return err
	}
	if err := OffHttp(); err != nil {
		return err
	}
	if err := OffSocks(); err != nil {
		return err
	}
	return nil
}

func SetIgnore(ignores []string) error {
	s, err := getNetworkInterface()
	if err != nil {
		return err
	}
	return set("proxybypassdomains", s, ignores...)
}

func ClearIgnore() error {
	s, err := getNetworkInterface()
	if err != nil {
		return err
	}
	return set("proxybypassdomains", s, "")
}

func GetIgnore() ([]string, error) {
	s, err := getNetworkInterface()
	if err != nil {
		return nil, err
	}
	m, err := get("proxybypassdomains", s)
	if err != nil {
		return nil, err
	}
	m = strings.TrimSpace(m)
	ignores := strings.Split(m, "\n")
	if len(ignores) != 0 && ignores[len(ignores)-1] == "" {
		ignores = ignores[:len(ignores)-1]
	}
	return ignores, nil
}

func OnHttps(addr Addr) error {
	s, err := getNetworkInterface()
	if err != nil {
		return err
	}
	err = set("securewebproxy", s, addr.Host, strconv.Itoa(addr.Port))
	if err != nil {
		return err
	}
	err = set("securewebproxystate", s, "on")
	if err != nil {
		return err
	}
	return nil
}

func OffHttps() error {
	s, err := getNetworkInterface()
	if err != nil {
		return err
	}
	err = set("securewebproxystate", s, "off")
	if err != nil {
		return err
	}
	return nil
}

func GetHttps() (*Addr, error) {
	s, err := getNetworkInterface()
	if err != nil {
		return nil, err
	}
	buf, err := get("securewebproxy", s)
	if err != nil {
		return nil, err
	}
	reader := textproto.NewReader(bufio.NewReader(bytes.NewBufferString(buf)))
	header, err := reader.ReadMIMEHeader()
	if err != nil && err != io.EOF {
		return nil, err
	}
	if header.Get("Enabled") == "Yes" {
		return ParseAddrPtr(fmt.Sprintf("%s:%s", header.Get("Server"), header.Get("Port"))), nil
	}
	return nil, nil
}

func OnHttp(addr Addr) error {
	s, err := getNetworkInterface()
	if err != nil {
		return err
	}
	err = set("webproxy", s, addr.Host, strconv.Itoa(addr.Port))
	if err != nil {
		return err
	}
	err = set("webproxystate", s, "on")
	if err != nil {
		return err
	}
	return nil
}

func OffHttp() error {
	s, err := getNetworkInterface()
	if err != nil {
		return err
	}
	err = set("webproxystate", s, "off")
	if err != nil {
		return err
	}
	return nil
}

func GetHttp() (*Addr, error) {
	s, err := getNetworkInterface()
	if err != nil {
		return nil, err
	}
	buf, err := get("webproxy", s)
	if err != nil {
		return nil, err
	}
	reader := textproto.NewReader(bufio.NewReader(bytes.NewBufferString(buf)))
	header, err := reader.ReadMIMEHeader()
	if err != nil && err != io.EOF {
		return nil, err
	}
	if header.Get("Enabled") == "Yes" {
		return ParseAddrPtr(fmt.Sprintf("%s:%s", header.Get("Server"), header.Get("Port"))), nil
	}
	return nil, nil
}

func OnSocks(addr Addr) error {
	s, err := getNetworkInterface()
	if err != nil {
		return err
	}
	err = set("socksfirewallproxy", s, addr.Host, strconv.Itoa(addr.Port))
	if err != nil {
		return err
	}
	err = set("socksfirewallproxystate", s, "on")
	if err != nil {
		return err
	}
	return nil
}

func OffSocks() error {
	s, err := getNetworkInterface()
	if err != nil {
		return err
	}
	err = set("socksfirewallproxystate", s, "off")
	if err != nil {
		return err
	}
	return nil
}

func GetSocks() (*Addr, error) {
	s, err := getNetworkInterface()
	if err != nil {
		return nil, err
	}
	buf, err := get("socksfirewallproxy", s)
	if err != nil {
		return nil, err
	}
	reader := textproto.NewReader(bufio.NewReader(bytes.NewBufferString(buf)))
	header, err := reader.ReadMIMEHeader()
	if err != nil && err != io.EOF {
		return nil, err
	}
	if header.Get("Enabled") == "Yes" {
		return ParseAddrPtr(fmt.Sprintf("%s:%s", header.Get("Server"), header.Get("Port"))), nil
	}
	return nil, nil
}

func set(key string, inter string, value ...string) error {
	_, err := sys.Command("networksetup", append([]string{"-set" + key, inter}, value...)...)
	return err
}

func get(key string, inter string) (string, error) {
	return sys.Command("networksetup", "-get"+key, inter)
}

func getNetworkInterface() (string, error) {
	buf, err := sys.Command("sh", "-c", "networksetup -listnetworkserviceorder | grep -B 1 $(route -n get default | grep interface | awk '{print $2}')")
	if err != nil {
		return "", err
	}
	reader := textproto.NewReader(bufio.NewReader(bytes.NewBufferString(buf)))
	reg := regexp.MustCompile(`^\(\d+\)\s(.*)$`)
	device, err := reader.ReadLine()
	if err != nil {
		return "", err
	}
	match := reg.FindStringSubmatch(device)
	if len(match) <= 1 {
		return "", fmt.Errorf("unable to get network interface")
	}
	return match[1], nil
}
