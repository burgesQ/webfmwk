package main

import (
	"net/http"

	w "github.com/burgesQ/webfmwk/v3"
)

func main() {
	// create server
	s := w.InitServer()

	s.GET("/hello/{id}", func(c w.IContext) {
		c.JSONBlob(http.StatusOK, []byte(`{ "id": "`+c.GetVar("id")+`" }`))
	})

	// start asynchronously on :4242
	go func() {
		s.Start(":4242")
	}()

	// ctrl+c is handled internaly
	defer s.WaitAndStop()
}
