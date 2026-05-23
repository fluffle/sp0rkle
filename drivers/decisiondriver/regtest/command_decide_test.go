//go:build integration

package regtest

import (
	"regexp"
	"testing"
)

func testDecide(t *testing.T) {
	t.Run("space_delimited", func(t *testing.T) {
		_, err := h.CommandAndExpect("decide apple banana cherry",
			h.Regex(regexp.MustCompile(`(apple|banana|cherry)$`)))
		if err != nil {
			t.Errorf("expected one of the options, got: %v", err)
		}
	})

	t.Run("pipe_delimited", func(t *testing.T) {
		_, err := h.CommandAndExpect("decide red|green|blue",
			h.Regex(regexp.MustCompile(`(red|green|blue)$`)))
		if err != nil {
			t.Errorf("expected one of the piped options, got: %v", err)
		}
	})

	t.Run("double_quoted", func(t *testing.T) {
		_, err := h.CommandAndExpect(`decide "foo bar" "baz qux" "hello world"`,
			h.Regex(regexp.MustCompile(`(foo bar|baz qux|hello world)$`)))
		if err != nil {
			t.Errorf("expected one of the quoted options, got: %v", err)
		}
	})

	t.Run("single_quoted", func(t *testing.T) {
		_, err := h.CommandAndExpect(`decide 'foo bar' 'baz qux' 'hello world'`,
			h.Regex(regexp.MustCompile(`(foo bar|baz qux|hello world)$`)))
		if err != nil {
			t.Errorf("expected one of the single-quoted options, got: %v", err)
		}
	})

	t.Run("single_value", func(t *testing.T) {
		_, err := h.CommandAndExpect("decide onlychoice", h.Contains("onlychoice"))
		if err != nil {
			t.Errorf("expected 'onlychoice', got: %v", err)
		}
	})

	t.Run("empty_input", func(t *testing.T) {
		line, err := h.CommandAndExpect("choose", h.Contains("need something to choose from"))
		if err != nil {
			t.Errorf("expected response for empty decide, got: %v", err)
		}
		if line == nil {
			t.Error("expected response line")
		}
	})

	t.Run("unbalanced_quotes", func(t *testing.T) {
		_, err := h.CommandAndExpect(`choose "cheese" "ham`, h.Contains("can't decide"))
		if err != nil {
			t.Errorf("expected error for unbalanced quotes, got: %v", err)
		}
	})

	t.Run("case_insensitive", func(t *testing.T) {
		_, err := h.CommandAndExpect("DECIDE foo bar",
			h.Regex(regexp.MustCompile(`(foo|bar)$`)))
		if err != nil {
			t.Errorf("expected case insensitive command, got: %v", err)
		}
	})
}
