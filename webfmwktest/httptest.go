package webfmwktest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/burgesQ/webfmwk/v4"
)

// testing method ... Maybe create a new package ?

type Expected struct {
	Code int
	Body string
}

// GetAndTest perfome a get request on the webfmwk using the httptest package
func GetAndTest(t *testing.T, h webfmwk.HandlerFunc, e Expected) {
	t.Helper()

	var check = func(ctx string, err error) {
		if err != nil {
			t.Fatal(ctx + err.Error())
		}
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter
	req, err := http.NewRequest("GET", "/tests", nil)
	check("cannot create http requesst : ", err)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response
	rr := httptest.NewRecorder()
	// We create a webfmwk server to wrapper the provided handler
	s := webfmwk.InitServer()
	// We wrap the webfmwk handler into the http one
	handler := http.HandlerFunc(s.CustomHandler(h))

	handler.ServeHTTP(rr, req)

	// we assert the content
	AssertStatusCode(t, rr, e.Code)
	AssertBody(t, rr, e.Body)
}
