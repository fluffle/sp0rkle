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

	"github.com/fluffle/golog/logging"
	"github.com/fluffle/goirc/client"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/collections/urls"
	"github.com/fluffle/sp0rkle/util"
	"github.com/fluffle/sp0rkle/util/bson"
)

const shortenPath string = "/s/"
const cachePath string = "/c/"
const autoShortenLimit int = 120
const maxCacheSize = 1 << 22 // 4MB

// TODO: this should be in a conf namespace.
var badUrlStrings = []string{
	"4chan",
}

var urlCacheDir *string = flag.String("url_cache_dir",
	util.JoinPath(os.Getenv("HOME"), ".sp0rkle"),
	"Path to store cached content under.")

type Driver struct {
	uc       *urls.Collection
	lastseen map[string]bson.ObjectId
}


func New(b *bot.Bot, uc *urls.Collection) *Driver {
	d := &Driver{uc: uc, lastseen: map[string]bson.ObjectId{}}

	if err := os.MkdirAll(*urlCacheDir, 0700); err != nil {
		logging.Fatal("Couldn't create URL cache dir: %v", err)
	}

	b.Handle(d.urlScan, client.PRIVMSG)

	b.Command(d.find, "urlfind", "urlfind <regex>  -- "+
		"searches for previously mentioned URLs matching <regex>")
	b.Command(d.find, "url find", "url find <regex>  -- "+
		"searches for previously mentioned URLs matching <regex>")
	b.Command(d.find, "urlsearch", "urlsearch <regex>  -- "+
		"searches for previously mentioned URLs matching <regex>")
	b.Command(d.find, "url search", "url search <regex>  -- "+
		"searches for previously mentioned URLs matching <regex>")

	b.Command(d.find, "randurl", "randurl  -- displays a random URL")
	b.Command(d.find, "random url", "random url  -- displays a random URL")

	b.Command(d.shorten, "shorten that", "shorten that  -- "+
		"shortens the last mentioned URL.")
	b.Command(d.shorten, "shorten", "shorten <url>  -- shortens <url>")

	b.Command(d.cache, "cache that", "cache that  -- "+
		"caches the last mentioned URL.")
	b.Command(d.cache, "cache", "cache <url>  -- caches <url>")
	b.Command(d.cache, "save that", "save that  -- "+
		"caches the last mentioned URL.")
	b.Command(d.cache, "save", "save <url>  -- caches <url>")

	// This serves "shortened" urls
	http.Handle(shortenPath, http.StripPrefix(shortenPath,
		http.HandlerFunc(d.shortenedServer)))

	// This serves "cached" urls
	http.Handle(cachePath, http.StripPrefix(cachePath,
		http.FileServer(http.Dir(*urlCacheDir))))

	return d
}

func (d *Driver) Encode(url string) string {
	// We shorten/cache a url with it's base-64 encoded CRC32 hash
	crc := crc32.ChecksumIEEE([]byte(url))
	crcb := make([]byte, 4)
	for i := range 4 {
		crcb[i] = byte((crc >> uint32(i)) & 0xff)
	}
	// Avoid collisions in shortened URLs
	for range 10 {
		// Since we're always encoding exactly 4 bytes (32 bits)
		// resulting in 5 1/3 bytes of encoded data, we can drop
		// the two padding equals signs for brevity.
		s := (base64.URLEncoding.EncodeToString(crcb))[:6]
		cached := d.uc.GetCached(s)
		shortened := d.uc.GetShortened(s)
		if !(cached.Exists() || shortened.Exists()) {
			return s
		}
		crcb[rand.Intn(4)]++
	}
	logging.Warn("Collided ten times while encoding URL.")
	return ""
}

func (d *Driver) Shorten(u *urls.Url) error {
	u.Shortened = d.Encode(u.Url)
	if err := d.uc.Put(u); err != nil {
		return err
	}
	return nil
}

func (d *Driver) Cache(u *urls.Url) error {
	u.CachedAs = d.Encode(u.Url)
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
	if err := d.uc.Put(u); err != nil {
		return err
	}
	return nil
}
