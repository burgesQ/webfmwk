package logging

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/burgesQ/gommon/webtest"
	"github.com/burgesQ/webfmwk/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const _testPort = ":6670"

func TestHandler(t *testing.T) {
	s, e := webfmwk.InitServer(webfmwk.CheckIsUp(),
		webfmwk.SetPrefix("/api"),
		webfmwk.WithHandlers(Handler),
	)

	require.Nil(t, e)

	t.Cleanup(func() { require.Nil(t, s.ShutAndWait()) })

	s.GET("/testing", func(c webfmwk.Context) error {
		return c.JSONOk(json.RawMessage(`{}`))
	})

	go s.Start(_testPort)
	<-s.IsReady()

	webtest.RequestAndTestAPI(t, "http://127.0.0.1"+_testPort+"/api/testing",
		func(t *testing.T, resp *http.Response) {
			t.Helper()
			assert.Contains(t, resp.Header, HeaderRequestID)
		})
}
