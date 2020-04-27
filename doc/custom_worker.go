package main

import (
	"time"

	"github.com/burgesQ/webfmwk/v4"
	"github.com/burgesQ/webfmwk/v4/log"
)

func custom_worker() {
	var (
		s  = webfmwk.InitServer()
		wl = s.GetLauncher()
	)

	// register /test
	s.GET("/test", func(c webfmwk.Context) error {
		return c.JSONOk("ok")
	})

	// register extra eorker
	wl.Start("custom worker", func() error {
		time.Sleep(10 * time.Second)
		log.Debugf("done")
		return nil
	})

	// start asynchronously on :4242
	s.Start(":4242")

	// ctrl+c is handled internaly
	defer s.WaitAndStop()
}
