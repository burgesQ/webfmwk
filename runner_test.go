package webfmwk

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/burgesQ/gommon/webtest"
	"github.com/stretchr/testify/require"
)

func TestRunner(t *testing.T) {
	s, e := InitServer(CheckIsUp())

	require.Nil(t, e)
	t.Cleanup(func() { require.Nil(t, s.ShutAndWait()) })

	s.GET("/test", func(c Context) error {
		return c.JSONOk(json.RawMessage(`{"value":"test"}`))
	})

	go s.Run(Address{Addr: ":6661"}, Address{Addr: ":6662"})

	<-s.isReady

	webtest.RequestAndTestAPI(t, "http://127.0.0.1:6661/test",
		func(t *testing.T, resp *http.Response) {
			t.Helper()
			webtest.StatusCode(t, http.StatusOK, resp)
		})

	webtest.RequestAndTestAPI(t, "http://127.0.0.1:6662/test",
		func(t *testing.T, resp *http.Response) {
			t.Helper()
			webtest.StatusCode(t, http.StatusOK, resp)
		})
}
