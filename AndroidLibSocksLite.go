package AndroidLibSocksLite

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/armon/go-socks5"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const CoreVersion = "1.0.0"

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Port     int    `json:"port"`
}

type ClientStatus struct {
	Username string `json:"username"`
	Port     int    `json:"port"`
	Running  bool   `json:"running"`
}

type ProxyLog struct {
	Username  string    `json:"username"`
	URL       string    `json:"url"`
	Timestamp time.Time `json:"timestamp"`
}

type Metrics struct {
	TotalUsers   int `json:"total_users"`
	ActiveUsers  int `json:"active_users"`
	TotalServers int `json:"total_servers"`
}

var (
	isCoreRunning   bool
	servers         map[int]*socks5.Server
	listeners       map[int]net.Listener
	userCredentials map[int]User
	httpServer      *http.Server
	mutex           sync.RWMutex
	proxyLogs       []ProxyLog
	logMutex        sync.Mutex
)

// Initialize logging with levels for better readability.
func logInfo(message string) {
	log.Printf("[INFO] %s\n", message)
}

func logError(message string) {
	log.Printf("[ERROR] %s\n", message)
}

// StartSocksServers initializes multiple SOCKS5 servers from JSON input data,
// and binds both SOCKS5 and HTTP servers to the specified `host`.
func StartSocksServers(host string, jsonData string) error {
	mutex.Lock()
	defer mutex.Unlock()

	if isCoreRunning {
		logInfo("Core is already running.")
		return fmt.Errorf("core is already running")
	}

	logInfo(fmt.Sprintf("Starting Core version %s on host %s...", CoreVersion, host))
	var users []User
	err := json.Unmarshal([]byte(jsonData), &users)
	if err != nil {
		logError(fmt.Sprintf("Core initialization failed: error parsing JSON data: %v", err))
		return fmt.Errorf("error parsing JSON data: %v", err)
	}

	servers = make(map[int]*socks5.Server)
	listeners = make(map[int]net.Listener)
	userCredentials = make(map[int]User)
	for _, user := range users {
		credMap := socks5.StaticCredentials{user.Username: user.Password}
		auth := socks5.UserPassAuthenticator{Credentials: credMap}

		conf := &socks5.Config{
			AuthMethods: []socks5.Authenticator{auth},
		}

		server, err := socks5.New(conf)
		if err != nil {
			logError(fmt.Sprintf("Error creating SOCKS5 server for user %s: %v", user.Username, err))
			continue
		}

		addr := fmt.Sprintf("%s:%d", host, user.Port)
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			logError(fmt.Sprintf("Error creating listener on %s: %v", addr, err))
			continue
		}

		servers[user.Port] = server
		listeners[user.Port] = listener
		userCredentials[user.Port] = user

		go func(s *socks5.Server, l net.Listener, username string) {
			for {
				conn, err := l.Accept()
				if err != nil {
					logError(fmt.Sprintf("Error accepting connection on %s:%d: %v", host, user.Port, err))
					return
				}

				wrappedConn := &loggingConn{Conn: conn, username: username}
				go func() {
					defer func(conn net.Conn) {
						if err := conn.Close(); err != nil {
							logError(fmt.Sprintf("Error closing connection: %v", err))
						}
					}(conn)

					if err := s.ServeConn(wrappedConn); err != nil {
						logError(fmt.Sprintf("Error serving SOCKS5 connection: %v", err))
					}
				}()
			}
		}(server, listener, user.Username)

		logInfo(fmt.Sprintf("User %s server started on %s:%d", user.Username, host, user.Port))
	}

	isCoreRunning = true
	logInfo(fmt.Sprintf("Core version %s started successfully on host %s.", CoreVersion, host))
	go startHTTPServer(host) // Start HTTP server on the specified host

	// Graceful shutdown handling
	go func() {
		shutdownChan := make(chan os.Signal, 1)
		signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)
		<-shutdownChan
		logInfo("Shutting down servers...")
		if err := Shutdown(); err != nil {
			logError(fmt.Sprintf("Error during shutdown: %v", err))
		}
	}()

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
	logRequest(c.username, destAddr)
	return c.Conn.Write(b)
}

// logRequest adds a log entry for each request that goes through the SOCKS5 server.
func logRequest(username, url string) {
	logMutex.Lock()
	defer logMutex.Unlock()

	// Keep only the last 100 logs
	if len(proxyLogs) >= 100 {
		proxyLogs = proxyLogs[1:]
	}

	proxyLogs = append(proxyLogs, ProxyLog{
		Username:  username,
		URL:       url,
		Timestamp: time.Now(),
	})
}

// ProxyLogsAPI returns the proxy logs, optionally filtered by username.
func ProxyLogsAPI(w http.ResponseWriter, r *http.Request) {
	logMutex.Lock()
	defer logMutex.Unlock()

	username := r.URL.Query().Get("username")
	var filteredLogs []ProxyLog

	for _, logEntry := range proxyLogs {
		if username == "" || logEntry.Username == username {
			filteredLogs = append(filteredLogs, logEntry)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(filteredLogs); err != nil {
		logError(fmt.Sprintf("Error encoding logs JSON response: %v", err))
		writeJSONError(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// startHTTPServer initializes and starts the HTTP server on the specified `host` for managing server status, metrics, and logs.
func startHTTPServer(host string) {
	http.HandleFunc("/getClientStatus", GetClientStatus)
	http.HandleFunc("/getMetrics", GetMetricsAPI)
	http.HandleFunc("/getProxyLogs", ProxyLogsAPI)
	httpServer = &http.Server{Addr: fmt.Sprintf("%s:8080", host)}
	logInfo(fmt.Sprintf("HTTP server started on %s:8080.", host))
	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		logError(fmt.Sprintf("HTTP server error: %v", err))
	}
}

// Shutdown gracefully shuts down the HTTP server.
func Shutdown() error {
	if httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return httpServer.Shutdown(ctx)
	}
	return nil
}

// GetClientStatus retrieves status information for each user’s SOCKS5 server.
func GetClientStatus(w http.ResponseWriter, r *http.Request) {
	mutex.RLock()
	defer mutex.RUnlock()

	var statuses []ClientStatus
	for port, user := range userCredentials {
		status := ClientStatus{
			Username: user.Username,
			Port:     port,
			Running:  servers[port] != nil && listeners[port] != nil,
		}
		statuses = append(statuses, status)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(statuses); err != nil {
		logError(fmt.Sprintf("Error encoding JSON response: %v", err))
		writeJSONError(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// GetMetricsAPI returns the system metrics as JSON.
func GetMetricsAPI(w http.ResponseWriter, r *http.Request) {
	mutex.RLock()
	defer mutex.RUnlock()

	metrics := Metrics{
		TotalUsers:   len(userCredentials),
		ActiveUsers:  len(servers),
		TotalServers: len(listeners),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(metrics); err != nil {
		logError(fmt.Sprintf("Error encoding metrics JSON response: %v", err))
		writeJSONError(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// writeJSONError sends a JSON-formatted error response.
func writeJSONError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
