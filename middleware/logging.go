// Package middleware implement some basic middleware for the webfmwk
// middleware provides a convenient mechanism for filtering HTTP requests
// entering the application. It returns a new handler which may perform various
// operations and should finish by calling the next HTTP handler.
package middleware

import (
	"net/http"

  webfmwk "github.com/burgesQ/webfmwk/v2"
)

// Logging log information about the newly receive request
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := webfmwk.GetLogger()
		log.Infof("[+] (%s): [%s]%s", r.RemoteAddr, r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
