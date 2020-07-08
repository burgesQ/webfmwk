package webfmwk

import (
	"net/http"
	"testing"

	"github.com/burgesQ/gommon/assert"
)

type customContext struct {
	Context
	Value string
}

type testSerial struct {
	A string `json:"test"`
}

func TestUseCase(t *testing.T) {
	var s = InitServer(
		CheckIsUp(), SetPrefix("/api"),
		WithHandlers(func(next HandlerFunc) HandlerFunc {
			return HandlerFunc(func(c Context) error {
				cc := customContext{c, "turlu"}
				return next(cc)
			})
		}),
	)

	// declare routes
	s.GET("/hello", func(c Context) error {
		return c.JSONBlob(http.StatusOK, []byte(`{ "message": "hello world" }`))
	})

	s.GET("/routes", func(c Context) error {
		return c.JSON(http.StatusOK, &testSerial{"hello"})
	})

	s.GET("/hello/{who}", func(c Context) error {
		var content = `{ "message": "hello ` + c.GetVar("who") + `" }`
		return c.JSONBlob(http.StatusOK, []byte(content))
	})

	s.GET("/testquery", func(c Context) error {
		return c.JSONOk(c.GetQueries())
	})

	s.GET("/testContext", func(c Context) error {
		return c.JSONBlob(http.StatusOK, []byte(`{ "message": "hello `+
			c.(customContext).Value+`" }`))
	})

	s.POST("/world", func(c Context) error {
		anonymous := struct {
			FirstName string `json:"first_name,omitempty" validate:"required"`
			LastName  string `json:"last_name,omitempty"  validate:"required"`
		}{}

		if e := c.FetchContent(&anonymous); e != nil {
			return e
		}

		return c.JSONCreated(anonymous)
	})

	defer stopServer(t, s)
	go s.Start(_testPort)
	<-s.isReady

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

				assert.RequestAndTestAPI(t, _testAddr+test.url, func(t *testing.T, resp *http.Response) {
					if test.header {
						for _, testVal := range []string{"Content-Type", "Accept", "Produce"} {
							assert.Header(t, resp, testVal, jsonEncode)
						}
					}
					if test.bodyDiffer {
						assert.BodyDiffere(t, resp, test.expectedBody)
					} else {
						assert.Body(t, resp, test.expectedBody)
					}
					assert.StatusCode(t, resp, test.expectedSC)
				})

			case _pushNTest:
				assert.PushAndTestAPI(t, _testAddr+test.url, []byte(string(`{"first_name":"jean"}`)),
					func(t *testing.T, resp *http.Response) {
						assert.Body(t, resp, test.expectedBody)
						assert.StatusCode(t, resp, test.expectedSC)
					})
			}

		})
	}

}
