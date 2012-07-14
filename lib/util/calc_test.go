package util

import (
	"math"
	"strings"
	"testing"
	"unicode"
)

func TestCalc(t *testing.T) {
	res, err := Calc("2+2")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if res != 4 {
		t.Errorf("2+2 is apparently %f", res)
	}
}

func makets(data ...float64) *tokenStack {
	ts := ts(len(data))
	for _, n := range data {
		ts.push(&token{T_NUM,"", n})
	}
	return ts
}

func TestFunctionisers(t *testing.T) {
	f1 := func(x float64) float64 { return 0.5*x }
	f2 := func(x,y float64) float64 { return 2*x/y }

	ff1 := functionise1(f1)
	ff2 := functionise2(f2)

	// Exercise f1 with a range of good inputs.
	for i, v := range []float64{2, 1e6, 6.75e-3, math.Pi} {
		ts := makets(v)
		if err := ff1.exec(ts); err != nil {
			t.Errorf("f1(%d): Unexpected error result from function.", i)
		} else if tok, err := ts.pop(); err != nil {
			t.Errorf("f1(%d): Stack not updated correctly with result", i)
		} else if exp := f1(v); tok.numval != exp {
			t.Errorf("f1(%d): Func result differed: %f != %f", i, tok.numval, exp)
		}
	}

	// Same for f2
	for i, v := range [][]float64{{2,4}, {1e6,2.9e4}, {6.75,48.1}, {math.Pi,math.E}} {
		ts := makets(v...)
		if err := ff2.exec(ts); err != nil {
			t.Errorf("f2(%d): Unexpected error result from function.", i)
		} else if tok, err := ts.pop(); err != nil {
			t.Errorf("f2(%d): Stack not updated correctly with result", i)
		} else if exp := f2(v[0],v[1]); tok.numval != exp {
			t.Errorf("f2(%d): Func result differed: %f != %f", i, tok.numval, exp)
		}
	}

	// Now try to break things
	zerots := ts(0)
	if err := ff1.exec(zerots); err == nil {
		t.Errorf("ff1 seemed happy with empty stack")
	}
	if err := ff2.exec(zerots); err == nil {
		t.Errorf("ff2 seemed happy with empty stack")
	}
	onets := makets(1)
	if err := ff2.exec(onets); err == nil {
		t.Errorf("ff2 seemed happy with undersized stack")
	}

	// Lastly, check they're not breaking the stack
	fourts := makets(1,2,3,4)
	ff1.exec(fourts) // should pop 4, push 2
	if len(*fourts) != 4 || (*fourts)[3].numval != 2 || (*fourts)[2].numval != 3 {
		t.Errorf("ff1 changed stack size unexpectedly")
	}
	ff2.exec(fourts) // should pop 2,3, push 3
	if len(*fourts) != 3 || (*fourts)[2].numval != 3 || (*fourts)[1].numval != 2 {
		t.Errorf("ff2 changed stack size unexpectedly")
	}
}

// Lexer tests

// This takes care of exercising peek, next, scan, and rewind
func TestLexerLowLevelFuncs(t *testing.T) {
	// Mmm. Unicodey. 42 bytes, 22 chars.
	l := &lexer{input: "This ïş æ 💩♽⛤ “Ŧəšţ”™."}

	// First, peek.
	if l.peek() != 'T' {
		t.Errorf("Lexer appears to be starting in the wrong place")
	}
	// Advance 5 bytes to ï (\u00ef, 0xC3 0xAf)
	l.pos += 5
	if l.peek() != 'ï' {
		t.Errorf("Lexer not decoding two-byte unicode chars")
	}
	// Advance another byte to the middle of ï
	l.pos += 1
	if l.peek() != 0 {
		t.Errorf("Didn't get EOF from bad unicode")
	}

	// Advance to POO, PILE OF!
	l.pos = strings.Index(l.input, "💩")
	if l.peek() != '💩' {
		t.Errorf("Lexer can't decode shit")
	}

	// For the next three chars, make sure peek() and next() are in sync
	for i := 0; i < 3; i++ {
		if string(l.peek()) != l.next() {
			t.Errorf("Peek and next don't agree")
		}
	}

	// We should be at the space before “Ŧəšţ” now.
	if l.next() != " " {
		t.Errorf("Lexer seems out of sync with reality")
	}
	l.rewind()
	if l.next() != " " {
		t.Errorf("Lexer still seems out of sync with reality")
	}
	l.next() // skip opening quote.

	// Test scanning Ŧəšţ
	if l.scan(unicode.IsLetter) != "Ŧəšţ" {
		t.Errorf("Scanning for letters didn't retrieve string")
	}
	l.rewind()
	if l.next() != "Ŧ" {
		t.Errorf("Rewinding scan didn't put lexer in correct place")
	}
}

