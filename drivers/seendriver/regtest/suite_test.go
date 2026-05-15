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
	t.Run("seen", testSeen)
}

func TestHandlers(t *testing.T) {
	t.Run("smoke", testSmoke)
	t.Run("privmsg", testRecordPrivmsg)
	t.Run("join", testRecordJoin)
	t.Run("nick", testRecordNick)
	t.Run("kick", testRecordKick)
}
