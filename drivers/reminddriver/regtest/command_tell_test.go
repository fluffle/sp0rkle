//go:build integration

package regtest

import (
	"testing"
)

func testTell(t *testing.T) {
	t.Run("missing_args", testTellMissingArgs)
	t.Run("self_tell", testTellSelf)
	t.Run("tell_me", testTellMe)
}

func testAsk(t *testing.T) {
	t.Run("missing_args", testAskMissingArgs)
	t.Run("ask_other", testAskOther)
}

func testTellMissingArgs(t *testing.T) {
	_, err := h.CommandAndExpect("tell", h.Contains("Tell who what?"))
	if err != nil {
		t.Errorf("expected 'Tell who what?' response, got: %v", err)
	}
}

func testTellSelf(t *testing.T) {
	_, err := h.CommandAndExpect("tell " + h.Me().Nick + " something", h.Contains("You're a dick"))
	if err != nil {
		t.Errorf("expected 'You're a dick' response, got: %v", err)
	}
}

func testTellMe(t *testing.T) {
	_, err := h.CommandAndExpect("tell me something", h.Contains("You're a dick"))
	if err != nil {
		t.Errorf("expected 'You're a dick' response for 'tell me', got: %v", err)
	}
}

func testAskMissingArgs(t *testing.T) {
	_, err := h.CommandAndExpect("ask", h.Contains("Tell who what?"))
	if err != nil {
		t.Errorf("expected 'Tell who what?' response for 'ask', got: %v", err)
	}
}

func testAskOther(t *testing.T) {
	_, err := h.CommandAndExpect("ask otheruser a secret", h.Contains("okay, i'll tell"))
	if err != nil {
		t.Errorf("expected ask acknowledgement, got: %v", err)
	}
}
