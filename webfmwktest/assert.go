package webfmwktest

import (
	"net/http/httptest"
	"testing"

	"github.com/burgesQ/gommon/assert"
)

func AssertBody(t *testing.T, rr *httptest.ResponseRecorder, expected string) {
	t.Helper()
	assert.StringEqual(t, rr.Body.String(), expected)
}

func AssertBodyDiffer(t *testing.T, rr *httptest.ResponseRecorder, expected string) {
	t.Helper()
	assert.StringNotEqual(t, rr.Body.String(), expected)
}

// AssertStatusCode assert the status code of the response
func AssertStatusCode(t *testing.T, rr *httptest.ResponseRecorder, expected int) {
	t.Helper()
	assert.IntEqual(t, rr.Code, expected)
}
