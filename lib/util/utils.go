package util

// Random utility functions that are useful in various places.

import (
	"strings"
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

// Does this string look like a URL to you?
// This should be fairly conservative, I hope:
//   s starts with http:// or https:// and contains no spaces
func LooksURLish(s string) bool {
	return ((strings.HasPrefix(s, "http://") ||
		strings.HasPrefix(s, "https://")) &&
		strings.Index(s, " ") == -1)
}
