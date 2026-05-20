//go:build integration

package regtest

import (
	"fmt"
	"testing"
	"time"

	"github.com/fluffle/goirc/client"
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
	w := h.Expect(h.Contains("You last went for a smoke"))
	h.Privmsg(h.Channel, "t0kez0r")
	_, err := w.Wait()
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
	wasMe := func(line *client.Line) bool {
		return line.Nick == h.Me().Nick && line.Args[0] == h.Channel
	}
	w := h.ExpectEvent(client.PART, regtest.PatternFunc(wasMe))
	h.Part(h.Channel)
	w.MustWait(t)

	w = h.ExpectEvent(client.JOIN, regtest.PatternFunc(wasMe))
	h.Join(h.Channel)
	w.MustWait(t)

	if err := seenCmd(h.Contains("joining "+h.Channel)); err != nil {
		t.Errorf("expected join response, got: %v", err)
	}
}

func testRecordNick(t *testing.T) {
	prev := h.Me().Nick
	next := "testNickPleaseIgnore"

	fwdMatch := func(line *client.Line) bool {
		return line.Nick == prev && line.Args[0] == next
	}
	w := h.ExpectEvent(client.NICK, regtest.PatternFunc(fwdMatch))
	h.Nick(next)
	w.MustWait(t)

	revMatch := func(line *client.Line) bool {
		return line.Nick == next && line.Args[0] == prev
	}
	w = h.ExpectEvent(client.NICK, regtest.PatternFunc(revMatch))
	h.Nick(prev)
	w.MustWait(t)

	if err := seenCmd(h.Contains(fmt.Sprintf("changing their nick to '%s'", next))); err != nil {
		t.Errorf("expected seen response, got: %v", err)
	}
}

func testRecordKick(t *testing.T) {
	t.Skip("requires second client for kick testing")
}
