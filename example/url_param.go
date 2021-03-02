package main

import (
	"net/http"

	"github.com/burgesQ/webfmwk/v5"
)

// curl -X GET 127.0.0.1:4242/hello/world
// {"content":"hello world"}
// curl -X GET 127.0.0.1:4242/acticles/abc
// 404
// curl -X GET 127.0.0.1:4242/acticles/
// 404
// curl -X GET 127.0.0.1:4242/acticles/0
// {"content":"is is 0"}
// for more see https://pkg.go.dev/github.com/fasthttp/router
func url_param() *webfmwk.Server {
	// init server
	var s = webfmwk.InitServer()

	// expose /hello/name
	s.GET("/hello/{name}", func(c webfmwk.Context) error {
		return c.JSONBlob(http.StatusOK, []byte(`{ "content": "hello `+c.GetVar("name")+`" }`))
	})

	// expose /acticles/01
	s.GET("/acticles/{id:[0-9]+}", func(c webfmwk.Context) error {
		return c.JSONBlob(http.StatusOK, []byte(`{ "content": "id is `+c.GetVar("id")+`" }`))
	})

	return s
}
