package datetime

import (
	"testing"
//	"time"
)

func TestParse(t *testing.T) {
	if ret := Parse("3 years ago"); ret.IsZero() || ret.Second() != 3600 {
		t.Errorf("oawww %#v %s.", ret, ret)
	}
}
