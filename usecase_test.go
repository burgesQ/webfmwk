package webfmwk

import (
	"net/http"
	"testing"
	"time"

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

	defer func(s Server) {
		s.Shutdown(*s.GetContext())
		s.WaitAndStop()
	}(s)

	// add middleware TODO: check headers
	// s.AddMiddleware(m.Security)

	// set url prefix
	s.SetPrefix("/api")

	// set custom context
	s.SetCustomContext(func(c *Context) IContext {
		cctx := &customContext{*c, "turlu"}
		return cctx
	})

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
		cc := c.(*customContext)
		c.JSONBlob(http.StatusOK, []byte(`{ "message": "hello `+cc.Value+`" }`))
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

	go s.Start(":4242")
	time.Sleep(50 * time.Millisecond)

	// request each routes
	z.RequestAndTestAPI(t, "/api/hello",
		func(t *testing.T, resp *http.Response) {
			for _, testVal := range []string{"Content-Type", "Accept", "Produce"} {
				z.AssertHeader(t, resp, testVal, jsonEncode)
			}
			z.AssertBody(t, resp, `{"message":"hello world"}`)
			z.AssertStatusCode(t, resp, http.StatusOK)
		})
	z.RequestAndTestAPI(t, "/api/routes",
		func(t *testing.T, resp *http.Response) {
			z.AssertBody(t, resp, `{"test":"hello"}`)
			z.AssertStatusCode(t, resp, http.StatusOK)
		})
	z.RequestAndTestAPI(t, "/api/hello/you",
		func(t *testing.T, resp *http.Response) {
			z.AssertBody(t, resp, `{"message":"hello you"}`)
			z.AssertStatusCode(t, resp, http.StatusOK)
		})
	z.RequestAndTestAPI(t, "/api/testquery?pjson=1",
		func(t *testing.T, resp *http.Response) {
			z.AssertBodyDiffere(t, resp, `{"pjson":["1"]}`)
			z.AssertStatusCode(t, resp, http.StatusOK)
		})
	z.RequestAndTestAPI(t, "/api/testContext",
		func(t *testing.T, resp *http.Response) {
			z.AssertBodyDiffere(t, resp, `{"message":"hello tutu"}`)
			z.AssertStatusCode(t, resp, http.StatusOK)
		})
	z.PushAndTestAPI(t, "/api/world", []byte(string(`{"first_name":"jean"}`)),
		func(t *testing.T, resp *http.Response) {
			z.AssertBody(t, resp, `{"first_name":"jean"}`)
			z.AssertStatusCode(t, resp, http.StatusCreated)
		})
}
