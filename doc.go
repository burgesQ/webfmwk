/*
Package webfmwk implements a json API server ready

Hello world example:
	package main

	import (
		"github.com/burgesQ/webfmwk/v5"
	)

	// Handler
	func hello(c webfmwk.Context) error {
		return c.JSONOk("Hello, world!")
	}

	func main() {
		// Echo instance
		s := webfmwk.InitServer(webfmwk.WithCtrlC())

		// Routes
		s.GET("/hello", hello)

		// ctrl+c is handled internaly
		defer s.WaitAndStop()

		// start server on :4242
		s.Start(":4242")
	}

Some extra feature are available like : tls, custom handler/logger/context, redoc support and more ...
Find other examples at https://github.com/burgesQ/webfmwk/tree/master/example.
*/
package webfmwk
