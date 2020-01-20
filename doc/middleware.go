package main

import (
	w "github.com/burgesQ/webfmwk/v3"
	m "github.com/burgesQ/webfmwk/v3/middleware"
)

func main() {

	// init server w/ ctrl+c support
	s := w.InitServer(true)

	s.AddMiddleware(m.Logging)
	s.AddMiddleware(m.Security)

	s.GET("/test", func(c w.IContext) error {
		return c.JSONOk("ok")
	})

	// start asynchronously on :4242
	go func() {
		s.Start(":4242")
	}()

	// ctrl+c is handled internaly
	defer s.WaitAndStop()
}