func TestNumber(t *testing.T) {
	tests := []struct{
		i string   // input
		o float64  // output
		p int      // expected value of lexer.pos afterwards
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
		{"0xf00", 0, 1},   // Hex not supported yet
		{"0b010", 0, 1},   // Binary not supported yet
		{"1foo", 1, 1},    // Stops at first non-digit
		{"०१२", 0, 9},     // 012 in devanagari digits
		{"໘໔໓", 0, 9},     // 843 in lao digits
		// I guess those poor devanagari etc. are SOL ;-(
	}

	for i, tc := range tests {
		l := &lexer{input: tc.i}
		if o := l.number(); o != tc.o {
			t.Errorf("number(%d): '%s' result %f != %f", i, tc.i, o, tc.o)
		}
		if l.pos != tc.p {
			t.Errorf("number(%d): '%s' pos %d != %d", i, tc.i, l.pos, tc.p)
		}
	}
}

type tt struct {
	i string     // input
	k tokenKind  // token.kind
	s string     // token.strval
	n float64    // token.numval
}

func TestToken(t *testing.T) {
	tests := []tt{
		{"",        T_EOF,   "",      0},
		{"       ", T_EOF,   "",      0},
		{"+",       T_OP,    "+",     0},
// We can't test '-' as a standalone operator because the lexer
// assumes that it's the unary minus at the beginning of a number.
//		{"-",       T_OP,    "-",     0},
		{"*",       T_OP,    "*",     0},
		{"/",       T_OP,    "/",     0},
		{"**",      T_OP,    "**",    0},
		{"^",       T_OP,    "^",     0},
		{"%",       T_OP,    "%",     0},
		{"(",       T_LPAR,  "(",     0},
		{")",       T_RPAR,  ")",     0},
		{",",       T_COMMA, ",",     0},
		{"1234.5",  T_NUM,   "",      1234.5},
		{"-1234.5", T_NUM,   "",      -1234.5},
		{"pie",     T_NFI,   "pie",   0},
		{"-pie",    T_NFI,   "-pie",  0},
		{"&",       T_NFI,   "&",     0},
	}
	// Test all the fucntions are correctly recognised
	for fun, _ := range functionMap {
		tests = append(tests, tt{fun, T_FUNC, fun, 0})
	}
	// Test all the constants are correctly recognised
	for con, val := range ConstMap {
		tests = append(tests, tt{con, T_NUM, "", val})
		tests = append(tests, tt{"-" + con, T_NUM, "", -val})
	}

	for i, tc := range tests {
		l := &lexer{input: tc.i}
		tok := l.token()
		if tok.kind != tc.k {
			t.Errorf("token(%d) '%s' kind mismatch, %d != %d",
				i, tc.i, tok.kind, tc.k)
		}
		if tok.strval != tc.s {
			t.Errorf("token(%d) '%s' str mismatch, %s != %s",
				i, tc.i, tok.strval, tc.s)
		}
		if tok.numval != tc.n {
			t.Errorf("token(%d) '%s' num mismatch, %f != %f",
				i, tc.i, tok.numval, tc.n)
		}
	}
}

func TestTokenMinus(t *testing.T) {
	// Special minus testing
	tests := []struct{
		i string // input
		n int    // expect only the nth token to be a minus
	}{
		{"-e - -3", 2},
		{"atan2(3, -2) - 1", 7},
		{"cos(-pi-2)", 4},
	}
	for i, tc := range tests {
		l := &lexer{input: tc.i}
		tok := &token{T_NUM, "", 0} // start things off ...
		for j := 1; tok.kind != T_EOF; j++ {
			tok = l.token()
			if j != tc.n && tok.kind == T_OP && tok.strval == "-" {
				t.Errorf("token(%d) '%s' unexpected - operator at token %d",
					i, tc.i, j)
			}
			if j == tc.n && tok.kind != T_OP && tok.strval != "-" {
				t.Errorf("token(%d) '%s' unexpected %#v token",
					i, tc.i, tok)
			}
		}
	}
}

// Parser tests


