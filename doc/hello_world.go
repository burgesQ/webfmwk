package main

import (
	"net/http"

	"github.com/burgesQ/webfmwk/v3"
)

// curl -X GET 127.0.0.1:4242/hello
// { "message": "hello world" }
func main() {
	// create server
	s := webfmwk.InitServer()

	// expose /hello
	s.GET("/hello", func(c webfmwk.IContext) {
		c.JSONBlob(http.StatusOK, []byte(`{ "message": "hello world" }`))
	})

	// start asynchronously on :4242
	s.Start(":4242")

	// ctrl+c is handled internaly
	defer s.WaitAndStop()
}
