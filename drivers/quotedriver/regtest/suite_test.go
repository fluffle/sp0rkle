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
	t.Run("qadd", testQadd)
	t.Run("quote_add", testQuoteAdd)
	t.Run("add_quote", testAddQuote)
	t.Run("qdel", testQdel)
	t.Run("quote_del", testQuoteDel)
	t.Run("del_quote", testDelQuote)
	t.Run("del_invalid_id", testDelInvalidID)
	t.Run("del_nonexistent", testDelNonexistent)
	t.Run("fetch", testFetch)
	t.Run("lookup", testLookup)
	t.Run("lookup_no_match", testLookupNoMatch)
}
