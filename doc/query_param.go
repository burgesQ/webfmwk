package main

import (
	"net/http"

	"github.com/burgesQ/webfmwk/v4"
)

// curl -i -X GET "127.0.0.1:4242/hello?pretty"
// {
//   "pretty": [
//     ""
// 		]
// }
// curl -i -X GET "127.0.0.1:4242/hello?prete"
// {"prete":[""]}%
func main() {
	var s = webfmwk.InitServer()

	// expose /hello
	s.GET("/hello", func(c webfmwk.IContext) {
		c.JSON(http.StatusOK, c.GetQueries())
	})

	// start asynchronously on :4242
	s.Start(":4242")

	// ctrl+c is handled internaly
	defer s.WaitAndStop()
}
