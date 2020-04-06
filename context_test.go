package webfmwk

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/burgesQ/webfmwk/v3/log"
	z "github.com/burgesQ/webfmwk/v3/testing"
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
	handlerRoute func(c IContext), handlerTest z.HandlerForTest) {
	var s = InitServer(CheckIsUp())

	t.Log("init server...")
	defer stopServer(t, s)

	s.POST(route, handlerRoute)
	s.Start(_testPort)
	<-s.isReady
	t.Log("server inited")

	z.PushAndTestAPI(t, _testAddr+routeReq, content, handlerTest)
}

func wrapperGet(t *testing.T, route, routeReq string,
	handlerRoute func(c IContext), handlerTest z.HandlerForTest) {
	var s = InitServer(CheckIsUp())

	t.Log("init server...")
	defer stopServer(t, s)

	s.GET(route, handlerRoute)
	s.Start(_testPort)
	<-s.isReady
	t.Log("server inited")

	z.RequestAndTestAPI(t, _testAddr+routeReq, handlerTest)
}

func TestParam(t *testing.T) {
	wrapperGet(t, "/test/{id}", "/test/tutu", func(c IContext) {
		id := c.GetVar("id")
		if id != "tutu" {
			t.Errorf("error fetching the url param : [%s] expected [tutu]", id)
		}
		c.JSONOk(id)
	}, func(t *testing.T, resp *http.Response) {
		z.AssertBody(t, resp, `"tutu"`)
		z.AssertStatusCode(t, resp, http.StatusOK)
	})
}

func TestQuery(t *testing.T) {
	var (
		c = Context{
			query: map[string][]string{"test": {"ok"}},
		}
		v, ok = c.GetQuery("test")
	)

	z.AssertTrue(t, ok)
	z.AssertStringEqual(t, v, "ok")

	v, ok = c.GetQuery("undef")
	z.AssertFalse(t, ok)
	z.AssertStringEqual(t, v, "")
}

func TestLogger(t *testing.T) {
	var (
		c      = Context{}
		logger = log.GetLogger()
	)

	c.SetLogger(logger)
	z.AssertTrue(t, logger == c.GetLogger())
	z.AssertTrue(t, logger == GetLogger())

}

func TestContext(t *testing.T) {
	var (
		ctx context.Context
		c   = Context{
			ctx: ctx,
		}
	)
	z.AssertTrue(t, ctx == c.GetContext())
}

func TestRequestID(t *testing.T) {
	var ctx = Context{}

	ctx.SetRequestID("testing")
	z.AssertStringEqual(t, ctx.GetRequestID(), "testing")
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
			wrapperPost(t, "/test", "/test", test.payload, func(c IContext) {
				var anonymous = struct {
					FirstName string `json:"first_name,omitempty" validate:"required"`
				}{}
				c.FetchContent(&anonymous)
				c.Validate(anonymous)
				c.JSON(http.StatusCreated, anonymous)
			}, func(t *testing.T, resp *http.Response) {
				switch test.t {

				case _ok:
					z.AssertBody(t, resp, `{"first_name":"tutu"}`)
					z.AssertStatusCode(t, resp, http.StatusCreated)

				case _unprocessable:
					z.AssertStatusCode(t, resp, http.StatusUnprocessableEntity)
				}
			})
		})
	}
}

func TestCheckHeader(t *testing.T) {
	wrapperPost(t, "/test", "/test", []byte(`{}`), func(c IContext) {
		c.JSONBlob(200, []byte(hBody))
	}, func(t *testing.T, resp *http.Response) {
		z.AssertBody(t, resp, hBody)
		z.AssertStatusCode(t, resp, http.StatusOK)
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
	s.POST("/test", func(c IContext) {
		c.JSONBlob(http.StatusOK, []byte(hBody))
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
				z.AssertStatusCode(t, resp, http.StatusNotAcceptable)
				z.AssertBody(t, resp, `{"error":"Content-Type is not application/json"}`)
			case _noValue:
				z.AssertStatusCode(t, resp, http.StatusNotAcceptable)
				z.AssertBody(t, resp, `{"error":"Missing Content-Type header"}`)
			case _noHeader:
				z.AssertStatusCode(t, resp, http.StatusNotAcceptable)
			}

		})
	}

}

func TestJSONBlobPretty(t *testing.T) {
	wrapperGet(t, "/test", "/test?pretty", func(c IContext) {
		c.JSONBlob(http.StatusOK, []byte(hBody))
	}, func(t *testing.T, resp *http.Response) {
		z.AssertBodyDiffere(t, resp, hBody)
		z.AssertStatusCode(t, resp, http.StatusOK)
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
			fn         func(c IContext, ret interface{})
		}{
			"blob": {http.StatusOK, func(c IContext, ret interface{}) {
				c.JSONBlob(http.StatusOK, []byte(hBody))
			}},
			"ok": {http.StatusOK, func(c IContext, ret interface{}) {
				c.JSONOk(ret)
			}},
			"created": {http.StatusCreated, func(c IContext, ret interface{}) {
				c.JSONCreated(ret)
			}},
			"accepted": {http.StatusAccepted, func(c IContext, ret interface{}) {
				c.JSONAccepted(ret)
			}},
			"no content": {http.StatusNoContent, func(c IContext, ret interface{}) {
				c.JSONNoContent()
			}},
			"bad request": {http.StatusBadRequest, func(c IContext, ret interface{}) {
				c.JSONBadRequest(ret)
			}},
			"unauthorized": {http.StatusUnauthorized, func(c IContext, ret interface{}) {
				c.JSONUnauthorized(ret)
			}},
			"forbiden": {http.StatusForbidden, func(c IContext, ret interface{}) {
				c.JSONForbiden(ret)
			}},
			"notFound": {http.StatusNotFound, func(c IContext, ret interface{}) {
				c.JSONNotFound(ret)
			}},
			"conflict": {http.StatusConflict, func(c IContext, ret interface{}) {
				c.JSONConflict(ret)
			}},
			"unprocessable": {http.StatusUnprocessableEntity, func(c IContext, ret interface{}) {
				c.JSONUnprocessable(ret)
			}},
			"internalError": {http.StatusInternalServerError, func(c IContext, ret interface{}) {
				c.JSONInternalError(ret)
			}},
			"notImplemented": {http.StatusNotImplemented, func(c IContext, ret interface{}) {
				c.JSONNotImplemented(ret)
			}},
		}
	)

	defer stopServer(t, s)

	// load custom endpoints
	for n, t := range tests {
		var fn = t.fn
		s.GET("/"+n, func(c IContext) {
			fn(c, ret)
		})
	}

	s.Start(_testPort)
	<-s.isReady

	// s.DumpRoutes()

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {

			z.RequestAndTestAPI(t, _testAddr+"/"+name, func(t *testing.T, resp *http.Response) {

				if test.expectedOP != http.StatusNoContent {
					z.AssertBody(t, resp, hBody)
				}
				z.AssertStatusCode(t, resp, test.expectedOP)
			})
		})
	}
}
