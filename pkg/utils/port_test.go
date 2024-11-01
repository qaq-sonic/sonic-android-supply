package utils

import "testing"

func TestFindAvailablePort(t *testing.T) {
	port, err := FindAvailablePort(1024)
	if err != nil {
		t.Errorf("Error finding available port: %s", err)
	}
	if port < 1024 || port > 65535 {
		t.Errorf("Invalid port number: %d", port)
	}
}
