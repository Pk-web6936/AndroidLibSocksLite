package socks

import (
	"testing"
	"net"
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
	// Generate a random port
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	// Get the port number
	port := listener.Addr().(*net.TCPAddr).Port

	user := User{Username: "test", Password: "test", Port: port}
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
