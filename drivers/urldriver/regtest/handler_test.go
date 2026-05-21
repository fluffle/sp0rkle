//go:build integration

package regtest

import (
	"strings"
	"testing"
	"time"
)

func testUrlScanLongUrlAutoShortens(t *testing.T) {
	// URLs longer than autoShortenLimit (120 chars) should be auto-shortened
	longUrl := "https://example.com/"
	for i := 0; i < 150; i++ {
		longUrl += "a"
	}
	w := h.Expect(h.Regex(regexAutoShortened))
	h.Privmsg(h.Channel, longUrl)
	_, err := w.Wait()
	if err != nil {
		t.Errorf("expected auto-shortened url response, got: %v", err)
	}
}

func testUrlScanDuplicateUrl(t *testing.T) {
	t.Skip("urlScan handler requires 2h since first mention before responding; impossible to test in harness")
	// First mention — recorded silently
	h.Privmsg(h.Channel, "https://example.com/duplicate-test")
	time.Sleep(100 * time.Millisecond)

	// Second mention — expect "that URL first mentioned by"
	w := h.Expect(h.Regex(regexFirstMentioned))
	h.Privmsg(h.Channel, "https://example.com/duplicate-test")
	_, err := w.Wait()
	if err != nil {
		t.Errorf("expected duplicate url response, got: %v", err)
	}
}

func testUrlScanMultipleUrls(t *testing.T) {
	// Multiple URLs in one message — both recorded silently
	h.Privmsg(h.Channel, "visit https://example.com/multi-one and https://example.com/multi-two for great peen sizes")
	time.Sleep(100 * time.Millisecond)

	// Verify both were recorded via urlfind
	for _, u := range []string{"one", "two"} {
		t.Run("multi_"+u, func(t *testing.T) {
			_, err := h.CommandAndExpect("urlfind multi-"+u, h.Contains("/multi-"+u))
			if err != nil {
				t.Errorf("expected to find scheme URL, got: %v", err)
			}
		})
	}
}

func testUrlScanSupportedSchemes(t *testing.T) {
	var urls = []string{
		"http://example.com/http-scheme-test",
		"https://example.com/https-scheme-test",
	}
	for _, u := range urls {
		scheme := strings.Split(u, ":")[0]
		t.Run("scheme_"+scheme, func(t *testing.T) {
			h.Privmsg(h.Channel, u)
			time.Sleep(100 * time.Millisecond)

			_, err := h.CommandAndExpect("urlfind "+scheme+"-scheme-test", h.Contains("/"+scheme+"-scheme-test"))
			if err != nil {
				t.Errorf("expected to find %s URL, got: %v", scheme, err)
			}
		})
	}
}

func testUrlScanUnsupportedSchemes(t *testing.T) {
	var urls = []string{
		"ftp://example.com/ftp-scheme-test",
		"smb://HOST.DOMAIN/smb-scheme-test",
	}
	for _, u := range urls {
		scheme := strings.Split(u, ":")[0]
		t.Run("scheme_"+scheme, func(t *testing.T) {
			h.Privmsg(h.Channel, u)
			time.Sleep(100 * time.Millisecond)

			_, err := h.CommandAndExpect("urlfind "+scheme+"-scheme-test", h.Contains("No urls matching"))
			if err != nil {
				t.Errorf("expected not to find %s URL, got: %v", scheme, err)
			}
		})
	}
}

func testUrlScanUrlInMiddleOfText(t *testing.T) {
	h.Privmsg(h.Channel, "hey guys check out https://example.com/middle-of-sentence please")
	time.Sleep(100 * time.Millisecond)

	_, err := h.CommandAndExpect("urlfind middle", h.Contains("middle-of-sentence"))
	if err != nil {
		t.Errorf("expected to find middle of sentence URL, got: %v", err)
	}
}
