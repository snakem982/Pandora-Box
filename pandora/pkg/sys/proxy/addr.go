package sys

import (
	"fmt"
	"strconv"
	"strings"
)

type Addr struct {
	Host string
	Port int
}

func (a Addr) String() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}

func ParseAddr(s string) Addr {
	tmp := strings.Split(s, ":")
	var (
		host = tmp[0]
		port int
	)
	if len(tmp) > 1 {
		port, _ = strconv.Atoi(tmp[1])
	}
	return Addr{
		Host: host,
		Port: port,
	}
}

func ParseAddrPtr(s string) *Addr {
	addr := ParseAddr(s)
	return &addr
}
