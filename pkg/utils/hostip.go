package utils

import (
	"fmt"
	"net"
)

func GetHostIP() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	for _, iface := range interfaces {
		// 排除名为 "docker0" 的接口
		if iface.Name == "docker0" {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		for _, addr := range addrs {
			// 检查地址类型并确保不是回环地址
			if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
				if ipNet.IP.To4() != nil {
					// fmt.Printf("Interface: %s, IP Address: %s\n", iface.Name, ipNet.IP.String())
					return ipNet.IP.String()
				}
			}
		}
	}
	return ""
}
