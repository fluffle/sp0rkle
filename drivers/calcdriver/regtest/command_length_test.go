//go:build integration

package regtest

import (
	"regexp"
	"testing"
)

func testLength(t *testing.T) {
	t.Run("length_simple", func(t *testing.T) {
		_, err := h.CommandAndExpect("length hello", h.Contains("'hello' is 5 characters long"))
		if err != nil {
			t.Errorf("expected 'hello is 5 characters', got: %v", err)
		}
	})

	t.Run("empty", func(t *testing.T) {
		_, err := h.CommandAndExpect("length ", h.Contains("'' is still longer"))
		if err != nil {
			t.Errorf("expected empty string length 0, got: %v", err)
		}
	})

	t.Run("length_with_spaces", func(t *testing.T) {
		_, err := h.CommandAndExpect("length hello world", h.Contains("'hello world' is 11 characters long"))
		if err != nil {
			t.Errorf("expected 'hello world is 11 characters', got: %v", err)
		}
	})

	t.Run("length_multibyte", func(t *testing.T) {
		_, err := h.CommandAndExpect("length café", h.Regex(regexp.MustCompile(`'café' is \d+ characters long`)))
		if err != nil {
			t.Errorf("expected multi-byte length, got: %v", err)
		}
	})

	t.Run("length_case_insensitive", func(t *testing.T) {
		_, err := h.CommandAndExpect("LENGTH test", h.Contains("'test' is 4 characters long"))
		if err != nil {
			t.Errorf("expected case insensitive command, got: %v", err)
		}
	})

	t.Run("length_single_char", func(t *testing.T) {
		_, err := h.CommandAndExpect("length x", h.Contains("'x' is 1 characters long"))
		if err != nil {
			t.Errorf("expected 'x is 1 characters', got: %v", err)
		}
	})

	t.Run("length_long_string", func(t *testing.T) {
		_, err := h.CommandAndExpect("length the quick brown fox jumps over the lazy dog", h.Contains("is 43 characters long"))
		if err != nil {
			t.Errorf("expected 43 characters, got: %v", err)
		}
	})
}
