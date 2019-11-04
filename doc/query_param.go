package main

import (
	"net/http"

	w "github.com/burgesQ/webfmwk/v2"
	"github.com/burgesQ/webfmwk/v2/log"
)

func main() {
	// create server
	s := w.InitServer(true)

	s.GET("/hello", func(c w.IContext) {
		var (
			queries   = c.GetQueries()
			pjson, ok = c.GetQuery("pjson")
		)
		if ok {
			log.Errorf("%#v", pjson)
		}
		c.JSON(http.StatusOK, queries)
	})

	// start asynchronously on :4242
	go func() {
		s.Start(":4242")
	}()

	// ctrl+c is handled internaly
	defer s.WaitAndStop()
}
