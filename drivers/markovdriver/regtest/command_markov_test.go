//go:build integration

package regtest

import (
	"testing"
	"time"
)

func testMarkovMe(t *testing.T) {
	_, err := h.CommandAndExpect("markov me", h.Contains("I'll markov you"))
	if err != nil {
		t.Errorf("expected markov me response, got: %v", err)
	}
}

func testDontMarkovMe(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		// Enable and record some data so ClearTag succeeds
		mustEnable(t)
		h.Privmsg(h.Channel, "some data to record")
		time.Sleep(100 * time.Millisecond)
		_, err := h.CommandAndExpect("don't markov me", h.Contains("Sure, bro, I'll stop"))
		if err != nil {
			t.Errorf("expected disable response with bro, got: %v", err)
		}
		_, err = h.CommandAndExpect("markov "+h.Me().Nick, h.Contains("not recording"))
		if err != nil {
			t.Errorf("expected not recording response, got: %v", err)
		}
	})

	t.Run("with_bro", func(t *testing.T) {
		// Enable and record some data so ClearTag succeeds
		mustEnable(t)
		h.Privmsg(h.Channel, "another data to record")
		time.Sleep(100 * time.Millisecond)
		_, err := h.CommandAndExpect("don't markov me, bro", h.Contains("Sure, bro, I'll stop"))
		if err != nil {
			t.Errorf("expected disable response with bro, got: %v", err)
		}
	})

	t.Run("no_data", func(t *testing.T) {
		// When no markov data exists, ClearTag returns "bucket not found"
		_, err := h.CommandAndExpect("don't markov me", h.Contains("Failed to clear tag"))
		if err != nil {
			t.Errorf("expected error response when no data, got: %v", err)
		}
	})
}

func testMarkovNick(t *testing.T) {
	t.Run("no_args", func(t *testing.T) {
		_, err := h.CommandAndExpect("markov", h.Contains("Be who?"))
		if err != nil {
			t.Errorf("expected 'Be who?' response, got: %v", err)
		}
	})

	t.Run("self_as_bot", func(t *testing.T) {
		_, err := h.CommandAndExpect("markov "+h.BotNick, h.Contains("Ha, you're funny"))
		if err != nil {
			t.Errorf("expected self-as-bot response, got: %v", err)
		}
	})

	t.Run("not_recording_self", func(t *testing.T) {
		_, err := h.CommandAndExpect("markov "+h.Me().Nick, h.Contains("not recording markov data"))
		if err != nil {
			t.Errorf("expected not recording response for self, got: %v", err)
		}
	})

	t.Run("not_recording_other", func(t *testing.T) {
		_, err := h.CommandAndExpect("markov someonewhodoesnotexist", h.Contains("Not recording markov data"))
		if err != nil {
			t.Errorf("expected not recording response for other, got: %v", err)
		}
	})

	t.Run("success_after_data", func(t *testing.T) {
		mustEnable(t)
		defer mustDisable(t)

		// Feed some data so markov chains exist
		h.Privmsg(h.Channel, "hello world test")
		time.Sleep(100 * time.Millisecond)
		h.Privmsg(h.Channel, "hello again world")
		time.Sleep(100 * time.Millisecond)
		h.Privmsg(h.Channel, "test hello world")
		time.Sleep(100 * time.Millisecond)

		_, err := h.CommandAndExpect("markov "+h.Me().Nick, h.Contains("would say:"))
		if err != nil {
			t.Errorf("expected markov output, got: %v", err)
		}

	})
}
