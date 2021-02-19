package webfmwk

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/burgesQ/gommon/assert"
	"github.com/burgesQ/gommon/log"
)

var (
	hBody      = `{"message":"nul"}`
	jsonEncode = "application/json; charset=UTF-8"
	_testPort  = ":6666"
	_testAddr  = "http://127.0.0.1" + _testPort
)

func stopServer(t *testing.T, s *Server) {
	var ctx = s.GetContext()

	ctx.Done()
	s.Shutdown(ctx)
	s.WaitAndStop()
	Shutdown(ctx)
	t.Log("server closed")
}

func wrapperPost(t *testing.T, route, routeReq string, content []byte,
	handlerRoute HandlerFunc, handlerTest assert.HandlerForTest) {
	var s = InitServer(CheckIsUp())

	t.Log("init server...")
	defer stopServer(t, s)

	s.POST(route, handlerRoute)
	s.Start(_testPort)
	<-s.isReady
	t.Log("server inited")

	assert.PushAndTestAPI(t, _testAddr+routeReq, content, handlerTest)
}

func wrapperGet(t *testing.T, route, routeReq string,
	handlerRoute HandlerFunc, handlerTest assert.HandlerForTest) {
	var s = InitServer(CheckIsUp())

	t.Log("init server...")
	defer stopServer(t, s)

	s.GET(route, handlerRoute)
	s.Start(_testPort)
	<-s.isReady
	t.Log("server inited")

	assert.RequestAndTestAPI(t, _testAddr+routeReq, handlerTest)
}

func TestParam(t *testing.T) {
	wrapperGet(t, "/test/{id}", "/test/tutu", func(c Context) error {
		id := c.GetVar("id")
		if id != "tutu" {
			t.Errorf("error fetching the url param : [%s] expected [tutu]", id)
		}
		return c.JSONOk(id)
	}, func(t *testing.T, resp *http.Response) {
		assert.Body(t, resp, `"tutu"`)
		assert.StatusCode(t, resp, http.StatusOK)
	})
}

func TestQuery(t *testing.T) {
	var (
		c = icontext{
			query: map[string][]string{"test": {"ok"}},
		}
		v, ok = c.GetQuery("test")
	)

	assert.True(t, ok)
	assert.StringEqual(t, v, "ok")

	v, ok = c.GetQuery("undef")
	assert.False(t, ok)
	assert.StringEqual(t, v, "")
}

func TestLogger(t *testing.T) {
	var (
		c      = icontext{}
		logger = log.GetLogger()
	)

	c.SetLogger(logger)
	assert.True(t, logger == c.GetLogger())
	assert.True(t, logger == GetLogger())

}

func TestContext(t *testing.T) {
	var (
		ctx context.Context
		c   = icontext{
			ctx: ctx,
		}
	)
	assert.True(t, ctx == c.GetContext())
}

func TestFetchContent(t *testing.T) {
	const (
		_ok = iota
		_unprocessable
	)

	var tests = map[string]struct {
		payload []byte
		t       int
	}{
		"fetch content":               {[]byte(`{"first_name": "tutu"}`), _ok},
		"fetch content unprocessable": {[]byte(`{"first_name": tutu"}`), _unprocessable},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			wrapperPost(t, "/test", "/test", test.payload, func(c Context) error {
				var anonymous = struct {
					FirstName string `json:"first_name,omitempty" validate:"required"`
				}{}
				if e := c.FetchContent(&anonymous); e != nil {
					return e
				} else if e := c.Validate(anonymous); e != nil {
					return e
				}
				return c.JSON(http.StatusCreated, anonymous)
			}, func(t *testing.T, resp *http.Response) {
				switch test.t {

				case _ok:
					assert.Body(t, resp, `{"first_name":"tutu"}`)
					assert.StatusCode(t, resp, http.StatusCreated)

				case _unprocessable:
					assert.StatusCode(t, resp, http.StatusUnprocessableEntity)
				}
			})
		})
	}
}

func TestCheckHeader(t *testing.T) {
	wrapperPost(t, "/test", "/test", []byte(`{}`), func(c Context) error {
		return c.JSONBlob(200, []byte(hBody))
	}, func(t *testing.T, resp *http.Response) {
		assert.Body(t, resp, hBody)
		assert.StatusCode(t, resp, http.StatusOK)
	})
}

