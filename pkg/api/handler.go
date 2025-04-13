package api

import (
	"AndroidLibSocksLite/pkg/logging"
	"AndroidLibSocksLite/pkg/socks"
	"encoding/json"
	"net/http"
)

// StartHTTPServer initializes and starts the HTTP server on the specified host.
func StartHTTPServer(host string) {
	http.HandleFunc("/getClientStatus", getClientStatus)
	http.HandleFunc("/shutdown", shutdownServers)

	logging.LogInfo("Starting HTTP server on " + host + ":8080")
	if err := http.ListenAndServe(host+":8080", nil); err != nil {
		logging.LogError("Failed to start HTTP server: " + err.Error())
	}
}

// getClientStatus handles requests to retrieve the current client status.
func getClientStatus(w http.ResponseWriter, r *http.Request) {
	logging.LogInfo("getClientStatus called")
	status := map[string]interface{}{
		"isCoreRunning": socks.IsCoreRunning(),
		"activeServers": len(socks.ActiveServers()),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(status); err != nil {
		logging.LogError("Failed to encode response: " + err.Error())
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// shutdownServers handles requests to shut down all servers.
func shutdownServers(w http.ResponseWriter, r *http.Request) {
	logging.LogInfo("shutdownServers called")
	if err := socks.Shutdown(); err != nil {
		logging.LogError("Failed to shutdown servers: " + err.Error())
		http.Error(w, "Failed to shutdown", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "All servers shut down successfully"}); err != nil {
		logging.LogError("Failed to encode shutdown response: " + err.Error())
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
