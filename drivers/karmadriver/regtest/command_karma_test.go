//go:build integration

package regtest

import (
	"strings"
	"testing"
	"time"
)

func testKarma(t *testing.T) {
	t.Run("no_args", func(t *testing.T) {
		_, err := h.CommandAndExpect("karma", h.Contains("chameeleeoooooonnn"))
		if err != nil {
			t.Errorf("expected response for karma without args, got: %v", err)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		_, err := h.CommandAndExpect("karma nonexistentthing",
			h.Contains("No karma found for 'nonexistentthing'"))
		if err != nil {
			t.Errorf("expected 'not found' response, got: %v", err)
		}
	})

	t.Run("after_upvote", func(t *testing.T) {
		h.Privmsg(h.Channel, "something++")
		time.Sleep(100 * time.Millisecond)
		_, err := h.CommandAndExpect("karma something",
			h.Contains("'something' has a karma of 1 after 1 votes. Last upvoted by "+h.Me().Nick))
		if err != nil {
			t.Errorf("expected 'something' karma of 1 after 1 vote, got: %v", err)
		}
	})

	t.Run("case_insensitive", func(t *testing.T) {
		h.Privmsg(h.Channel, "UPPERCASE++")
		time.Sleep(100 * time.Millisecond)
		_, err := h.CommandAndExpect("karma UPPERCASE",
			h.Contains("'UPPERCASE' has a karma of 1 after 1 votes. Last upvoted by "+h.Me().Nick))
		if err != nil {
			t.Errorf("expected 'UPPERCASE' karma of 1 after 1 vote, got: %v", err)
		}
	})

	t.Run("with_multiple_votes", func(t *testing.T) {
		h.Privmsg(h.Channel, "multi++")
		time.Sleep(100 * time.Millisecond)
		h.Privmsg(h.Channel, "multi++")
		time.Sleep(100 * time.Millisecond)
		h.Privmsg(h.Channel, "multi--")
		time.Sleep(100 * time.Millisecond)
		line, err := h.CommandAndExpect("karma multi",
			h.Contains("'multi' has a karma of 1 after 3 votes."))
		if err != nil || line == nil {
			t.Fatalf("expected 'multi' karma of 1 after 3 votes with both up/downvoter, got: %v", err)
		}
		wants := []string{
			"Last upvoted by "+h.Me().Nick,
			"Last downvoted by "+h.Me().Nick,
		}
		for _, want := range wants {
			if !strings.Contains(line.Text(), want) {
				t.Errorf("expected 'multi' karma to contain %q", want)
				t.Logf("original line: %s", line.Text())
			}
		}
	})
}
