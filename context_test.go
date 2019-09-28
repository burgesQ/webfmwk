package webfmwk

import (
	"bytes"
	"net/http"
	"testing"
	"time"

	z "github.com/burgesQ/webfmwk/testing"
)

var (
	h_body      = `{"message":"nul"}`
	json_encode = "application/json; charset=UTF-8"
)

func wrapperPost(t *testing.T,
	route, routeReq string,
	content []byte,
	handlerRoute func(c IContext) error,
	handlerTest z.HandlerForTest) {
	s := InitServer(false)
	defer s.WaitAndStop()
	defer s.Shutdown(*s.GetContext())

	s.POST(route, handlerRoute)

	go func() {
		if e := s.Start(":4242"); e != nil {
			t.Fatalf("error while booting the server : %s", e.Error())
		}
	}()
	time.Sleep(1 * time.Second)

	z.PushAndTestAPI(t, routeReq, content, handlerTest)
}

func wrapperGet(t *testing.T,
	route, routeReq string,
	handlerRoute func(c IContext) error,
	handlerTest z.HandlerForTest) {
	s := InitServer(false)
	defer s.WaitAndStop()
	defer s.Shutdown(*s.GetContext())

	s.GET(route, handlerRoute)

	go func() {
		if e := s.Start(":4242"); e != nil {
			t.Fatalf("error while booting the server : %s", e.Error())
		}
	}()
	time.Sleep(1 * time.Second)

	z.RequestAndTestAPI(t, routeReq, handlerTest)
}

func TestParam(t *testing.T) {
	wrapperGet(t, "/test/{id}", "/test/tutu", func(c IContext) error {
		id := c.GetVar("id")
		if id != "tutu" {
			t.Errorf("error fetching the url param : [%s] expected [tutu]", id)
		}
		return c.JSONOk(id)
	}, func(t *testing.T, resp *http.Response) {
		z.AssertBody(t, resp, `"tutu"`)
		z.AssertStatusCode(t, resp, http.StatusOK)
	})
}

func TestFetchContentUnprocessable(t *testing.T) {
	wrapperPost(t, "/test", "/test", []byte(`{"first_name": tutu"}`), func(c IContext) error {

		anonymous := struct {
			FirstName string `json:"first_name,omitempty" validate:"required"`
		}{}

		if err := c.FetchContent(&anonymous); err != nil {
			return c.JSONUnprocessable(AnonymousError{err.Error()})
		} else if err = c.Validate(anonymous); err != nil {
			return c.JSONUnprocessable(AnonymousError{err.Error()})
		}

		return c.JSON(http.StatusCreated, anonymous)
	}, func(t *testing.T, resp *http.Response) {
		z.AssertStatusCode(t, resp, http.StatusUnprocessableEntity)
	})
}

func TestFetchContent(t *testing.T) {
	wrapperPost(t, "/test", "/test", []byte(`{"first_name": "tutu"}`), func(c IContext) error {

		anonymous := struct {
			FirstName string `json:"first_name,omitempty" validate:"required"`
		}{}

		if err := c.FetchContent(&anonymous); err != nil {
			return err
		}

		return c.JSON(http.StatusCreated, anonymous)
	}, func(t *testing.T, resp *http.Response) {
		z.AssertBody(t, resp, `{"first_name":"tutu"}`)
		z.AssertStatusCode(t, resp, http.StatusCreated)
	})
}

func TestCheckHeaderNoHeader(t *testing.T) {
	s := InitServer(false)
	defer s.WaitAndStop()
	defer s.Shutdown(*s.GetContext())

	s.POST("/test", func(c IContext) error {
		return c.JSONBlob(http.StatusOK, []byte(h_body))
	})

	go func() {
		if e := s.Start(":4242"); e != nil {
			t.Fatalf("error while booting the server : %s", e.Error())
		}
	}()
	time.Sleep(1 * time.Second)

	url := "http://127.0.0.1:4242" + "/test"
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(h_body)))
	client := &http.Client{}

	if resp, err := client.Do(req); err != nil {
		t.Fatalf("error requesting the api : %s", err.Error())
	} else {
		z.AssertStatusCode(t, resp, http.StatusNotAcceptable)
	}
}

func TestCheckHeaderWrongHeader(t *testing.T) {
	s := InitServer(false)
	defer s.WaitAndStop()
	defer s.Shutdown(*s.GetContext())

	s.POST("/test", func(c IContext) error {
		return c.JSONBlob(http.StatusOK, []byte(h_body))
	})

	go func() {
		if e := s.Start(":4242"); e != nil {
			t.Fatalf("error while booting the server : %s", e.Error())
		}
	}()
	time.Sleep(1 * time.Second)

	url := "http://127.0.0.1:4242" + "/test"
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(h_body)))
	req.Header.Set("Content-Type", "")

	client := &http.Client{}
	if resp, err := client.Do(req); err != nil {
		t.Fatalf("error requesting the api : %s", err.Error())
	} else {
		z.AssertStatusCode(t, resp, http.StatusNotAcceptable)
	}
}

