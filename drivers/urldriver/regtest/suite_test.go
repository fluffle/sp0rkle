//go:build integration

package regtest

import (
	"context"
	"os"
	"testing"

	"github.com/fluffle/sp0rkle/regtest"
)

var h *regtest.Harness

func TestMain(m *testing.M) {
	var err error
	if h, err = regtest.Start(context.Background()); err != nil {
		panic("Start: " + err.Error())
	}
	code := m.Run()
	if err = h.Stop(); err != nil {
		panic("Stop: " + err.Error())
	}
	os.Exit(code)
}

func TestCommands(t *testing.T) {
	t.Run("urlfind", testUrlFind)
	t.Run("urlfind_no_match", testUrlFindNoMatch)
	t.Run("randurl", testRandurl)
	t.Run("shorten_with_url", testShortenWithUrl)
	t.Run("shorten_that", testShortenThat)
	t.Run("shorten_already_shortened", testShortenThatAlreadyShortened)
	t.Run("shorten_not_urlish", testShortenNotUrlish)
	t.Run("shorten_with_extra_args", testShortenWithExtraArgs)
	t.Run("cache_not_urlish", testCacheNotUrlish)
	t.Run("cache_bad_url", testCacheBadUrlSubstring)
	t.Run("find_case_insensitive", testFindCaseInsensitive)
}

func TestHandlers(t *testing.T) {
	t.Run("scan_long_url_auto_shortens", testUrlScanLongUrlAutoShortens)
	t.Run("scan_duplicate_url", testUrlScanDuplicateUrl)
	t.Run("scan_multiple_urls", testUrlScanMultipleUrls)
	t.Run("scan_supported_schemes", testUrlScanSupportedSchemes)
	t.Run("scan_unsupported_schemes", testUrlScanUnsupportedSchemes)
	t.Run("scan_url_in_middle", testUrlScanUrlInMiddleOfText)
}
