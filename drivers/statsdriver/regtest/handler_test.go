//go:build integration

package regtest

import (
	"regexp"
	"strconv"
	"testing"
	"time"
)

func testRecordPrivmsg(t *testing.T) {
	h.Privmsg(h.Channel, "recorded message")
	time.Sleep(100 * time.Millisecond)
	_, err := h.CommandAndExpect("lines "+h.Me().Nick, h.Contains("has said"))
	if err != nil {
		t.Errorf("expected stats after PRIVMSG, got: %v", err)
	}
}

func testRecordAction(t *testing.T) {
	h.Action(h.Channel, "wave hello")
	time.Sleep(100 * time.Millisecond)
	_, err := h.CommandAndExpect("lines "+h.Me().Nick, h.Contains("has said"))
	if err != nil {
		t.Errorf("expected stats after ACTION, got: %v", err)
	}
}

func testRecordMultiple(t *testing.T) {
	// Record current line count before sending new messages
	rx := regexp.MustCompile(`(\d+) lines`)
	line, err := h.CommandAndExpect("lines "+h.Me().Nick, h.Regex(rx))
	if err != nil {
		t.Fatalf("failed to get current stats: %v", err)
	}
	matches := rx.FindStringSubmatch(line.Text())
	if len(matches) < 2 {
		t.Fatalf("could not parse line count from response: %s", line.Text())
	}
	before, err := strconv.Atoi(matches[1])
	if err != nil {
		t.Fatalf("could not convert line count: %v", err)
	}

	h.Privmsg(h.Channel, "message one")
	time.Sleep(100 * time.Millisecond)
	h.Privmsg(h.Channel, "message two")
	time.Sleep(100 * time.Millisecond)
	h.Privmsg(h.Channel, "message three")
	time.Sleep(100 * time.Millisecond)

	line, err = h.CommandAndExpect("lines "+h.Me().Nick, h.Regex(rx))
	if err != nil {
		t.Fatalf("failed to get updated stats: %v", err)
	}
	matches = rx.FindStringSubmatch(line.Text())
	if len(matches) < 2 {
		t.Fatalf("could not parse line count from response: %s", line.Text())
	}
	after, err := strconv.Atoi(matches[1])
	if err != nil {
		t.Fatalf("could not convert line count: %v", err)
	}

	// The difference is 4: 3 PRIVMSGs + 1 PRIVMSG from the "lines" query itself
	if after-before != 4 {
		t.Errorf("expected 4 more lines (3 messages + 1 query, before=%d, after=%d), got %d", before, after, after-before)
	}
}
