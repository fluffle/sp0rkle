package factdriver

import (
	"flag"
	"testing"
	"time"

	"github.com/fluffle/goirc/client"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/util"
)

func TestIdentifiers(t *testing.T) {
	tests := []string{
		"nothing to see here",
		"just a $nick",
		"lots of $nick $nick $nick",
		"$nick $chan $username $user $host $time $date",
	}
	// Static timestamp for great testing justice, no "local" time here kthx.
	flag.Set("timezone", "UTC")
	ts := time.Unix(1234567890, 0).UTC()
	ctx := &bot.Context{
		Line: &client.Line{
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
		"tester #test tests tests goirc.github.com 23:31:30 23:31:30, Friday 13 February 2009 UTC",
	}
	for i, s := range tests {
		ret := id_replacer(s, ctx, ts)
		if ret != expected[i] {
			t.Errorf("Expected: %s\nGot: %s\n", expected[i], ret)
		}
	}
}

func TestExtractRx(t *testing.T) {
	tests := []struct {
		in, out string
	}{
		{"", ""},
		{"/foo", ""},
		{"foo", "foo"},
		{"foo/bar", "foo"},
		{"foo\\/bar", "foo\\/bar"},
		{"foo\\/bar/", "foo\\/bar"},
		{"foo\\\\/bar", "foo\\\\"},
		{"foo\\\\\\/bar", "foo\\\\\\/bar"},
		{"foo\\", "foo\\"},
		{"foo\\\\", "foo\\\\"},
		{"\\/", "\\/"},
		{"\\//", "\\/"},
		{"\\\\/", "\\\\"},
		{"\\\\//", "\\\\"},
		{"\\/foo", "\\/foo"},
	}
	for i, test := range tests {
		l := &util.Lexer{Input: test.in}
		if o := extractRx(l, '/'); o != test.out {
			t.Errorf("extractRx(%d) '%s': exp '%s' got '%s'",
				i, test.in, test.out, o)
		}
	}
}
