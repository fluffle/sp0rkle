package bot

import (
	"strings"
	"testing"
)

func TestMatchCaseInsensitive(t *testing.T) {
	cs := newCommandSet()
	// newCommandSet already adds "help"
	cs.Add(&command{help: "remind help"}, "remind")
	cs.Add(&command{help: "remind list help"}, "remind list")

	tests := []struct {
		input          string
		expectedPrefix string
	}{
		{"remind me", "remind"},
		{"REMIND me", "remind"},
		{"REmInD me", "remind"},
		{"remind list", "remind list"},
		{"REMIND LIST", "remind list"},
		{"Remind List 1", "remind list"},
		{"help remind", "help"},
		{"HELP remind", "help"},
	}

	for _, test := range tests {
		r, ln := cs.match(test.input)
		if r == nil {
			t.Errorf("Input %q failed to match", test.input)
			continue
		}
		if ln != len(test.expectedPrefix) {
			t.Errorf("Input %q matched with wrong length. Expected %d, got %d", test.input, len(test.expectedPrefix), ln)
		}
		if r.Help() != cs.set[test.expectedPrefix].Help() {
			t.Errorf("Input %q matched wrong command. Expected help %q, got %q", test.input, cs.set[test.expectedPrefix].Help(), r.Help())
		}
	}
}

func TestPossibleCaseInsensitive(t *testing.T) {
	cs := newCommandSet()
	cs.Add(&command{help: "remind help"}, "remind")

	tests := []struct {
		input    string
		expected []string
	}{
		{"remind", []string{"remind"}},
		{"REMIND", []string{"remind"}},
		{"REm", []string{"remind"}},
	}

	for _, test := range tests {
		poss := cs.possible(test.input)
		found := false
		for _, p := range poss {
			if p == test.expected[0] {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Input %q failed to find %q in %v", test.input, test.expected[0], poss)
		}
	}
}

func TestCasingPreservation(t *testing.T) {
	// This tests the logic in Handle that cuts off the command
	cs := newCommandSet()
	cs.Add(&command{help: "remind help"}, "remind")

	input := "REmInD me Put the CAKE in the oven"
	expectedArgs1 := "me Put the CAKE in the oven"

	r, ln := cs.match(input)
	if r == nil {
		t.Fatalf("Failed to match")
	}

	// Simulate the logic in Handle:
	// ctx.Args[1] = strings.Join(strings.Fields(ctx.Args[1][ln:]), " ")
	result := strings.Join(strings.Fields(input[ln:]), " ")

	if result != expectedArgs1 {
		t.Errorf("Casing not preserved or unexpected result.\nExpected: %q\nGot:      %q", expectedArgs1, result)
	}
}
