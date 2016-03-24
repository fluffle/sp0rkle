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
		{"a++ b-- c++", []kt{{"c", true}, {"b", false}, {"a", true}}},
		{" (a b c)-- d++ ", []kt{{"d", true}, {"a b c", false}}},
		{"(a b c)++", []kt{{"a b c", true}}},
		{"a b c)++", []kt{}},
		{"++", []kt{}},
		{"foo -- bar", []kt{}},
		{"foo ()-- bar", []kt{}},
		{"++.++++++", []kt{}},
		{"http://foo.bar/some-url--with-hyphens", []kt{}},
		// Brainfuck has poor karma :-)
		{"++++++++[>++++[>++>+++>+++>+<<<<-]>+>+>->>+[<]<-]>>.>---.+++++++..+++.>>.<-.<.+++.----", []kt{}},
		// ... as does Awk.
		{"BEGIN{n=0}{if(t[$2,$3]){o[n]=$0;n++};", []kt{}},
	}
	for i, test := range tests {
		o := karmaThings(test.in)
		if !reflect.DeepEqual(o, test.out) {
			t.Errorf("KarmaThings test %d: %s\nExpected: %#v\nGot: %#v\n",
				i, test.in, test.out, o)
		}
	}
}
