package util

import (
	"context"
	"os"
	"os/signal"

	"github.com/burgesQ/webfmwk/log"
)

// Clean exit from infinit loop
func ExitHandler(ctx context.Context, sig ...os.Signal) {
	c := make(chan os.Signal)
	signal.Notify(c, sig...)

	select {
	case <-ctx.Done():
		log.Infof("Context canceled, tchao!")
		return
	case s := <-c:
		log.Infof("captured %v, exiting...", s)
		ctx.Done()
		return
	}
}
