package webfmwk

import (
	"net/http"
	"testing"

	"github.com/burgesQ/gommon/webtest"
	"github.com/stretchr/testify/require"
)

func TestJSONResponse(t *testing.T) {
	var (
		s, e = InitServer(CheckIsUp())
		ret  = struct {
			Message string `json:"message"`
		}{"nul"}
		tests = map[string]struct {
			fn         func(c Context, ret interface{}) error
			expectedOP int
		}{
			"blob": {expectedOP: http.StatusOK, fn: func(c Context, ret interface{}) error {
				return c.JSONBlob(http.StatusOK, []byte(hBody))
			}},
			"ok": {expectedOP: http.StatusOK, fn: func(c Context, ret interface{}) error {
				return c.JSONOk(ret)
			}},
			"created": {expectedOP: http.StatusCreated, fn: func(c Context, ret interface{}) error {
				return c.JSONCreated(ret)
			}},
			"accepted": {expectedOP: http.StatusAccepted, fn: func(c Context, ret interface{}) error {
				return c.JSONAccepted(ret)
			}},
			"noContent": {expectedOP: http.StatusNoContent, fn: func(c Context, ret interface{}) error {
				return c.JSONNoContent()
			}},
			"badRequest": {expectedOP: http.StatusBadRequest, fn: func(c Context, ret interface{}) error {
				return c.JSONBadRequest(ret)
			}},
			"unauthorized": {expectedOP: http.StatusUnauthorized, fn: func(c Context, ret interface{}) error {
				return c.JSONUnauthorized(ret)
			}},
			"forbidden": {expectedOP: http.StatusForbidden, fn: func(c Context, ret interface{}) error {
				return c.JSONForbidden(ret)
			}},
			"notFound": {expectedOP: http.StatusNotFound, fn: func(c Context, ret interface{}) error {
				return c.JSONNotFound(ret)
			}},
			"conflict": {expectedOP: http.StatusConflict, fn: func(c Context, ret interface{}) error {
				return c.JSONConflict(ret)
			}},
			"unprocessable": {expectedOP: http.StatusUnprocessableEntity, fn: func(c Context, ret interface{}) error {
				return c.JSONUnprocessable(ret)
			}},
			"internalError": {expectedOP: http.StatusInternalServerError, fn: func(c Context, ret interface{}) error {
				return c.JSONInternalError(ret)
			}},
			"notImplemented": {expectedOP: http.StatusNotImplemented, fn: func(c Context, ret interface{}) error {
				return c.JSONNotImplemented(ret)
			}},
		}
	)

	require.Nil(t, e)
	t.Cleanup(func() { stopServer(s) })

	// load custom endpoints
	for n, t := range tests {
		fn := t.fn
		s.GET("/"+n, func(c Context) error {
			return fn(c, ret)
		})
	}

	go s.Start(_testPort)
	<-s.isReady

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			webtest.RequestAndTestAPI(t, _testAddr+"/"+name,
				func(t *testing.T, resp *http.Response) {
					t.Helper()

					if test.expectedOP != http.StatusNoContent {
						webtest.Body(t, hBody, resp)
					}
					webtest.StatusCode(t, test.expectedOP, resp)
				})
		})
	}
}
