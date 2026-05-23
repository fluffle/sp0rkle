//go:build integration

package regtest

import (
	"regexp"
	"testing"

	"github.com/fluffle/goirc/client"
)

// ----- insert handler -----

func testInsert(t *testing.T) {
	t.Run("colon_equals", func(t *testing.T) {
		_, err := h.CommandAndExpect("regtest_insert_colon_key := insert_value", h.Contains("I now know 1 things about"))
		if err != nil {
			t.Errorf("expected insert confirmation, got: %v", err)
		}
	})

	t.Run("colon_is", func(t *testing.T) {
		_, err := h.CommandAndExpect("regtest_insert_is_key :is is_value", h.Contains("I now know 1 things about"))
		if err != nil {
			t.Errorf("expected insert confirmation for :is, got: %v", err)
		}

		// Verify :is creates "key is value" format
		_, err = h.CommandAndExpect("literal regtest_insert_is_key", h.Contains("regtest_insert_is_key is is_value"))
		if err != nil {
			t.Errorf("expected 'key is value' format for :is, got: %v", err)
		}
	})

	t.Run("multiple_values", func(t *testing.T) {
		// First value
		_, err := h.CommandAndExpect("regtest_multi_key := first_value", h.Contains("I now know 1 things about"))
		if err != nil {
			t.Errorf("expected first insert, got: %v", err)
		}

		// Second value
		_, err = h.CommandAndExpect("regtest_multi_key := second_value", h.Regex(regexp.MustCompile(`I now know 2 things about`)))
		if err != nil {
			t.Errorf("expected second insert with count 2, got: %v", err)
		}
	})

	t.Run("randomwoot", func(t *testing.T) {
		createFactoid(t, "randomwoot", "woot test")
		_, err := h.CommandAndExpect("regtest_woot_key := some value", h.Contains("woot test, I now know 1 things about"))
		if err != nil {
			t.Errorf("expected randomwoot, got: %v", err)
		}
	})
}

// ----- lookup handler -----

func testLookup(t *testing.T) {
	t.Run("privmsg_lookup", func(t *testing.T) {
		// Create a factoid
		key := "regtest_lookup_privmsg_key"
		_, err := h.CommandAndExpect(key+" := lookup_value", h.Contains("I now know 1 things about"))
		if err != nil {
			t.Fatalf("failed to create factoid: %v", err)
		}

		// Trigger lookup via PRIVMSG
		if !triggerLookupFast(t, key, "lookup_value", 10) {
			t.Errorf("expected lookup response on PRIVMSG")
		}
	})

	t.Run("action_lookup", func(t *testing.T) {
		// Create a factoid that responds with an action, trigger via an action
		key := "regtest_lookup_action_key"
		_, err := h.CommandAndExpect(key+" := <me>action_lookup_value", h.Contains("I now know 1 things about"))
		if err != nil {
			t.Fatalf("failed to create factoid: %v", err)
		}

		// Trigger lookup via ACTION
		w := h.ExpectEvent(client.ACTION, h.Contains("action_lookup_value"))
		h.Action(h.Channel, "regtest_lookup_action_key")
		if _, err := w.Wait(); err != nil {
			t.Errorf("expected lookup response on ACTION")
		}
	})

	t.Run("action_lookup_strips_bot_nick", func(t *testing.T) {
		// Create a factoid that responds with an action, trigger via an action
		key := "regtest_lookup_action_strip"
		_, err := h.CommandAndExpect(key+" := action_strip_value", h.Contains("I now know 1 things about"))
		if err != nil {
			t.Fatalf("failed to create factoid: %v", err)
		}

		// Trigger lookup via ACTION
		w := h.ExpectEvent(client.ACTION, h.Contains("action_lookup_value"))
		h.Action(h.Channel, "regtest_lookup_action_key "+h.BotNick)
		if _, err := w.Wait(); err != nil {
			t.Errorf("expected lookup response on ACTION with bot nick present")
		}
	})
}

func testLookupRecurse(t *testing.T) {
	t.Run("one_pointer", func(t *testing.T) {
		createFactoid(t, "recurse_value", "recursion successful")
		ptr := createFactoid(t, "ptr1", "*{recurse_value}")
		mustTriggerLookup(t, ptr, "recursion successful")
	})

	t.Run("many_pointers", func(t *testing.T) {
		createFactoid(t, "ptr2", "*{ptr1} 2")
		createFactoid(t, "ptr3", "*{ ptr2 } 3")
		ptr := createFactoid(t, "ptr4", "*ptr3 4")
		mustTriggerLookup(t, ptr, "recursion successful 2 3 4")
	})

	t.Run("loop", func(t *testing.T) {
		createFactoid(t, "ptr5", "pool *ptr6")
		loop := createFactoid(t, "ptr6", "loop *ptr5")
		mustTriggerLookup(t, loop, "loop pool [circular reference]")
	})
}

