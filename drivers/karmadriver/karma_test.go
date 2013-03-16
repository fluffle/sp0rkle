package karmadriver

import (
	"reflect"
	"testing"
)

func TestKarmaThings(t *testing.T) {
	tests := []struct {
		in  string
		out []kt
	}{
		{"no karma here", []kt{}},
		{"karma++", []kt{{"karma", true}}},
		{"karma--", []kt{{"karma", false}}},
		{"some funky karma-- text", []kt{{"karma", false}}},
		{"a++ b-- c++", []kt{{"a", true}, {"b", false}, {"c", true}}},
		{" (a b c)-- d++ ", []kt{{"a b c", false}, {"d", true}}},
		{"(a b c)++", []kt{{"a b c", true}}},
		{"a b c)++", []kt{}},
		{"++", []kt{}},
		{"foo -- bar", []kt{}},
		{"foo ()-- bar", []kt{}},
	}
	for i, test := range tests {
		o := karmaThings(test.in)
		if !reflect.DeepEqual(o, test.out) {
			t.Errorf("KarmaThings test %d: %s\nExpected: %#v\nGot: %#v\n",
				i, test.in, test.out, o)
		}
	}
}
