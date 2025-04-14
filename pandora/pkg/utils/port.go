package utils

import (
	"net"
	"strconv"
)

// IsPortAvailable 检测端口是否被占用
func IsPortAvailable(port int) bool {
	address := ":" + strconv.Itoa(port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return false // 端口被占用
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {

		}
	}(listener)
	return true // 端口可用
}

// GetRandomPort 获取一个随机可用端口
func GetRandomPort() (int, error) {
	listener, err := net.Listen("tcp", ":0") // 监听端口 0
	if err != nil {
		return 0, err
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {

		}
	}(listener)
	return listener.Addr().(*net.TCPAddr).Port, nil
}
