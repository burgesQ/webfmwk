package main

import (
	"net/http"

	"github.com/burgesQ/webfmwk/v4"
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
	var s = webfmwk.InitServer()

	s.POST("/hello", func(c webfmwk.IContext) {
		var out = Payload{}

		// process query params
		c.DecodeQP(&out.qp)
		c.Validate(out.qp)

		// process payload
		c.FetchContent(&out.content)
		c.Validate(out.content)

		c.JSON(http.StatusOK, out)
	})

	// start asynchronously on :4242
	s.Start(":4244")

	// ctrl+c is handled internaly
	defer s.WaitAndStop()
}
