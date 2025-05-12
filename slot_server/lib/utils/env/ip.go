package env

import (
	"errors"
	"net"
	"strings"
)

var (
	IP     string
	IPLast string
)

func init() {
	ip, err := GetLocalIp()
	if err != nil {
		ip = net.IPv4(127, 0, 0, 1)
	}
	IP = ip.String()
	IPLast = strings.Split(IP, ".")[3]
}

// GetLocalIp 获取本地ip地址
func GetLocalIp() (net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP, nil
			}
		}
	}
	return nil, errors.New("can not find the client ip address")
}
