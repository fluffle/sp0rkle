package regtest

import (
	"fmt"
	"testing"
	"time"

	"github.com/fluffle/sp0rkle/regtest"
)

func testIgnore(t *testing.T) {
	t.Run("basic", testIgnoreBasic)
	t.Run("case_insensitive", testIgnoreCaseInsensitive)
}

func testUnignore(t *testing.T) {
	t.Run("basic", testUnignoreBasic)
}

func testIgnoreBasic(t *testing.T) {
	h.Privmsg(h.Channel(), fmt.Sprintf("%s: ignore testnick", h.BotNick()))
	_, err := h.Expect(
		regtest.Contains("I'll ignore 'testnick'"),
		2*time.Second,
	)
	if err != nil {
		t.Errorf("expected ignore response, got: %v", err)
	}
}

func testIgnoreCaseInsensitive(t *testing.T) {
	h.Privmsg(h.Channel(), fmt.Sprintf("%s: IGNORE TESTNICK", h.BotNick()))
	_, err := h.Expect(
		regtest.Contains("I'll ignore 'testnick'"),
		2*time.Second,
	)
	if err != nil {
		t.Errorf("expected case-insensitive ignore response, got: %v", err)
	}
}

func testUnignoreBasic(t *testing.T) {
	h.Privmsg(h.Channel(), fmt.Sprintf("%s: unignore testnick", h.BotNick()))
	_, err := h.Expect(
		regtest.Contains("No longer ignoring"),
		2*time.Second,
	)
	if err != nil {
		t.Errorf("expected unignore response, got: %v", err)
	}
}
