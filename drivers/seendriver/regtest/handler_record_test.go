package regtest

import (
	"testing"
	"time"

	"github.com/fluffle/sp0rkle/regtest"
)

func testSmoke(t *testing.T) {
	h.Privmsg(h.Channel(), "schmoke?")
	_, err := h.Expect(
		regtest.Contains("smoke"),
		2*time.Second,
	)
	if err != nil {
		t.Errorf("expected smoke response, got: %v", err)
	}
}

func testRecordPrivmsg(t *testing.T) {
	h.Privmsg(h.Channel(), "this is a test message for recording")
	time.Sleep(100 * time.Millisecond)

	h.Privmsg(h.Channel(), "seen "+h.Me().Nick)
	_, err := h.Expect(
		regtest.Contains("in "+h.Channel()),
		2*time.Second,
	)
	if err != nil {
		t.Errorf("expected recorded privmsg, got: %v", err)
	}
}

func testRecordJoin(t *testing.T) {
	h.Privmsg(h.Channel(), "seen "+h.BotNick())
	_, err := h.Expect(
		regtest.Contains("joining"),
		2*time.Second,
	)
	if err != nil {
		t.Errorf("expected join record, got: %v", err)
	}
}

func testRecordNick(t *testing.T) {
	t.Skip("requires second client for nick change")
}

func testRecordKick(t *testing.T) {
	t.Skip("requires second client for kick testing")
}
