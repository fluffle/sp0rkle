//go:build integration

package regtest

import (
	"testing"
	"time"
)

func testInsult(t *testing.T) {
	t.Run("no_data", func(t *testing.T) {
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

	t.Run("learn_insult_then_insult", func(t *testing.T) {
		// Learn an insult first
		_, err := h.CommandAndExpect("learn insult you are a terrible person", h.Contains("Ta"))
		if err != nil {
			t.Fatalf("failed to learn insult: %v", err)
		}
		time.Sleep(100 * time.Millisecond)

		// Now insult should work
		_, err = h.CommandAndExpect("insult someone", h.Contains("someone:"))
		if err != nil {
			t.Errorf("expected insult response, got: %v", err)
		}
	})

	t.Run("insult_me", func(t *testing.T) {
		_, err := h.CommandAndExpect("insult me", h.Contains(h.Me().Nick+":"))
		if err != nil {
			t.Errorf("expected insult with harness nick, got: %v", err)
		}
	})

	t.Run("insult_no_target", func(t *testing.T) {
		// Without target, bot just replies with the sentence
		_, err := h.CommandAndExpect("insult", h.Contains("terrible person"))
		if err != nil {
			t.Errorf("expected insult without target, got: %v", err)
		}
	})
}
