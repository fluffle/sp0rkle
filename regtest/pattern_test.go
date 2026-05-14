package regtest

import (
	"regexp"
	"testing"

	"github.com/fluffle/goirc/client"
)

func TestExactPattern(t *testing.T) {
	p := Exact("hello world")
	line := client.ParseLine("PRIVMSG #chan :hello world")
	if !p.Match(line) {
		t.Error("Exact should match matching line")
	}
	line2 := client.ParseLine("PRIVMSG #chan :hello world!")
	if p.Match(line2) {
		t.Error("Exact should not match non-matching line")
	}
}

func TestContainsPattern(t *testing.T) {
	p := Contains("hello")
	line := client.ParseLine("PRIVMSG #chan :hello world")
	if !p.Match(line) {
		t.Error("Contains should match line with substring")
	}
	line2 := client.ParseLine("PRIVMSG #chan :goodbye world")
	if p.Match(line2) {
		t.Error("Contains should not match line without substring")
	}
}

func TestRegexPattern(t *testing.T) {
	p := Regex(regexp.MustCompile(`^hello.*world$`))
	line := client.ParseLine("PRIVMSG #chan :hello world")
	if !p.Match(line) {
		t.Error("Regex should match matching line")
	}
	line2 := client.ParseLine("PRIVMSG #chan :hello")
	if p.Match(line2) {
		t.Error("Regex should not match non-matching line")
	}
}
