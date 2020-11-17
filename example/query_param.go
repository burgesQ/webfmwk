package main

import (
	"github.com/burgesQ/webfmwk/v4"
)

// QueryParam
// see post_content.go for validation
type QueryParam struct {
	Pretty bool `json:"pretty" schema:"pretty"`
	Val    int  `schema:"val" json:"val" validate:"gte=1"`
}

// curl -i -X GET "127.0.0.1:4242/hello?pretty"
// {
//   "pretty": [
//     ""
// 		],
//    val: 0
// }
// curl -i -X GET "127.0.0.1:4242/hello?prete"
// {"pretty":false,"val":0}%
func query_param() {
	var s = webfmwk.InitServer()

	// expose /query_param
	s.GET("/hello", func(c webfmwk.Context) error {
		qp := &QueryParam{}

		if e := c.DecodeQP(qp); e != nil {
			return e
		}

		return c.JSONOk(qp)
	})

	// start asynchronously on :4242
	s.Start(":4242")

	// ctrl+c is handled internaly
	defer s.WaitAndStop()
}
