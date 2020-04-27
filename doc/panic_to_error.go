package main

import (
	"github.com/burgesQ/webfmwk/v4"
	"github.com/burgesQ/webfmwk/v4/handler"
)

func panic_to_error() {
	var s = webfmwk.InitServer(
		webfmwk.WithCtrlC(),
		webfmwk.WithHandlers(handler.Recover, handler.Logging, handler.RequestID),
	)

	// expose /no_panic
	s.GET("/no_panic", func(c webfmwk.Context) error {
		return webfmwk.NewErrorHandled(500, webfmwk.NewAnonymousError("error by return"))
	})

	// expose /no_panic
	s.GET("/panic", func(c webfmwk.Context) error {
		panic(webfmwk.NewErrorHandled(422, webfmwk.NewAnonymousError("user not logged")))
		return nil
	})

	// start asynchronously on :4242
	s.Start(":4242")

	// ctrl+c is handled internaly
	defer s.WaitAndStop()
}
