/*
Package webfmwk implements an minimalist Go web framework

Example:
  package main

  import (
    w "github.com/burgesQ/webfmwk/v4"
  )

  type Context struct {
    webfmwk.IContext
    content string
  }

  // Handler
  func hello(c w.IContext) error {
    return c.JSONOk("Hello, World!")
  }

  func main() {
    // Echo instance
    s := w.InitServer(
     webfmwk.EnableCheckIsUp()
		webfmwk.WithCORS(),
		webfmwk.WithLogger(log.GetLogger()),
		webfmwk.WithMiddlewars(
			middleware.Logging,
			middleware.Security),
		webfmwk.WithCustomContext(func(c *webfmwk.Context) webfmwk.IContext {
			return &Context{
				Context:  *c,
				content: "testing",
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
