package main

import (
	"github.com/fluffle/goirc/client"
	"lib/util"
	"rand"
	"testing"
	"time"
)

func TestIdentifiers(t *testing.T) {
	tests := []string{
		"nothing to see here",
		"just a $nick",
		"lots of $nick $nick $nick",
		"$nick $chan $username $user $host $time $date",
	}
	// Static timestamp for great testing justice, no "local" time here kthx.
	ts := time.SecondsToUTC(1234567890)
	line := &client.Line{
		Nick: "tester", Ident: "tests", Host: "goirc.github.com",
		Src: "tester!tests@goirc.github.com", Cmd: "PRIVMSG",
		Raw: ":tester!tests@goirc.github.com PRIVMSG #test :I love testing.",
		Args: []string{"#test", "I love testing."}, Time: ts,
	}
	expected := []string{
		"nothing to see here",
		"just a tester",
		"lots of tester tester tester",
		"tester #test tests tests goirc.github.com 23:31:30 Fri Feb 13 23:31:30 2009",
	}
	for i, s := range tests {
		ret := id_replacer(s, line, ts)
		if ret != expected[i] {
			t.Errorf("Expected: %s\nGot: %s\n", expected[i], ret)
		}
	}
}


// deterministic randomizer
var mytestrand *rand.Rand = util.NewRand(42)

func TestRand(t *testing.T) {
	tests := []string{
		"no plugin value",
		"this has a <plugin=rand 100>% chance of working",
		"<plugin=rand 40-50> should be between 40 and 50",
		"there's a <plugin=rand 60-100 %.3f%%> chance of accurate random",
		"<plugin=rand dongs> shouldn't work, but <plugin=rand 20> should",
	}
	expected := []string{
		"no plugin value",
		"this has a 37% chance of working",
		"41 should be between 40 and 50",
		"there's a 84.164% chance of accurate random",
		"0 shouldn't work, but 1 should",
	}
	for i, s := range tests {
		ret := rand_replacer(s, mytestrand)
		if ret != expected[i] {
			t.Fail()
		}
	}
}
