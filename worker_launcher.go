package webfmwk

import (
	"context"
	"net/http"
	"sync"
)

// WorkerLauncher hold the different workers
type WorkerLauncher struct {
	wg     *sync.WaitGroup
	cancel context.CancelFunc
}

// to use as a factory as the fields are unexported
func CreateWorkerLauncher(wg *sync.WaitGroup, cancel context.CancelFunc) WorkerLauncher {
	return WorkerLauncher{wg, cancel}
}

func (l *WorkerLauncher) run(name string, fn func() error) {
	logger.Debugf("%s: starting", name)

	if err := fn(); err != nil && err != http.ErrServerClosed {
		logger.Errorf("%s (%T): %s", name, err, err)
	} else {
		logger.Infof("%s: done", name)
	}

	l.cancel()
	l.wg.Done()
}

// launch a worker who will be wait & kill at the same time than the others
func (l *WorkerLauncher) Start(name string, fn func() error) {
	l.wg.Add(1)

	go l.run(name, fn)
}
