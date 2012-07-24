package util

import (
	"strings"
	"testing"
	"unicode"
)
	
// This takes care of exercising peek, next, scan, and rewind
func TestLexerLowLevelFuncs(t *testing.T) {
	// Mmm. Unicodey. 42 bytes, 22 chars.
	l := &Lexer{Input: "This ïş æ 💩♽⛤ “Ŧəšţ”™."}

	// First, peek.
	if l.Peek() != 'T' {
		t.Errorf("Lexer appears to be starting in the wrong place")
	}
	// Advance 5 bytes to ï (\u00ef, 0xC3 0xAf)
	l.pos += 5
	if l.Peek() != 'ï' {
		t.Errorf("Lexer not decoding two-byte unicode chars")
	}
	// Advance another byte to the middle of ï
	l.pos += 1
	if l.Peek() != 0 {
		t.Errorf("Didn't get EOF from bad unicode")
	}

	// Advance to POO, PILE OF!
	l.pos = strings.Index(l.Input, "💩")
	if l.Peek() != '💩' {
		t.Errorf("Lexer can't decode shit")
	}

	// For the next three chars, make sure peek() and next() are in sync
	for i := 0; i < 3; i++ {
		if string(l.Peek()) != l.Next() {
			t.Errorf("Peek and next don't agree")
		}
	}

	// We should be at the space before “Ŧəšţ” now.
	if l.Next() != " " {
		t.Errorf("Lexer seems out of sync with reality")
	}
	l.Rewind()
	if l.Next() != " " {
		t.Errorf("Lexer still seems out of sync with reality")
	}
	l.Next() // skip opening quote.

	// Test scanning Ŧəšţ
	if l.Scan(unicode.IsLetter) != "Ŧəšţ" {
		t.Errorf("Scanning for letters didn't retrieve string")
	}
	l.Rewind()
	if l.Next() != "Ŧ" {
		t.Errorf("Rewinding scan didn't put lexer in correct place")
	}
}

func TestScanBadEOFHandling(t *testing.T) {
	l := &Lexer{Input: "alongstringwithnospaces"}
	s := l.Scan(func(r rune) bool {
		if r == ' ' { return false }
		return true
	})
	if s != "alongstringwithnospaces" {
		t.Errorf("Scan failed to handle func with no EOF detection")
	}
}	

func TestNumber(t *testing.T) {
	tests := []struct {
		i string  // input
		o float64 // output
		p int     // expected value of lexer.pos afterwards
	}{
		// GOOD CASES
		{"0", 0, 1},
		{"-1", -1, 2},
		{"1.25", 1.25, 4},
		{"-12345.6789", -12345.6789, 11},
		{"1e6", 1e6, 3},
		{"-1.23e45", -1.23e45, 8},
		{"1.23e-45", 1.23e-45, 8},
		// BAD CASES
		{"1e999", 0, 5},   // > MaxFloat
		{"1e-999", 0, 6},  // < MinFloat
		{"NaN", 0, 0},     // should result in ParseFloat("")
		{"a123.45", 0, 0}, //   ""
		// UGLY CASES
		{"0xf00", 0, 1}, // Hex not supported yet
		{"0b010", 0, 1}, // Binary not supported yet
		{"1foo", 1, 1},  // Stops at first non-digit
		{"०१२", 0, 9},   // 012 in devanagari digits
		{"໘໔໓", 0, 9},   // 843 in lao digits
		// I guess those poor devanagari etc. are SOL ;-(
	}

	for i, tc := range tests {
		l := &Lexer{Input: tc.i}
		if o := l.Number(); o != tc.o {
			t.Errorf("number(%d): '%s' result %f != %f", i, tc.i, o, tc.o)
		}
		if l.pos != tc.p {
			t.Errorf("number(%d): '%s' pos %d != %d", i, tc.i, l.pos, tc.p)
		}
	}
}
