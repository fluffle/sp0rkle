//go:build integration

package regtest

import (
	"regexp"
	"testing"
)

func testRand(t *testing.T) {
	t.Run("basic_0_to_10", func(t *testing.T) {
		_, err := h.CommandAndExpect("rand 10", h.Regex(regexp.MustCompile(`-?[0-9]+$`)))
		if err != nil {
			t.Errorf("expected integer response, got: %v", err)
		}
	})

	t.Run("range_5_to_10", func(t *testing.T) {
		_, err := h.CommandAndExpect("rand 5-10", h.Regex(regexp.MustCompile(`[5-9]$|^10$`)))
		if err != nil {
			t.Errorf("expected integer 5-10, got: %v", err)
		}
	})

	t.Run("format_two_decimals", func(t *testing.T) {
		_, err := h.CommandAndExpect("rand 100 %.2f", h.Regex(regexp.MustCompile(`\d+\.\d{2}$`)))
		if err != nil {
			t.Errorf("expected two decimal places, got: %v", err)
		}
	})

	t.Run("format_scientific", func(t *testing.T) {
		_, err := h.CommandAndExpect("rand 1000000 %e", h.Regex(regexp.MustCompile(`\d+e[+-]?\d+$`)))
		if err != nil {
			t.Errorf("expected scientific notation, got: %v", err)
		}
	})

	t.Run("negative_range", func(t *testing.T) {
		_, err := h.CommandAndExpect("rand -9-9", h.Regex(regexp.MustCompile(`-?[0-9]$`)))
		if err != nil {
			t.Errorf("expected integer in -9 to 9, got: %v", err)
		}
	})

	t.Run("float_range", func(t *testing.T) {
		_, err := h.CommandAndExpect("rand 0.5-1.5 %.2f", h.Regex(regexp.MustCompile(`[01]\.\d{2}$`)))
		if err != nil {
			t.Errorf("expected float 0.50-1.50, got: %v", err)
		}
	})

	t.Run("no_args", func(t *testing.T) {
		_, err := h.CommandAndExpect("rand", h.Contains("0"))
		if err != nil {
			t.Errorf("expected '0' for empty rand, got: %v", err)
		}
	})

	t.Run("invalid_input", func(t *testing.T) {
		_, err := h.CommandAndExpect("rand abc", h.Contains("0"))
		if err != nil {
			t.Errorf("expected '0' for invalid input, got: %v", err)
		}
	})

	t.Run("case_insensitive", func(t *testing.T) {
		_, err := h.CommandAndExpect("RAND 10", h.Regex(regexp.MustCompile(`-?[0-9]+$`)))
		if err != nil {
			t.Errorf("expected case insensitive command, got: %v", err)
		}
	})
}
