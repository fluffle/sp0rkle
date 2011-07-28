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
