package main

import (
	"github.com/burgesQ/webfmwk/v5"
	"github.com/burgesQ/webfmwk/v5/handler/logging"
	"github.com/burgesQ/webfmwk/v5/handler/security"
)

// Handlers implement webfmwk.Handler methods
// Check the server logs
//
// curl -i -X GET 127.0.0.1:4242/test
// Accept: application/json; charset=UTF-8
// Content-Type: application/json; charset=UTF-8
// Produce: application/json; charset=UTF-8
// Strict-Transport-Security: max-age=3600; includeSubDomains
// X-Content-Type-Options: nosniff
// X-Xss-Protection: 1; mode=block
// Date: Mon, 06 Apr 2020 14:58:44 GMT
// Content-Length: 4
func handlers() *webfmwk.Server {
	// init server w/ ctrl+c support and middlewares
	var s = webfmwk.InitServer(
		webfmwk.WithCtrlC(),
		webfmwk.WithHandlers(logging.Handler),
	)

	// expose /test
	s.GET("/test", security.Handler(func(c webfmwk.Context) error {
		return c.JSONOk("ok")
	}))

	return s
}
