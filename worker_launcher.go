package webfmwk

import (
	"context"
	"net/http"
	"sync"
)

type WorkerLauncher struct {
	wg     *sync.WaitGroup
	cancel context.CancelFunc
}

// to use as a factory as the fields are unexported
func CreateWorkerLauncher(wg *sync.WaitGroup, cancel context.CancelFunc) WorkerLauncher {
	return WorkerLauncher{wg, cancel}
}

// launch a worker who will be wait & kill at the same time than the others
func (l *WorkerLauncher) Start(name string, fn func() error) {
	l.wg.Add(1)
	go func(n string) {
		logger.Debugf("%s: starting", n)
		if err := fn(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("%s (%T): %s", n, err, err)
		} else {
			logger.Infof("%s: done", n)
		}
		l.cancel()
		l.wg.Done()
	}(name)
}
