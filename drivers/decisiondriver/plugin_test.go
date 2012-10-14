package decisiondriver

import (
	"github.com/fluffle/sp0rkle/util"
	"reflect"
	"testing"
)

func TestSplitDelimitedString(t *testing.T) {
	tests := []string{
		"singlevalue",
		"AAA BBB CCC",
		"DDD | EEE",
		"FFF | GGG | HHH",
		"spam | spam and sausage | eggs | ham | spam eggs and spam",
		`"spam" "spam and sausage" "eggs" "ham" "spam spam spam spam eggs and spam"`,
		`mixed "quoting styles" 'just to' cause problems`,
		`"cheese" "ham`,
		`'cheese' 'carrots' 'sausage'`,
		`"foo bar" "foo's bar" "something with spaces in it"`,
		`"foobar" "bar" "cheese"`,
	}
	expected := [][]string{
		[]string{"singlevalue"},
		[]string{"AAA", "BBB", "CCC"},
		[]string{"DDD ", " EEE"},
		[]string{"FFF ", " GGG ", " HHH"},
		[]string{"spam ", " spam and sausage ", " eggs ", " ham ", " spam eggs and spam"},
		[]string{"spam", "spam and sausage", "eggs", "ham", "spam spam spam spam eggs and spam"},
		[]string{"mixed", "quoting styles", "just to", "cause", "problems"},
		[]string{},
		[]string{"cheese", "carrots", "sausage"},
		[]string{"foo bar", "foo's bar", "something with spaces in it"},
		[]string{"foobar", "bar", "cheese"},
	}
	for i, s := range tests {
		ret := splitDelimitedString(s)
		// We don't trim space from the possible choices *in* choices()
		// as it's only really necessary to do so for the one chosen.
		if !reflect.DeepEqual(expected[i], ret) {
			t.Errorf("Test: %s\nExpected: %#v\nGot: %#v\n\n", s, expected[i], ret)
		}
	}
}
