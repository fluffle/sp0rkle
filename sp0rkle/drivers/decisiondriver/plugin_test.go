package decisiondriver

import (
	"rand"
	"testing"
	"lib/util"
)

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
			t.Errorf("Expected: %s\nGot: %s\n", expected[i], ret)
		}
	}
}

func TestDecide(t *testing.T) {
	tests := []string{
		"singlevalue",
		"sausage eggs spam",
		"carrots | ham",
		"spam | spam and sausage | eggs | ham | spam spam spam spam eggs and spam",
		"\"spam\" \"spam and sausage\" \"eggs\" \"ham\" \"spam spam spam spam eggs and spam\"",
		"\"cheese\" \"ham",
		"'cheese' 'carrots' 'sausage'",
	}
	expected := []string{
		"error: see help for details",
		"sausage",
		"carrots",
		"spam",
		"spam",
		"Error: unmatched quotes",
		"cheese",
	}
	for i, s := range tests {
		ret := rand_decider(s, mytestrand)
		if ret != expected[i] {
			t.Errorf("Expected: %s\nGot: %s\n\n",expected[i], ret)
		}
	}
}

