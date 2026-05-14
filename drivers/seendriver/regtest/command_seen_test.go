package regtest

import (
	"fmt"
	"testing"
	"time"

	"github.com/fluffle/sp0rkle/regtest"
)

func testSeen(t *testing.T) {
	t.Run("have_seen_privmsg", testHaveSeenPrivmsg)
	t.Run("have_seen_nick", testHaveSeenNick)
	t.Run("not_seen", testNotSeen)
}

func testHaveSeenPrivmsg(t *testing.T) {
	h.Privmsg(h.Channel(), "hello there")
	time.Sleep(100 * time.Millisecond)

	h.Privmsg(h.Channel(), fmt.Sprintf("%s: seen %s", h.BotNick(), h.Me().Nick))
	_, err := h.Expect(
		regtest.Contains("I last saw "+h.Me().Nick+" on"),
		2*time.Second,
	)
	if err != nil {
		t.Errorf("expected seen response, got: %v", err)
	}
}

func testHaveSeenNick(t *testing.T) {
	h.Privmsg(h.Channel(), "hello")
	time.Sleep(100 * time.Millisecond)

	h.Privmsg(h.Channel(), fmt.Sprintf("%s: seen %s", h.BotNick(), h.Me().Nick))
	_, err := h.Expect(
		regtest.Contains("I last saw "+h.Me().Nick+" on"),
		2*time.Second,
	)
	if err != nil {
		t.Errorf("expected seen response, got: %v", err)
	}
}

func testNotSeen(t *testing.T) {
	h.Privmsg(h.Channel(), fmt.Sprintf("%s: seen nonexistentuser", h.BotNick()))
	_, err := h.Expect(
		regtest.Contains("Haven't seen nonexistentuser before"),
		2*time.Second,
	)
	if err != nil {
		t.Errorf("expected 'not seen' response, got: %v", err)
	}
}
