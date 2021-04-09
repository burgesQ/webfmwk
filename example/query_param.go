package main

import (
	"github.com/burgesQ/webfmwk/v5"
)

// QueryParam
// see post_content.go for validation
type qp struct {
	Pretty bool   `json:"pretty" schema:"pretty"`
	Val    int    `schema:"val" json:"val" validate:"gte=1"`
	Smth   string `schema:"else" json:"else"`
}

// curl -i -X GET "127.0.0.1:4242/hello?pretty"
// {
//   "pretty": [
//     ""
//		],
//    val: 0
// }
// curl -i -X GET "127.0.0.1:4242/hello?prete"
// {"pretty":false,"val":0}%
func queryParam() *webfmwk.Server {
	var s = webfmwk.InitServer()

	// expose /query_param
	s.GET("/hello", func(c webfmwk.Context) error {
		qp := &qp{}

		if e := c.DecodeQP(qp); e != nil {
			return e
		}

		return c.JSONOk(qp)
	})

	return s
}
