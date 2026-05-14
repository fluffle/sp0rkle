package regtest

import (
	"os"
	"testing"

	"github.com/fluffle/sp0rkle/regtest"
)

var h *regtest.Harness

func TestMain(m *testing.M) {
	var err error
	if h, err = regtest.Start(); err != nil {
		panic("Start: " + err.Error())
	}
	code := m.Run()
	if err = regtest.Stop(); err != nil {
		panic("Stop: " + err.Error())
	}
	os.Exit(code)
}

func TestCommands(t *testing.T) {
	t.Run("karma", testKarma)
}

func testKarma(t *testing.T) { t.Skip("karma driver test") }
