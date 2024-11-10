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

const CoreVersion = "1.0.0"

package socks

func CheckCoreVersion() string {
return CoreVersion 
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
		credMap := socks5.StaticCredentials{user.Username: user.Password}
		auth := socks5.UserPassAuthenticator{Credentials: credMap}
		conf := &socks5.Config{AuthMethods: []socks5.Authenticator{auth}}

		server, err := socks5.New(conf)
		if err != nil {
			logging.LogError(fmt.Sprintf("Error creating SOCKS5 server for user %s: %v", user.Username, err))
			continue
		}

		addr := fmt.Sprintf("%s:%d", host, user.Port)
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			logging.LogError(fmt.Sprintf("Error creating listener on %s: %v", addr, err))
			continue
		}

		servers[user.Port] = server
		listeners[user.Port] = listener
		userCredentials[user.Port] = user

		go func(s *socks5.Server, l net.Listener, username string) {
			for {
				conn, err := l.Accept()
				if err != nil {
					logging.LogError(fmt.Sprintf("Error accepting connection: %v", err))
					return
				}

				wrappedConn := &loggingConn{Conn: conn, username: username}
				go s.ServeConn(wrappedConn)
			}
		}(server, listener, user.Username)

		logging.LogInfo(fmt.Sprintf("User %s server started on %s", user.Username, addr))
	}

	isCoreRunning = true
	logging.LogInfo("Core started successfully.")

	// Handle shutdown
	go func() {
		shutdownChan := make(chan os.Signal, 1)
		signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)
		<-shutdownChan
		logging.LogInfo("Shutting down servers...")
		if err := Shutdown(); err != nil {
			logging.LogError(fmt.Sprintf("Error during shutdown: %v", err))
		}
	}()

	return nil
}

// Shutdown gracefully shuts down all SOCKS5 servers.
func Shutdown() error {
	mutex.Lock()
	defer mutex.Unlock()
	for port, listener := range listeners {
		listener.Close()
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
