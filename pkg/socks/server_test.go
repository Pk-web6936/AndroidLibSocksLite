package socks

import (
    "testing"
    "net"
    "syscall"
)

func TestStartServer(t *testing.T) {
    // Generate a random port
    listener, err := net.Listen("tcp", ":0")
    if err != nil {
        t.Fatal(err)
    }
    defer listener.Close()

    // Convert net.Listener to net.TCPListener
    tcpListener, ok := listener.(*net.TCPListener)
    if !ok {
        t.Fatal("Failed to convert listener to TCPListener")
    }

    // Get the port number
    port := listener.Addr().(*net.TCPAddr).Port

    // Set SO_REUSEADDR (not necessary in Go, but for demonstration)
    file, err := tcpListener.File()
    if err != nil {
        t.Fatal(err)
    }
    defer file.Close()

    syscall.SetsockoptInt(int(file.Fd()), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)

    user := User{Username: "test", Password: "test", Port: port}
    if err := startServer("localhost", user); err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
}
