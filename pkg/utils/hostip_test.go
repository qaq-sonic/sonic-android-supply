package utils

import "testing"

func TestGetHostIP(t *testing.T) {
	ip := GetHostIP()
	expectedIP := "10.160.74.183"
	if ip != expectedIP {
		t.Errorf("Expected IP %s, but got %s", expectedIP, ip)
	}
}
