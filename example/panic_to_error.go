package main

import (
	"crypto/rand"
	"math/big"

	"github.com/burgesQ/webfmwk/v4"
	"github.com/burgesQ/webfmwk/v4/handler"
)

func panic_to_error() {
	var (
		s = webfmwk.InitServer(
			webfmwk.WithCtrlC(),
			webfmwk.WithHandlers(handler.Recover, handler.Logging, handler.RequestID),
		)
	)
	// expose /panic
	s.GET("/no_panic", func(c webfmwk.Context) error {
		if n, _ := rand.Int(rand.Reader, big.NewInt(1000)); n.Mod(n, big.NewInt(2)) == big.NewInt(0) {
			panic(webfmwk.NewErrorHandled(422, webfmwk.NewError("error by panic")))
		}
		return webfmwk.NewErrorHandled(500, webfmwk.NewError("error by return"))
	})

	// start asynchronously on :4242
	s.Start(":4242")

	// ctrl+c is handled internaly
	defer s.WaitAndStop()
}
