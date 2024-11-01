package utils

import (
	"fmt"
	"net"
)

func FindAvailablePort(startPort int) (int, error) {
	for port := startPort; port <= 65535; port++ {
		address := fmt.Sprintf(":%d", port)
		listener, err := net.Listen("tcp", address)
		if err == nil {
			// Port is available
			listener.Close()
			return port, nil
		}
	}
	return 0, fmt.Errorf("no available ports found")
}
