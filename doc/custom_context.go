package main

import "github.com/burgesQ/webfmwk/v3"

// customContext extend the webfmwk.Context
type customContext struct {
	webfmwk.Context
	val string
}

// curl -X GET 127.0.0.1:4242/test
// {"content":"42"}
func main() {
	// init server w/ ctrl+c support and custom context options
	var s = webfmwk.InitServer(
		webfmwk.WithCtrlC(),
		webfmwk.WithCustomContext(func(c *webfmwk.Context) webfmwk.IContext {
			return &customContext{*c, "42"}
		}))

	// expose /test
	s.GET("/test", func(c webfmwk.IContext) {
		c.JSONOk(webfmwk.NewResponse(c.(*customContext).val))
	})

	// start asynchronously on :4242
	s.Start(":4242")

	// ctrl+c is handled internaly
	defer s.WaitAndStop()
}
