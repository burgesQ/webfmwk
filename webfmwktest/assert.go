package webfmwktest

import (
	"net/http/httptest"
	"testing"

	z "github.com/burgesQ/gommon/testing"
)

func AssertBody(t *testing.T, rr *httptest.ResponseRecorder, expected string) {
	t.Helper()
	z.AssertStringEqual(t, rr.Body.String(), expected)
}

func AssertBodyDiffer(t *testing.T, rr *httptest.ResponseRecorder, expected string) {
	t.Helper()
	z.AssertStringNotEqual(t, rr.Body.String(), expected)
}

// AssertStatusCode assert the status code of the response
func AssertStatusCode(t *testing.T, rr *httptest.ResponseRecorder, expected int) {
	t.Helper()
	z.AssertIntEqual(t, rr.Code, expected)
}
