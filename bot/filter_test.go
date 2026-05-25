package bot

import (
	"strings"
	"testing"

	"github.com/fluffle/goirc/client"
	"github.com/fluffle/sp0rkle/db/conf"
)

func TestFilterPipeline(t *testing.T) {
	p := &FilterPipeline{}
	p.AddFunc(
		func (line *client.Line) bool {
			return strings.Contains(line.Raw, "PRIVMSG")
		},
		func (line *client.Line) bool {
			return strings.Contains(line.Raw, "hello")
		},
	)

	tests := []struct {
		line string
		want bool
	} {
		{"", false},
		{"neither", false},
		{"PRIVMSG :only one", false},
		{"OTHER :hello", false},
		{"PRIVMSG #channel :hello", true},
	}

	for _, test := range tests {
		got := p.ShouldProcess(&client.Line{Raw: test.line})
		if got != test.want {
			t.Errorf("ShouldProcess(%q) = %t, want %t", test.line, got, test.want)
		}
	}
}

func TestNickIgnoreFilter(t *testing.T) {
	ns := conf.InMem("ignore")
	// Ignore the nick "ignored" for testing purposes.
	ns.String("ignored", "ignored")
	p := &FilterPipeline{}
	p.Add(&nickIgnoreFilter{ns: ns})

	tests := []struct {
		name      string
		nick      string
		want      bool
	}{
		{
			name:     "permit message from normal nick",
			nick:     "somedude",
			want:     true,
		},
		{
			name:     "ignore message from ignored nick",
			nick:     "ignored",
			want:     false,
		},
		{
			name:     "validate case sensitivity",
			nick:     "iGNoReD",
			want:     false,
		},
		{
			name:     "permit message empty nick",
			nick:     "",
			want:     true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			line := &client.Line{Nick: test.nick}
			got := p.ShouldProcess(line)
			if got != test.want {
				t.Errorf("NamespaceFilter(%q) = %t, want %t", test.nick, got, test.want)
			}
		})
	}
}
