//go:build integration

package regtest

import (
	"testing"
)

func testOrd(t *testing.T) {
	t.Run("ord_uppercase", func(t *testing.T) {
		_, err := h.CommandAndExpect("ord A", h.Contains("ord(A) is 65, U+0041"))
		if err != nil {
			t.Errorf("expected ord(A) = 65, got: %v", err)
		}
	})

	t.Run("ord_lowercase", func(t *testing.T) {
		_, err := h.CommandAndExpect("ord z", h.Contains("ord(z) is 122, U+007A"))
		if err != nil {
			t.Errorf("expected ord(z) = 122, got: %v", err)
		}
	})

	t.Run("ord_digit", func(t *testing.T) {
		_, err := h.CommandAndExpect("ord 5", h.Contains("ord(5) is 53, U+0035"))
		if err != nil {
			t.Errorf("expected ord(5) = 53, got: %v", err)
		}
	})

	t.Run("ord_multibyte_utf8", func(t *testing.T) {
		_, err := h.CommandAndExpect("ord é", h.Contains("ord(é) is 233, U+00E9, '0xc3 0xa9'"))
		if err != nil {
			t.Errorf("expected ord(é) = 233, got: %v", err)
		}
	})

	t.Run("ord_percent", func(t *testing.T) {
		_, err := h.CommandAndExpect("ord %", h.Contains("ord(%) is 37, U+0025"))
		if err != nil {
			t.Errorf("expected ord(%%) = 37, got: %v", err)
		}
	})

	t.Run("ord_case_insensitive", func(t *testing.T) {
		_, err := h.CommandAndExpect("ORD A", h.Contains("ord(A) is 65, U+0041"))
		if err != nil {
			t.Errorf("expected case insensitive command, got: %v", err)
		}
	})

	t.Run("ord_multichar", func(t *testing.T) {
		_, err := h.CommandAndExpect("ord hello", h.Contains("ord(h) is 104, U+0068"))
		if err != nil {
			t.Errorf("expected ord(hello) = 104 (first rune), got: %v", err)
		}
	})

	t.Run("ord_emoji", func(t *testing.T) {
		_, err := h.CommandAndExpect("ord 😀", h.Contains("ord(😀) is 128512, U+1F600, '0xf0 0x9f 0x98 0x80'"))
		if err != nil {
			t.Errorf("expected ord(😀) = 128512, got: %v", err)
		}
	})

	t.Run("ord_bad_utf8", func(t *testing.T) {
		_, err := h.CommandAndExpect("ord \x91\x80\x80\x80", h.Contains("utf8 rune"))
		if err != nil {
			t.Errorf("expected bad rune, got: %v", err)
		}
	})
}
