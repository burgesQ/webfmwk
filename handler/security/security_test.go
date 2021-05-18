package security

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/burgesQ/gommon/webtest"
	"github.com/burgesQ/webfmwk/v5"
)

const (
	_testPort = ":6668"
)

func TestHandler(t *testing.T) {
	var (
		s = webfmwk.InitServer(webfmwk.CheckIsUp(),
			webfmwk.SetPrefix("/api"),
			webfmwk.WithHandlers(Handler),
		)
	)

	t.Cleanup(func() {
		var ctx = s.GetContext()

		ctx.Done()
		s.Shutdown()
		s.WaitAndStop()
		webfmwk.Shutdown()
	})

	s.GET("/testing", func(c webfmwk.Context) error {
		// never reach
		return c.JSONOk(json.RawMessage(`{}`))
	})

	go s.Start(_testPort)
	<-s.IsReady()

	webtest.RequestAndTestAPI(t, "http://127.0.0.1"+_testPort+"/api/testing",
		func(t *testing.T, resp *http.Response) {
			webtest.Headers(t, resp,
				[2]string{headerProtection, headerProtectionV},
				[2]string{headerSecu, headerSecuV},
				[2]string{headerOption, headerOptionV})
		})
}
