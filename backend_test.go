package simplelb

import (
	"net/url"
	"testing"
)

func TestBackend(t *testing.T) {
	t.Run("SetAlive and IsAlive", func(t *testing.T) {
		u, _ := url.Parse("http://localhost:8080")
		b := &Backend{URL: u}

		if b.IsAlive() {
			t.Error("Expected backend to be initially down")
		}

		b.SetAlive(true)
		if !b.IsAlive() {
			t.Error("Expected backend to be alive")
		}

		b.SetAlive(false)
		if b.IsAlive() {
			t.Error("Expected backend to be down")
		}
	})
}
