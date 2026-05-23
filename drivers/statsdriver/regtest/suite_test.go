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
	t.Run("lines_self", testLinesSelf)
	t.Run("lines_other", testLinesOther)
	t.Run("lines_unknown", testLinesUnknown)
	t.Run("stats_alias", testStatsAlias)
	t.Run("topten", testTopten)
	t.Run("top10_alias", testTop10Alias)
}

func TestHandlers(t *testing.T) {
	t.Run("record_privmsg", testRecordPrivmsg)
	t.Run("record_action", testRecordAction)
	t.Run("record_multiple", testRecordMultiple)
}
