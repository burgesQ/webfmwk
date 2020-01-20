package main

import (
	w "github.com/burgesQ/webfmwk/v3"
)

type customContext struct {
	w.Context
	customVal string
}

func main() {
	// init server w/ ctrl+c support
	s := w.InitServer(true)

	s.SetCustomContext(func(c *w.Context) w.IContext {
		ctx := &customContext{*c, "42"}
		return ctx
	})

	s.GET("/test", func(c w.IContext) {
		ctx := c.(*customContext)
		c.JSONOk(ctx.customVal)
	})

	// start asynchronously on :4242
	go func() {
		s.Start(":4244")
	}()

	// ctrl+c is handled internaly
	defer s.WaitAndStop()
}
