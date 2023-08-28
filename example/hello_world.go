package main

import (
	"net/http"

	"github.com/burgesQ/webfmwk/v6"
)

// curl -X GET 127.0.0.1:4242/hello
// { "message": "hello world" }
func helloWorld() *webfmwk.Server {
	s := webfmwk.InitServer()

	// expose /hello
	s.GET("/hello", func(c webfmwk.Context) error {
		return c.JSONBlob(http.StatusOK, []byte(`{ "message": "hello world" }`))
	})

	return s
}
