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
	t.Run("qadd", testQadd)
	t.Run("qdel", testQdel)
}

func testQadd(t *testing.T) { t.Skip("qadd driver test") }
func testQdel(t *testing.T) { t.Skip("qdel driver test") }
