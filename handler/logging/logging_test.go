package logging

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/burgesQ/webfmwk/v5"
	"github.com/stretchr/testify/assert"
)

const _testPort = ":6670"

func TestHandler(t *testing.T) {
	var (
		s = webfmwk.InitServer(webfmwk.CheckIsUp(),
			webfmwk.DisableKeepAlive(),
			webfmwk.SetPrefix("/api"),
			webfmwk.WithHandlers(Handler),
		)
	)

	t.Log("init server...")
	defer func() {
		var ctx = s.GetContext()
		t.Log("closing server ...")
		ctx.Done()
		s.Shutdown()
		s.WaitAndStop()
		webfmwk.Shutdown()
		t.Log("server closed")
	}()

	s.GET("/testing", func(c webfmwk.Context) error {
		return c.JSONOk(json.RawMessage(`{}`))
	})

	go s.Start(_testPort)
	<-s.IsReady()
	t.Log("server inited")

	// req
	resp, err := http.Get("http://127.0.0.1" + _testPort + "/api/testing")
	if err != nil {
		t.Errorf("error requesting the api : %s", err.Error())
	}
	defer resp.Body.Close()

	assert.Contains(t, resp.Header, HeaderRequestID)
}
