package main

import (
	"github.com/burgesQ/webfmwk/v4"
)

var (
	_routes = webfmwk.RoutesPerPrefix{
		"/api/v1": {
			{
				Verbe: "GET",
				Path:  "/test",
				Name:  "test v1",
				Handler: func(c webfmwk.Context) error {
					return c.JSONOk("v1 ok")
				},
			},
		},
		"/api/v2": {
			{
				Verbe: "GET",
				Path:  "/test",
				Name:  "test v2",
				Handler: func(c webfmwk.Context) error {
					return c.JSONOk("v2 ok")
				},
			},
		},
	}
)

// curl -X GET 127.0.0.1:4242/api/v1/test
// "v1 ok"
// curl -X GET 127.0.0.1:4242/api/v2/test
// "v2 ok"
func routes() {
	var s = webfmwk.InitServer()

	// register routes object
	s.RouteApplier(_routes)

	// start asynchronously on :4242
	s.Start(":4242")

	// ctrl+c is handled internaly
	defer s.WaitAndStop()
}
