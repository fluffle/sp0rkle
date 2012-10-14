package urldriver

import (
	"net/http"
)

func shortenedServer(rw http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "" {
		http.NotFound(rw, req)
	}
	if u := uc.GetShortened(req.URL.Path); u != nil {
		rw.Header().Set("Location", u.Url)
		rw.WriteHeader(302)
		return
	}
	http.NotFound(rw, req)
}
