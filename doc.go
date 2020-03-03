/*
Package webfmwk implements an minimalist Go web framework

Example:
  package main

  import (
    w "github.com/burgesQ/webfmwk/v3"
  )
  // Handler
  func hello(c w.IContext) error {
    return c.JSONOk("Hello, World!")
  }

  func main() {
    // Echo instance
    s := w.InitServer(
     webfmwk.EnableCheckIsUp()
		// webfmwk.WithCORS(),
		webfmwk.WithLogger(log.GetLogger()),
		webfmwk.WithMiddlewars(
			middleware.Logging,
			middleware.Security),
		webfmwk.WithCustomContext(func(c *webfmwk.Context) webfmwk.IContext {
			return &server.CustomContext{
				Context:  *c,
				T:        tmpl,
				Chans:    x.chans,
				MaxEntry: maxEntry,
				RCList:   listRC,
				LCList:   listLC,
				StatAddr: stat,
			}
		}))

    // Routes
    s.GET("/hello", hello)

    // start server on :4242
    go func() {
      s.Start(":4242")
    }()

    // ctrl+c is handled internaly
    defer s.WaitAndStop()
  }

Learn more at https://github.com/burgesQ/webfmwk
*/
package webfmwk
