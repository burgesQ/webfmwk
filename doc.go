/*
Package webfmwk implements an minimalist Go web framework

Example:
  package main

  import (
	  w "github.com/burgesQ/webfmwk/v2"
  )
  // Handler
  func hello(c w.IContext) error {
	  return c.JSONOk("Hello, World!")
  }

  func main() {
	  // Echo instance
	  s := w.InitServer(true)

  	// Routes
	  s.GET("/hello", hello)

    // start server on :4242
	  go func() {
		  s.Start(":4242")
	  }()

  	// ctrl+c is handled internaly
	  defer s.WaitAndStop()
  }

Learn more at https://github.com/burgesQ/echo
*/
package webfmwk
