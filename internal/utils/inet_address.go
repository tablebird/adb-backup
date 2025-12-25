package utils

import (
	"fmt"
	"net"
	"os"
)

func GetLocalHostIP() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("get hostname error：%w", err)
	}

	addrs, err := net.LookupHost(hostname)
	if err != nil {
		return "", fmt.Errorf("lookup IP error：%w", err)
	}

	for _, addr := range addrs {
		ip := net.ParseIP(addr)
		if ip != nil && ip.To4() != nil {
			return addr, nil
		}
	}

	if len(addrs) > 0 {
		return addrs[0], nil
	}

	return "", fmt.Errorf("can not get host IP")
}
