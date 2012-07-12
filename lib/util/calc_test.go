package util

import (
	"testing"
)

func TestCalc(t *testing.T) {
	res, err := Calc("2+2")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if res != 4 {
		t.Errorf("2+2 is apparently %f", res)
	}
}
