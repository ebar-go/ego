package helper

import (
	"net"
	"errors"
)

// GetLocalIp 获取本地IP
func GetLocalIp() (string, error) {
	addressItems, err := net.InterfaceAddrs()

	if err != nil {
		return "", err
	}

	for _, address := range addressItems {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}

		}
	}

	return "", errors.New("failed to get local address")
}
