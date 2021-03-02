package main

import (
	"crypto/rand"
	"math/big"

	"github.com/burgesQ/webfmwk/v5"
	"github.com/burgesQ/webfmwk/v5/handler"
)

func panic_to_error() *webfmwk.Server {
	var (
		s = webfmwk.InitServer(
			webfmwk.WithCtrlC(),
			webfmwk.WithHandlers(handler.Recover, handler.Logging),
		)
	)
	// expose /panic
	s.GET("/no_panic", func(c webfmwk.Context) error {
		if n, _ := rand.Int(rand.Reader, big.NewInt(1000)); n.Mod(n, big.NewInt(2)) == big.NewInt(0) {
			panic(webfmwk.NewErrorHandled(422, webfmwk.NewError("error by panic")))
		}
		return webfmwk.NewErrorHandled(500, webfmwk.NewError("error by return"))
	})

	return s
}
