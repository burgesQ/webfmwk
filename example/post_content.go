package main

import (
	"net/http"

	w "github.com/burgesQ/webfmwk/v2"
)

type (
	Content struct {
		Name string `schema:"name" json:"name" validate:"omitempty"`
		Age  int    `schema:"age" json:"age" validate:"gte=1"`
	}

	QueryParam struct {
		pjson bool `schema:"pjson" json:"pjson"`
		ok    int  `schema:"pjson" json:"pjson" validate:"gte=1"`
	}

	Payload struct {
		content Content    `json:"content"`
		qp      QueryParam `json:"query_param"`
	}
)

func main() {
	// create server
	s := w.InitServer(true)

	s.POST("/hello", func(c w.IContext) {

		out := Payload{}

		c.FetchContent(&out.content)
		c.Validate(out.content)

		c.DecodeQP(&out.qp)
		c.Validate(out.qp)

		c.JSON(http.StatusOK, out)
	})

	// start asynchronously on :4242
	go func() {
		s.Start(":4244")
	}()

	// ctrl+c is handled internaly
	defer s.WaitAndStop()
}
