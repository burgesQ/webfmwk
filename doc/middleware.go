package main

import (
	"net/http"

	"github.com/burgesQ/webfmwk/v4"
	"github.com/burgesQ/webfmwk/v4/log"
)

// Middleware implement http.Handler methods
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
		webfmwk.WithMiddlewares(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				log.Infof("[%s] %s", r.Method, r.RequestURI)
				next.ServeHTTP(w, r)
			})
		}))

	// expose /test
	s.GET("/test", func(c webfmwk.IContext) {
		c.JSONOk("ok")
	})

	// start asynchronously on :4242
	s.Start(":4242")

	// ctrl+c is handled internaly
	defer s.WaitAndStop()
}
