package webfmwk

import (
	"context"
	"errors"
	"net/http"
	"sync"
)

// WorkerLauncher hold the different workers and wait for them to finish before exiting.
type WorkerLauncher struct {
	wg     *sync.WaitGroup
	cancel context.CancelFunc
}

// CreateWorkerLauncher initialize and return a WorkerLauncher instance.
func CreateWorkerLauncher(wg *sync.WaitGroup, cancel context.CancelFunc) WorkerLauncher {
	return WorkerLauncher{wg, cancel}
}

// Start launch a worker job.
func (l *WorkerLauncher) Start(name string, fn func() error) {
	l.wg.Add(1)

	go l.run(name, fn)
}

func (l *WorkerLauncher) run(name string, fn func() error) {
	logger.Debugf("%s: starting", name)

	if err := fn(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Errorf("%s (%T): %s", name, err, err)
	} else {
		loggerMu.Lock()
		logger.Infof("%s: done", name)
		loggerMu.Unlock()
	}

	l.cancel()
	l.wg.Done()
}
