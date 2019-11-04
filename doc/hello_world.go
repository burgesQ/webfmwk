package main

import (
	"net/http"

	w "github.com/burgesQ/webfmwk/v2"
)

func main() {
	// create server
	s := w.InitServer(true)

	s.GET("/hello", func(c w.IContext) {
		c.JSONBlob(http.StatusOK, []byte(`{ "message": "hello world" }`))
	})

	// start asynchronously on :4242
	go func() {
		s.Start(":4242")
	}()

	// ctrl+c is handled internaly
	defer s.WaitAndStop()
}
