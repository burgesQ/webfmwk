package main

import (
	"net/http"

	w "github.com/burgesQ/webfmwk/v3"
)

type (
	// Content hold the body of the request
	Content struct {
		Name string `schema:"name" json:"name" validate:"omitempty"`
		Age  int    `schema:"age" json:"age" validate:"gte=1"`
	}

	// QueryParam hold the query params
	QueryParam struct {
		PJSON bool `schema:"pjson" json:"pjson"`
		Val   int  `schema:"val" json:"val" validate:"gte=1"`
	}

	// Payload hold the output of the endpoint
	Payload struct {
		Content Content    `json:"content"`
		QP      QueryParam `json:"query_param"`
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
