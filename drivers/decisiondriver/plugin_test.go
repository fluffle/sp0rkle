package decisiondriver

import (
	"reflect"
	"testing"
)

func TestSplitDelimitedString(t *testing.T) {
	tests := []struct {
		in   string
		want []string
		err  error
	}{
		{"singlevalue", []string{"singlevalue"}, nil},
		{"AAA BBB CCC", []string{"AAA", "BBB", "CCC"}, nil},
		{"DDD | EEE", []string{"DDD ", " EEE"}, nil},
		{"FFF | GGG | HHH", []string{"FFF ", " GGG ", " HHH"}, nil},
		{"spam | spam and sausage | eggs | ham | spam eggs and spam",
			[]string{"spam ", " spam and sausage ", " eggs ", " ham ", " spam eggs and spam"}, nil},
		{`"spam" "spam and sausage" "eggs" "ham" "spam spam spam spam eggs and spam"`,
			[]string{"spam", "spam and sausage", "eggs", "ham", "spam spam spam spam eggs and spam"}, nil},
		{`mixed "quoting styles" 'just to' cause problems`,
			[]string{"mixed", "quoting styles", "just to", "cause", "problems"}, nil},
		{`'cheese' 'carrots' 'sausage'`, []string{"cheese", "carrots", "sausage"}, nil},
		{`"foo bar" "foo's bar" "something with spaces in it"`,
			[]string{"foo bar", "foo's bar", "something with spaces in it"}, nil},
		{`"foobar" "bar" "cheese"`, []string{"foobar", "bar", "cheese"}, nil},
		// Error condition.
		{`"cheese" "ham`, nil, ErrUnbalanced},
	}
	for _, s := range tests {
		got, err := splitDelimitedString(s.in)
		// We don't trim space from the possible choices *in* choices()
		// as it's only really necessary to do so for the one chosen.
		if err != s.err || !reflect.DeepEqual(s.want, got) {
			t.Errorf("splitDelimitedString(%s) = (%v, %v), want (%v, %v)\n", s.in, got, err, s.want, s.err)
		}
	}
}
