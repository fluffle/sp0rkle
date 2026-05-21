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
	t.Run("karma", testKarma)
}

func TestHandlers(t *testing.T) {
	t.Run("record_privmsg", testRecordPrivmsg)
	t.Run("record_action", testRecordAction)
	t.Run("record_multiple", testRecordMultiple)
	t.Run("record_bracketed", testRecordBracketed)
	t.Run("record_nothing", testRecordNothing)
}
