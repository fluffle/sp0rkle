package flightdriver

import (
	"testing"
	"time"

	"github.com/fluffle/sp0rkle/bot"
)

func TestStalkLogic(t *testing.T) {
	fp = &flightPoller{
		tracking: make(map[string]*flightInfo),
	}

	fn := "AA123"
	ch := bot.Chan("#channel")
	me := "sp0rkle"
	key := me + ":" + string(ch) + ":" + fn

	fp.tracking[key] = &flightInfo{
		FlightNum: fn,
		Target:    ch,
		Me:        me,
		StartTime: time.Now(),
	}

	if len(fp.tracking) != 1 {
		t.Errorf("Expected 1 tracked flight, got %d", len(fp.tracking))
	}

	info, ok := fp.tracking[key]
	if !ok {
		t.Errorf("Flight %s not found in tracking", key)
	}
	if info.FlightNum != fn {
		t.Errorf("Expected flight number %s, got %s", fn, info.FlightNum)
	}
}

func TestFormatDelay(t *testing.T) {
	tests := []struct {
		in  interface{}
		out string
	}{
		{nil, ""},
		{0.0, ""},
		{0, ""},
		{15.0, "15 mins"},
		{20, "20 mins"},
	}

	for _, tc := range tests {
		got := formatDelay(tc.in)
		if got != tc.out {
			t.Errorf("formatDelay(%v) = %q; want %q", tc.in, got, tc.out)
		}
	}
}
