package webfmwk

import (
	"bytes"
	"net/http"
	"testing"
	"time"

	z "github.com/burgesQ/webfmwk/v2/testing"
)

var (
	hBody      = `{"message":"nul"}`
	jsonEncode = "application/json; charset=UTF-8"
)

func wrapperPost(t *testing.T, route, routeReq string, content []byte,
	handlerRoute func(c IContext), handlerTest z.HandlerForTest) {

	s := InitServer(false)
	// TODO:	s.SetLogLevel(-1)
	defer func() {
		s.Shutdown(*s.GetContext())
		s.WaitAndStop()
	}()

	s.POST(route, handlerRoute)
	go func() {
		if e := s.Start(":4242"); e != nil {
			t.Fatalf("error while booting the server : %s", e.Error())
		}
	}()
	time.Sleep(50 * time.Millisecond)

	z.PushAndTestAPI(t, routeReq, content, handlerTest)
}

func wrapperGet(t *testing.T, route, routeReq string,
	handlerRoute func(c IContext), handlerTest z.HandlerForTest) {

	s := InitServer(false)
	defer func() {
		s.Shutdown(*s.GetContext())
		s.WaitAndStop()
	}()

	s.GET(route, handlerRoute)
	go func() {
		if e := s.Start(":4242"); e != nil {
			t.Fatalf("error while booting the server : %s", e.Error())
		}
	}()
	time.Sleep(50 * time.Millisecond)

	z.RequestAndTestAPI(t, routeReq, handlerTest)
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

func TestFetchContentUnprocessable(t *testing.T) {

	wrapperPost(t, "/test", "/test", []byte(`{"first_name": tutu"}`), func(c IContext) {

		anonymous := struct {
			FirstName string `json:"first_name,omitempty" validate:"required"`
		}{}

		c.FetchContent(&anonymous)
		c.Validate(anonymous)

		c.JSON(http.StatusCreated, anonymous)
	}, func(t *testing.T, resp *http.Response) {
		z.AssertStatusCode(t, resp, http.StatusUnprocessableEntity)
	})
}

func TestFetchContent(t *testing.T) {
	wrapperPost(t, "/test", "/test", []byte(`{"first_name": "tutu"}`), func(c IContext) {

		anonymous := struct {
			FirstName string `json:"first_name,omitempty" validate:"required"`
		}{}

		c.FetchContent(&anonymous)
		c.JSON(http.StatusCreated, anonymous)
	}, func(t *testing.T, resp *http.Response) {
		z.AssertBody(t, resp, `{"first_name":"tutu"}`)
		z.AssertStatusCode(t, resp, http.StatusCreated)
	})
}

func TestCheckHeaderNoHeader(t *testing.T) {
	s := InitServer(false)
	defer func() {
		s.Shutdown(*s.GetContext())
		s.WaitAndStop()
	}()

	s.POST("/test", func(c IContext) {
		c.JSONBlob(http.StatusOK, []byte(hBody))
	})

	go func() {
		if e := s.Start(":4242"); e != nil {
			t.Fatalf("error while booting the server : %s", e.Error())
		}
	}()
	time.Sleep(50 * time.Millisecond)

	url := "http://127.0.0.1:4242" + "/test"
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(hBody)))
	client := &http.Client{}

	if resp, err := client.Do(req); err != nil {
		t.Fatalf("error requesting the api : %s", err.Error())
	} else {
		z.AssertStatusCode(t, resp, http.StatusNotAcceptable)
	}
}

