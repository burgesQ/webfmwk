package webfmwk

import (
	"net/http"
	"strings"
	"testing"
	"time"

	z "github.com/burgesQ/webfmwk/testing"
	"github.com/gorilla/mux"
)

func TestSetRouter(t *testing.T) {
	s := InitServer(false)
	defer s.WaitAndStop()
	defer s.Shutdown(*s.GetContext())

	s.SetPrefix("/api")
	s.GET("/test", func(c IContext) error { return nil })

	r := s.SetRouter()

	if err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		var (
			pathTemplate, _ = route.GetPathTemplate()
			methods, _      = route.GetMethods()
			method          = strings.Join(methods, ",")
		)

		if method != "GET" || pathTemplate != "/api/test" {
			t.Errorf("Router Routing wrongly created : [%s](%s)", methods, pathTemplate)
		}
		return nil
	}); err != nil {
		t.Errorf("Router wrongly created : %s", err.Error())
	}
}

func TestHandleParam(t *testing.T) {
	s := InitServer(false)
	defer s.WaitAndStop()
	defer s.Shutdown(*s.GetContext())

	s.GET("/test/{id}", func(c IContext) error {

		if val, ok := c.GetQuery("pjson"); !ok || val != "1" {
			t.Errorf("query Param wrongly decoded %s", val)
		} else if c.GetVar("id") != "toto" {
			t.Errorf("URL Param wrongly decoded")
		}

		return nil
	})

	go func() {
		if e := s.Start(":4242"); e != nil {
			t.Fatalf("error while booting the server : %s", e.Error())
		}
	}()
	time.Sleep(50 * time.Millisecond)

	z.RequestAndTestAPI(t, "/test/toto?pjson=1", func(t *testing.T, resp *http.Response) {})
}
