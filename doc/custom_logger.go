package main

import (
	w "github.com/burgesQ/webfmwk/v3"
	"github.com/burgesQ/webfmwk/v3/log"
)

// GetLogger return a log.ILog interface
var logger = log.GetLogger()

func main() {
	var s = w.InitServer(webfmwk.WithLogger(logger))

	s.GET("/test", func(c w.IContext) error {
		return c.JSONOk("ok")
	})

	// start asynchronously on :4242
	go func() {
		s.StartTLS(":4242", TLSConfig{
			Cert:     "/path/to/cert",
			Key:      "/path/to/key",
			Insecure: true,
		})
	}()

	// ctrl+c is handled internaly
	defer s.WaitAndStop()
}
