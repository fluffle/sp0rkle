//go:build integration

package regtest

import (
	"regexp"
	"testing"
)

func testQadd(t *testing.T) {
	_, err := h.CommandAndExpect("qadd the quick brown fox", h.Contains("Quote added succesfully"))
	if err != nil {
		t.Errorf("expected success response, got: %v", err)
	}
}

func testQuoteAdd(t *testing.T) {
	_, err := h.CommandAndExpect("quote add lazy dogs are fun", h.Contains("Quote added succesfully"))
	if err != nil {
		t.Errorf("expected success response, got: %v", err)
	}
}

func testAddQuote(t *testing.T) {
	_, err := h.CommandAndExpect("add quote another fine day", h.Contains("Quote added succesfully"))
	if err != nil {
		t.Errorf("expected success response, got: %v", err)
	}
}

func testQdel(t *testing.T) {
	// Add a quote, extract its ID, then delete it
	line, err := h.CommandAndExpect("qadd temp quote for deletion", h.Regex(regexp.MustCompile(`Quote added succesfully, id #(\d+)\.`)))
	if err != nil {
		t.Fatalf("failed to add quote: %v", err)
	}
	matched := regexp.MustCompile(`id #(\d+)\.`).FindStringSubmatch(line.Text())
	if len(matched) < 2 {
		t.Fatalf("could not extract QID from: %s", line.Text())
	}

	_, err = h.CommandAndExpect("qdel #"+matched[1], h.Contains("I forgot quote"))
	if err != nil {
		t.Errorf("expected delete success, got: %v", err)
	}

	// Fetch the quote by ID to validate deleted.
	_, err = h.CommandAndExpect("quote #"+matched[1], h.Contains("No quote found"))
	if err != nil {
		t.Errorf("expected no quote, got: %v", err)
	}

}

func testQuoteDel(t *testing.T) {
	line, err := h.CommandAndExpect("qadd temp quote for del2", h.Regex(regexp.MustCompile(`Quote added succesfully, id #(\d+)\.`)))
	if err != nil {
		t.Fatalf("failed to add quote: %v", err)
	}
	matched := regexp.MustCompile(`id #(\d+)\.`).FindStringSubmatch(line.Text())
	if len(matched) < 2 {
		t.Fatalf("could not extract QID from: %s", line.Text())
	}

	_, err = h.CommandAndExpect("quote del "+matched[1], h.Contains("I forgot quote"))
	if err != nil {
		t.Errorf("expected delete success, got: %v", err)
	}
}

func testDelQuote(t *testing.T) {
	line, err := h.CommandAndExpect("qadd temp quote for del3", h.Regex(regexp.MustCompile(`Quote added succesfully, id #(\d+)\.`)))
	if err != nil {
		t.Fatalf("failed to add quote: %v", err)
	}
	matched := regexp.MustCompile(`id #(\d+)\.`).FindStringSubmatch(line.Text())
	if len(matched) < 2 {
		t.Fatalf("could not extract QID from: %s", line.Text())
	}

	_, err = h.CommandAndExpect("del quote #"+matched[1], h.Contains("I forgot quote"))
	if err != nil {
		t.Errorf("expected delete success, got: %v", err)
	}
}

func testDelInvalidID(t *testing.T) {
	_, err := h.CommandAndExpect("qdel notanumber", h.Contains("doesn't look like a quote id"))
	if err != nil {
		t.Errorf("expected error for invalid ID, got: %v", err)
	}
}

func testDelNonexistent(t *testing.T) {
	_, err := h.CommandAndExpect("quote del 1234567", h.Contains("No quote found"))
	if err != nil {
		t.Errorf("expected no quote, got: %v", err)
	}
}

func testFetch(t *testing.T) {
	// Add a quote, fetch it by ID, and test invalid ID
	line, err := h.CommandAndExpect("qadd fetch me if you can", h.Regex(regexp.MustCompile(`Quote added succesfully, id #(\d+)\.`)))
	if err != nil {
		t.Fatalf("failed to add quote: %v", err)
	}
	matched := regexp.MustCompile(`id #(\d+)\.`).FindStringSubmatch(line.Text())
	if len(matched) < 2 {
		t.Fatalf("could not extract QID from: %s", line.Text())
	}

	// Fetch the quote by ID
	_, err = h.CommandAndExpect("quote #"+matched[1], h.Contains("fetch me if you can"))
	if err != nil {
		t.Errorf("expected quote content, got: %v", err)
	}

	// Fetch with invalid ID
	_, err = h.CommandAndExpect("quote #abc", h.Contains("doesn't look like a quote id"))
	if err != nil {
		t.Errorf("expected error for invalid ID, got: %v", err)
	}
}

func testLookup(t *testing.T) {
	// Lookup a quote matching "quick" (we added "the quick brown fox" earlier)
	_, err := h.CommandAndExpect("quote quick", h.Regex(regexp.MustCompile(`#\d+:.*quick.*fox`)))
	if err != nil {
		t.Errorf("expected matching quote, got: %v", err)
	}
}

func testLookupNoMatch(t *testing.T) {
	_, err := h.CommandAndExpect("quote xyzzy999nonexistent42", h.Contains("No quotes matching"))
	if err != nil {
		t.Errorf("expected no match response, got: %v", err)
	}
}
