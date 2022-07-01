package webfmwk

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/burgesQ/gommon/webtest"
)

func TestRunner(t *testing.T) {
	var s = InitServer(CheckIsUp())

	t.Cleanup(func() { stopServer(s) })

	s.GET("/test", func(c Context) error {
		return c.JSONOk(json.RawMessage(`{"value":"test"}`))
	})

	go s.Run(Address{Addr: ":6661"}, Address{Addr: ":6662"})

	<-s.isReady

	webtest.RequestAndTestAPI(t, "http://127.0.0.1:6661/test",
		func(t *testing.T, resp *http.Response) {
			webtest.StatusCode(t, http.StatusOK, resp)
		})

	webtest.RequestAndTestAPI(t, "http://127.0.0.1:6662/test",
		func(t *testing.T, resp *http.Response) {
			webtest.StatusCode(t, http.StatusOK, resp)
		})
}
