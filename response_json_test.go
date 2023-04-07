package webfmwk

import (
	"net/http"
	"testing"

	"github.com/burgesQ/gommon/webtest"
)

func TestJSONResponse(t *testing.T) {
	var (
		s   = InitServer(CheckIsUp())
		ret = struct {
			Message string `json:"message"`
		}{"nul"}
		tests = map[string]struct {
			expectedOP int
			fn         func(c Context, ret interface{}) error
		}{
			"blob": {http.StatusOK, func(c Context, ret interface{}) error {
				return c.JSONBlob(http.StatusOK, []byte(hBody))
			}},
			"ok": {http.StatusOK, func(c Context, ret interface{}) error {
				return c.JSONOk(ret)
			}},
			"created": {http.StatusCreated, func(c Context, ret interface{}) error {
				return c.JSONCreated(ret)
			}},
			"accepted": {http.StatusAccepted, func(c Context, ret interface{}) error {
				return c.JSONAccepted(ret)
			}},
			"noContent": {http.StatusNoContent, func(c Context, ret interface{}) error {
				return c.JSONNoContent()
			}},
			"badRequest": {http.StatusBadRequest, func(c Context, ret interface{}) error {
				return c.JSONBadRequest(ret)
			}},
			"unauthorized": {http.StatusUnauthorized, func(c Context, ret interface{}) error {
				return c.JSONUnauthorized(ret)
			}},
			"forbidden": {http.StatusForbidden, func(c Context, ret interface{}) error {
				return c.JSONForbidden(ret)
			}},
			"notFound": {http.StatusNotFound, func(c Context, ret interface{}) error {
				return c.JSONNotFound(ret)
			}},
			"conflict": {http.StatusConflict, func(c Context, ret interface{}) error {
				return c.JSONConflict(ret)
			}},
			"unprocessable": {http.StatusUnprocessableEntity, func(c Context, ret interface{}) error {
				return c.JSONUnprocessable(ret)
			}},
			"internalError": {http.StatusInternalServerError, func(c Context, ret interface{}) error {
				return c.JSONInternalError(ret)
			}},
			"notImplemented": {http.StatusNotImplemented, func(c Context, ret interface{}) error {
				return c.JSONNotImplemented(ret)
			}},
		}
	)

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
