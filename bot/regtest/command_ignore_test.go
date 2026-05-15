//go:build integration

package regtest

import (
	"testing"
)

func testIgnore(t *testing.T) {
	t.Run("ignore", testIgnoreBasic)
	t.Run("unignore", testUnignoreBasic)
	t.Run("case_insensitive", testIgnoreCaseInsensitive)
}

func testIgnoreBasic(t *testing.T) {
	_, err := h.CommandAndExpect("ignore testnick", h.Contains("I'll ignore 'testnick'"))
	if err != nil {
		t.Errorf("expected ignore response, got: %v", err)
	}
}

func testIgnoreCaseInsensitive(t *testing.T) {
	_, err := h.CommandAndExpect("IGNORE TESTNICK", h.Contains("I'll ignore 'testnick'"))
	if err != nil {
		t.Errorf("expected case-insensitive ignore response, got: %v", err)
	}
}

func testUnignoreBasic(t *testing.T) {
	_, err := h.CommandAndExpect("unignore testnick", h.Contains("No longer ignoring"))
	if err != nil {
		t.Errorf("expected unignore response, got: %v", err)
	}
}
