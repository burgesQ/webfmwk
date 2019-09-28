package webfmwk

import (
	"net/http"
	"testing"
	"time"

	m "github.com/burgesQ/webfmwk/middleware"
	z "github.com/burgesQ/webfmwk/testing"
)

type customContext struct {
	Context
	Value string
}

type testSerial struct {
	A string `json:"test"`
}

// implem ctx interface ?
func (c *customContext) GetName() string {
	return "custom context"
}

func TestUseCase(t *testing.T) {
	s := InitServer(false)
	defer s.WaitAndStop()
	defer s.Shutdown(*s.GetContext())

	// add middleware TODO: check headers
	s.AddMiddleware(m.Security)

	// set url prefix
	s.SetPrefix("/api")

	// set custom context
	s.SetCustomContext(func(c *Context) IContext {
		cctx := &customContext{*c, "turlu"}
		return cctx
	})

	// declare routes
	s.GET("/hello", func(c IContext) error { return c.JSONBlob(http.StatusOK, []byte(`{ "message": "hello world" }`)) })
	s.GET("/routes", func(c IContext) error { return c.JSON(http.StatusOK, &testSerial{"hello"}) })
	s.GET("/hello/{who}", func(c IContext) error {
		var content = `{ "message": "hello ` + c.GetVar("who") + `" }`
		return c.JSONBlob(http.StatusOK, []byte(content))
	})
	s.GET("/testquery", func(c IContext) error { return c.JSONOk(c.GetQueries()) })
	s.GET("/testContext", func(c IContext) error {
		cc := c.(*customContext)
		return c.JSONBlob(http.StatusOK, []byte(string(`{ "message": "hello `+cc.Value+`" }`)))
	})
	s.POST("/world", func(c IContext) error {
		anonymous := struct {
			FirstName string `json:"first_name,omitempty" validate:"required"`
			LastName  string `json:"last_name,omitempty"  validate:"required"`
		}{}

		// check body handle the error management, so no return needed
		if err := c.FetchContent(&anonymous); err != nil {
			return c.JSONUnprocessable(err)
		}
		return c.JSONCreated(anonymous)
	})

	go func() {
		if e := s.Start(":4242"); e != nil {
			t.Fatalf("error while booting the server : %s", e.Error())
		}
	}()
	time.Sleep(1 * time.Second)

	// request each routes
	z.RequestAndTestAPI(t, "/api/hello", func(t *testing.T, resp *http.Response) bool {
		for _, test_val := range []string{"Content-Type", "Accept", "Produce"} {
			if !z.TestHeader(t, resp, test_val, json_encode) {
				return false
			}
		}
		return z.TestBody(t, resp, `{"message":"hello world"}`) && z.TestStatusCode(t, resp, http.StatusOK)
	})
	z.RequestAndTestAPI(t, "/api/routes", func(t *testing.T, resp *http.Response) bool {
		return z.TestBody(t, resp, `{"test":"hello"}`) && z.TestStatusCode(t, resp, http.StatusOK)
	})
	z.RequestAndTestAPI(t, "/api/hello/you", func(t *testing.T, resp *http.Response) bool {
		return z.TestBody(t, resp, `{"message":"hello you"}`) && z.TestStatusCode(t, resp, http.StatusOK)
	})
	z.RequestAndTestAPI(t, "/api/testquery?pjson=1", func(t *testing.T, resp *http.Response) bool {
		return z.TestBodyDiffere(t, resp, `{"pjson":["1"]}`) && z.TestStatusCode(t, resp, http.StatusOK)
	})
	z.RequestAndTestAPI(t, "/api/testContext", func(t *testing.T, resp *http.Response) bool {
		return z.TestBodyDiffere(t, resp, `{"message":"hello tutu"}`) && z.TestStatusCode(t, resp, http.StatusOK)
	})
	z.PushAndTestAPI(t, "/api/world", []byte(string(`{"first_name":"jean"}`)), func(t *testing.T, resp *http.Response) bool {
		return z.TestBody(t, resp, `{"first_name":"jean"}`) && z.TestStatusCode(t, resp, http.StatusCreated)
	})
}
