package socks

import (
	"testing"
        "net"
        "syscall"
)

func TestStartSocksServers(t *testing.T) {
    // Test starting servers with valid JSON data
    jsonData := `[{"username": "test", "password": "test", "port": 8080}]`
    if err := StartSocksServers("localhost", jsonData); err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
}

func TestStartSocksServersInvalidJSON(t *testing.T) {
    // Test starting servers with invalid JSON data
    jsonData := `Invalid JSON`
    if err := StartSocksServers("localhost", jsonData); err == nil {
        t.Errorf("Expected error, got nil")
    }
}
