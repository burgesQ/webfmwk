package main

import (
	"github.com/burgesQ/webfmwk/v3"
	"github.com/burgesQ/webfmwk/v3/log"
)

// GetLogger return a log.ILog interface
var logger = log.GetLogger()

// curl -X GET 127.0.0.1:4242/test
// "ok"
func main() {
	// init server w/ ctrl+c support and custom logger options
	var s = webfmwk.InitServer(
		webfmwk.WithLogger(logger),
		webfmwk.WithCtrlC())

	// expose /test
	s.GET("/test", func(c webfmwk.IContext) error {
		c.JSONOk("ok")
	})

	// start asynchronously on :4242
	s.Start(":4242")

	// ctrl+c is handled internaly
	defer s.WaitAndStop()
}
