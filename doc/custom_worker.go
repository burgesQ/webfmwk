package main

import (
	"time"

	w "github.com/burgesQ/webfmwk/v3"
	"github.com/burgesQ/webfmwk/v3/log"
)

func main() {
	log.SetLogLevel(log.LogDEBUG)

	// init server w/ ctrl+c support
	var (
		s  = w.InitServer()
		wl = s.GetLauncher()
	)

	s.GET("/test", func(c w.IContext) {
		c.JSONOk("ok")
	})

	wl.Start("custom worker", func() error {
		time.Sleep(10 * time.Second)
		log.Debugf("done")
		return nil
	})

	// start asynchronously on :4242
	go func() {
		s.Start(":4242")
	}()

	// ctrl+c is handled internaly
	defer s.WaitAndStop()
}
