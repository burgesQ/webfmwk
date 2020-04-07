package main

import (
	"github.com/burgesQ/webfmwk/v3"
	"github.com/burgesQ/webfmwk/v3/handler"
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
func main() {
	// init server w/ ctrl+c support and middlewares
	s := webfmwk.InitServer(
		webfmwk.WithCtrlC(),
		webfmwk.WithHandlers(handler.Logging))

	// expose /test
	s.GET("/test", handler.Security(func(c webfmwk.IContext) {
		c.JSONOk("ok")
	}))

	// start asynchronously on :4242
	s.Start(":4242")

	// ctrl+c is handled internaly
	s.WaitAndStop()
}
