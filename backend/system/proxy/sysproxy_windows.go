//go:build windows

package proxy

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/metacubex/mihomo/log"
	"io"
	"net/textproto"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"syscall"
)

func init() {
	log.Infoln("system proxy use windows")
}

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
	return set("ProxyOverride", "REG_SZ", strings.Join(ignores, ";"))
}

func ClearIgnore() error {
	return set("ProxyOverride", "REG_SZ", "")
}

func GetIgnore() ([]string, error) {
	m, err := get("ProxyOverride")
	if err != nil {
		return nil, err
	}
	ignores := m["ProxyOverride"]
	if ignores == "" {
		return []string{}, nil
	}
	return strings.Split(ignores, ";"), nil
}

func OnHttps(addr Addr) error {
	b, err := getProxy()
	if err != nil {
		return err
	}
	m := map[string]string{}
	if b {
		m, err = getAllProxy()
		if err != nil {
			return err
		}
	}
	m["https"] = addr.String()
	return setAllProxy(m)
}

func OffHttps() error {
	b, err := getProxy()
	if err != nil {
		return err
	}
	m := map[string]string{}
	if b {
		m, err = getAllProxy()
		if err != nil {
			return err
		}
	}
	delete(m, "https")
	return setAllProxy(m)
}

func GetHttps() (*Addr, error) {
	b, err := getProxy()
	if err != nil {
		return nil, err
	}
	if !b {
		return nil, nil
	}
	m, err := getAllProxy()
	if err != nil {
		return nil, err
	}
	return ParseAddrPtr(m["https"]), nil
}

func OnHttp(addr Addr) error {
	b, err := getProxy()
	if err != nil {
		return err
	}
	m := map[string]string{}
	if b {
		m, err = getAllProxy()
		if err != nil {
			return err
		}
	}
	m["http"] = addr.String()
	return setAllProxy(m)
}

func OffHttp() error {
	b, err := getProxy()
	if err != nil {
		return err
	}
	m := map[string]string{}
	if b {
		m, err = getAllProxy()
		if err != nil {
			return err
		}
	}
	delete(m, "http")
	return setAllProxy(m)
}

func GetHttp() (*Addr, error) {
	b, err := getProxy()
	if err != nil {
		return nil, err
	}
	if !b {
		return nil, nil
	}
	m, err := getAllProxy()
	if err != nil {
		return nil, err
	}
	return ParseAddrPtr(m["http"]), nil
}

func OnSocks(addr Addr) error {
	b, err := getProxy()
	if err != nil {
		return err
	}
	m := map[string]string{}
	if b {
		m, err = getAllProxy()
		if err != nil {
			return err
		}
	}
	m["socks"] = addr.String()
	return setAllProxy(m)
}

func OffSocks() error {
	b, err := getProxy()
	if err != nil {
		return err
	}
	m := map[string]string{}
	if b {
		m, err = getAllProxy()
		if err != nil {
			return err
		}
	}
	delete(m, "socks")
	return setAllProxy(m)
}

func GetSocks() (*Addr, error) {
	b, err := getProxy()
	if err != nil {
		return nil, err
	}
	if !b {
		return nil, nil
	}
	m, err := getAllProxy()
	if err != nil {
		return nil, err
	}
	return ParseAddrPtr(m["socks"]), nil
}

func setAllProxy(m map[string]string) error {
	list := make([]string, 0, len(m))
	for key, item := range m {
		if item == "" {
			continue
		}
		list = append(list, strings.Join([]string{key, item}, "="))
	}
	sort.Strings(list)
	err := set("ProxyServer", "REG_SZ", strings.Join(list, ";"))
	if err != nil {
		return err
	}
	return useProxy(len(list) != 0)
}

func getAllProxy() (map[string]string, error) {
	m, err := get("ProxyServer")
	if err != nil {
		return nil, err
	}
	list := strings.Split(m["ProxyServer"], ";")
	proxy := map[string]string{}
	for _, item := range list {
		n := strings.SplitN(item, "=", 2)
		if len(n) == 1 {
			proxy["http"] = item
			proxy["https"] = item
			proxy["socks"] = item
			break
		}
		proxy[n[0]] = n[1]
	}
	return proxy, nil
}

const settingPath = `HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Internet Settings`

func set(key string, typ string, value string) error {
	_, err := Command(`reg`, `add`, settingPath, `/v`, key, `/t`, typ, `/d`, value, `/f`)
	return err
}

func get(keys ...string) (map[string]string, error) {
	buf, err := Command(`reg`, `query`, settingPath)
	if err != nil {
		return nil, err
	}
	return getFrom(buf, settingPath, keys...)
}

func del(key string) error {
	_, err := Command(`reg`, `delete`, settingPath, `/v`, key, `/f`)
	return err
}

func strBool(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

func useProxy(b bool) error {
	return set("ProxyEnable", "REG_DWORD", strBool(b))
}

func getProxy() (bool, error) {
	m, err := get("ProxyEnable", "REG_DWORD")
	if err != nil {
		return false, err
	}
	i, err := strconv.ParseInt(m["ProxyEnable"], 0, 0)
	if err != nil {
		return false, err
	}
	return i != 0, nil
}

func getFrom(data string, path string, keys ...string) (map[string]string, error) {
	m := map[string]string{}
	index := strings.Index(data, path)
	if index == -1 {
		return m, nil
	}
	data = data[index+len(path):]
	reader := textproto.NewReader(bufio.NewReader(bytes.NewBufferString(data)))
	reader.ReadLine()
	for len(m) != len(keys) {
		row, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		if row == "" {
			break
		}
		row = strings.TrimSpace(row)
		s := strings.SplitN(row, "    ", 3)
		key := s[0]
		skip := true
		for _, k := range keys {
			if k == key {
				skip = false
				break
			}
		}
		if skip {
			continue
		}
		val := ""
		if len(s) == 3 {
			val = s[2]
		}
		m[key] = val
	}
	return m, nil
}

func Command(name string, arg ...string) (string, error) {
	c := exec.Command(name, arg...)
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := c.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%q: %w: %q", strings.Join(append([]string{name}, arg...), " "), err, out)
	}
	return strings.TrimSpace(string(out)), nil
}
