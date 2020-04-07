// Package handler implement some basic handler for the webfmwk
// handler provides a convenient mechanism for filtering HTTP requests
// entering the application. It returns a new handler which may perform various
// operations and should finish by calling the next HTTP handler.
package handler

import "github.com/burgesQ/webfmwk/v4"

// Security append few security headers
func Security(next webfmwk.HandlerFunc) webfmwk.HandlerFunc {
	return webfmwk.HandlerFunc(func(c webfmwk.IContext) {
		c.SetHeaders(webfmwk.Header{"X-XSS-Protection", "1; mode=block"},
			webfmwk.Header{"X-Content-Type-Options", "nosniff"},
			webfmwk.Header{"Strict-Transport-Security", "max-age=3600; includeSubDomains"})
		next(c)
	})
}
