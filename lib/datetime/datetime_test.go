package datetime

import (
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	if ret := Parse("2pm PDT 15th feb"); ret.IsZero() || ret.Second() != 3600 {
		t.Errorf("oawww %#v %s.", ret, ret)
	}
}
