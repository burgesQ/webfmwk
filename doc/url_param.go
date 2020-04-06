package main

import (
	"net/http"

	"github.com/burgesQ/webfmwk/v3"
)

// curl -X GET 127.0.0.1:4242/hello/world
// {"content":"hello world"}
func main() {
	// init server
	var s = webfmwk.InitServer()

	// expose /hello/name
	s.GET("/hello/{name}", func(c webfmwk.IContext) {
		c.JSONBlob(http.StatusOK, []byte(`{ "content": "hello `+c.GetVar("name")+`" }`))
	})

	// start asynchronously on :4242
	s.Start(":4242")

	// ctrl+c is handled internaly
	defer s.WaitAndStop()
}
