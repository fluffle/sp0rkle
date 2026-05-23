//go:build integration

package regtest

import (
	"context"
	"os"
	"regexp"
	"testing"
	"time"

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
	t.Run("edit", testEdit)
	t.Run("forget", testForget)
	t.Run("delete", testDelete)
	t.Run("info", testInfo)
	t.Run("literal", testLiteral)
	t.Run("replace", testReplace)
	t.Run("search", testSearch)
}

func TestHandlers(t *testing.T) {
	t.Run("insert", testInsert)
	t.Run("lookup", testLookup)
	t.Run("recurse", testLookupRecurse)
}

// createFactoid creates a factoid and returns the key.
func createFactoid(t *testing.T, key, value string) string {
	_, err := h.CommandAndExpect(
		key+" := "+value,
		h.Regex(regexp.MustCompile(`I now know \d+ things about '`+regexp.QuoteMeta(key)+`'`)))
	if err != nil {
		t.Fatalf("failed to create factoid '%s': %v", key, err)
	}
	return key
}

// triggerLookupFast tries to trigger a lookup with short timeouts.
// Returns true if the bot responded with the expected value.
func triggerLookupFast(t *testing.T, key, expected string, retries int) bool {
	for i := 0; i < retries; i++ {
		w := h.Expect(h.Contains(expected))
		w.Timeout = 100 * time.Millisecond
		h.Privmsg(h.Channel, key)
		if _, err := w.Wait(); err == nil {
			return true
		}
		time.Sleep(50 * time.Millisecond)
	}
	return false
}

// mustTriggerLookup wraps triggerLookupFast and calls t.Fatal if it fails
func mustTriggerLookup(t *testing.T, key, expected string) {
	if !triggerLookupFast(t, key, expected, 1) {
		t.Fatalf("look up key %q with expected value %q failed", key, expected)
	}
}

