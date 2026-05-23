//go:build integration

package regtest

import (
	"testing"
	"time"
)

func testRecordMarkovPrivmsg(t *testing.T) {
	mustEnable(t)
	defer mustDisable(t)

	// Send unaddressed public PRIVMSG — should be recorded
	h.Privmsg(h.Channel, "recording this message")
	time.Sleep(100 * time.Millisecond)

	// Verify by generating a markov sentence
	// Because each test clears markov state it should repeat what we said.
	_, err := h.CommandAndExpect("markov "+h.Me().Nick, h.Contains("would say: recording this message"))
	if err != nil {
		t.Errorf("expected markov output after recording, got: %v", err)
	}
}

func testRecordMarkovAction(t *testing.T) {
	t.Skip("markov $nick only uses SENTENCE_START, can't retrieve data")
	// otherwise...
	mustEnable(t)
	defer mustDisable(t)

	// Send unaddressed public ACTION — should be recorded
	h.Action(h.Channel, "does a funny dance")
	time.Sleep(100 * time.Millisecond)

	// Verify by generating a markov sentence
	_, err := h.CommandAndExpect("markov "+h.Me().Nick, h.Contains("would say: does a funny dance"))
	if err != nil {
		t.Errorf("expected markov output after action, got: %v", err)
	}
}

func testRecordMarkovNotEnabled(t *testing.T) {
	// Send unaddressed public PRIVMSG — should NOT be recorded
	h.Privmsg(h.Channel, "this should not be recorded at all")
	time.Sleep(100 * time.Millisecond)

	// Verify markov data is not recorded — should say "not recording"
	_, err := h.CommandAndExpect("markov "+h.Me().Nick, h.Contains("not recording markov data"))
	if err != nil {
		t.Errorf("expected not recording response, got: %v", err)
	}
}

func testRecordMarkovAddressed(t *testing.T) {
	mustEnable(t)

	// Send addressed message — should NOT be recorded (addressed messages are excluded)
	h.Command("this should not be recorded")
	time.Sleep(100 * time.Millisecond)

	// Markov enabled but no data should result in failure
	_, err := h.CommandAndExpect("markov "+h.Me().Nick, h.Contains("!SENTENCE_START"))
	if err != nil {
		t.Errorf("expected markov failure, got: %v", err)
	}

	// Disabling markov with nothing logged should result in failure
	_, err = h.CommandAndExpect("don't markov me", h.Contains("bucket not found"))
	if err != nil {
		t.Errorf("expected disable response, got: %v", err)
	}
}
