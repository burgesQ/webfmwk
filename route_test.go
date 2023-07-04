package webfmwk

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	s, e := InitServer(CheckIsUp(), SetPrefix(_testPrefix))
	require.Nil(t, e)
	t.Cleanup(func() { require.Nil(t, s.ShutdownAndWait()) })

	s.GET(_testURL, _emptyController)

	all := s.GetRouter().List()
	assert.Contains(t, all["GET"], _testPrefix+_testURL)
	assert.Contains(t, all["GET"], _testPrefix+"/ping")
	assert.True(t, len(all["GET"]) == 2, "only 2 routes should be loaded")
}

func TestAddRoutes(t *testing.T) {
	s, e := InitServer(CheckIsUp())
	require.Nil(t, e)
	t.Cleanup(func() { require.Nil(t, s.ShutdownAndWait()) })

	asserter := assert.New(t)

	s.AddRoutes(Route{
		Path:    _testURI,
		Verbe:   _testVerbe,
		Handler: _emptyController,
	})

	t.Log("ensuring route path and verbe are persisted")
	{
		asserter.Equal(_testURI, s.meta.routes[s.meta.prefix][0].Path)
		asserter.Equal(_testVerbe, s.meta.routes[s.meta.prefix][0].Verbe)
	}

	s.AddRoutes(Routes{
		{
			Path:    _testURI + "1",
			Verbe:   POST,
			Handler: _emptyController,
		},
		{
			Path:    _testURI2,
			Verbe:   _testVerbe,
			Handler: _emptyController,
		},
	}...)

	t.Log("ensuring route path and verbe are persisted in correct order")
	{
		asserter.Equal(_testURI+"1", s.meta.routes[s.meta.prefix][1].Path)
		asserter.Equal(POST, s.meta.routes[s.meta.prefix][1].Verbe)
		asserter.Equal(_testURI2, s.meta.routes[s.meta.prefix][2].Path)
		asserter.Equal(_testVerbe, s.meta.routes[s.meta.prefix][2].Verbe)
	}
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
			s, e := InitServer(CheckIsUp())
			require.Nil(t, e)
			t.Cleanup(func() { require.Nil(t, s.ShutdownAndWait()) })

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

			assert.Equal(t, _testURL, s.meta.routes[s.meta.prefix][0].Path)
			assert.Equal(t, testVerb, s.meta.routes[s.meta.prefix][0].Verbe)
		})
	}
}

// TODO: func TestRouteApplier(t *testing.T) {}

func TestHandleParam(t *testing.T) {
	wrapperGet(t, "/test/{id}", "/test/toto?pretty=1", func(c Context) error {
		assert.Equal(t, []byte("1"), c.GetQuery().Peek("pretty"))
		assert.Equal(t, "toto", c.GetVar("id"))

		return c.JSONNoContent()
	}, func(t *testing.T, resp *http.Response) { t.Helper() })
}