func TestCheckHeaderError(t *testing.T) {
	const (
		_xml = iota
		_noHeader
		_noValue
	)

	var (
		s     = InitServer(CheckIsUp())
		tests = map[string]struct {
			headerValue string
			noHeader    bool
			t           int
		}{
			"xml value": {"application/xml", false, _xml},
			"no value":  {"", false, _noValue},
			"no header": {"", true, _noHeader},
		}
	)

	defer stopServer(t, s)
	s.POST("/test", func(c Context) error {
		return c.JSONBlob(http.StatusOK, []byte(hBody))
	})
	s.Start(_testPort)
	<-s.isReady

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// TODO: wrap that in the test fmwk
			var (
				url    = "http://127.0.0.1" + _testPort + "/test"
				req, _ = http.NewRequest("POST", url, bytes.NewBuffer([]byte(hBody)))
				client = &http.Client{}
			)

			if !test.noHeader {
				req.Header.Set("Content-Type", test.headerValue)
			}

			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("error requesting the api : %s", err.Error())
			}
			defer resp.Body.Close()

			switch test.t {
			case _xml:
				assert.StatusCode(t, resp, http.StatusNotAcceptable)
				assert.Body(t, resp, `{"status":406,"message":"Content-Type is not application/json"}`)
			case _noValue:
				assert.StatusCode(t, resp, http.StatusNotAcceptable)
				assert.Body(t, resp, `{"status":406,"message":"Missing Content-Type header"}`)
			case _noHeader:
				assert.StatusCode(t, resp, http.StatusNotAcceptable)
			}

		})
	}

}

func TestJSONBlobPretty(t *testing.T) {
	wrapperGet(t, "/test", "/test?pretty", func(c Context) error {
		return c.JSONBlob(http.StatusOK, []byte(hBody))
	}, func(t *testing.T, resp *http.Response) {
		assert.BodyDiffere(t, resp, hBody)
		assert.StatusCode(t, resp, http.StatusOK)
	})
}

func TestJSONResponse(t *testing.T) {
	// log.SetLogLevel(log.LogDEBUG)
	var (
		s   = InitServer(CheckIsUp(), WithLogger(log.GetLogger()))
		ret = struct {
			Message string `json:"message"`
		}{"nul"}
		tests = map[string]struct {
			expectedOP int
			fn         func(c Context, ret interface{}) error
		}{
			"blob": {http.StatusOK, func(c Context, ret interface{}) error {
				return c.JSONBlob(http.StatusOK, []byte(hBody))
			}},
			"ok": {http.StatusOK, func(c Context, ret interface{}) error {
				return c.JSONOk(ret)
			}},
			"created": {http.StatusCreated, func(c Context, ret interface{}) error {
				return c.JSONCreated(ret)
			}},
			"accepted": {http.StatusAccepted, func(c Context, ret interface{}) error {
				return c.JSONAccepted(ret)
			}},
			"no content": {http.StatusNoContent, func(c Context, ret interface{}) error {
				return c.JSONNoContent()
			}},
			"bad request": {http.StatusBadRequest, func(c Context, ret interface{}) error {
				return c.JSONBadRequest(ret)
			}},
			"unauthorized": {http.StatusUnauthorized, func(c Context, ret interface{}) error {
				return c.JSONUnauthorized(ret)
			}},
			"forbiden": {http.StatusForbidden, func(c Context, ret interface{}) error {
				return c.JSONForbiden(ret)
			}},
			"notFound": {http.StatusNotFound, func(c Context, ret interface{}) error {
				return c.JSONNotFound(ret)
			}},
			"conflict": {http.StatusConflict, func(c Context, ret interface{}) error {
				return c.JSONConflict(ret)
			}},
			"unprocessable": {http.StatusUnprocessableEntity, func(c Context, ret interface{}) error {
				return c.JSONUnprocessable(ret)
			}},
			"internalError": {http.StatusInternalServerError, func(c Context, ret interface{}) error {
				return c.JSONInternalError(ret)
			}},
			"notImplemented": {http.StatusNotImplemented, func(c Context, ret interface{}) error {
				return c.JSONNotImplemented(ret)
			}},
		}
	)

	defer stopServer(t, s)

	// load custom endpoints
	for n, t := range tests {
		var fn = t.fn
		s.GET("/"+n, func(c Context) error {
			return fn(c, ret)
		})
	}

	s.Start(_testPort)
	<-s.isReady

	// s.DumpRoutes()

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {

			assert.RequestAndTestAPI(t, _testAddr+"/"+name, func(t *testing.T, resp *http.Response) {

				if test.expectedOP != http.StatusNoContent {
					assert.Body(t, resp, hBody)
				}
				assert.StatusCode(t, resp, test.expectedOP)
			})
		})
	}
}
