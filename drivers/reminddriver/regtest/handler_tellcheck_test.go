//go:build integration

package regtest

import (
	"testing"
)

func testTellCheck(t *testing.T) {
	t.Run("nick", testTellForNick)
	t.Run("join", testTellForJoin)
	t.Run("privmsg", testTellForPrivmsg)
	t.Run("action", testTellForAction)
}

func testTellForNick(t *testing.T) {
	_, err := h.CommandAndExpect("tell otheruser hello there", h.Contains("okay, i'll tell"))
	if err != nil {
		t.Errorf("expected tell acknowledgement, got: %v", err)
	}
	// If we change our nick to "otheruser" we should be told a thing.
	w := h.Expect(h.Contains("asked me to tell you hello there"))
	old := h.MustRenick(t, "otheruser")
	if _, err = w.Wait(); err != nil {
		t.Errorf("expected to be told a thing, got: %v", err)
	}
	_ = h.MustRenick(t, old)
}

func testTellForJoin(t *testing.T) {
	_, err := h.CommandAndExpect("tell otheruser hello there", h.Contains("okay, i'll tell"))
	if err != nil {
		t.Errorf("expected tell acknowledgement, got: %v", err)
	}
	// If we change our nick to "otheruser" we should be told a thing.
	// But we must part before changing nick so the bot does not know!
	h.MustPart(t)
	old := h.MustRenick(t, "otheruser")
	w := h.Expect(h.Contains("asked me to tell you hello there"))
	h.MustJoin(t)
	if _, err = w.Wait(); err != nil {
		t.Errorf("expected to be told a thing, got: %v", err)
	}
	_ = h.MustRenick(t, old)
}

func testTellForPrivmsg(t *testing.T) {
	t.Skip("needs a second nick already present")
}

func testTellForAction(t *testing.T) {
	t.Skip("needs a second nick already present")
}
