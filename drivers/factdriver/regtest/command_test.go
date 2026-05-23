//go:build integration

package regtest

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

// ----- chance -----

func testChance(t *testing.T) {
	key := createFactoid(t, "regtest_chance_key", "chance_value")

	t.Run("set_chance", func(t *testing.T) {
		// Trigger lookup to set lastSeen.
		mustTriggerLookup(t, key, "chance_value")

		// Set chance to 50% — response reports old/new chance
		_, err := h.CommandAndExpect(
			"chance of that is 50%",
			h.Regex(regexp.MustCompile(`was at \d+% chance, now is at 50%`)))
		if err != nil {
			t.Errorf("expected chance update response, got: %v", err)
		}
	})

	t.Run("modify_chance", func(t *testing.T) {
		// Float chance: 0.75 means 75%
		// Trigger lookup: chance of not getting a response at 50% in 20 rolls is
		// 1 in 1M (1/2^20), so hopefully this won't be _too_ flaky :D
		if !triggerLookupFast(t, key, "chance_value", 20) {
			t.Fatalf("could not look up %q in 20 attempts, go buy a lottery ticket", key)
		}
		_, err := h.CommandAndExpect(
			"chance of that is 0.75",
			h.Regex(regexp.MustCompile(`was at 50% chance, now is at 75%`)))
		if err != nil {
			t.Errorf("expected float chance update to 75 pct, got: %v", err)
		}
	})

	t.Run("chance_too_high", func(t *testing.T) {
		_, err := h.CommandAndExpect("chance of that is 150%", h.Contains("outside possible chance ranges"))
		if err != nil {
			t.Errorf("expected out-of-range error, got: %v", err)
		}
	})

	t.Run("chance_too_low", func(t *testing.T) {
		_, err := h.CommandAndExpect("chance of that is 0%", h.Contains("outside possible chance ranges"))
		if err != nil {
			t.Errorf("expected out-of-range error for zero, got: %v", err)
		}
	})

	t.Run("negative_chance", func(t *testing.T) {
		_, err := h.CommandAndExpect("chance of that is -5%", h.Contains("outside possible chance ranges"))
		if err != nil {
			t.Errorf("expected out-of-range error for negative, got: %v", err)
		}
	})

	t.Run("nan_chance", func(t *testing.T) {
		_, err := h.CommandAndExpect("chance of that is notanumber", h.Contains("didn't look like"))
		if err != nil {
			t.Errorf("expected parse error for non-numeric, got: %v", err)
		}
	})
}

// ----- edit (that =~) -----

