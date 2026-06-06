//go:build integration

package regtest

import (
	"testing"
)

func testInsult(t *testing.T) {
	t.Run("no_data", func(t *testing.T) {
		t.Skip("requires http server in regtest harness")
		// No insult data learned yet, so markov error
		_, err := h.CommandAndExpect("insult someone", h.Contains("markov error"))
		if err != nil {
			t.Errorf("expected markov error for empty insult, got: %v", err)
		}
	})

	t.Run("insult_yourself", func(t *testing.T) {
		_, err := h.CommandAndExpect("insult yourself", h.Contains("Ha, you're funny"))
		if err != nil {
			t.Errorf("expected yourself response, got: %v", err)
		}
	})

	t.Run("insult_bot", func(t *testing.T) {
		_, err := h.CommandAndExpect("insult "+h.BotNick, h.Contains("Ha, you're funny"))
		if err != nil {
			t.Errorf("expected bot response, got: %v", err)
		}
	})

	t.Run("insult_me", func(t *testing.T) {
		t.Skip("requires http server in regtest harness")
		_, err := h.CommandAndExpect("insult me", h.Contains(h.Me().Nick+":"))
		if err != nil {
			t.Errorf("expected insult with harness nick, got: %v", err)
		}
	})

	t.Run("insult_no_target", func(t *testing.T) {
		t.Skip("requires http server in regtest harness")
		_, err := h.CommandAndExpect("insult", h.Contains("TODO"))
		if err != nil {
			t.Errorf("expected insult without target, got: %v", err)
		}
	})
}
