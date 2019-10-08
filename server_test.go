package webfmwk

import (
	"net/http"
	"testing"

	"github.com/gorilla/mux"
)

var emptyController = func(c IContext) {}

func TestSetPrefix(t *testing.T) {
	s := InitServer(false)
	defer func() {
		s.Shutdown(*s.GetContext())
		s.WaitAndStop()
	}()

	s.SetPrefix("/api")
	s.GET("/test", emptyController)

	r := s.SetRouter()

	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, _ := route.GetPathTemplate()

		if pathTemplate != "/api/test" {
			t.Errorf("Router prefix wrongly applied : (%s)", pathTemplate)
		}
		return nil
	})
}

func TestGetLauncher(t *testing.T) {
	s := InitServer(false)
	defer func() {
		s.Shutdown(*s.GetContext())
		s.WaitAndStop()
	}()

	if s.GetLauncher() == nil {
		t.Errorf("Launcher wrongly created : %v.", s.launcher)
	}
}

func TestGetContext(t *testing.T) {
	s := InitServer(false)
	defer func() {
		s.Shutdown(*s.GetContext())
		s.WaitAndStop()
	}()

	if s.GetContext() == nil {
		t.Errorf("Context wrongly created : %v.", s.ctx)
	}
}

func TestAddMiddleware(t *testing.T) {
	s := InitServer(false)
	defer func() {
		s.Shutdown(*s.GetContext())
		s.WaitAndStop()
	}()

	s.AddMiddleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		})
	})

	if len(s.middlewares) != 1 {
		t.Errorf("Middleware wrongly saved : %v.", s.middlewares)
	}
}

func TestAddRoute(t *testing.T) {
	s := InitServer(false)
	defer func() {
		s.Shutdown(*s.GetContext())
		s.WaitAndStop()
	}()

	s.AddRoute(Route{
		Pattern: "/test/1",
		Method:  "GET",
		Handler: emptyController,
	})

	if !(s.routes[0].Pattern == "/test/1" && s.routes[0].Method == "GET") {
		t.Errorf("Routes wrongly saved : %v.", s.routes[0])
	}
}

func TestAddRoutes(t *testing.T) {
	s := InitServer(false)
	defer func() {
		s.Shutdown(*s.GetContext())
		s.WaitAndStop()
	}()

	s.AddRoutes([]Route{
		Route{
			Pattern: "/test/1",
			Method:  "GET",
			Handler: emptyController,
		},
		Route{
			Pattern: "/test/2",
			Method:  "GET",
			Handler: emptyController,
		},
	})

	if !(s.routes[0].Pattern == "/test/1" && s.routes[1].Pattern == "/test/2") {
		t.Errorf("Routes wrongly saved : %v.", s.routes[0])
	}
}

func TestGET(t *testing.T) {
	s := InitServer(false)
	defer func() {
		s.Shutdown(*s.GetContext())
		s.WaitAndStop()
	}()

	s.GET("/test", emptyController)

	if !(s.routes[0].Pattern == "/test" && s.routes[0].Method == "GET") {
		t.Errorf("Routes wrongly saved : %v.", s.routes[0])
	}
}

func TestDELETE(t *testing.T) {
	s := InitServer(false)
	defer func() {
		s.Shutdown(*s.GetContext())
		s.WaitAndStop()
	}()

	s.DELETE("/test", emptyController)

	if !(s.routes[0].Pattern == "/test" && s.routes[0].Method == "DELETE") {
		t.Errorf("Routes wrongly saved : %v.", s.routes[0])
	}
}

func TestPOST(t *testing.T) {
	s := InitServer(false)
	defer func() {
		s.Shutdown(*s.GetContext())
		s.WaitAndStop()
	}()

	s.POST("/test", emptyController)

	if !(s.routes[0].Pattern == "/test" && s.routes[0].Method == "POST") {
		t.Errorf("Routes wrongly saved : %v.", s.routes[0])
	}
}

func TestPUT(t *testing.T) {
	s := InitServer(false)
	defer func() {
		s.Shutdown(*s.GetContext())
		s.WaitAndStop()
	}()
	s.PUT("/test", emptyController)

	if !(s.routes[0].Pattern == "/test" && s.routes[0].Method == "PUT") {
		t.Errorf("Routes wrongly saved : %v.", s.routes[0])
	}
}

func TestPATCH(t *testing.T) {
	s := InitServer(false)
	defer func() {
		s.Shutdown(*s.GetContext())
		s.WaitAndStop()
	}()

	s.PATCH("/test", emptyController)

	if !(s.routes[0].Pattern == "/test" && s.routes[0].Method == "PATCH") {
		t.Errorf("Routes wrongly saved : %v.", s.routes[0])
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
		defer func() {
			s.Shutdown(*s.GetContext())
			s.WaitAndStop()
		}()
		if s.GetLauncher() == nil || s.GetContext() == nil {
			t.Errorf("Error while creating the server entity")
		}
	})
}