func TestCheckHeaderWrongHeader(t *testing.T) {
	s := InitServer(false)
	defer func() {
		s.Shutdown(*s.GetContext())
		s.WaitAndStop()
	}()
	s.POST("/test", func(c IContext) {
		c.JSONBlob(http.StatusOK, []byte(hBody))
	})

	go func() {
		if e := s.Start(":4242"); e != nil {
			t.Fatalf("error while booting the server : %s", e.Error())
		}
	}()
	time.Sleep(50 * time.Millisecond)

	url := "http://127.0.0.1:4242" + "/test"
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(hBody)))
	req.Header.Set("Content-Type", "")

	client := &http.Client{}
	if resp, err := client.Do(req); err != nil {
		t.Fatalf("error requesting the api : %s", err.Error())
	} else {
		z.AssertStatusCode(t, resp, http.StatusNotAcceptable)
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

func TestJSONBlobPretty(t *testing.T) {
	wrapperGet(t, "/test", "/test?pjson", func(c IContext) {
		c.JSONBlob(http.StatusOK, []byte(hBody))
	}, func(t *testing.T, resp *http.Response) {
		z.AssertBodyDiffere(t, resp, hBody)
		z.AssertStatusCode(t, resp, http.StatusOK)
	})
}

func TestJSONBlob(t *testing.T) {
	wrapperGet(t, "/test", "/test", func(c IContext) {
		c.JSONBlob(http.StatusOK, []byte(hBody))
	}, func(t *testing.T, resp *http.Response) {
		for _, testVal := range []string{"Content-Type", "Accept", "Produce"} {
			z.AssertHeader(t, resp, testVal, jsonEncode)
		}
		z.AssertBody(t, resp, hBody)
		z.AssertStatusCode(t, resp, http.StatusOK)
	})
}

func TestJSONNotImplemented(t *testing.T) {
	wrapperGet(t, "/test", "/test", func(c IContext) {
		ret := struct {
			Message string `json:"message"`
		}{"nul"}
		c.JSONNotImplemented(ret)
	}, func(t *testing.T, resp *http.Response) {
		z.AssertBody(t, resp, hBody)
		z.AssertStatusCode(t, resp, http.StatusNotImplemented)
	})
}

func TestJSONCreated(t *testing.T) {
	wrapperGet(t, "/test", "/test", func(c IContext) {
		ret := struct {
			Message string `json:"message"`
		}{"nul"}
		c.JSONCreated(ret)
	}, func(t *testing.T, resp *http.Response) {
		z.AssertBody(t, resp, hBody)
		z.AssertStatusCode(t, resp, http.StatusCreated)
	})
}

func TestJSONUnprocessable(t *testing.T) {
	wrapperGet(t, "/test", "/test", func(c IContext) {
		ret := struct {
			Message string `json:"message"`
		}{"nul"}
		c.JSONUnprocessable(ret)
	}, func(t *testing.T, resp *http.Response) {
		z.AssertBody(t, resp, hBody)
		z.AssertStatusCode(t, resp, http.StatusUnprocessableEntity)
	})
}

func TestJSONOk(t *testing.T) {
	wrapperGet(t, "/test", "/test", func(c IContext) {
		ret := struct {
			Message string `json:"message"`
		}{"nul"}
		c.JSONOk(ret)
	}, func(t *testing.T, resp *http.Response) {
		z.AssertBody(t, resp, hBody)
		z.AssertStatusCode(t, resp, http.StatusOK)
	})
}

func TestJSONNotFound(t *testing.T) {
	wrapperGet(t, "/test", "/test", func(c IContext) {
		ret := struct {
			Message string `json:"message"`
		}{"nul"}
		c.JSONNotFound(ret)
	}, func(t *testing.T, resp *http.Response) {
		z.AssertBody(t, resp, hBody)
		z.AssertStatusCode(t, resp, http.StatusNotFound)
	})
}

func TestJSONConflict(t *testing.T) {
	wrapperGet(t, "/test", "/test", func(c IContext) {
		ret := struct {
			Message string `json:"message"`
		}{"nul"}
		c.JSONConflict(ret)
	}, func(t *testing.T, resp *http.Response) {
		z.AssertBody(t, resp, hBody)
		z.AssertStatusCode(t, resp, http.StatusConflict)
	})
}

func TestJSONInternalError(t *testing.T) {
	wrapperGet(t, "/test", "/test", func(c IContext) {
		ret := struct {
			Message string `json:"message"`
		}{"nul"}
		c.JSONInternalError(ret)
	}, func(t *testing.T, resp *http.Response) {
		z.AssertBody(t, resp, hBody)
		z.AssertStatusCode(t, resp, http.StatusInternalServerError)
	})
}
