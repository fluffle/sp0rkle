package util

// Random utility functions that are useful in various places.

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode"
)

func RemovePrefixedNick(text, nick string) (string, bool) {
	if HasPrefixedNick(text, nick) {
		text = strings.TrimSpace(text[len(nick)+1:])
		return text, true
	}
	return text, false
}

func HasPrefixedNick(text, nick string) bool {
	prefixed := false
	if len(text) <= len(nick) {
		return false
	}
	if strings.HasPrefix(strings.ToLower(text), strings.ToLower(nick)) {
		switch text[len(nick)] {
		// This is nicer than an if statement :-)
		// We only cut off the nick if it's followed by one of these chars
		// and an optional space, to indicate that it was prefixed to the text.
		case ':', ';', ',', '>', '-':
			prefixed = true
		}
	}
	return prefixed
}

// Removes mIRC-style colours from a string.
// These colours match the following BNF notation:
//   colour ::= idchar | idchar colnum | idchar colnum "," colnum
//   idchar ::= "\003"
//   colnum ::= digit | digit digit
//   digit  ::= "0" "1" "2" "3" "4" "5" "6" "7" "8" "9"
func RemoveColours(s string) string {
	for {
		i := strings.Index(s, "\003")
		if i == -1 {
			break
		}
		j := i + 1 // end of colour sequence
		c := -1    // comma position, if found
	L:
		for {
			// Who needs regex anyway.
			// util.BenchmarkRemoveColours    1000000  1936 ns/op
			// util.BenchmarkRemoveColoursRx    50000 41497 ns/op
			switch {
			case c != -1 && (j-c) > 2:
				break L
			case s[j] == ',':
				c = j
				j++
			case c == -1 && (j-i) > 2:
				break L
			case s[j] >= '0' && s[j] <= '9':
				j++
			default:
				break L
			}
		}
		s = s[:i] + s[j:]
	}
	return s
}

func RemoveFormatting(s string) string {
	return strings.Map(func(c rune) rune {
		switch c {
		case '\002', '\025':
			// \002 == bold, \025 == underline
			return -1
		}
		return c
	}, s)
}

var prefixes []string = []string{
	"o*k+", "see", "u(h+m*|m+)", "hey", "actually", "ooo+",
	"we+ll+", "iirc", "but", "and", "or", "eh", `\.+`,
	"like", "o+h+", "y(e+a+h*|e+h+|a+h+)", "yup", "lol",
	"wow", "h+m+", "e+r+", "[ha][ha]+", "[he][he][he]+",
}
var prefixrx *regexp.Regexp = regexp.MustCompile(
	"^((" + strings.Join(prefixes, "|") + "),? *)+ ")

func RemovePrefixes(s string) string {
	if idx := prefixrx.FindStringIndex(s); idx != nil {
		return s[idx[1]:]
	}
	return s
}

// Apply a set of strings tests to a source string s
func _any(f func(string, string) bool, s string, l []string) bool {
	for _, i := range l {
		if f(s, i) {
			return true
		}
	}
	return false
}

// Returns true if string begins with any of prefixes.
// NOTE: Does prefix comparisons against strings.ToLower(*s)!
func HasAnyPrefix(s string, prefixes []string) bool {
	return _any(strings.HasPrefix, strings.ToLower(s), prefixes)
}

// Returns true if string contains any of indexes.
// NOTE: Does index searches against strings.ToLower(*s)!
func ContainsAny(s string, indexes []string) bool {
	return _any(strings.Contains, strings.ToLower(s), indexes)
}

// Does this string look like a URL to you?
// This should be fairly conservative, I hope:
//   s starts with http:// or https:// and contains no spaces
func LooksURLish(s string) bool {
	return ((strings.HasPrefix(s, "http://") ||
		strings.HasPrefix(s, "https://")) &&
		strings.Index(s, " ") == -1)

}

func ApplyPluginFunction(val, plugin string, f func(string) string) string {
	plstart := fmt.Sprintf("<plugin=%s", plugin)
	for {
		// Work out the indices of the plugin start and end.
		ps := strings.Index(val, plstart)
		if ps == -1 {
			break
		}
		pe := strings.Index(val[ps:], ">")
		if pe == -1 {
			// No closing '>', so abort
			break
		}
		pe += ps
		// Mid is where the plugin args start.
		mid := ps + len(plstart)
		// And if there *are* args we should skip the leading space
		for val[mid] == ' ' { mid++ }
		val = val[:ps] + f(val[mid:pe]) + val[pe+1:]
	}
	return val
}

func FactPointer(val string) (key string, start, end int) {
	// A pointer looks like *key or *{key with optional spaces}
	// In the former case key must be alphanumeric
	if start = strings.Index(val, "*"); start == -1 || start + 1 == len(val) {
		return "", -1, -1
	}
	if val[start+1] == '{' {
		end = strings.Index(val[start:], "}") + start + 1
		// TrimSpace since it's not possible to have a fact key that
		// starts/ends with a space, but someone *could* write *{ foo }
		key = strings.TrimSpace(val[start+2:end-1])
	} else {
		// util.Lexer helps find the next char that isn't alphabetical
		l := &Lexer{Input: val}
		l.Pos(start+1)
		key = l.Scan(func (r rune) bool {
			if unicode.IsLetter(r) || unicode.IsNumber(r) {
				return true
			}
			return false
		})
		end = l.Pos()
		// Special case handling because *pointer might be *emphasis*
		// perlfu's designer has a lot to answer for :-/
		if l.Peek() == '*' {
			return "", -1, -1
		}
	}
	return
}

func JoinPath(items ...string) string {
	return strings.Join(items, string(os.PathSeparator))
}

func TimeSince(t time.Time) string {
	s := ""
	sec := int(time.Since(t)/time.Second)
	times := []struct{
		d int
		s string
	}{
		{31536000, "y"}, // a year is 365 days, natch.
		{604800, "w"},
		{86400, "d"},
		{3600, "h"},
		{60, "m"},
		{1, "s"},
	}
	for _, v := range times {
		if div := sec / v.d; div > 0 {
			s = fmt.Sprintf("%s%d%s ", s, div, v.s)
		}
		sec = sec % v.d
	}
	if len(s) > 0 {
		return s[:len(s)-1]
	}
	return ""
}
