//go:build integration

package regtest

import (
	"testing"
)

// We can't necessarily guarantee test ordering, or at least should not
// _depend_ on ordering between tests, so we can only verifiably test snooze
// and actual reminder firing in a single test that does things in order.
//
// This comment is hilariously hypocritical given the inter-test dependencies
// present in command_remind_test.go
func testSnooze(t *testing.T) {
	// first test snooze _without_ a previous expired reminder
	_, err := h.CommandAndExpect("snooze", h.Contains("No record of an expired reminder"))
	if err != nil {
		t.Errorf("expected no snooze record message, got: %v", err)
	}

	// Now create a reminder that will fire almost instantly and
	// verify that it does indeed fire
	_, err = h.CommandAndExpect("remind me water plants in 1 second", h.Contains("okay, i'll remind you"))
	if err != nil {
		t.Errorf("expected reminder ack, got: %v", err)
	}
	_, err = h.Expect(h.Contains("you asked me to remind you water plants")).Wait()
	if err != nil {
		t.Errorf("expected reminder to fire, got: %v", err)
	}

	// now test snooze with a previous expired reminder and
	// verify that it fires a second time
	_, err = h.CommandAndExpect("snooze 1 second", h.Contains("okay, i'll remind you"))
	if err != nil {
		t.Errorf("expected snooze ack, got: %v", err)
	}
	_, err = h.Expect(h.Contains("you asked me to remind you water plants")).Wait()
	if err != nil {
		t.Errorf("expected snoozed reminder to fire, got: %v", err)
	}
}

