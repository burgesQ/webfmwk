package main

import (
	"net/http"

	"github.com/burgesQ/webfmwk/v5"
	"github.com/burgesQ/webfmwk/v5/handler/logging"
)

// curl -X GET 127.0.0.1:4242/hello
// { "message": "hello world" }
func logMe() *webfmwk.Server {
	s := webfmwk.InitServer(webfmwk.WithHandlers(logging.Handler))

	// expose /hello
	s.GET("/logMe", func(c webfmwk.Context) error {
		c.GetLogger().Errorf("yo")

		return c.JSONBlob(http.StatusOK, []byte(`{ "message": "hello world" }`))
	})

	return s
}
