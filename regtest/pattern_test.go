package regtest

import (
	"regexp"
	"testing"

	"github.com/fluffle/goirc/client"
)

func TestExactPattern(t *testing.T) {
	h := &Harness{BotNick: "test"}
	p := h.Exact("hello world")
	line := client.ParseLine(":test!u@h PRIVMSG #chan :hello world")
	if !p.Match(line) {
		t.Error("Exact should match matching line")
	}
	line2 := client.ParseLine(":test!u@h PRIVMSG #chan :hello world!")
	if p.Match(line2) {
		t.Error("Exact should not match non-matching line")
	}
	line3 := client.ParseLine(":fail!u@h PRIVMSG #chan :hello world")
	if p.Match(line3) {
		t.Error("Exact should not match matching line from wrong nick")
	}
}

func TestContainsPattern(t *testing.T) {
	h := &Harness{BotNick: "test"}
	p := h.Contains("hello")
	line := client.ParseLine(":test!u@h PRIVMSG #chan :hello world")
	if !p.Match(line) {
		t.Error("Contains should match line with substring")
	}
	line2 := client.ParseLine(":test!u@h PRIVMSG #chan :goodbye world")
	if p.Match(line2) {
		t.Error("Contains should not match line without substring")
	}
	line3 := client.ParseLine(":fail!u@h PRIVMSG #chan :hello world")
	if p.Match(line3) {
		t.Error("Contains should not match matching line from wrong nick")
	}
}

func TestRegexPattern(t *testing.T) {
	h := &Harness{BotNick: "test"}
	p := h.Regex(regexp.MustCompile(`^hello.*world$`))
	line := client.ParseLine(":test!u@h PRIVMSG #chan :hello irc world")
	if !p.Match(line) {
		t.Error("Regex should match matching line")
	}
	line2 := client.ParseLine(":test!u@h PRIVMSG #chan :hello world irc")
	if p.Match(line2) {
		t.Error("Regex should not match non-matching line")
	}
	line3 := client.ParseLine(":fail!u@h PRIVMSG #chan :hello irc world")
	if p.Match(line3) {
		t.Error("Regex should not match matching line from wrong nick")
	}
}
