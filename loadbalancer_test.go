package simplelb

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestLb(t *testing.T) {
	t.Run("Lb with available peer", func(t *testing.T) {
		serverPool = ServerPool{}
		backendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer backendServer.Close()

		backends := []string{backendServer.URL}
		SetupServerPool(backends)

		req := httptest.NewRequest("GET", "http://localhost:3030", nil)
		rr := httptest.NewRecorder()
		Lb(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})

	t.Run("Lb with no available peer", func(t *testing.T) {
		serverPool = ServerPool{}
		req := httptest.NewRequest("GET", "http://localhost:3030", nil)
		rr := httptest.NewRecorder()
		Lb(rr, req)
		if rr.Code != http.StatusServiceUnavailable {
			t.Errorf("Expected status code %d, got %d", http.StatusServiceUnavailable, rr.Code)
		}
	})

	t.Run("GetAttemptsFromContext", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://localhost:3030", nil)
		ctx := context.WithValue(req.Context(), Attempts, 3)
		req = req.WithContext(ctx)
		attempts := GetAttemptsFromContext(req)
		if attempts != 3 {
			t.Errorf("Expected 3 attempts, got %d", attempts)
		}
	})

	t.Run("GetRetryFromContext", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://localhost:3030", nil)
		ctx := context.WithValue(req.Context(), Retry, 2)
		req = req.WithContext(ctx)
		retries := GetRetryFromContext(req)
		if retries != 2 {
			t.Errorf("Expected 2 retries, got %d", retries)
		}
	})

	t.Run("Lb with proxy error", func(t *testing.T) {
		serverPool = ServerPool{}
		backendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// This will be closed immediately to simulate a connection error
		}))
		backendUrl, _ := url.Parse(backendServer.URL)
		backendServer.Close() // Close the server to trigger a proxy error

		backends := []string{backendUrl.String()}
		SetupServerPool(backends)

		req := httptest.NewRequest("GET", "http://localhost:3030", nil)
		rr := httptest.NewRecorder()

		Lb(rr, req)

		if rr.Code != http.StatusServiceUnavailable {
			t.Errorf("Expected status code %d, got %d", http.StatusServiceUnavailable, rr.Code)
		}

		// Check if the backend is marked as down
		var backendAlive bool
		for _, b := range serverPool.backends {
			if b.URL.String() == backendUrl.String() {
				backendAlive = b.IsAlive()
				break
			}
		}
		if backendAlive {
			t.Error("Expected backend to be marked as down")
		}
	})
}
