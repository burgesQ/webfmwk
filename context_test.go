package webfmwk

import (
	"context"
	"log/slog"
	"net/http"
	"testing"

	"github.com/burgesQ/gommon/webtest"
	"github.com/stretchr/testify/assert"
)

var (
	hBody      = `{"message":"nul"}`
	jsonEncode = "application/json; charset=UTF-8"
)

func TestParam(t *testing.T) {
	wrapperGet(t, "/test/{id}", "/test/tutu", func(c Context) error {
		id := c.GetVar("id")
		if id != "tutu" {
			t.Errorf("error fetching the url param : [%s] expected [tutu]", id)
		}

		return c.JSONOk(id)
	}, func(t *testing.T, resp *http.Response) {
		t.Helper()
		webtest.Body(t, `"tutu"`, resp)
		webtest.StatusCode(t, http.StatusOK, resp)
	})
}

func TestLogger(t *testing.T) {
	var (
		c      = icontext{}
		logger = slog.Default()
	)

	c.SetStructuredLogger(logger)
	assert.True(t, logger == c.GetStructuredLogger(), "context logger should be the setted one")
	// assert.True(t, logger == GetLogger(), "default logger should be the setted one")
}

func TestContext(t *testing.T) {
	var (
		ctx context.Context
		c   = icontext{
			ctx: ctx,
		}
	)
	assert.True(t, ctx == c.GetContext())
}

func TestFetchContent(t *testing.T) {
	const (
		_ok = iota
		_unprocessable
	)

	tests := map[string]struct {
		payload []byte
		t       int
	}{
		"fetch content":               {[]byte(`{"first_name": "tutu"}`), _ok},
		"fetch content unprocessable": {[]byte(`{"first_name": tutu"}`), _unprocessable},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			wrapperPost(t, "/test", "/test", test.payload, func(c Context) error {
				anonymous := struct {
					FirstName string `json:"first_name,omitempty" validate:"required"`
				}{}
				if e := c.FetchContent(&anonymous); e != nil {
					return e
				} else if e := c.Validate(anonymous); e != nil {
					return e
				}

				return c.JSON(http.StatusCreated, anonymous)
			}, func(t *testing.T, resp *http.Response) {
				t.Helper()

				switch test.t {
				case _ok:
					webtest.Body(t, `{"first_name":"tutu"}`, resp)
					webtest.StatusCode(t, http.StatusCreated, resp)

				case _unprocessable:
					webtest.StatusCode(t, http.StatusUnprocessableEntity, resp)
				}
			})
		})
	}
}

func TestJSONBlobPretty(t *testing.T) {
	wrapperGet(t, "/test", "/test?pretty", func(c Context) error {
		return c.JSONBlob(http.StatusOK, []byte(hBody))
	}, func(t *testing.T, resp *http.Response) {
		t.Helper()
		webtest.BodyDiffere(t, hBody, resp)
		webtest.StatusCode(t, http.StatusOK, resp)
	})
}
