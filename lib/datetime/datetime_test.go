package datetime

import (
	"testing"
//	"time"
)

func TestParse(t *testing.T) {
	if ret, _ := Parse("this sunday"); true {
		t.Errorf("oawww %#v %s.", ret, ret)
	}
}
