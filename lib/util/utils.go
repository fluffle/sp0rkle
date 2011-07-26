package util

// Random utility functions that are useful in various places.

import (
	"strings"
)

func RemovePrefixedNick(text, nick string) string {
	if strings.HasPrefix(strings.ToLower(text), strings.ToLower(nick)) {
		l := len(nick)
		switch text[l] {
		// This is nicer than an if statement :-)
		// We only cut off the nick if it's followed by one of these chars
		// and an optional space, to indicate that it was prefixed to the text.
		case ':', ';', ',', '>', '-':
			l++
			text = strings.TrimSpace(text[l:])
		}
	}
	return text
}
