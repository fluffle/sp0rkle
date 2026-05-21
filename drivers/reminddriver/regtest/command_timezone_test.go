//go:build integration

package regtest

import (
	"testing"
)

func testZone(t *testing.T) {
	t.Run("valid_zone", testZoneValid)
	t.Run("invalid_zone", testZoneInvalid)
	t.Run("no_zone_arg", testZoneNoArg)
}

func testUnzone(t *testing.T) {
	_, err := h.CommandAndExpect("forget my timezone", h.Contains("forgotten where you live"))
	if err != nil {
		t.Errorf("expected timezone forgotten response, got: %v", err)
	}
}

func testZoneValid(t *testing.T) {
	_, err := h.CommandAndExpect("my timezone is America/New_York", h.Contains(`in "America/New_York"`))
	if err != nil {
		t.Errorf("expected new york set response, got: %v", err)
	}

	// Overwrite zone with one that won't have DST naming problems.
	_, err = h.CommandAndExpect("my timezone is UTC", h.Contains(`in "UTC"`))
	if err != nil {
		t.Errorf("expected timezone set response, got: %v", err)
	}

	// Trigger a reminder in 1s and expect it to have UTC zone.
	_, err = h.CommandAndExpect("remind me timezones in 1 second", h.Contains("UTC"))
	if err != nil {
		t.Errorf("expected reminder ack, got: %v", err)
	}
	_, err = h.Expect(h.Contains("you asked me to remind you timezones")).Wait()
	if err != nil {
		t.Errorf("expected reminder to fire, got: %v", err)
	}
}

func testZoneInvalid(t *testing.T) {
	_, err := h.CommandAndExpect("my timezone is NotARealZone", h.Contains("Don't recognise"))
	if err != nil {
		t.Errorf("expected invalid timezone response, got: %v", err)
	}
}

func testZoneNoArg(t *testing.T) {
	_, err := h.CommandAndExpect("my timezone is", h.Contains("fat"))
	if err != nil {
		t.Errorf("expected fat timezone response, got: %v", err)
	}
}
