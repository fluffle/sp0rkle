package datetime

import (
	"testing"
)

func TestParse(t *testing.T) {
	if ret := Parse("11:24 pm march 4, 2030"); ret == nil || ret.Seconds() != 3600 {
		t.Errorf("oawww %T(%#v).", ret, ret)
	}
}
