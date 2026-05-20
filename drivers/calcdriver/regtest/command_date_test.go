//go:build integration

package regtest

import (
	"regexp"
	"testing"
	"time"
)

// Matches datetime.TimeFormat
var dateRx = regexp.MustCompile(`\d{2}:\d{2}:\d{2}, \w+ \d{1,2} \w+ \d{4} [A-Z]{3}`)

func testDate(t *testing.T) {
	t.Run("current_time_no_args", func(t *testing.T) {
		_, err := h.CommandAndExpect("date", h.Regex(dateRx))
		if err != nil {
			t.Errorf("expected to match date regex, got: %v", err)
		}
	})

	t.Run("date_with_midday", func(t *testing.T) {
		_, err := h.CommandAndExpect("date midday", h.Contains("12:00:00"))
		if err != nil {
			t.Errorf("expected midday time in response, got: %v", err)
		}
	})

	t.Run("date_with_specific_date", func(t *testing.T) {
		_, err := h.CommandAndExpect("date 2024-01-01", h.Contains("1 January 2024"))
		if err != nil {
			t.Errorf("expected date output for Jan 1 2024, got: %v", err)
		}
	})

	t.Run("date_with_tomorrow", func(t *testing.T) {
		// Might be a bit flaky around midnights, especially during TZ changes.
		// Ehhhhhh do I care? No.
		expect := time.Now().Add(24*time.Hour).Format("Monday 2 January 2006")
		_, err := h.CommandAndExpect("date tomorrow", h.Contains(expect))
		if err != nil {
			t.Errorf("expected date output for tomorrow, got: %v", err)
		}
	})

	t.Run("date_with_yesterday", func(t *testing.T) {
		expect := time.Now().Add(-24*time.Hour).Format("Monday 2 January 2006")
		_, err := h.CommandAndExpect("date yesterday", h.Contains(expect))
		if err != nil {
			t.Errorf("expected date output for yesterday, got: %v", err)
		}
	})

	t.Run("date_with_next_day", func(t *testing.T) {
		_, err := h.CommandAndExpect("date next monday", h.Contains("Monday"))
		if err != nil {
			t.Errorf("expected monday in response, got: %v", err)
		}
	})

	t.Run("date_with_timezone", func(t *testing.T) {
		_, err := h.CommandAndExpect("date in UTC", h.Contains("UTC"))
		if err != nil {
			t.Errorf("expected UTC timezone in response, got: %v", err)
		}
	})

	t.Run("date_with_timezone_in_prefix", func(t *testing.T) {
		_, err := h.CommandAndExpect("date 12:00 in UTC", h.Regex(regexp.MustCompile(`12:00:00, .* UTC`)))
		if err != nil {
			t.Errorf("expected UTC timezone with time, got: %v", err)
		}
	})

	t.Run("date_with_iso_format", func(t *testing.T) {
		_, err := h.CommandAndExpect("date 2024-06-15", h.Contains("15 June"))
		if err != nil {
			t.Errorf("expected 15 June in response, got: %v", err)
		}
	})

	t.Run("date_with_time_only", func(t *testing.T) {
		_, err := h.CommandAndExpect("date 14:30", h.Contains("14:30:00"))
		if err != nil {
			t.Errorf("expected 14:30 in response, got: %v", err)
		}
	})

	t.Run("date_invalid_midday_word", func(t *testing.T) {
		// "noon" is not a recognized date token (only "midday" is)
		_, err := h.CommandAndExpect("date noon", h.Contains("Couldn't parse"))
		if err != nil {
			t.Errorf("expected parse error for 'noon', got: %v", err)
		}
	})

	t.Run("date_invalid", func(t *testing.T) {
		_, err := h.CommandAndExpect("date notadate xyz", h.Contains("Couldn't parse"))
		if err != nil {
			t.Errorf("expected parse error for invalid date, got: %v", err)
		}
	})

	t.Run("date_next_month", func(t *testing.T) {
		_, err := h.CommandAndExpect("date next march", h.Contains("March"))
		if err != nil {
			t.Errorf("expected March in response, got: %v", err)
		}
	})

	t.Run("date_with_offset", func(t *testing.T) {
		expect := time.Now().Add(5*24*time.Hour).Format("Monday 2 January 2006")
		_, err := h.CommandAndExpect("date 5 days", h.Contains(expect))
		if err != nil {
			t.Errorf("expected 5 days in the future in response, got: %v", err)
		}
	})

	t.Run("date_case_insensitive", func(t *testing.T) {
		_, err := h.CommandAndExpect("DATE", h.Regex(dateRx))
		if err != nil {
			t.Errorf("expected case insensitive command, got: %v", err)
		}
	})

	t.Run("date_midnight", func(t *testing.T) {
		_, err := h.CommandAndExpect("date midnight", h.Contains("00:00:00"))
		if err != nil {
			t.Errorf("expected midnight time in response, got: %v", err)
		}
	})

	t.Run("date_next_week", func(t *testing.T) {
		expect := time.Now().Format("Monday")
		_, err := h.CommandAndExpect("date next week", h.Contains(expect))
		if err != nil {
			t.Errorf("expected same day of week in response, got: %v", err)
		}
	})

	t.Run("date_today", func(t *testing.T) {
		expect := time.Now().Format("Monday")
		_, err := h.CommandAndExpect("date today", h.Contains(expect))
		if err != nil {
			t.Errorf("expected same day of week in response, got: %v", err)
		}
	})

	t.Run("date_iso_datetime", func(t *testing.T) {
		_, err := h.CommandAndExpect("date 2024-01-15 10:30", h.Contains("10:30:00, Monday 15 January 2024"))
		if err != nil {
			t.Errorf("expected exact datetime in response, got: %v", err)
		}
	})

	t.Run("date_with_am_pm", func(t *testing.T) {
		_, err := h.CommandAndExpect("date 12pm", h.Contains("12:00:00"))
		if err != nil {
			t.Errorf("expected 12 in response, got: %v", err)
		}
	})

	t.Run("date_with_fortnight", func(t *testing.T) {
		expect := time.Now().Format("Monday")
		_, err := h.CommandAndExpect("date in a fortnight", h.Contains(expect))
		if err != nil {
			t.Errorf("expected same day in response, got: %v", err)
		}
	})

	t.Run("date_with_week", func(t *testing.T) {
		expect := time.Now().Format("Monday")
		_, err := h.CommandAndExpect("date in a week", h.Contains(expect))
		if err != nil {
			t.Errorf("expected same day in response, got: %v", err)
		}
	})
}
