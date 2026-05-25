package urldriver

import (
	"net/http"
)

func (d *Driver) shortenedServer(rw http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "" {
		http.NotFound(rw, req)
	}
	if u := d.uc.GetShortened(req.URL.Path); u != nil {
		rw.Header().Set("Location", u.Url)
		rw.WriteHeader(302)
		return
	}
	http.NotFound(rw, req)
}
