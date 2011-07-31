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

func BenchmarkRemoveColours(b *testing.B) {
	teststr := "has \00312,144\00312 sets \0032,4of\003 colours"
	for i := 0; i < b.N; i++ {
		RemoveColours(teststr)
	}
}

/*
func RemoveColoursRx(s string) string {
	rx := regexp.MustCompile("\003([0-9][0-9]?(,[0-9][0-9]?)?)?")
	s = rx.ReplaceAllString(s, "")
	return s
}

func BenchmarkRemoveColoursRx(b *testing.B) {
	teststr := "has \00312,144\00312 sets \0032,4of\003 colours"
	for i := 0; i < b.N; i++ {
		RemoveColoursRx(teststr)
	}
}
*/

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

func TestRemovePrefixes(t *testing.T) {
	tests := []string{
		"postfix",
		"oook, postfix",
		"see postfix",
		"uhhm postfix",
		"ummmmm postfix",
		"hey, postfix",
		"actually postfix",
		"oooooo postfix",
		"well, postfix",
		"iirc postfix",
		"but postfix",
		"and postfix",
		"or, postfix",
		"eh, postfix",
		".... postfix",
		"like, postfix",
		"ooohhhh postfix",
		"yeeaaaa postfix",
		"yeehhhh postfix",
		"yahhhhh postfix",
		"yup, postfix",
		"lol, postfix",
		"wow, postfix",
		"hhmmmm postfix",
		"eeerr, postfix",
		"hahaha postfix",
		"heh postfix",
		"hey, like actually iirc postfix",
		"... like haha and oooo, but iirc uhhm well, actually yup postfix",
	}
	for _, s := range tests {
		ret := RemovePrefixes(s)
		if ret != "postfix" {
			t.Errorf("Expected: postfix\nGot: %s\n", ret)
		}
	}
}
