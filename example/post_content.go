package main

import (
	"net/http"

	"github.com/burgesQ/webfmwk/v5"
)

type (
	// Content hold the body of the request
	Content struct {
		Name string `schema:"name" json:"name" validate:"omitempty"`
		Age  int    `schema:"age" json:"age" validate:"gte=1"`
	}

	// Payload hold the output of the endpoint
	// QueryParam is imported from query_param.go file
	Payload struct {
		Content    Content    `json:"content"`
		QueryParam QueryParam `json:"query_param"`
	}
)

func postContent() *webfmwk.Server {
	var s = webfmwk.InitServer()

	s.POST("/post", func(c webfmwk.Context) error {
		var out = Payload{}

		// process query params
		if e := c.DecodeQP(&out.QueryParam); e != nil {
			return e
		} else if e := c.Validate(out.QueryParam); e != nil {
			return e
		}

		// process payload
		if e := c.FetchContent(&out.Content); e != nil {
			return e
		} else if e := c.Validate(out.Content); e != nil {
			return e
		}

		return c.JSON(http.StatusOK, out)
	})

	return s
}
