package logging

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/burgesQ/gommon/webtest"
	"github.com/burgesQ/webfmwk/v5"
	"github.com/stretchr/testify/assert"
)

const _testPort = ":6670"

func TestHandler(t *testing.T) {
	s := webfmwk.InitServer(webfmwk.CheckIsUp(),
		webfmwk.SetPrefix("/api"),
		webfmwk.WithHandlers(Handler),
	)

	t.Cleanup(func() {
		ctx := s.GetContext()

		ctx.Done()
		s.Shutdown()
		s.WaitAndStop()
		webfmwk.Shutdown()
	})

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
