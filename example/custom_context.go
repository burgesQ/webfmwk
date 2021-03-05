package main

import "github.com/burgesQ/webfmwk/v5"

// customContext extend the webfmwk.Context
type ctx struct {
	webfmwk.Context
	val string
}

func loadContext(next webfmwk.HandlerFunc) webfmwk.HandlerFunc {
	return webfmwk.HandlerFunc(func(c webfmwk.Context) error {
		cc := ctx{c, "val"}
		return next(cc)
	})
}

// To use a custom context, create a struct that
// extend webfmwk.Context.
// Then, load it the earliest possible in the handler call chain
//
// curl -X GET 127.0.0.1:4242/test
// {"content":"42"}
func customContext() *webfmwk.Server {
	// init server w/ ctrl+c support and custom context options
	var s = webfmwk.InitServer(
		webfmwk.WithCtrlC(),
		webfmwk.WithHandlers(loadContext),
	)

	// expose /test
	s.GET("/test", func(c webfmwk.Context) error {
		return c.JSONOk(webfmwk.NewResponse(c.(*ctx).val))
	})

	return s
}
