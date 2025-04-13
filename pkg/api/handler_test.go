package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetClientStatus(t *testing.T) {
	req, err := http.NewRequest("GET", "/getClientStatus", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	http.HandlerFunc(getClientStatus).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestShutdownServers(t *testing.T) {
	req, err := http.NewRequest("GET", "/shutdown", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	http.HandlerFunc(shutdownServers).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}
