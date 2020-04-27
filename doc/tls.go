package main

import (
	"github.com/burgesQ/webfmwk/v4"
)

// TODO: curl with HTTPS
func tls() {
	// init server w/ ctrl+c support
	var s = webfmwk.InitServer(webfmwk.WithCtrlC())

	// expose /test
	s.GET("/test", func(c webfmwk.Context) error {
		return c.JSONOk("ok")
	})

	// start asynchronously on :4242
	s.StartTLS(":4242", webfmwk.TLSConfig{
		Cert:     "/path/to/cert",
		Key:      "/path/to/key",
		Insecure: false,
	})

	// ctrl+c is handled internaly
	defer s.WaitAndStop()
}
