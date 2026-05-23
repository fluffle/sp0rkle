//go:build integration

package regtest

import (
	"testing"
)

func testLearn(t *testing.T) {
	t.Run("no_args", func(t *testing.T) {
		_, err := h.CommandAndExpect("learn", h.Contains("can't learn from you"))
		if err != nil {
			t.Errorf("expected error for learn without args, got: %v", err)
		}
	})

	t.Run("tag_only", func(t *testing.T) {
		_, err := h.CommandAndExpect("learn onlytag", h.Contains("can't learn from you"))
		if err != nil {
			t.Errorf("expected error for learn with tag only, got: %v", err)
		}
	})

	t.Run("success", func(t *testing.T) {
		_, err := h.CommandAndExpect("learn testtag this is a test sentence", h.Contains("Ta"))
		if err != nil {
			t.Errorf("expected success response, got: %v", err)
		}
	})
}
