package webfmwk

import (
	"context"
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
	logger.Debugf("%s: starting", name)
	go func() {
		if err := fn(); err != nil {
			logger.Errorf("%s: %s", name, err)
		} else {
			logger.Infof("%s: done", name)
		}
		l.cancel()
		l.wg.Done()
	}()
}
