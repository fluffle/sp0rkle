//go:build integration

package regtest

import (
	"context"
	"os"
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
	t.Run("markov_me", testMarkovMe)
	t.Run("dont_markov_me", testDontMarkovMe)
	t.Run("markov_nick", testMarkovNick)
	t.Run("insult", testInsult)
	t.Run("learn", testLearn)
}

func TestHandlers(t *testing.T) {
	t.Run("record_markov_privmsg", testRecordMarkovPrivmsg)
	t.Run("record_markov_action", testRecordMarkovAction)
	t.Run("record_markov_not_enabled", testRecordMarkovNotEnabled)
	t.Run("record_markov_addressed", testRecordMarkovAddressed)
}

func mustEnable(t *testing.T) {
	_, err := h.CommandAndExpect("markov me", h.Contains("I'll markov you"))
	if err != nil {
		t.Fatalf("failed to enable markov: %v", err)
	}
}

func mustDisable(t *testing.T) {
	_, err := h.CommandAndExpect("don't markov me", h.Contains("Sure, bro, I'll stop"))
	if err != nil {
		t.Fatalf("expected disable response, got: %v", err)
	}
}
