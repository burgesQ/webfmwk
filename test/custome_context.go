package main

import (
	w "github.com/burgesQ/webfmwk"
	"github.com/burgesQ/webfmwk/log"
)

type customContext struct {
	w.Context
	customVal string
}

func main() {

	// init logging
	log.SetLogLevel(log.LOG_DEBUG)
	log.Init(log.LOGGER_STDOUT | log.LOGFORMAT_LONG)

	// init server w/ ctrl+c support
	s := w.InitServer(true)

	s.SetCustomContext(func(c *w.Context) w.IContext {
		ctx := &customContext{*c, "42"}
		return ctx
	})

	s.GET("/test", func(c w.IContext) error {
		ctx := c.(*custom Context)
		return c.JSONOk(ctx.customVal)
	})

	// start asynchronously on :4242
	go func() {
		s.Start(":4242")
	}()

	// ctrl+c is handled internaly
	defer s.WaitAndStop()
}
