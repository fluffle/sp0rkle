//go:build integration

package regtest

import (
	"regexp"
	"testing"
	"time"
)

var (
	regexUrl       = regexp.MustCompile(`https?://[^\s]+`)
	regexShortened = regexp.MustCompile(`shortened to https://[^/]+/s/[A-Za-z0-9_-]{6}`)
	regexCached    = regexp.MustCompile(`cached as https://[^/]+/c/[A-Za-z0-9_-]{6}`)
	regexAlreadyShortened = regexp.MustCompile(`already shortened as https://[^/]+/s/[A-Za-z0-9_-]{6}`)
	regexAlreadyCached = regexp.MustCompile(`already cached as https://[^/]+/c/[A-Za-z0-9_-]{6}`)
	regexFirstMentioned = regexp.MustCompile(`first mentioned by`)
	regexAutoShortened = regexp.MustCompile(`URL shortened as https://[^/]+/s/[A-Za-z0-9_-]{6}`)
	regexBadSubstring = regexp.MustCompile(`bad substring`)
)

func testUrlFind(t *testing.T) {
	// First populate the URL collection
	h.Privmsg(h.Channel, "check out https://example.com/foo for details")
	time.Sleep(100 * time.Millisecond)

	// Now use various aliases of the find command
	for _, cmd := range []string{"urlfind", "url find", "urlsearch", "url search"} {
		t.Run("find_"+cmd, func(t *testing.T) {
			_, err := h.CommandAndExpect(cmd+" example.com", h.Contains("example.com"))
			if err != nil {
				t.Errorf("command %q: expected to find URL, got: %v", cmd, err)
			}
		})
	}
}

func testUrlFindNoMatch(t *testing.T) {
	_, err := h.CommandAndExpect("url find nonexistent.com", h.Contains("No urls matching"))
	if err != nil {
		t.Errorf("expected not to find URL, got: %v", err)
	}
}

func testFindCaseInsensitive(t *testing.T) {
	h.Privmsg(h.Channel, "http://example.com/uppercase")
	time.Sleep(100 * time.Millisecond)

	_, err := h.CommandAndExpect("URLfind example.com", h.Contains("example.com"))
	if err != nil {
		t.Errorf("expected case-insensitive command to work, got: %v", err)
	}
}

func testRandurl(t *testing.T) {
	// First populate the URL collection
	h.Privmsg(h.Channel, "see https://example.com/rand for fun stuff")
	time.Sleep(100 * time.Millisecond)

	// Now use various aliases of the randurl command
	for _, cmd := range []string{"randurl", "random url"} {
		t.Run("randurl_"+cmd, func(t *testing.T) {
			_, err := h.CommandAndExpect(cmd, h.Contains("example.com"))
			if err != nil {
				t.Errorf("command %q: expected a random URL, got: %v", cmd, err)
			}
		})
	}
}

func testShortenWithUrl(t *testing.T) {
	_, err := h.CommandAndExpect("shorten https://example.com/very/long/path",
		h.Regex(regexShortened))
	if err != nil {
		t.Errorf("expected shortened url response, got: %v", err)
	}
}

func testShortenThat(t *testing.T) {
	// First, record a URL via urlScan (silent)
	h.Privmsg(h.Channel, "https://example.com/shortenme")
	time.Sleep(100 * time.Millisecond)

	// Now shorten the last seen URL
	_, err := h.CommandAndExpect("shorten that",
		h.Regex(regexShortened))
	if err != nil {
		t.Errorf("expected shorten that response, got: %v", err)
	}
}

func testShortenThatAlreadyShortened(t *testing.T) {
	// First, record and shorten a URL
	h.Privmsg(h.Channel, "https://example.com/alreadyshortened")
	time.Sleep(100 * time.Millisecond)
	_, err := h.CommandAndExpect("shorten that",
		h.Regex(regexShortened))
	if err != nil {
		t.Errorf("expected shorten that response, got: %v", err)
	}

	// Now try to shorten the already shortened URL
	_, err = h.CommandAndExpect("shorten https://example.com/alreadyshortened",
		h.Regex(regexAlreadyShortened))
	if err != nil {
		t.Errorf("expected already shortened response, got: %v", err)
	}
}

func testShortenNotUrlish(t *testing.T) {
	_, err := h.CommandAndExpect("shorten not-a-url",
		h.Contains("doesn't look URLish"))
	if err != nil {
		t.Errorf("expected 'not URLish' response, got: %v", err)
	}
}

func testShortenWithExtraArgs(t *testing.T) {
	_, err := h.CommandAndExpect("shorten https://example.com/urlwithargs plus extra text",
		h.Regex(regexShortened))
	if err != nil {
		t.Errorf("expected shortened url response, got: %v", err)
	}
}

func testCacheNotUrlish(t *testing.T) {
	_, err := h.CommandAndExpect("cache not-a-url",
		h.Contains("doesn't look URLish"))
	if err != nil {
		t.Errorf("expected 'not URLish' response, got: %v", err)
	}
}

func testCacheBadUrlSubstring(t *testing.T) {
	_, err := h.CommandAndExpect("cache https://4chan.org/example",
		h.Regex(regexBadSubstring))
	if err != nil {
		t.Errorf("expected bad substring error, got: %v", err)
	}
}

func testCacheThat(t *testing.T) {
	t.Skip("Need local HTTP server in harness to test caching.")
}
