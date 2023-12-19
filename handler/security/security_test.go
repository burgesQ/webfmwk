package security

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/burgesQ/gommon/webtest"
	"github.com/burgesQ/webfmwk/v6"
	"github.com/stretchr/testify/require"
)

const (
	_testPort = ":6668"
)

func TestHandler(t *testing.T) {
	s, e := webfmwk.InitServer(webfmwk.CheckIsUp(),
		webfmwk.SetPrefix("/api"),
		webfmwk.WithHandlers(Handler),
	)

	require.Nil(t, e)

	t.Cleanup(func() { require.Nil(t, s.ShutdownAndWait()) })

	s.GET("/testing", func(c webfmwk.Context) error {
		// never reach
		return c.JSONOk(json.RawMessage(`{}`))
	})

	go require.Nil(t, s.Start(_testPort))
	<-s.IsReady()

	webtest.RequestAndTestAPI(t, "http://127.0.0.1"+_testPort+"/api/testing",
		func(t *testing.T, resp *http.Response) {
			t.Helper()
			webtest.Headers(t, resp,
				[2]string{headerProtection, headerProtectionV},
				[2]string{headerSecu, headerSecuV},
				[2]string{headerOption, headerOptionV})
		})
}
