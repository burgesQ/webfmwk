package main

import (
	w "github.com/burgesQ/webfmwk"
	"github.com/burgesQ/webfmwk/log"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Answer struct {
	message string `json:"message"`
}

// @Summary hello world
// @Description Return a simple greeting
// @Param pjson query bool false "return a pretty JSON"
// @Success 200 {object} db.Reply
// @Produce application/json
// @Router /hello [get]
func hello(c w.IContext) error {
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
func main() {

	// init logging
	log.SetLogLevel(log.LogDebug)
	log.Init(log.LoggerSTDOUT | log.LogFormatLong)

	// init server w/ ctrl+c support
	s := w.InitServer(true)

	s.RegisterDocHandler(httpSwagger.WrapHandler)

	s.GET("/test", func(c w.IContext) error {
		return c.JSONOk("ok")
	})

	// start asynchronously on :4242
	go func() {
		s.Start(":4242")
	}()

	// ctrl+c is handled internaly
	defer s.WaitAndStop()
}
