package main

import (
	"github.com/burgesQ/webfmwk/v3"
)

// TODO: curl with HTTPS
func main() {
	// init server w/ ctrl+c support
	var s = webfmwk.InitServer(webfmwk.WithCtrlC())

	// expose /test
	s.GET("/test", func(c webfmwk.IContext) error {
		c.JSONOk("ok")
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
