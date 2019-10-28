package main

import (
	"github.com/burgesQ/webfmwk/v2"
)

var (
	routes = webfmwk.RoutesPerPrefix{
		"/api/v1": {
			{
				Verbe: "GET",
				Path:  "/test",
				Name:  "test v1",
				Handler: func(c webfmwk.IContext) {
					c.JSONOk("v1 ok")
				},
			},
		},
		"/api/v2": {
			{
				Verbe: "GET",
				Path:  "/test",
				Name:  "test v2",
				Handler: func(c webfmwk.IContext) {
					c.JSONOk("v2 ok")
				},
			},
		},
	}
)

func main() {

	s := webfmwk.InitServer(true)

	s.RouteApplier(routes)

	// start asynchronously on :4242
	go func() {
		s.Start(":4242")
	}()

	// ctrl+c is handled internaly
	defer s.WaitAndStop()

}
