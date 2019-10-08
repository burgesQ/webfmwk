package main

import (
	"net/http"

	w "github.com/burgesQ/webfmwk/v2"
)

type Content struct {
	Name string `schema:"name" json:"name" validate:"omitempty"`
	Age  int    `schema:"age" json:"age" vallidate:"gte=1"`
}

func main() {
	// create server
	s := w.InitServer(true)

	s.POST("/hello", func(c w.IContext) {
		data := Content{}
		c.FetchContent(&data)
		c.Validate(data)
		c.JSON(http.StatusOK, data)
	})

	// start asynchronously on :4242
	go func() {
		s.Start(":4242")
	}()

	// ctrl+c is handled internaly
	defer s.WaitAndStop()
}
