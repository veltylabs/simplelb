package simplelb

import (
	"net/url"
	"testing"
)

func TestServerPool(t *testing.T) {
	t.Run("AddBackend and GetNextPeer", func(t *testing.T) {
		serverPool = ServerPool{}
		sp := &ServerPool{}
		u1, _ := url.Parse("http://localhost:8080")
		b1 := &Backend{URL: u1}
		b1.SetAlive(true)
		sp.AddBackend(b1)

		u2, _ := url.Parse("http://localhost:8081")
		b2 := &Backend{URL: u2}
		b2.SetAlive(true)
		sp.AddBackend(b2)

		peer1 := sp.GetNextPeer()
		if peer1.URL.String() != u2.String() {
			t.Errorf("Expected peer %s, got %s", u1.String(), peer1.URL.String())
		}

		peer2 := sp.GetNextPeer()
		if peer2.URL.String() != u1.String() {
			t.Errorf("Expected peer %s, got %s", u2.String(), peer2.URL.String())
		}
	})

	t.Run("MarkBackendStatus", func(t *testing.T) {
		serverPool = ServerPool{}
		sp := &ServerPool{}
		u, _ := url.Parse("http://localhost:8080")
		b := &Backend{URL: u}
		sp.AddBackend(b)

		sp.MarkBackendStatus(u, true)
		if !b.IsAlive() {
			t.Error("Expected backend to be alive")
		}

		sp.MarkBackendStatus(u, false)
		if b.IsAlive() {
			t.Error("Expected backend to be down")
		}
	})

	t.Run("GetNextPeer with no available peer", func(t *testing.T) {
		serverPool = ServerPool{}
		sp := &ServerPool{}
		u1, _ := url.Parse("http://localhost:8080")
		b1 := &Backend{URL: u1}
		b1.SetAlive(false)
		sp.AddBackend(b1)

		peer := sp.GetNextPeer()
		if peer != nil {
			t.Error("Expected no peer, got one")
		}
	})
}
