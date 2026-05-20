//go:build integration

package regtest

import (
	"regexp"
	"testing"
)

func testChr(t *testing.T) {
	t.Run("chr_decimal", func(t *testing.T) {
		_, err := h.CommandAndExpect("chr 65", h.Contains("chr(65) is A, U+0041"))
		if err != nil {
			t.Errorf("expected chr(65) = A, got: %v", err)
		}
	})

	t.Run("chr_lowercase_a", func(t *testing.T) {
		_, err := h.CommandAndExpect("chr 97", h.Contains("chr(97) is a, U+0061"))
		if err != nil {
			t.Errorf("expected chr(97) = a, got: %v", err)
		}
	})

	t.Run("chr_hex", func(t *testing.T) {
		_, err := h.CommandAndExpect("chr 0x41", h.Contains("chr(0x41) is A"))
		if err != nil {
			t.Errorf("expected chr(0x41) = A, got: %v", err)
		}
	})

	t.Run("chr_octal", func(t *testing.T) {
		_, err := h.CommandAndExpect("chr 0o101", h.Contains("chr(0o101) is A"))
		if err != nil {
			t.Errorf("expected chr(0o101) = A, got: %v", err)
		}
	})

	t.Run("chr_unicode_syntax", func(t *testing.T) {
		// "u+0041" is translated to "0x0041" by the driver, output shows the transformed form
		_, err := h.CommandAndExpect("chr u+0041", h.Contains("chr(0x0041) is A"))
		if err != nil {
			t.Errorf("expected chr(0x0041) = A, got: %v", err)
		}
	})

	t.Run("chr_space", func(t *testing.T) {
		_, err := h.CommandAndExpect("chr 32", h.Contains("chr(32) is  ,"))
		if err != nil {
			t.Errorf("expected chr(32) = space, got: %v", err)
		}
	})

	t.Run("chr_null", func(t *testing.T) {
		_, err := h.CommandAndExpect("chr 0", h.Contains(`chr(0) is "\x00"`))
		if err != nil {
			t.Errorf("expected chr(10) = newline, got: %v", err)
		}
	})

	t.Run("chr_newline", func(t *testing.T) {
		_, err := h.CommandAndExpect("chr 10", h.Contains(`chr(10) is "\n"`))
		if err != nil {
			t.Errorf("expected chr(10) = newline, got: %v", err)
		}
	})

	t.Run("chr_multibyte_utf8", func(t *testing.T) {
		_, err := h.CommandAndExpect("chr 233", h.Contains("chr(233) is é, U+00E9, '0xc3 0xa9'"))
		if err != nil {
			t.Errorf("expected chr(233) = é, got: %v", err)
		}
	})

	t.Run("chr_emoji", func(t *testing.T) {
		_, err := h.CommandAndExpect("chr 0x1F600", h.Contains("chr(0x1f600) is 😀, U+1F600, '0xf0 0x9f 0x98 0x80'"))
		if err != nil {
			t.Errorf("expected chr(0x1F600) emoji, got: %v", err)
		}
	})

	t.Run("chr_invalid", func(t *testing.T) {
		_, err := h.CommandAndExpect("chr abc", h.Contains("Couldn't parse"))
		if err != nil {
			t.Errorf("expected parse error for invalid input, got: %v", err)
		}
	})

	t.Run("chr_empty", func(t *testing.T) {
		_, err := h.CommandAndExpect("chr", h.Contains("Couldn't parse"))
		if err != nil {
			t.Errorf("expected parse error for empty input, got: %v", err)
		}
	})

	t.Run("chr_case_insensitive", func(t *testing.T) {
		_, err := h.CommandAndExpect("CHR 65", h.Contains("chr(65) is A"))
		if err != nil {
			t.Errorf("expected case insensitive command, got: %v", err)
		}
	})

	t.Run("chr_negative", func(t *testing.T) {
		_, err := h.CommandAndExpect("chr -1", h.Contains("Don't be so negative."))
		if err != nil {
			t.Errorf("expected chr(-1) output, got: %v", err)
		}
	})

	t.Run("chr_large_value", func(t *testing.T) {
		_, err := h.CommandAndExpect("chr 0x10FFFF", h.Regex(regexp.MustCompile(`chr\(0x10ffff\)`)))
		if err != nil {
			t.Errorf("expected chr(0x10FFFF) output, got: %v", err)
		}
	})
}
