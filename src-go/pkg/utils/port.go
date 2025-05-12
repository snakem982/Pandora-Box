package utils

import (
	"net"
	"strconv"
)

// IsPortAvailable 检测地址端口是否被占用
func IsPortAvailable(host string, port int) error {
	address := host + ":" + strconv.Itoa(port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err // 端口被占用
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {

		}
	}(listener)

	return nil // 端口可用
}

// GetRandomPort 获取一个随机可用端口
func GetRandomPort(host string) (int, error) {
	listener, err := net.Listen("tcp", host+":0") // 监听端口 0
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
