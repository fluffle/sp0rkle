//go:build integration

package regtest

import (
	"testing"
	"time"
)

func testLinesSelf(t *testing.T) {
	h.Privmsg(h.Channel, "hello world")
	time.Sleep(100 * time.Millisecond)
	_, err := h.CommandAndExpect("lines", h.Contains("has said"))
	if err != nil {
		t.Errorf("expected stats response for self, got: %v", err)
	}
}

func testLinesOther(t *testing.T) {
	h.Privmsg(h.Channel, "hello from test0rkle")
	time.Sleep(100 * time.Millisecond)
	_, err := h.CommandAndExpect("lines "+h.Me().Nick, h.Contains("has said"))
	if err != nil {
		t.Errorf("expected stats response for specific nick, got: %v", err)
	}
}

func testLinesUnknown(t *testing.T) {
	_, err := h.CommandAndExpect("lines nobody", h.Contains(`not seen "nobody"`))
	if err != nil {
		t.Errorf("expected stats response for specific nick, got: %v", err)
	}
}

func testStatsAlias(t *testing.T) {
	h.Privmsg(h.Channel, "another message")
	time.Sleep(100 * time.Millisecond)
	_, err := h.CommandAndExpect("stats", h.Contains("has said"))
	if err != nil {
		t.Errorf("expected stats alias response, got: %v", err)
	}
}

func testTopten(t *testing.T) {
	h.Privmsg(h.Channel, "first message for topten")
	time.Sleep(100 * time.Millisecond)
	_, err := h.CommandAndExpect("topten", h.Contains("#1:"))
	if err != nil {
		t.Errorf("expected topten response with ranking, got: %v", err)
	}
}

func testTop10Alias(t *testing.T) {
	h.Privmsg(h.Channel, "message for top10 alias")
	time.Sleep(100 * time.Millisecond)
	_, err := h.CommandAndExpect("top10", h.Contains("#1:"))
	if err != nil {
		t.Errorf("expected top10 alias response with ranking, got: %v", err)
	}
}
