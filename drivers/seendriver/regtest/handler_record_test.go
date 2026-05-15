//go:build integration

package regtest

import (
	"testing"
	"time"

	"github.com/fluffle/sp0rkle/regtest"
)

func seenCmd(p regtest.Pattern) error {
	_, err := h.CommandAndExpect("seen "+h.Me().Nick, p)
	return err
}

func testSmoke(t *testing.T) {
	h.Privmsg(h.Channel, "schmoke?")
	time.Sleep(100 * time.Millisecond)
	if err := seenCmd(h.Contains("going for a smoke"));  err != nil {
		t.Errorf("expected smoke response, got: %v", err)
	}
	h.Privmsg(h.Channel, "t0kez0r")
	_, err := h.Expect(h.Contains("You last went for a smoke"))
	if err != nil {
		t.Errorf("expected second smoke respone, got: %v", err)
	}
}

func testRecordPrivmsg(t *testing.T) {
	h.Privmsg(h.Channel, "hello there")
	time.Sleep(100 * time.Millisecond)
	if err := seenCmd(h.Contains(`saying 'hello there'`)); err != nil {
		t.Errorf("expected privmsg response, got: %v", err)
	}
}

func testRecordJoin(t *testing.T) {
	h.Part(h.Channel)
	time.Sleep(100 * time.Millisecond)
	h.Join(h.Channel)
	time.Sleep(100 * time.Millisecond)
	if err := seenCmd(h.Contains("joining "+h.Channel)); err != nil {
		t.Errorf("expected join response, got: %v", err)
	}
}

func testRecordNick(t *testing.T) {
	old := h.Me().Nick
	h.Nick("testNickPleaseIgnore")
	time.Sleep(100 * time.Millisecond)
	h.Nick(old)
	time.Sleep(100 * time.Millisecond)

	if err := seenCmd(h.Contains(`changing their nick to 'testNickPleaseIgnore'`)); err != nil {
		t.Errorf("expected seen response, got: %v", err)
	}
}

func testRecordKick(t *testing.T) {
	t.Skip("requires second client for kick testing")
}
