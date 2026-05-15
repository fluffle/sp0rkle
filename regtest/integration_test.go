//go:build integration

package regtest

import (
	"context"
	"testing"
)

func TestFullLifecycle(t *testing.T) {
	ctx := context.Background()
	h, err := Start(ctx)
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	defer func() {
		if err := h.Stop(); err != nil {
			t.Logf("Stop error: %v", err)
		}
	}()

	// Verify the harness has a valid channel.
	if h.Channel == "" {
		t.Error("Channel is empty")
	}

	// Verify SendAndExpect can roundtrip a help to the bot
	_, err = h.CommandAndExpect("help", h.Contains("github.com/fluffle/sp0rkle/wiki"))
	if err != nil {
		t.Errorf("SendAndExpect err = %v, want nil", err)
	}
}
