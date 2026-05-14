package regtest

import (
	"regexp"
	"strings"

	"github.com/fluffle/goirc/client"
)

// Pattern matches an IRC line.
type Pattern interface {
	Match(line *client.Line) bool
}

// exactPattern matches the full command text.
type exactPattern struct {
	text string
}

func (e exactPattern) Match(line *client.Line) bool {
	return strings.TrimSpace(line.Text()) == e.text
}

// containsPattern matches if the command text contains the substring.
type containsPattern struct {
	text string
}

func (c containsPattern) Match(line *client.Line) bool {
	return strings.Contains(line.Text(), c.text)
}

// regexPattern matches via regexp.
type regexPattern struct {
	rx *regexp.Regexp
}

func (r regexPattern) Match(line *client.Line) bool {
	return r.rx.MatchString(line.Text())
}

// Exact returns a Pattern that matches the full command text.
func Exact(text string) Pattern {
	return exactPattern{text: text}
}

// Contains returns a Pattern that matches if the text contains the substring.
func Contains(text string) Pattern {
	return containsPattern{text: text}
}

// Regex returns a Pattern that matches via regexp.
func Regex(rx *regexp.Regexp) Pattern {
	return regexPattern{rx: rx}
}
