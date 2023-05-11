package webfmwk

import (
	"context"
	"sync"
)

type (
	Worker         func()
	WorkerLauncher interface{ Start(Worker) }
)

// WorkerLauncher hold the different workers and wait for them to finish before exiting.
type launcher struct {
	wg     *sync.WaitGroup
	cancel context.CancelFunc
}

// CreateWorkerLauncher initialize and return a WorkerLauncher instance.
func CreateWorkerLauncher(wg *sync.WaitGroup, cancel context.CancelFunc) WorkerLauncher {
	return &launcher{wg, cancel}
}

// Start launch a worker job.
func (l *launcher) Start(fn Worker) {
	l.wg.Add(1)

	go func() {
		fn()
		l.cancel()
		l.wg.Done()
	}()
}
