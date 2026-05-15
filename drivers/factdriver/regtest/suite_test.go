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
	t.Run("chance", testChance)
	t.Run("that_is", testThatIs)
	t.Run("forget", testForget)
	t.Run("info", testInfo)
}

func testChance(t *testing.T) { t.Skip("chance driver test") }
func testThatIs(t *testing.T) { t.Skip("that_is driver test") }
func testForget(t *testing.T) { t.Skip("forget driver test") }
func testInfo(t *testing.T)  { t.Skip("info driver test") }
