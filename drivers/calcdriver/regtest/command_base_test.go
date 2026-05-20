//go:build integration

package regtest

import (
	"regexp"
	"testing"
)

func testBase(t *testing.T) {
	t.Run("base_binary_to_decimal", func(t *testing.T) {
		_, err := h.CommandAndExpect("base 2to10 1010", h.Contains("1010 in base 2 is 10 in base 10"))
		if err != nil {
			t.Errorf("expected '1010 in base 2 is 10 in base 10', got: %v", err)
		}
	})

	t.Run("base_hex_to_decimal", func(t *testing.T) {
		_, err := h.CommandAndExpect("base 16to10 ff", h.Contains("ff in base 16 is 255 in base 10"))
		if err != nil {
			t.Errorf("expected 'ff in base 16 is 255 in base 10', got: %v", err)
		}
	})

	t.Run("base_decimal_to_binary", func(t *testing.T) {
		_, err := h.CommandAndExpect("base 10to2 42", h.Contains("42 in base 10 is 101010 in base 2"))
		if err != nil {
			t.Errorf("expected binary conversion, got: %v", err)
		}
	})

	t.Run("base_decimal_to_hex", func(t *testing.T) {
		_, err := h.CommandAndExpect("base 10to16 255", h.Contains("255 in base 10 is ff in base 16"))
		if err != nil {
			t.Errorf("expected hex conversion, got: %v", err)
		}
	})

	t.Run("base_base36", func(t *testing.T) {
		_, err := h.CommandAndExpect("base 10to36 35", h.Contains("35 in base 10 is z in base 36"))
		if err != nil {
			t.Errorf("expected base 36 conversion, got: %v", err)
		}
	})

	t.Run("base_octal_to_binary", func(t *testing.T) {
		_, err := h.CommandAndExpect("base 8to2 77", h.Contains("77 in base 8 is 111111 in base 2"))
		if err != nil {
			t.Errorf("expected octal to binary, got: %v", err)
		}
	})

	t.Run("base_invalid_format", func(t *testing.T) {
		_, err := h.CommandAndExpect("base invalid", h.Contains("Specify base as:"))
		if err != nil {
			t.Errorf("expected base format error, got: %v", err)
		}
	})

	t.Run("base_invalid_from", func(t *testing.T) {
		_, err := h.CommandAndExpect("base 1to10 5", h.Regex(regexp.MustCompile(`(?i)bad base`)))
		if err != nil {
			t.Errorf("expected invalid base error, got: %v", err)
		}
	})

	t.Run("base_invalid_to", func(t *testing.T) {
		_, err := h.CommandAndExpect("base 10to1 5", h.Regex(regexp.MustCompile(`(?i)bad base`)))
		if err != nil {
			t.Errorf("expected invalid base error, got: %v", err)
		}
	})

	t.Run("base_too_large_base", func(t *testing.T) {
		_, err := h.CommandAndExpect("base 10to37 5", h.Regex(regexp.MustCompile(`(?i)bad base`)))
		if err != nil {
			t.Errorf("expected base > 36 error, got: %v", err)
		}
	})

	t.Run("base_invalid_number", func(t *testing.T) {
		_, err := h.CommandAndExpect("base 10to16 zz", h.Regex(regexp.MustCompile(`(?i)Couldn't parse`)))
		if err != nil {
			t.Errorf("expected parse error for invalid number, got: %v", err)
		}
	})

	t.Run("base_case_insensitive", func(t *testing.T) {
		_, err := h.CommandAndExpect("BASE 10to2 1", h.Contains("1 in base 10 is 1 in base 2"))
		if err != nil {
			t.Errorf("expected case insensitive command, got: %v", err)
		}
	})

	t.Run("base_zero_to_any", func(t *testing.T) {
		_, err := h.CommandAndExpect("base 10to16 0", h.Contains("0 in base 10 is 0 in base 16"))
		if err != nil {
			t.Errorf("expected 0 conversion, got: %v", err)
		}
	})

	t.Run("base_negative_number", func(t *testing.T) {
		_, err := h.CommandAndExpect("base 10to2 -5", h.Contains("-5 in base 10 is -101 in base 2"))
		if err != nil {
			t.Errorf("expected negative number conversion, got: %v", err)
		}
	})
}
