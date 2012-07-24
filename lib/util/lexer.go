package util

import (
	"strconv"
	"unicode"
	"unicode/utf8"
)

// A basic string lexer that simplifies extracting substrings.
type Lexer struct {
	Input             string
	start, pos, width int
}

// peek() returns the utf8 rune that is at lexer.pos in the input string.
// It does not move input.pos; repeated peek()s will return the same rune.
func (l *Lexer) Peek() (r rune) {
	if l.pos >= len(l.Input) {
		l.width = 0
		return 0
	}
	r, l.width = utf8.DecodeRuneInString(l.Input[l.pos:])
	if r == utf8.RuneError {
		// Treat bad unicode as EOF.
		l.width = 0
		return 0
	}
	return r
}

// next() returns the utf8 rune (in string form, for convenience elsewhere)
// that is at lexer.pos in the input string, then advances lexer.pos past it.
func (l *Lexer) Next() string {
	l.start = l.pos
	r := l.Peek()
	l.pos += l.width
	return string(r)
}

// scan() returns the sequence of runes in the input string anchored
// at lexer.pos that the supplied function returns true for. Usefully,
// unicode.IsDigit et al. fit the required function signature ;-)
func (l *Lexer) Scan(f func(rune) bool) string {
	l.start = l.pos
	for f(l.Peek()) {
		l.pos += l.width
	}
	return l.Input[l.start:l.pos]
}

// rewind() undoes the last next() or scan() by resetting lexer.pos.
func (l *Lexer) Rewind() {
	l.pos = l.start
}

// number() is a higher-level function that extracts a number from the
// input beginning at lexer.pos. A number matches the following regex:
//     -?[0-9]+(.[0-9]+)?([eE]-?[0-9]+)?
func (l *Lexer) Number() float64 {
	s := l.pos // l.start is reset through the multiple scans
	if l.Peek() == '-' {
		l.pos += l.width
	}
	l.Scan(unicode.IsDigit)
	if l.Next() == "." {
		l.Scan(unicode.IsDigit)
	} else {
		l.Rewind()
	}
	if c := l.Next(); c == "e" || c == "E" {
		if l.Peek() == '-' {
			l.pos += l.width
		}
		l.Scan(unicode.IsDigit)
	} else {
		l.Rewind()
	}
	l.start = s
	n, err := strconv.ParseFloat(l.Input[s:l.pos], 64)
	if err != nil {
		// This might be a bad idea in the long run.
		return 0
	}
	return n
}

