package sys

import (
	"bytes"
	"fmt"
	sys "github.com/snakem982/pandora-box/pkg/sys/cmd"
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
	buf := bytes.NewBuffer(nil)
	buf.WriteString("[ ")
	for i, item := range ignores {
		if item == "" {
			continue
		}
		buf.WriteByte('\'')
		buf.WriteString(item)
		buf.WriteByte('\'')
		if len(ignores)-1 != i {
			buf.WriteString(", ")
		}
	}
	buf.WriteString(" ]")
	return set("", "ignore-hosts", buf.String())
}

func ClearIgnore() error {
	return set("", "ignore-hosts", "[]")
}

func GetIgnore() ([]string, error) {
	data, err := get("", "ignore-hosts")
	if err != nil {
		return nil, err
	}
	data = strings.TrimPrefix(data, "@as")
	data = strings.TrimSpace(data)
	data = strings.TrimPrefix(data, "[")
	data = strings.TrimSuffix(data, "]")
	data = strings.TrimSpace(data)
	if data == "" {
		return []string{}, nil
	}
	ignores := strings.Split(data, ",")
	for i := range ignores {
		item := ignores[i]
		item = strings.TrimSpace(item)
		item = strings.Trim(item, "'")
		ignores[i] = item
	}
	return ignores, nil
}

func OnHttps(addr Addr) error {
	err := set("https", "host", addr.Host)
	if err != nil {
		return err
	}
	err = set("https", "port", strconv.Itoa(addr.Port))
	if err != nil {
		return err
	}
	err = set("", "mode", "manual")
	if err != nil {
		return err
	}
	return nil
}

func OffHttps() error {
	err := reset("", "mode")
	if err != nil {
		return err
	}
	err = reset("https", "host")
	if err != nil {
		return err
	}
	err = reset("https", "port")
	if err != nil {
		return err
	}
	return nil
}

func GetHttps() (*Addr, error) {
	mode, err := get("", "mode")
	if err != nil {
		return nil, err
	}
	mode = strings.Trim(mode, "'")
	if mode != "manual" {
		return nil, nil
	}
	host, err := get("https", "host")
	if err != nil {
		return nil, err
	}
	port, err := get("https", "port")
	if err != nil {
		return nil, err
	}
	return ParseAddrPtr(fmt.Sprintf("%s:%s", host, port)), nil
}

func OnHttp(addr Addr) error {
	err := set("http", "host", addr.Host)
	if err != nil {
		return err
	}
	err = set("http", "port", strconv.Itoa(addr.Port))
	if err != nil {
		return err
	}
	err = set("", "mode", "manual")
	if err != nil {
		return err
	}
	return nil
}

func OffHttp() error {
	err := reset("", "mode")
	if err != nil {
		return err
	}
	err = reset("http", "host")
	if err != nil {
		return err
	}
	err = reset("http", "port")
	if err != nil {
		return err
	}
	return nil
}

func GetHttp() (*Addr, error) {
	mode, err := get("", "mode")
	if err != nil {
		return nil, err
	}
	if mode != "manual" {
		return nil, nil
	}
	host, err := get("http", "host")
	if err != nil {
		return nil, err
	}
	port, err := get("http", "port")
	if err != nil {
		return nil, err
	}
	return ParseAddrPtr(fmt.Sprintf("%s:%s", host, port)), nil
}

func OnSocks(addr Addr) error {
	err := set("socks", "host", addr.Host)
	if err != nil {
		return err
	}
	err = set("socks", "port", strconv.Itoa(addr.Port))
	if err != nil {
		return err
	}
	err = set("", "mode", "manual")
	if err != nil {
		return err
	}
	return nil
}

func OffSocks() error {
	err := reset("", "mode")
	if err != nil {
		return err
	}
	err = reset("socks", "host")
	if err != nil {
		return err
	}
	err = reset("socks", "port")
	if err != nil {
		return err
	}
	return nil
}

func GetSocks() (*Addr, error) {
	mode, err := get("", "mode")
	if err != nil {
		return nil, err
	}
	if mode != "manual" {
		return nil, nil
	}
	host, err := get("socks", "host")
	if err != nil {
		return nil, err
	}
	port, err := get("socks", "port")
	if err != nil {
		return nil, err
	}
	return ParseAddrPtr(fmt.Sprintf("%s:%s", host, port)), nil
}

const scheme = "org.gnome.system.proxy"

func reset(sub, key string) error {
	scheme := scheme
	if sub != "" {
		scheme = scheme + "." + sub
	}
	_, err := sys.Command("gsettings", "reset", scheme, key)
	return err
}

func get(sub, key string) (string, error) {
	scheme := scheme
	if sub != "" {
		scheme = scheme + "." + sub
	}
	out, err := sys.Command("gsettings", "get", scheme, key)
	if err != nil {
		return "", err
	}
	out = strings.Trim(out, "'")
	return out, nil
}

func set(sub, key string, val string) error {
	scheme := scheme
	if sub != "" {
		scheme = scheme + "." + sub
	}
	_, err := sys.Command("gsettings", "set", scheme, key, val)
	return err
}
