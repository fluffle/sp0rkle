//go:build integration

package regtest

import (
	"testing"
	"time"
)

func testRecordPrivmsg(t *testing.T) {
	h.Privmsg(h.Channel, "smoke++")
	time.Sleep(100 * time.Millisecond)
	_, err := h.CommandAndExpect("karma smoke",
		h.Contains("'smoke' has a karma of 1 after 1 votes. Last upvoted by "+h.Me().Nick))
	if err != nil {
		t.Errorf("expected 'smoke' karma of 1 after 1 vote, got: %v", err)
	}
}

func testRecordAction(t *testing.T) {
	h.Action(h.Channel, "testaction++")
	time.Sleep(100 * time.Millisecond)
	_, err := h.CommandAndExpect("karma testaction",
		h.Contains("'testaction' has a karma of 1 after 1 votes. Last upvoted by "+h.Me().Nick))
	if err != nil {
		t.Errorf("expected 'testaction' karma of 1 after 1 vote, got: %v", err)
	}
}

func testRecordMultiple(t *testing.T) {
	h.Privmsg(h.Channel, "alpha++ beta-- gamma++")
	time.Sleep(100 * time.Millisecond)
	_, err := h.CommandAndExpect("karma alpha",
		h.Contains("'alpha' has a karma of 1 after 1 votes. Last upvoted by "+h.Me().Nick))
	if err != nil {
		t.Errorf("expected 'alpha' karma of 1 after 1 vote, got: %v", err)
	}
	_, err = h.CommandAndExpect("karma beta",
		h.Contains("'beta' has a karma of -1 after 1 votes. Last downvoted by "+h.Me().Nick))
	if err != nil {
		t.Errorf("expected 'beta' karma of -1 after 1 vote, got: %v", err)
	}
	_, err = h.CommandAndExpect("karma gamma",
		h.Contains("'gamma' has a karma of 1 after 1 votes. Last upvoted by "+h.Me().Nick))
	if err != nil {
		t.Errorf("expected 'gamma' karma of 1 after 1 vote, got: %v", err)
	}
}

func testRecordBracketed(t *testing.T) {
	h.Privmsg(h.Channel, "(bracketed thing)++")
	time.Sleep(100 * time.Millisecond)
	_, err := h.CommandAndExpect("karma bracketed thing",
		h.Contains("'bracketed thing' has a karma of 1 after 1 votes. Last upvoted by "+h.Me().Nick))
	if err != nil {
		t.Errorf("expected 'bracketed thing' karma of 1 after 1 vote, got: %v", err)
	}
}

func testRecordNothing(t *testing.T) {
	h.Privmsg(h.Channel, "no karma markers here at all")
	time.Sleep(100 * time.Millisecond)
	_, err := h.CommandAndExpect("karma randomquery", h.Contains("No karma found for 'randomquery'"))
	if err != nil {
		t.Errorf("expected 'not found' response for random query, got: %v", err)
	}
}
