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
	t.Run("calc", testCalc)
	t.Run("date", testDate)
	t.Run("netmask", testNetmask)
	t.Run("chr", testChr)
	t.Run("ord", testOrd)
}

func testCalc(t *testing.T)    { t.Skip("calc driver test") }
func testDate(t *testing.T)    { t.Skip("date driver test") }
func testNetmask(t *testing.T) { t.Skip("netmask driver test") }
func testChr(t *testing.T)     { t.Skip("chr driver test") }
func testOrd(t *testing.T)     { t.Skip("ord driver test") }
