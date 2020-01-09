package webfmwk

import (
	"net/http"
	"testing"
	"time"

	// "gitlab.frafos.net/gommon/golib/log"
	z "github.com/burgesQ/webfmwk/v2/testing"
)

type customContext struct {
	Context
	Value string
}

type testSerial struct {
	A string `json:"test"`
}

func TestUseCase(t *testing.T) {
	s := InitServer(false)

	// log.Init(log.LOGFORMAT_LONG | log.LOGGER_STDOUT)
	// log.SetLogLevel(log.LOG_DEBUG)
	// s.SetLogger(log.GetLogger())

	// set custom context
	if s.SetCustomContext(func(c *Context) IContext {
		return &customContext{*c, "turlu"}
	}) == false {
		t.Errorf("cannot set the custom context")
	}

	// add middleware TODO: check headers
	// s.AddMiddleware(m.Security)

	// set url prefix
	s.SetPrefix("/api")

	// declare routes
	s.GET("/hello", func(c IContext) {
		c.JSONBlob(http.StatusOK, []byte(`{ "message": "hello world" }`))
	})

	s.GET("/routes", func(c IContext) {
		c.JSON(http.StatusOK, &testSerial{"hello"})
	})

	s.GET("/hello/{who}", func(c IContext) {
		var content = `{ "message": "hello ` + c.GetVar("who") + `" }`
		c.JSONBlob(http.StatusOK, []byte(content))
	})

	s.GET("/testquery", func(c IContext) {
		c.JSONOk(c.GetQueries())
	})

	s.GET("/testContext", func(c IContext) {
		c.JSONBlob(http.StatusOK, []byte(`{ "message": "hello `+
			c.(*customContext).Value+`" }`))
	})

	s.POST("/world", func(c IContext) {
		anonymous := struct {
			FirstName string `json:"first_name,omitempty" validate:"required"`
			LastName  string `json:"last_name,omitempty"  validate:"required"`
		}{}
		// check body handle the error management, so no return needed
		c.FetchContent(&anonymous)
		c.JSONCreated(anonymous)
	})

	defer func(s Server) {
		s.Shutdown(*s.GetContext())
		s.WaitAndStop()
	}(s)
	go s.Start(":4242")
	time.Sleep(50 * time.Millisecond)

	const (
		_reqNTest = iota
		_pushNTest
	)

	tests := map[string]struct {
		testType     int
		header       bool
		bodyDiffer   bool
		url          string
		expectedBody string
		expectedSC   int
	}{
		"hello world": {
			_reqNTest, true, false, "/api/hello", `{"message":"hello world"}`, http.StatusOK,
		},
		"simple fetch": {
			_reqNTest, false, false, "/api/routes", `{"test":"hello"}`, http.StatusOK,
		},
		"url params": {
			_reqNTest, false, false, "/api/hello/you", `{"message":"hello you"}`, http.StatusOK,
		},
		"query params": {
			_reqNTest, false, true, "/api/testquery?pretty=1", `{"pretty":["1"]}`, http.StatusOK,
		},
		"context": {
			_reqNTest, false, true, "/api/testContext", `{"message":"hello tutu"}`, http.StatusOK,
		},
		"push": {
			_pushNTest, false, false, "/api/world", `{"first_name":"jean"}`, http.StatusCreated,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {

			switch test.testType {
			case _reqNTest:

				z.RequestAndTestAPI(t, test.url, func(t *testing.T, resp *http.Response) {
					if test.header {
						for _, testVal := range []string{"Content-Type", "Accept", "Produce"} {
							z.AssertHeader(t, resp, testVal, jsonEncode)
						}
					}
					if test.bodyDiffer {
						z.AssertBodyDiffere(t, resp, test.expectedBody)
					} else {
						z.AssertBody(t, resp, test.expectedBody)
					}
					z.AssertStatusCode(t, resp, test.expectedSC)
				})

			case _pushNTest:
				z.PushAndTestAPI(t, test.url, []byte(string(`{"first_name":"jean"}`)),
					func(t *testing.T, resp *http.Response) {
						z.AssertBody(t, resp, test.expectedBody)
						z.AssertStatusCode(t, resp, test.expectedSC)
					})
			}

		})
	}

}
