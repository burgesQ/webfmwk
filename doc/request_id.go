package main

import (
	"github.com/burgesQ/webfmwk/v4"
)

// register the RequestID BEFORE the logger
func main() {
	// init server w/ ctrl+c support and middlewares
	var s = webfmwk.InitServer(webfmwk.WithHandlers(handler.Logging, handler.SetRequestID))

	// expose /test
	s.GET("/test", func(c webfmwk.IContext) {
		c.JSONOk("ok")
	})

	// start asynchronously on :4242
	s.Start(":4242")

	// ctrl+c is handled internaly
	defer s.WaitAndStop()
}
