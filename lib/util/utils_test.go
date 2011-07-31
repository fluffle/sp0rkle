package util

import (
	"testing"
)

// This is also implicitly testing HasPrefixedNick, I guess...
func TestHasPrefixedNick(t *testing.T) {
	tests := []string{
		"has no prefixed nick",
		"foo: has prefixed foo",
		"bar: has prefixed bar",
		"foo; has different prefix",
		"foo>has different prefix and no space",
		"foo-    has different prefix and lots of space",
		"foo! has wrong prefix char",
	}
	exp_str := []string{
		"has no prefixed nick",
		"has prefixed foo",
		"bar: has prefixed bar",
		"has different prefix",
		"has different prefix and no space",
		"has different prefix and lots of space",
		"foo! has wrong prefix char",
	}
	exp_bool := []bool{false, true, false, true, true, true, false}
	for i, s := range tests {
		r, b := RemovePrefixedNick(s, "foo")
		if r != exp_str[i] || b != exp_bool[i] {
			t.Errorf("Expected: %s, %t\nGot: %s, %t\n",
				exp_str[i], exp_bool[i], r, b)
		}
	}
}


func TestRemoveColours(t *testing.T) {
	tests := []string{
		"has no colours",
		"has \0035one colour",
		"has \00315one long colour\003 and a reset",
		"has \00312,4one background colour",
		"has \0036,12a different background colour\003 and a reset",
		"has \00312,14lots\00312 \0032,4of\003 colours",
		"has a colour \00313200 followed by a number",
		"has a background \0034,12300 followed by a number",
	}
	expected := []string{
		"has no colours",
		"has one colour",
		"has one long colour and a reset",
		"has one background colour",
		"has a different background colour and a reset",
		"has lots of colours",
		"has a colour 200 followed by a number",
		"has a background 300 followed by a number",
	}
	for i, s := range tests {
		ret := RemoveColours(s)
		if ret != expected[i] {
			t.Errorf("Expected: %s\nGot: %s\n", expected[i], ret)
		}
	}
}

func TestRemoveFormatting(t *testing.T) {
	tests := []string{
		"has no formatting",
		"has a \002BOLD\002 word",
		"has an \025underlined\025 word",
		"has \002lots\002 of \025\002formatted\002 words\002",
	}
	expected := []string{
		"has no formatting",
		"has a BOLD word",
		"has an underlined word",
		"has lots of formatted words",
	}
	for i, s := range tests {
		ret := RemoveFormatting(s)
		if ret != expected[i] {
			t.Errorf("Expected: %s\nGot: %s\n", expected[i], ret)
		}
	}
}
