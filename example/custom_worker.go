package main

import (
	"time"

	"github.com/burgesQ/gommon/log"
	"github.com/burgesQ/webfmwk/v5"
)

func customWorker() *webfmwk.Server {
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
	return s
}
