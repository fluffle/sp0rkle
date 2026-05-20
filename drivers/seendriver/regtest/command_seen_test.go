//go:build integration

package regtest

import (
	"fmt"
	"testing"
)

func testSeen(t *testing.T) {
	// Other variations of the seen command are tested in the handler tests.
	t.Run("not_seen", testNotSeen)
}

func testNotSeen(t *testing.T) {
	w := h.Expect(h.Contains("Haven't seen nonexistentuser before"))
	h.Privmsg(h.Channel, fmt.Sprintf("%s: seen nonexistentuser", h.BotNick))
	if _, err := w.Wait(); err != nil {
		t.Errorf("expected 'not seen' response, got: %v", err)
	}
}
