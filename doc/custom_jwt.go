package main

import (
	"fmt"

	"github.com/burgesQ/webfmwk/v3"
	"github.com/burgesQ/webfmwk/v3/jwt"
	"github.com/burgesQ/webfmwk/v3/log"
)

// GetLogger return a log.ILog interface
var logger = log.GetLogger()

type customContext struct {
	webfmwk.Context
	customVal string
}

// curl -X GET 127.0.0.1:4242/test
// {"error":"Missing Authorization Header"}
// curl -X GET "127.0.0.1:4242/test" -H "Authorization: Bearer invalid_value"
// {"error":"Forbidden"}
// curl -X GET "127.0.0.1:4242/test" -H "Authorization: Bearer `fetch token printed at server start`"
func main() {
	// to secure all endpoints use :
	// webfmwk.InitServer(webfmwk.WithHandlers(jwt.Handler))
	// or
	// webfmwk.InitServer(webfmwk.WithMiddlewares(jwt.Middleware))
	var (
		s        = webfmwk.InitServer()
		token, _ = jwt.GenToken("dev")
	)

	fmt.Printf("use %q as JWT token\n", token)

	// secure only that endpoint
	s.GET("/test", jwt.Handler(func(c webfmwk.IContext) {
		c.JSONOk("ok")
	}))

	// start asynchronously on :4242
	s.Start(":4242")

	s.WaitAndStop()
}
