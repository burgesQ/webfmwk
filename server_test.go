package webfmwk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLauncher(t *testing.T) {
	s := InitServer(CheckIsUp())

	defer stopServer(t, s)
	if s.GetLauncher() == nil {
		t.Errorf("Launcher wrongly created : %v.", s.launcher)
	}
}

func TestGetContext(t *testing.T) {
	s := InitServer(CheckIsUp())

	defer stopServer(t, s)

	if s.GetContext() == nil {
		t.Errorf("Context wrongly created : %v.", s.ctx)
	}
}

func TestAddHandlers(t *testing.T) {
	s := InitServer(CheckIsUp())
	defer stopServer(t, s)

	s.addHandlers(func(next HandlerFunc) HandlerFunc {
		return HandlerFunc(func(c Context) error {
			return nil
		})
	})

	assert.True(t, len(s.meta.handlers) == 1, "handler wrongly saved")
}

// // // TODO: TestStartTLS(t *testing.T)
// // // TODO: TestStart
// // // TODO: TestShutDown
// // // TODO: TestWaitAndStop
// // // TODO: TestExitHandler
