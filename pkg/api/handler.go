package api

import (
	"AndroidLibSocksLite/pkg/logging"
	"AndroidLibSocksLite/pkg/socks"
	"encoding/json"
	"net/http"
)

func StartHTTPServer(host string) {
	http.HandleFunc("/getClientStatus", getClientStatus)
	http.HandleFunc("/shutdown", shutdownServers)

	logging.LogInfo("Starting HTTP server on " + host + ":8080")
	http.ListenAndServe(host+":8080", nil)
}

func getClientStatus(w http.ResponseWriter, r *http.Request) {
	status := map[string]interface{}{
		"isCoreRunning": socks.IsCoreRunning(),
		"activeServers": len(socks.ActiveServers()),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func shutdownServers(w http.ResponseWriter, r *http.Request) {
	err := socks.Shutdown()
	if err != nil {
		http.Error(w, "Failed to shutdown", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("All servers shut down successfully"))
}
