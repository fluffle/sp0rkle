//go:build integration

package regtest

import (
	"fmt"
	"testing"
	"time"
)

// -- remind list

func testRemindList(t *testing.T) {
	// NOTE: ordering of these tests matters!
	t.Run("empty", testRemindListEmpty)
	t.Run("one", testRemindListOne)
	t.Run("too_many", testRemindListTooMany)
}

func testRemindListEmpty(t *testing.T) {
	_, err := h.CommandAndExpect("remind list", h.Contains("no reminders set"))
	if err != nil {
		t.Errorf("expected no reminders message, got: %v", err)
	}
}

func testRemindListOne(t *testing.T) {
	mustRemind(t, "test remind list 1")
	// Also expect the actual reminder list entry.
	w := h.Expect(h.Contains("1: you asked me to remind you test remind list 1"))
	_, err := h.CommandAndExpect("remind list", h.Contains("1 reminders set"))
	if err != nil {
		t.Errorf("expected list header response, got: %v", err)
	}
	if _, err = w.Wait(); err != nil {
		t.Errorf("expected list entry response, got: %v", err)
	}
}

func testRemindListTooMany(t *testing.T) {
	// Create 5 more reminders to push the total over 5,
	// triggering the public channel limit.
	for i := range 5 {
		mustRemind(t, fmt.Sprintf("test remind list %d", i+2))
		// Timestamp index entries are at millisecond precision.
		// If we create 2 reminders in the same millisecond, which
		// appears to be _actually possible_ on localhost loopback,
		// we lose the index for one as the other overwrites it!
		time.Sleep(10*time.Millisecond)
	}
	// Now listing publicly should say "ask me privately"
	_, err := h.CommandAndExpect("remind list", h.Contains("lots of reminders"))
	if err != nil {
		t.Errorf("expected 'lots of reminders' message, got: %v", err)
	}

	// Creating so many reminders causes problems for other tests, soo:
	mustClearReminders(t)
}

// -- remind set

func testRemindSet(t *testing.T) {
	t.Run("self", testRemindSetSelf)
	t.Run("other", testRemindSetOther)
	t.Run("missing_args", testRemindSetMissingArgs)
	t.Run("no_time", testRemindSetNoTime)
	t.Run("past_time", testRemindSetPastTime)
	t.Run("relative", testRemindSetRelative)
}

func testRemindSetSelf(t *testing.T) {
	_, err := h.CommandAndExpect("remind me pick up laundry in 1 hour", h.Contains("okay, i'll remind you"))
	if err != nil {
		t.Errorf("expected self-reminder ack, got: %v", err)
	}
}

func testRemindSetOther(t *testing.T) {
	_, err := h.CommandAndExpect("remind otheruser call mom at 5pm", h.Contains("okay, i'll remind otheruser"))
	if err != nil {
		t.Errorf("expected other-reminder ack, got: %v", err)
	}
}

func testRemindSetMissingArgs(t *testing.T) {
	_, err := h.CommandAndExpect("remind", h.Contains("forget something?"))
	if err != nil {
		t.Errorf("expected missing args response, got: %v", err)
	}
}

func testRemindSetNoTime(t *testing.T) {
	_, err := h.CommandAndExpect("remind me hello", h.Contains("You asked me to remind you hello."))
	if err != nil {
		t.Errorf("expected no-time response, got: %v", err)
	}
}

func testRemindSetPastTime(t *testing.T) {
	_, err := h.CommandAndExpect("remind me read that book at 1970-01-01", h.Contains("in the past"))
	if err != nil {
		t.Errorf("expected past time error, got: %v", err)
	}
}

func testRemindSetRelative(t *testing.T) {
	_, err := h.CommandAndExpect("remind me water plants in 30 minutes", h.Contains("okay, i'll remind you"))
	if err != nil {
		t.Errorf("expected relative time ack, got: %v", err)
	}
}

// -- remind del

func testRemindDel(t *testing.T) {
	// NOTE: More ordering concerns here:
	//   - reminders all cleared after list tests
	//   - 3 reminders added in set tests
	//   - 1 reminder deleted in deletes test
	t.Run("no_list_first", testRemindDelNoListFirst)
	t.Run("deletes", testRemindDelDeletes)
	t.Run("invalid_index", testRemindDelInvalidIndex)
}

func testRemindDelNoListFirst(t *testing.T) {
	_, err := h.CommandAndExpect("remind del 1", h.Contains("use 'remind list' first"))
	if err != nil {
		t.Errorf("expected 'use remind list first' response, got: %v", err)
	}
}

func testRemindDelDeletes(t *testing.T) {
	_, err := h.CommandAndExpect("remind list", h.Contains("You have 3 reminders"))
	if err != nil {
		t.Errorf("expected list response, got: %v", err)
	}
	_, err = h.CommandAndExpect("remind del 1", h.Contains("forget that one"))
	if err != nil {
		t.Errorf("expected del response, got: %v", err)
	}
}

func testRemindDelInvalidIndex(t *testing.T) {
	_, err := h.CommandAndExpect("remind list", h.Contains("You have 2 reminders"))
	if err != nil {
		t.Errorf("expected list response, got: %v", err)
	}
	_, err = h.CommandAndExpect("remind del abc", h.Contains("Invalid reminder index"))
	if err != nil {
		t.Errorf("expected invalid index error for abc, got: %v", err)
	}
	_, err = h.CommandAndExpect("remind del 99", h.Contains("Invalid reminder index"))
	if err != nil {
		t.Errorf("expected invalid index error for 99, got: %v", err)
	}
}
