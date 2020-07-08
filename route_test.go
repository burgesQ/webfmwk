package webfmwk

import (
	"net/http"
	"strings"
	"testing"

	"github.com/burgesQ/gommon/assert"
	"github.com/gorilla/mux"
)

const (
	_testPrefix = "/api"
	_testURL    = "/test"
	_testURI    = _testPrefix + _testURL
	_testURI2   = _testPrefix + _testURL + "/2"
	_testVerbe  = GET
)

var _emptyController = func(c Context) error {
	return nil
}

// TODO: func TestAddRoute(t *testing.T)  {}
// TODO: func TestAddRoutes(t *testing.T) {}

func TestSetPrefix(t *testing.T) {
	var s = InitServer(CheckIsUp(), SetPrefix(_testPrefix))
	defer stopServer(t, s)

	s.GET(_testURL, _emptyController)

	r := s.SetRouter()

	if e := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, _ := route.GetPathTemplate()

		if !(pathTemplate == _testURI || pathTemplate == _testPrefix) && pathTemplate != _testPrefix+_pingEndpoint {
			t.Errorf("route wrongly created : [%s]", pathTemplate)
		}

		return nil
	}); e != nil {
		t.Errorf("cannot walk routes : %s", e.Error())
	}
}

func TestAddRoutes(t *testing.T) {
	var s = InitServer(CheckIsUp())
	defer stopServer(t, s)

	s.AddRoutes(Route{
		Path:    _testURI,
		Verbe:   _testVerbe,
		Handler: _emptyController,
	})

	assert.StringEqual(t, s.meta.routes[s.meta.prefix][0].Path, _testURI)
	assert.StringEqual(t, s.meta.routes[s.meta.prefix][0].Verbe, _testVerbe)

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
	}...)

	assert.StringEqual(t, s.meta.routes[s.meta.prefix][1].Path, _testURI)
	assert.StringEqual(t, s.meta.routes[s.meta.prefix][1].Verbe, _testVerbe)
	assert.StringEqual(t, s.meta.routes[s.meta.prefix][2].Path, _testURI2)
	assert.StringEqual(t, s.meta.routes[s.meta.prefix][2].Verbe, _testVerbe)

}

func TestRouteMethod(t *testing.T) {
	const (
		_get = iota
		_delete
		_post
		_put
		_patch
	)

	tests := map[string]struct {
		reqType int
	}{
		"get":    {reqType: _get},
		"delete": {reqType: _delete},
		"post":   {reqType: _post},
		"put":    {reqType: _put},
		"patch":  {reqType: _patch},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			var s = InitServer(CheckIsUp())

			defer stopServer(t, s)

			testVerb := ""
			switch test.reqType {
			case _get:
				s.GET(_testURL, _emptyController)
				testVerb = GET
			case _delete:
				s.DELETE(_testURL, _emptyController)
				testVerb = DELETE
			case _post:
				s.POST(_testURL, _emptyController)
				testVerb = POST
			case _put:
				s.PUT(_testURL, _emptyController)
				testVerb = PUT
			case _patch:
				s.PATCH(_testURL, _emptyController)
				testVerb = PATCH
			}

			assert.StringEqual(t, s.meta.routes[s.meta.prefix][0].Path, _testURL)
			assert.StringEqual(t, s.meta.routes[s.meta.prefix][0].Verbe, testVerb)
		})
	}

}

func TestSetRouter(t *testing.T) {
	s := InitServer(CheckIsUp(), SetPrefix(_testPrefix))
	defer stopServer(t, s)

	s.GET(_testURL, func(c Context) error {
		return c.JSONNoContent()
	})

	r := s.SetRouter()

	if e := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		var (
			path, _   = route.GetPathTemplate()
			verbes, _ = route.GetMethods()
			verbe     = strings.Join(verbes, ",")
		)

		if path != _testURI && path != _testPrefix && path != _testPrefix+_pingEndpoint {
			t.Errorf("route wrongly created : [%s]", path)
			t.Errorf("[%s][%s][%s]", _testURI, _testPrefix, _pingEndpoint)
		}
		if verbe != "" {
			assert.StringEqual(t, verbe, GET)
		}
		return nil
	}); e != nil {
		t.Errorf("cannot walk routes : %s", e.Error())
	}
}

// TODO: func TestRouteApplier(t *testing.T) {}

func TestHandleParam(t *testing.T) {
	s := InitServer(CheckIsUp())
	defer stopServer(t, s)

	s.GET("/test/{id}", func(c Context) error {
		if val, ok := c.GetQuery("pretty"); !ok || val != "1" {
			t.Errorf("query Param wrongly decoded %s", val)
		} else if c.GetVar("id") != "toto" {
			t.Errorf("URL Param wrongly decoded")
		}
		return c.JSONNoContent()
	})

	s.Start(_testPort)
	<-s.isReady

	assert.RequestAndTestAPI(t, _testAddr+"/test/toto?pretty=1", func(t *testing.T, resp *http.Response) {})
}
