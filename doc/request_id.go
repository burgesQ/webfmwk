package main

import (
	"github.com/burgesQ/webfmwk/v4"
	"github.com/burgesQ/webfmwk/v4/handler"
)

// register the RequestID BEFORE the logger
func request_id() {
	// init server w/ ctrl+c support and middlewares
	var s = webfmwk.InitServer(webfmwk.WithCtrlC(),
		webfmwk.WithHandlers(handler.Logging, handler.RequestID))

	// expose /test
	s.GET("/test", func(c webfmwk.Context) error {
		return c.JSONOk("ok")
	})

	// start asynchronously on :4242
	s.Start(":4242")

	// ctrl+c is handled internaly
	defer s.WaitAndStop()
}
