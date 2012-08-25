package datetime

import (
	"testing"
//	"time"
)

func TestParse(t *testing.T) {
	if ret := Parse("March 20th"); ret.IsZero() || ret.Second() != 3600 {
		t.Errorf("oawww %#v %s.", ret, ret)
	}
}
