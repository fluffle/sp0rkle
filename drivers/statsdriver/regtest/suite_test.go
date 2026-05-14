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
	t.Run("lines", testLines)
	t.Run("stats", testStats)
	t.Run("topten", testTopten)
}

func testLines(t *testing.T) { t.Skip("lines driver test") }
func testStats(t *testing.T) { t.Skip("stats driver test") }
func testTopten(t *testing.T) { t.Skip("topten driver test") }
