package socks

import (
	"testing"
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

func TestStartServer(t *testing.T) {
	// Test starting a server with valid data
	user := User{Username: "test", Password: "test", Port: 8080}
	if err := startServer("localhost", user); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestShutdown(t *testing.T) {
	// Test shutting down servers
	if err := Shutdown(); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestIsCoreRunning(t *testing.T) {
	// Test checking if core is running
	if IsCoreRunning() {
		t.Errorf("Expected core not to be running")
	}
}

func TestActiveServers(t *testing.T) {
	// Test getting active servers
	servers := ActiveServers()
	if len(servers) != 0 {
		t.Errorf("Expected no active servers")
	}
}
