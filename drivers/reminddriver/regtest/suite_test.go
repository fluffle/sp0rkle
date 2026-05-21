//go:build integration

package regtest

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/fluffle/sp0rkle/regtest"
)

var h *regtest.Harness

func TestMain(m *testing.M) {
	var err error
	if h, err = regtest.Start(context.Background()); err != nil {
		panic("Start: " + err.Error())
	}
	code := m.Run()
	if err = h.Stop(); err != nil {
		panic("Stop: " + err.Error())
	}
	os.Exit(code)
}

func TestCommands(t *testing.T) {
	t.Run("tell", testTell)
	t.Run("ask", testAsk)
	t.Run("remind_list", testRemindList)
	t.Run("remind_set", testRemindSet)
	t.Run("remind_del", testRemindDel)
	t.Run("snooze", testSnooze)
	t.Run("zone", testZone)
	t.Run("unzone", testUnzone)
}

func TestHandlers(t *testing.T) {
	t.Run("tell_check", testTellCheck)
}

func mustRemind(t *testing.T, msg string) {
	_, err := h.CommandAndExpect("remind me "+msg+" in 1 hour", h.Contains("okay, i'll remind"))
	if err != nil {
		t.Fatalf("expected remind ack, got: %v", err)
	}
}

func mustClearReminders(t *testing.T) {
	for {
		w := h.Expect(h.Contains("You have"))
		// Always ask privately.
		h.Privmsg(h.BotNick, "remind list")
		line, err := w.Wait()
		if err != nil {
			t.Fatalf("expected remind list response, got: %v", err)
			return
		}
		if strings.Contains(line.Text(), "no reminders set") {
			// done!
			return
		}
		// Since we have to relist every time, just delete first one.
		w = h.Expect(h.Contains("forget that one"))
		h.Privmsg(h.BotNick, "remind del 1")
		_, err = w.Wait()
		if err != nil {
			t.Fatalf("expected remind del response, got: %v", err)
			return
		}
	}
}