func TestCheckHeader(t *testing.T) {
	wrapperPost(t, "/test", "/test", []byte(`{}`), func(c IContext) error {
		return c.JSONBlob(200, []byte(h_body))
	}, func(t *testing.T, resp *http.Response) {
		z.AssertBody(t, resp, h_body)
		z.AssertStatusCode(t, resp, http.StatusOK)
	})
}

func TestJSONBlobPretty(t *testing.T) {
	wrapperGet(t, "/test", "/test?pjson", func(c IContext) error {
		return c.JSONBlob(http.StatusOK, []byte(h_body))
	}, func(t *testing.T, resp *http.Response) {
		z.AssertBodyDiffere(t, resp, h_body)
		z.AssertStatusCode(t, resp, http.StatusOK)
	})
}

func TestJSONBlob(t *testing.T) {
	wrapperGet(t, "/test", "/test", func(c IContext) error {
		return c.JSONBlob(http.StatusOK, []byte(h_body))
	}, func(t *testing.T, resp *http.Response) {
		for _, test_val := range []string{"Content-Type", "Accept", "Produce"} {
			z.AssertHeader(t, resp, test_val, json_encode)
		}
		z.AssertBody(t, resp, h_body)
		z.AssertStatusCode(t, resp, http.StatusOK)
	})
}

func TestJSONNotImplemented(t *testing.T) {
	wrapperGet(t, "/test", "/test", func(c IContext) error {
		ret := struct {
			Message string `json:"message"`
		}{"nul"}
		return c.JSONNotImplemented(ret)
	}, func(t *testing.T, resp *http.Response) {
		z.AssertBody(t, resp, h_body)
		z.AssertStatusCode(t, resp, http.StatusNotImplemented)
	})
}

func TestJSONCreated(t *testing.T) {
	wrapperGet(t, "/test", "/test", func(c IContext) error {
		ret := struct {
			Message string `json:"message"`
		}{"nul"}
		return c.JSONCreated(ret)
	}, func(t *testing.T, resp *http.Response) {
		z.AssertBody(t, resp, h_body)
		z.AssertStatusCode(t, resp, http.StatusCreated)
	})
}

func TestJSONUnprocessable(t *testing.T) {
	wrapperGet(t, "/test", "/test", func(c IContext) error {
		ret := struct {
			Message string `json:"message"`
		}{"nul"}
		return c.JSONUnprocessable(ret)
	}, func(t *testing.T, resp *http.Response) {
		z.AssertBody(t, resp, h_body)
		z.AssertStatusCode(t, resp, http.StatusUnprocessableEntity)
	})
}

func TestJSONOk(t *testing.T) {
	wrapperGet(t, "/test", "/test", func(c IContext) error {
		ret := struct {
			Message string `json:"message"`
		}{"nul"}
		return c.JSONOk(ret)
	}, func(t *testing.T, resp *http.Response) {
		z.AssertBody(t, resp, h_body)
		z.AssertStatusCode(t, resp, http.StatusOK)
	})
}

func TestJSONNotFound(t *testing.T) {
	wrapperGet(t, "/test", "/test", func(c IContext) error {
		ret := struct {
			Message string `json:"message"`
		}{"nul"}
		return c.JSONNotFound(ret)
	}, func(t *testing.T, resp *http.Response) {
		z.AssertBody(t, resp, h_body)
		z.AssertStatusCode(t, resp, http.StatusNotFound)
	})
}

func TestJSONConflict(t *testing.T) {
	wrapperGet(t, "/test", "/test", func(c IContext) error {
		ret := struct {
			Message string `json:"message"`
		}{"nul"}
		return c.JSONConflict(ret)
	}, func(t *testing.T, resp *http.Response) {
		z.AssertBody(t, resp, h_body)
		z.AssertStatusCode(t, resp, http.StatusConflict)
	})
}

func TestJSONInternalError(t *testing.T) {
	wrapperGet(t, "/test", "/test", func(c IContext) error {
		ret := struct {
			Message string `json:"message"`
		}{"nul"}
		return c.JSONInternalError(ret)
	}, func(t *testing.T, resp *http.Response) {
		z.AssertBody(t, resp, h_body)
		z.AssertStatusCode(t, resp, http.StatusInternalServerError)
	})
}
