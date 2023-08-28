package main

import (
	"github.com/burgesQ/webfmwk/v6"
)

// TODO: curl with HTTPS
func tls() *webfmwk.Server {
	// init server w/ ctrl+c support
	var s = webfmwk.InitServer(webfmwk.WithCtrlC())

	// expose /test
	s.GET("/test", func(c webfmwk.Context) error {
		return c.JSONOk("ok")
	})

	return s
}
