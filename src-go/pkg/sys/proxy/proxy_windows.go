package sys

import (
	"bufio"
	"bytes"
	"io"
	"net/textproto"
	"strconv"
	"strings"

	sys "github.com/snakem982/pandora-box/pkg/sys/cmd"
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
	err := set("ProxyServer", "REG_SZ", addr.String())
	if err != nil {
		return err
	}

	return useProxy(true)
}

func OffHttps() error {
	err := useProxy(false)
	if err != nil {
		return err
	}

	return set("ProxyServer", "REG_SZ", "")
}

func OnHttp(addr Addr) error {
	return nil
}

func OffHttp() error {
	return nil
}

func OnSocks(addr Addr) error {
	return nil
}

func OffSocks() error {
	return nil
}

func GetHttp() (*Addr, error) {
	// 检查代理是否启用
	enabled, err := getProxy()
	if err != nil {
		return nil, err
	}
	if !enabled {
		// 如果代理未启用，返回 nil
		return nil, nil
	}

	// 获取代理服务器地址
	m, err := get("ProxyServer")
	if err != nil {
		return nil, err
	}
	addr := m["ProxyServer"]
	if addr == "" {
		return nil, nil
	}

	// 解析 HTTP 代理地址
	if strings.Contains(addr, "=") {
		// 如果 ProxyServer 包含多个协议的代理地址，提取 http= 部分
		parts := strings.Split(addr, ";")
		for _, part := range parts {
			if strings.HasPrefix(part, "http=") {
				addr = strings.TrimPrefix(part, "http=")
				break
			}
		}
	} else {
		// 如果只有一个代理地址，直接使用
		addr = strings.TrimSpace(addr)
	}

	// 返回解析后的地址
	return ParseAddrPtr(addr), nil
}

const settingPath = `HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Internet Settings`

func set(key string, typ string, value string) error {
	_, err := sys.Command(`reg`, `add`, settingPath, `/v`, key, `/t`, typ, `/d`, value, `/f`)
	return err
}

func get(keys ...string) (map[string]string, error) {
	buf, err := sys.Command(`reg`, `query`, settingPath)
	if err != nil {
		return nil, err
	}
	return getFrom(buf, settingPath, keys...)
}

func del(key string) error {
	_, err := sys.Command(`reg`, `delete`, settingPath, `/v`, key, `/f`)
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
	_, _ = reader.ReadLine()
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
