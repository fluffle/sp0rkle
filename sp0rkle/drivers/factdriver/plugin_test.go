package factdriver

import (
	"github.com/fluffle/goirc/client"
	"sp0rkle/base"
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
	ts := time.Unix(1234567890, 0).UTC()
	line := &base.Line{
		Line: client.Line{
			Nick: "tester", Ident: "tests", Host: "goirc.github.com",
			Src: "tester!tests@goirc.github.com", Cmd: "PRIVMSG",
			Raw:  ":tester!tests@goirc.github.com PRIVMSG #test :I love testing.",
			Args: []string{"#test", "I love testing."}, Time: ts,
		},
		Addressed: false,
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
