package urldriver

import (
	"encoding/base64"
	"flag"
	"fmt"
	"hash/crc32"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fluffle/goirc/client"
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/collections/urls"
	"github.com/fluffle/sp0rkle/util"
	"gopkg.in/mgo.v2/bson"
)

const shortenPath string = "/s/"
const cachePath string = "/c/"
const autoShortenLimit int = 120
const maxCacheSize = 1 << 22 // 4MB

var badUrlStrings = []string{
	"4chan",
}

var urlCacheDir *string = flag.String("url_cache_dir",
	util.JoinPath(os.Getenv("HOME"), ".sp0rkle"),
	"Path to store cached content under.")

var uc *urls.Collection

// Remember the last url seen on a per-channel basis
var lastseen = map[string]bson.ObjectId{}

func Init() {
	uc = urls.Init()

	if err := os.MkdirAll(*urlCacheDir, 0700); err != nil {
		logging.Fatal("Couldn't create URL cache dir: %v", err)
	}

	bot.Handle(urlScan, client.PRIVMSG)

	bot.Command(find, "urlfind", "urlfind <regex>  -- "+
		"searches for previously mentioned URLs matching <regex>")
	bot.Command(find, "url find", "url find <regex>  -- "+
		"searches for previously mentioned URLs matching <regex>")
	bot.Command(find, "urlsearch", "urlsearch <regex>  -- "+
		"searches for previously mentioned URLs matching <regex>")
	bot.Command(find, "url search", "url search <regex>  -- "+
		"searches for previously mentioned URLs matching <regex>")

	bot.Command(find, "randurl", "randurl  -- displays a random URL")
	bot.Command(find, "random url", "random url  -- displays a random URL")

	bot.Command(shorten, "shorten that", "shorten that  -- "+
		"shortens the last mentioned URL.")
	bot.Command(shorten, "shorten", "shorten <url>  -- shortens <url>")

	bot.Command(cache, "cache that", "cache that  -- "+
		"caches the last mentioned URL.")
	bot.Command(cache, "cache", "cache <url>  -- caches <url>")
	bot.Command(cache, "save that", "save that  -- "+
		"caches the last mentioned URL.")
	bot.Command(cache, "save", "save <url>  -- caches <url>")

	// This serves "shortened" urls
	http.Handle(shortenPath, http.StripPrefix(shortenPath,
		http.HandlerFunc(shortenedServer)))

	// This serves "cached" urls
	http.Handle(cachePath, http.StripPrefix(cachePath,
		http.FileServer(http.Dir(*urlCacheDir))))
}

func Encode(url string) string {
	// We shorten/cache a url with it's base-64 encoded CRC32 hash
	crc := crc32.ChecksumIEEE([]byte(url))
	crcb := make([]byte, 4)
	for i := 0; i < 4; i++ {
		crcb[i] = byte((crc >> uint32(i)) & 0xff)
	}
	// Avoid collisions in shortened URLs
	for i := 0; i < 10; i++ {
		// Since we're always encoding exactly 4 bytes (32 bits)
		// resulting in 5 1/3 bytes of encoded data, we can drop
		// the two padding equals signs for brevity.
		s := (base64.URLEncoding.EncodeToString(crcb))[:6]
		cached := uc.GetCached(s)
		shortened := uc.GetShortened(s)
		if !(cached.Exists() || shortened.Exists()) {
			return s
		}
		crcb[rand.Intn(4)]++
	}
	logging.Warn("Collided ten times while encoding URL.")
	return ""
}

func Shorten(u *urls.Url) error {
	u.Shortened = Encode(u.Url)
	if err := uc.Put(u); err != nil {
		return err
	}
	return nil
}

func Cache(u *urls.Url) error {
	u.CachedAs = Encode(u.Url)
	if u.CachedAs == "" {
		return fmt.Errorf("collided 10 times while encoding URL")
	}
	for _, s := range badUrlStrings {
		if strings.Index(u.Url, s) != -1 {
			return fmt.Errorf("url contains bad substring %q", s)
		}
	}
	// Try a HEAD req first to get Content-Length header.
	res, err := http.Head(u.Url)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("received non-200 response %q from server.", res.Status)
	}
	if size := res.Header.Get("Content-Length"); size != "" {
		if bytes, err := strconv.Atoi(size); err != nil {
			return fmt.Errorf("received unparseable content length %q from server: %v", size, err)
		} else if bytes > maxCacheSize {
			return fmt.Errorf("response too large (%d MB) to cache safely", bytes/1024/1024)
		}
	}
	res, err = http.Get(u.Url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.ContentLength > maxCacheSize {
		return fmt.Errorf("response too large (%d MB) to cache safely",
			res.ContentLength/1024/1024)
	}
	fh, err := os.OpenFile(util.JoinPath(*urlCacheDir, u.CachedAs),
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(0600))
	defer fh.Close()
	if err != nil {
		return err
	}
	if _, err := io.Copy(fh, res.Body); err != nil {
		return err
	}
	u.CacheTime = time.Now()
	u.MimeType = res.Header.Get("Content-Type")
	if err := uc.Put(u); err != nil {
		return err
	}
	return nil
}
