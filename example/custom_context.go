package main

import "github.com/burgesQ/webfmwk/v4"

// customContext extend the webfmwk.Context
type Context struct {
	webfmwk.Context
	val string
}

func loadContext(next webfmwk.HandlerFunc) webfmwk.HandlerFunc {
	return webfmwk.HandlerFunc(func(c webfmwk.Context) error {
		cc := Context{c, "val"}
		return next(cc)
	})
}

// To use a custom context, create a struct that
// extend webfmwk.Context.
// Then, load it the earliest possible in the handler call chain
//
// curl -X GET 127.0.0.1:4242/test
// {"content":"42"}
func custom_context() {
	// init server w/ ctrl+c support and custom context options
	var s = webfmwk.InitServer(
		webfmwk.WithCtrlC(),
		webfmwk.WithHandlers(loadContext),
	)

	// expose /test
	s.GET("/test", func(c webfmwk.Context) error {
		return c.JSONOk(webfmwk.NewResponse(c.(*Context).val))
	})

	// start asynchronously on :4242
	s.Start(":4242")

	// ctrl+c is handled internaly
	defer s.WaitAndStop()
}
