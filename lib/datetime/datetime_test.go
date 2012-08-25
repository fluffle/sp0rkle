package datetime

import (
	"testing"
//	"time"
)

func TestParse(t *testing.T) {
	if ret, _ := Parse("1 sunday"); true {
		t.Errorf("oawww %#v %s.", ret, ret)
	}
}
