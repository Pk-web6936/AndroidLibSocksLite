package socks

import (
	"AndroidLibSocksLite/pkg/logging"
	"encoding/json"
	"fmt"
	"github.com/armon/go-socks5"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const coreVersion = "1.0.1"

func CheckCoreVersion() string {
	return coreVersion
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Port     int    `json:"port"`
}

var (
	isCoreRunning   bool
	servers         map[int]*socks5.Server
	listeners       map[int]net.Listener
	userCredentials map[int]User
	mutex           sync.RWMutex
)

// StartSocksServers initializes multiple SOCKS5 servers
func StartSocksServers(host string, jsonData string) error {
	mutex.Lock()
	defer mutex.Unlock()

	if isCoreRunning {
		logging.LogInfo("Core is already running.")
		return fmt.Errorf("core is already running")
	}

	var users []User
	err := json.Unmarshal([]byte(jsonData), &users)
	if err != nil {
		logging.LogError(fmt.Sprintf("Error parsing JSON data: %v", err))
		return fmt.Errorf("error parsing JSON data: %v", err)
	}

	servers = make(map[int]*socks5.Server)
	listeners = make(map[int]net.Listener)
	userCredentials = make(map[int]User)

	for _, user := range users {
		if err := startServer(host, user); err != nil {
			logging.LogError(fmt.Sprintf("Failed to start server for user %s: %v", user.Username, err))
			// Consider rolling back changes if some servers fail to start
		}
	}

	isCoreRunning = true
	logging.LogInfo("Core started successfully.")

	// Handle shutdown
	go handleShutdown()

	return nil
}

func startServer(host string, user User) error {
	credMap := socks5.StaticCredentials{user.Username: user.Password}
	auth := socks5.UserPassAuthenticator{Credentials: credMap}
	conf := &socks5.Config{AuthMethods: []socks5.Authenticator{auth}}

	server, err := socks5.New(conf)
	if err != nil {
		return fmt.Errorf("error creating SOCKS5 server: %v", err)
	}

	addr := fmt.Sprintf("%s:%d", host, user.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("error creating listener: %v", err)
	}

	servers[user.Port] = server
	listeners[user.Port] = listener
	userCredentials[user.Port] = user

	go acceptConnections(server, listener, user.Username)

	logging.LogInfo(fmt.Sprintf("User %s server started on %s", user.Username, addr))
	return nil
}

func acceptConnections(server *socks5.Server, listener net.Listener, username string) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			logging.LogError(fmt.Sprintf("Error accepting connection: %v", err))
			return
		}

		wrappedConn := &loggingConn{Conn: conn, username: username}
		go server.ServeConn(wrappedConn)
	}
}

func handleShutdown() {
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)
	<-shutdownChan
	logging.LogInfo("Shutting down servers...")
	if err := Shutdown(); err != nil {
		logging.LogError(fmt.Sprintf("Error during shutdown: %v", err))
	}
}

// Shutdown gracefully shuts down all SOCKS5 servers.
func Shutdown() error {
	mutex.Lock()
	defer mutex.Unlock()
	for port, listener := range listeners {
		if err := listener.Close(); err != nil {
			logging.LogError(fmt.Sprintf("Error closing listener on port %d: %v", port, err))
		}
		delete(servers, port)
		delete(listeners, port)
	}
	isCoreRunning = false
	logging.LogInfo("All servers shut down.")
	return nil
}

// loggingConn is a wrapper around net.Conn that logs each request's destination address.
type loggingConn struct {
	net.Conn
	username string
}

// Write logs the destination address before forwarding the request.
func (c *loggingConn) Write(b []byte) (int, error) {
	destAddr := c.RemoteAddr().String()
	logging.LogInfo(fmt.Sprintf("User %s connected to %s", c.username, destAddr))
	return c.Conn.Write(b)
}

// IsCoreRunning returns whether the core is running.
func IsCoreRunning() bool {
	return isCoreRunning
}

// ActiveServers returns the list of active servers.
func ActiveServers() map[int]*socks5.Server {
	return servers
}