func testEdit(t *testing.T) {
	key := createFactoid(t, "regtest_edit_key", "hello world")

	// NOTE: Ordering is important for these subtests.
	t.Run("basic_replace", func(t *testing.T) {
		mustTriggerLookup(t, key, "hello world")
		_, err := h.CommandAndExpect("that =~ s/hello/byebye/", h.Contains("byebye world"))
		if err != nil {
			t.Errorf("expected edit response with 'byebye world', got: %v", err)
		}
	})

	t.Run("alternate_delimiters", func(t *testing.T) {
		mustTriggerLookup(t, key, "byebye world")
		_, err := h.CommandAndExpect("that =~ s#world#http://some.host/#", h.Contains("byebye http://some.host/"))
		if err != nil {
			t.Errorf("expected edit response with 'byebye http://some.host/', got: %v", err)
		}
	})

	t.Run("escaped_delimiters", func(t *testing.T) {
		mustTriggerLookup(t, key, "byebye http://some.host/")
		_, err := h.CommandAndExpect(`that =~ s/http:\/\/some.host\//C:\\Program Files\\/`, h.Contains(`byebye C:\Program Files\`))
		if err != nil {
			t.Errorf(`expected edit response with 'byebye C:\Program Files\', got: %v`, err)
		}
	})

	t.Run("no_last_seen", func(t *testing.T) {
		_, err := h.CommandAndExpect("that =~ s#byebye#i forgor#", h.Contains("forgotten"))
		if err != nil {
			t.Errorf("expected edit failure with 'forgotten', got: %v", err)
		}
	})

	t.Run("invalid_regex", func(t *testing.T) {
		_, err := h.CommandAndExpect("that =~ s/[invalid/x/", h.Contains("Couldn't compile regex"))
		if err != nil {
			t.Errorf("expected regex compile error, got: %v", err)
		}
	})

	t.Run("wrong_sub_char", func(t *testing.T) {
		_, err := h.CommandAndExpect("that =~ g/hello/bye/", h.Contains("fool"))
		if err != nil {
			t.Errorf("expected syntax error for wrong delimiter, got: %v", err)
		}
	})

	t.Run("unbalanced_delimiters", func(t *testing.T) {
		_, err := h.CommandAndExpect(`that =~ s/C:\Program Files\/NO DATA/`, h.Contains(`parse regex`))
		if err != nil {
			t.Errorf(`expected edit response with 'byebye C:\Program Files\', got: %v`, err)
		}
	})
}

// ----- forget -----

func testForget(t *testing.T) {
	key := createFactoid(t, "regtest_forget_key", "forget_value")

	t.Run("deletes", func(t *testing.T) {
		mustTriggerLookup(t, key, "forget_value")
		// Forget the factoid
		_, err := h.CommandAndExpect("forget that", h.Contains("I forgot that"))
		if err != nil {
			t.Errorf("expected forget response, got: %v", err)
		}

		// Verify it's gone
		_, err = h.CommandAndExpect("fact info "+key, h.Contains("I don't know anything"))
		if err != nil {
			t.Errorf("expected factoid to be deleted, got: %v", err)
		}
	})

	t.Run("no_last_seen", func(t *testing.T) {
		_, err := h.CommandAndExpect("forget that", h.Contains("already forgotten"))
		if err != nil {
			t.Errorf("expected already forgotten error, got: %v", err)
		}
	})
}

// ----- delete -----

func testDelete(t *testing.T) {
	t.Run("no_last_seen", func(t *testing.T) {
		_, err := h.CommandAndExpect("delete that", h.Contains("already forgotten"))
		if err != nil {
			t.Errorf("expected already forgotten error, got: %v", err)
		}
	})

	key := createFactoid(t, "regtest_delete_key", "delete_value")

	t.Run("deletes", func(t *testing.T) {
	mustTriggerLookup(t, key, "delete_value")
		// Delete the factoid
		_, err := h.CommandAndExpect("delete that", h.Contains("I forgot that"))
		if err != nil {
			t.Errorf("expected delete response, got: %v", err)
		}

		// Verify it's gone
		_, err = h.CommandAndExpect("fact info "+key, h.Contains("I don't know anything"))
		if err != nil {
			t.Errorf("expected factoid to be deleted, got: %v", err)
		}
	})
}

// ----- info -----

func testInfo(t *testing.T) {
	key := createFactoid(t, "regtest_info_key", "info_value")

	// Info on existing key
	t.Run("pristine_key", func(t *testing.T) {
		line, err := h.CommandAndExpect("fact info "+key, h.Contains("1 things about"))
		if err != nil {
			t.Errorf("expected info response for existing key, got: %v", err)
		}

		// Validate some individual bits of info behaviour more carefully.
		text := line.Text()
		date := `\d\d:\d\d:\d\d, \w+ \d{1,2} \w+ \d{4} [A-Z]{3}`
		nick := h.Me().Nick
		wantrx := []*regexp.Regexp{
			regexp.MustCompile(`A factoid for 'regtest_info_key'`),
			regexp.MustCompile(`created on `+date+` by `+nick+`,`),
			regexp.MustCompile(`modified on `+date+` by `+nick+`,`),
			regexp.MustCompile(`accessed on `+date+` by `+nick+`.`),
			regexp.MustCompile(`modified 0 times and accessed 0 times.`),
		}
		for _, rx := range wantrx {
			if !rx.MatchString(text) {
				t.Errorf("info text %q did not match %q", text, rx)
			}
		}
	})

	// Access the factoid, edit the value, add a second value, and try again.
	t.Run("modified_key", func(t *testing.T) {
		mustTriggerLookup(t, key, "info_value")
		_, err := h.CommandAndExpect("that =~ s/_/ /", h.Contains("info value"))
		if err != nil {
			t.Errorf("expected edit response with 'info value', got: %v", err)
		}
		key = createFactoid(t, "regtest_info_key", "other value")

		line, err := h.CommandAndExpect("fact info "+key, h.Contains("2 things about"))
		if err != nil {
			t.Errorf("expected info response for existing key, got: %v", err)
		}

		// Validate modification and access counters increased.
		if !strings.Contains(line.Text(), `modified 1 times and accessed 1 times.`) {
			t.Errorf("info text %q did not show modification / access increments", line.Text())
		}
	})

	t.Run("missing_key", func(t *testing.T) {
		_, err := h.CommandAndExpect("fact info nonexistingskey12345", h.Contains("I don't know anything about"))
		if err != nil {
			t.Errorf("expected not found for non-existent key, got: %v", err)
		}
	})
}

// ----- literal -----

func testLiteral(t *testing.T) {
	key := createFactoid(t, "regtest_literal_key", "literal_value")

	t.Run("single_value", func(t *testing.T) {
		_, err := h.CommandAndExpect("literal "+key, h.Contains("[100%] literal_value"))
		if err != nil {
			t.Errorf("expected literal response with value, got: %v", err)
		}
	})

	t.Run("two_values", func(t *testing.T) {
		// Modify chance of value #1, add #2
		mustTriggerLookup(t, key, "literal_value")
		_, err := h.CommandAndExpect(
			"chance of that is 50%",
			h.Regex(regexp.MustCompile(`was at \d+% chance, now is at 50%`)))
		if err != nil {
			t.Errorf("expected chance update response, got: %v", err)
		}
		key = createFactoid(t, "regtest_literal_key", "oi $nick")

		// Literal on existing key — literal sends raw Privmsg, not ReplyN
		// so $nick should be repeated verbatim. Expect both lines, no ordering.
		w1 := h.Expect(h.Contains("[ 50%] literal_value"))
		w2 := h.Expect(h.Contains("[100%] oi $nick"))
		h.Command("literal "+key)
		_, err1 := w1.Wait()
		_, err2 := w2.Wait()
		if err1 != nil || err2 != nil {
			t.Errorf("expected literal response with value, got: %v / %v", err1, err2)
		}
	})

	t.Run("too_many_values", func(t *testing.T) {
		for i := range 10 {
			createFactoid(t, key, fmt.Sprintf("value %d", i))
		}
		// Big factoids have to be asked about privately to avoid channel spam.
		_, err := h.CommandAndExpect("literal "+key, h.Contains("know too much about"))
		if err != nil {
			t.Errorf("expected know too much response, got: %v", err)
		}
	})

	t.Run("missing_key", func(t *testing.T) {
		// Literal on non-existent key
		_, err := h.CommandAndExpect("literal nonexistingskey12345", h.Contains("I don't know anything about"))
		if err != nil {
			t.Errorf("expected not found for non-existent key, got: %v", err)
		}
	})

	// Delete last factoid to reset lastSeen state for other tests.
	_, err := h.CommandAndExpect("delete that", h.Contains("I forgot that"))
	if err != nil {
		t.Errorf("expected delete response, got: %v", err)
	}

}

// ----- replace -----

func testReplace(t *testing.T) {
	t.Run("no_last_seen", func(t *testing.T) {
		_, err := h.CommandAndExpect("replace that with a thing", h.Contains("already forgotten"))
		if err != nil {
			t.Errorf("expected already forgotten error, got: %v", err)
		}
	})

	key := createFactoid(t, "regtest_replace_key", "old_value")

	t.Run("no_last_seen", func(t *testing.T) {
		mustTriggerLookup(t, key, "old_value")
		_, err := h.CommandAndExpect("replace that with new_value", h.Contains("new_value"))
		if err != nil {
			t.Errorf("expected replace response with new value, got: %v", err)
		}

		// Verify the replacement via lookup
		mustTriggerLookup(t, key, "new_value")
	})
}

// ----- search -----

func testSearch(t *testing.T) {
	key := createFactoid(t, "regtest_search_key", "search_value")

	t.Run("existing_key", func(t *testing.T) {
		_, err := h.CommandAndExpect("fact search "+regexp.QuoteMeta(key), h.Contains(key))
		if err != nil {
			t.Errorf("expected search to find key, got: %v", err)
		}
	})

	t.Run("existing_key", func(t *testing.T) {
		_, err := h.CommandAndExpect("fact search xyznonexistentkey12345", h.Contains("couldn't think of anything"))
		if err != nil {
			t.Errorf("expected no results for non-existent search, got: %v", err)
		}
	})

	t.Run("regex_key", func(t *testing.T) {
		createFactoid(t, "regtest_search_key2", "another_value")
		_, err := h.CommandAndExpect(`fact search r[^_]+_s.arch_.*`, h.Contains("found 2 keys matching"))
		if err != nil {
			t.Errorf("expected search to find multiple keys, got: %v", err)
		}
	})
}
