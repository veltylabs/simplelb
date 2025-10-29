package simplelb

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	t.Run("HealthCheck with alive backend", func(t *testing.T) {
		serverPool = ServerPool{}
		// Create a mock backend server
		backendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer backendServer.Close()

		backendUrl, _ := url.Parse(backendServer.URL)
		b := &Backend{URL: backendUrl}
		sp := &ServerPool{}
		sp.AddBackend(b)
		serverPool = *sp

		sp.HealthCheck()

		if !b.IsAlive() {
			t.Error("Expected backend to be alive")
		}
	})

	t.Run("HealthCheck with dead backend", func(t *testing.T) {
		serverPool = ServerPool{}
		// Create a mock backend server that is closed
		backendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		backendServer.Close()

		backendUrl, _ := url.Parse(backendServer.URL)
		b := &Backend{URL: backendUrl}
		sp := &ServerPool{}
		sp.AddBackend(b)
		serverPool = *sp

		sp.HealthCheck()

		if b.IsAlive() {
			t.Error("Expected backend to be down")
		}
	})
}
