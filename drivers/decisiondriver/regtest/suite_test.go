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
	t.Run("rand", testRand)
	t.Run("decide", testDecide)
	t.Run("choose", testChoose)
}

func testRand(t *testing.T)    { t.Skip("rand driver test") }
func testDecide(t *testing.T)  { t.Skip("decide driver test") }
func testChoose(t *testing.T)  { t.Skip("choose driver test") }
