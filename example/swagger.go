package main

import (
	"github.com/burgesQ/webfmwk/v5"
	"github.com/burgesQ/webfmwk/v5/docs"
)

// TODO: form & payload schema
// TODO: for & payload validation?

// Answer implement the http response
type Answer struct {
	Message string `json:"message"`
}

// @Summary hello world
// @Description Return a simple greeting
// @Param pjson query bool false "return a pretty JSON"
// @Success 200 {object} db.Reply
// @Produce application/json
// @Router /hello [get]
func hello(c webfmwk.Context) error {
	return c.JSONOk(Answer{"ok"})
}

// @title hello world API
// @version 1.0
// @description This is an simple API
// @termsOfService https://www.youtube.com/watch?v=DLzxrzFCyOs
// @contact.name Quentin Burgess
// @contact.url github.com/burgesQ
// @contact.email quentin@frafos.com
// @license.name GFO
// @host localhost:4242
func swagger() *webfmwk.Server {
	// init server w/ ctrl+c support, prefix and APIDoc.
	// register prefix BEFORE api doc.
	s := webfmwk.InitServer(
		webfmwk.SetPrefix("/api"),
		// webfmwk.WithDocHandlers(httpSwagger.WrapHandler),
		webfmwk.WithDocHandlers(docs.GetRedocHandler(nil)),

		webfmwk.WithCtrlC())

	// register /test
	s.GET("/hello", hello)

	return s
}
