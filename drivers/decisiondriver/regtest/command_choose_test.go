//go:build integration

package regtest

import (
	"regexp"
	"testing"
)

func testChoose(t *testing.T) {
	t.Run("space_delimited", func(t *testing.T) {
		_, err := h.CommandAndExpect("choose alpha beta gamma",
			h.Regex(regexp.MustCompile(`(alpha|beta|gamma)$`)))
		if err != nil {
			t.Errorf("expected one of the options, got: %v", err)
		}
	})

	t.Run("pipe_delimited", func(t *testing.T) {
		_, err := h.CommandAndExpect("choose x|y|z",
			h.Regex(regexp.MustCompile(`(x|y|z)$`)))
		if err != nil {
			t.Errorf("expected one of the piped options, got: %v", err)
		}
	})

	t.Run("double_quoted", func(t *testing.T) {
		_, err := h.CommandAndExpect(`choose "one two" "three four"`,
			h.Regex(regexp.MustCompile(`(one two|three four)$`)))
		if err != nil {
			t.Errorf("expected one of the quoted options, got: %v", err)
		}
	})

	t.Run("single_value", func(t *testing.T) {
		_, err := h.CommandAndExpect("choose onlyone", h.Contains("onlyone"))
		if err != nil {
			t.Errorf("expected 'onlyone', got: %v", err)
		}
	})

	t.Run("empty_input", func(t *testing.T) {
		line, err := h.CommandAndExpect("choose", h.Contains("need something to choose from"))
		if err != nil {
			t.Errorf("expected response for empty choose, got: %v", err)
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
		_, err := h.CommandAndExpect("CHOOSE foo bar",
			h.Regex(regexp.MustCompile(`(foo|bar)$`)))
		if err != nil {
			t.Errorf("expected case insensitive command, got: %v", err)
		}
	})
}
