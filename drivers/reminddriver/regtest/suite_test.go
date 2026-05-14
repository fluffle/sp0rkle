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
	t.Run("tell", testTell)
	t.Run("ask", testAsk)
	t.Run("remind_list", testRemindList)
	t.Run("remind_del", testRemindDel)
	t.Run("remind_set", testRemindSet)
}

func testTell(t *testing.T)    { t.Skip("tell driver test") }
func testAsk(t *testing.T)     { t.Skip("ask driver test") }
func testRemindList(t *testing.T) { t.Skip("remind_list driver test") }
func testRemindDel(t *testing.T) { t.Skip("remind_del driver test") }
func testRemindSet(t *testing.T) { t.Skip("remind_set driver test") }
