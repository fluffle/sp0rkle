//go:build integration

package regtest

import (
	"os"
	"testing"
	"time"
)

func TestFullLifecycle(t *testing.T) {
	if os.Getenv("REGTEST_SERVER") == "" {
		t.Skip("REGTEST_SERVER not set, skipping integration test")
	}

	h, err := Start()
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	defer func() {
		if err := Stop(); err != nil {
			t.Logf("Stop error: %v", err)
		}
	}()

	// Verify the harness has a valid channel.
	if h.Channel() == "" {
		t.Error("Channel() returned empty string")
	}

	// Verify SendAndExpect times out gracefully (no bot registered).
	_, err = h.SendAndExpect("test", Exact("test"), 100*time.Millisecond)
	if err == nil {
		t.Error("SendAndExpect should timeout with no bot responding")
	}
}
