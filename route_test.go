package webfmwk

import (
	"net/http"
	"strings"
	"testing"
	"time"

	z "github.com/burgesQ/webfmwk/v2/testing"
	"github.com/gorilla/mux"
)

const (
	_testPrefix = "/api"
	_testURL    = "/test"
	_testURI    = _testPrefix + _testURL
	_testURI2   = _testPrefix + _testURL + "/2"
	_testVerbe  = GET
)

var _emptyController = func(c IContext) {}

// TODO: func TestAddRoute(t *testing.T)  {}
// TODO: func TestAddRoutes(t *testing.T) {}

func TestSetPrefix(t *testing.T) {
	s := InitServer(false)

	defer func(s Server) {
		s.Shutdown(*s.GetContext())
		s.WaitAndStop()
	}(s)

	s.SetPrefix(_testPrefix)
	s.GET(_testURL, _emptyController)

	r := s.SetRouter()

	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, _ := route.GetPathTemplate()

		if !(pathTemplate == _testURI || pathTemplate == _testPrefix) {
			t.Errorf("route wrongly created : [%s]", pathTemplate)
		}

		return nil
	})
}

func TestAddRoute(t *testing.T) {
	s := InitServer(false)

	defer func(s Server) {
		s.Shutdown(*s.GetContext())
		s.WaitAndStop()
	}(s)

	s.AddRoute(Route{
		Path:    _testURI,
		Verbe:   _testVerbe,
		Handler: _emptyController,
	})

	z.AssertStringEqual(t, s.routes[s.prefix][0].Path, _testURI)
	z.AssertStringEqual(t, s.routes[s.prefix][0].Verbe, _testVerbe)
}

func TestAddRoutes(t *testing.T) {
	s := InitServer(false)

	defer func(s Server) {
		s.Shutdown(*s.GetContext())
		s.WaitAndStop()
	}(s)

	s.AddRoutes(Routes{
		{
			Path:    _testURI,
			Verbe:   _testVerbe,
			Handler: _emptyController,
		},
		{
			Path:    _testURI2,
			Verbe:   _testVerbe,
			Handler: _emptyController,
		},
	})

	z.AssertStringEqual(t, s.routes[s.prefix][0].Path, _testURI)
	z.AssertStringEqual(t, s.routes[s.prefix][0].Verbe, _testVerbe)
	z.AssertStringEqual(t, s.routes[s.prefix][1].Path, _testURI2)
	z.AssertStringEqual(t, s.routes[s.prefix][1].Verbe, _testVerbe)
}

func TestGET(t *testing.T) {
	s := InitServer(false)

	defer func(s Server) {
		s.Shutdown(*s.GetContext())
		s.WaitAndStop()
	}(s)

	s.GET(_testURL, _emptyController)

	z.AssertStringEqual(t, s.routes[s.prefix][0].Path, _testURL)
	z.AssertStringEqual(t, s.routes[s.prefix][0].Verbe, GET)
}

func TestDELETE(t *testing.T) {
	s := InitServer(false)

	defer func(s Server) {
		s.Shutdown(*s.GetContext())
		s.WaitAndStop()
	}(s)

	s.DELETE(_testURL, _emptyController)

	z.AssertStringEqual(t, s.routes[s.prefix][0].Path, _testURL)
	z.AssertStringEqual(t, s.routes[s.prefix][0].Verbe, DELETE)
}

func TestPOST(t *testing.T) {
	s := InitServer(false)

	defer func(s Server) {
		s.Shutdown(*s.GetContext())
		s.WaitAndStop()
	}(s)

	s.POST(_testURL, _emptyController)

	z.AssertStringEqual(t, s.routes[s.prefix][0].Path, _testURL)
	z.AssertStringEqual(t, s.routes[s.prefix][0].Verbe, POST)
}

func TestPUT(t *testing.T) {
	s := InitServer(false)

	defer func(s Server) {
		s.Shutdown(*s.GetContext())
		s.WaitAndStop()
	}(s)

	s.PUT(_testURL, _emptyController)

	z.AssertStringEqual(t, s.routes[s.prefix][0].Path, _testURL)
	z.AssertStringEqual(t, s.routes[s.prefix][0].Verbe, PUT)
}

func TestPATCH(t *testing.T) {
	s := InitServer(false)

	defer func(s Server) {
		s.Shutdown(*s.GetContext())
		s.WaitAndStop()
	}(s)

	s.PATCH(_testURL, _emptyController)

	z.AssertStringEqual(t, s.routes[s.prefix][0].Path, _testURL)
	z.AssertStringEqual(t, s.routes[s.prefix][0].Verbe, PATCH)
}

func TestSetRouter(t *testing.T) {
	s := InitServer(false)
	defer s.WaitAndStop()
	defer s.Shutdown(*s.GetContext())

	s.SetPrefix(_testPrefix)
	s.GET(_testURL, func(c IContext) { c.JSONNoContent() })

	r := s.SetRouter()

	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		var (
			path, _   = route.GetPathTemplate()
			verbes, _ = route.GetMethods()
			verbe     = strings.Join(verbes, ",")
		)

		if !(path == _testURI || path == _testPrefix) {
			t.Errorf("route wrongly created : [%s]", path)
		}
		if verbe != "" {
			z.AssertStringEqual(t, verbe, GET)
		}
		return nil
	})
}

// TODO: func TestRouteApplier(t *testing.T) {}

func TestHandleParam(t *testing.T) {
	s := InitServer(false)
	defer s.WaitAndStop()
	defer s.Shutdown(*s.GetContext())

	s.GET("/test/{id}", func(c IContext) {
		if val, ok := c.GetQuery("pjson"); !ok || val != "1" {
			t.Errorf("query Param wrongly decoded %s", val)
		} else if c.GetVar("id") != "toto" {
			t.Errorf("URL Param wrongly decoded")
		}
		c.JSONNoContent()
	})

	go s.Start(":4242")
	time.Sleep(50 * time.Millisecond)

	z.RequestAndTestAPI(t, "/test/toto?pjson=1", func(t *testing.T, resp *http.Response) {})
}
