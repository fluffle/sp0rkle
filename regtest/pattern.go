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

// PatternFunc allows a bare func to implement Pattern
type PatternFunc func(line *client.Line) bool

// Match implements Pattern for a PatternFunc
func (pf PatternFunc) Match(line *client.Line) bool { return pf(line) }

// exactPattern matches the full command text.
type exactPattern struct {
	nick string
	text string
}

func (e exactPattern) Match(line *client.Line) bool {
	return line.Nick == e.nick && strings.TrimSpace(line.Text()) == e.text
}

// containsPattern matches if the command text contains the substring.
type containsPattern struct {
	nick string
	text string
}

func (c containsPattern) Match(line *client.Line) bool {
	return line.Nick == c.nick && strings.Contains(line.Text(), c.text)
}

// regexPattern matches via regexp.
type regexPattern struct {
	nick string
	rx *regexp.Regexp
}

func (r regexPattern) Match(line *client.Line) bool {
	return line.Nick == r.nick && r.rx.MatchString(line.Text())
}

type cmdPattern struct {
	who, what string
}

func (c cmdPattern) Match(line *client.Line) bool {
	return line.Nick == c.who && line.Args[0] == c.what
}

// Exact returns a Pattern that matches a full response from the bot.
func (h *Harness) Exact(text string) Pattern {
	return exactPattern{nick: h.BotNick, text: text}
}

// Contains returns a Pattern that matches if the bot response contains the substring.
func (h *Harness) Contains(text string) Pattern {
	return containsPattern{nick: h.BotNick, text: text}
}

// Regex returns a Pattern that matches via regexp.
func (h *Harness) Regex(rx *regexp.Regexp) Pattern {
	return regexPattern{nick: h.BotNick, rx: rx}
}
