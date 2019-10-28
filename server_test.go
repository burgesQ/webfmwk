package webfmwk

import (
	"net/http"
	"testing"
)

func TestGetLauncher(t *testing.T) {
	s := InitServer(false)

	defer func(s Server) {
		s.Shutdown(*s.GetContext())
		s.WaitAndStop()
	}(s)

	if s.GetLauncher() == nil {
		t.Errorf("Launcher wrongly created : %v.", s.launcher)
	}
}

func TestGetContext(t *testing.T) {
	s := InitServer(false)

	defer func(s Server) {
		s.Shutdown(*s.GetContext())
		s.WaitAndStop()
	}(s)

	if s.GetContext() == nil {
		t.Errorf("Context wrongly created : %v.", s.ctx)
	}
}

func TestAddMiddleware(t *testing.T) {
	s := InitServer(false)

	defer func(s Server) {
		s.Shutdown(*s.GetContext())
		s.WaitAndStop()
	}(s)

	s.AddMiddleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		})
	})

	if len(s.middlewares) != 1 {
		t.Errorf("Middleware wrongly saved : %v.", s.middlewares)
	}
}

// // TODO: TestStartTLS(t *testing.T)
// // TODO: TestStart
// // TODO: TestShutDown
// // TODO: TestWaitAndStop
// // TODO: TestExitHandler

func TestInitServer(t *testing.T) {
	t.Run("simple init server", func(t *testing.T) {
		s := InitServer(false)

		defer func(s Server) {
			s.Shutdown(*s.GetContext())
			s.WaitAndStop()
		}(s)

		if s.GetLauncher() == nil || s.GetContext() == nil {
			t.Errorf("Error while creating the server entity")
		}
	})
}
