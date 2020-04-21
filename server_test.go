package webfmwk

import (
	"net/http"
	"testing"
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

func TestAddMiddleware(t *testing.T) {
	s := InitServer(CheckIsUp())
	defer stopServer(t, s)

	s.addMiddlewares(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		})
	})

	if len(s.meta.middlewares) != 1 {
		t.Errorf("Middleware wrongly saved : %v.", s.meta.middlewares)
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

	if len(s.meta.handlers) != 1 {
		t.Errorf("Middleware wrongly saved : %v.", s.meta.handlers)
	}
}

// // TODO: TestStartTLS(t *testing.T)
// // TODO: TestStart
// // TODO: TestShutDown
// // TODO: TestWaitAndStop
// // TODO: TestExitHandler

func TestInitServer(t *testing.T) {
	t.Run("simple init server", func(t *testing.T) {
		s := InitServer(CheckIsUp())

		defer stopServer(t, s)

		if s.GetLauncher() == nil || s.GetContext() == nil {
			t.Errorf("Error while creating the server entity")
		}
	})
}
