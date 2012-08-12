package urldriver

import (
	"net/http"
	"os"
)

func (ud *urlDriver) RegisterHttpHandlers() {
	if err := os.MkdirAll(*urlCacheDir, 0700); err != nil {
		ud.l.Fatal("Couldn't create URL cache dir: %v", err)
	}
	// This serves "shortened" urls 
	http.Handle(shortenPath, http.StripPrefix(shortenPath, ud))
	// This serves "cached" urls
	http.Handle(cachePath, http.FileServer(http.Dir(*urlCacheDir)))
}

func (ud *urlDriver) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if u := ud.GetShortened(req.URL.Path); u != nil {
		rw.Header().Set("Location", u.Url)
		rw.WriteHeader(302)
		return
	}
	http.NotFound(rw, req)
}
