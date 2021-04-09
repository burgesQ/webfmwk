package recover

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/burgesQ/webfmwk/v5"
	"github.com/stretchr/testify/assert"
)

const _testPort = ":6671"

func TestHandler(t *testing.T) {
	var (
		s = webfmwk.InitServer(webfmwk.CheckIsUp(),
			webfmwk.DisableKeepAlive(),
			webfmwk.SetPrefix("/api"),
			webfmwk.WithHandlers(Handler),
		)
	)

	// t.Log("init server...")
	t.Cleanup(func() {
		var ctx = s.GetContext()
		// t.Log("closing server ...")
		ctx.Done()
		s.Shutdown()
		s.WaitAndStop()
		webfmwk.Shutdown()
		// t.Log("server closed")
	})

	s.GET("/testing/string", func(c webfmwk.Context) error {
		panic("some fatal")
		// never reached
		return c.JSONOk(json.RawMessage(`{}`))
	})

	s.GET("/testing/error", func(c webfmwk.Context) error {
		panic(webfmwk.NewForbidden(webfmwk.NewError("some fatal error")))
		// never reached
		return c.JSONOk(json.RawMessage(`{}`))
	})

	go s.Start(_testPort)
	<-s.IsReady()
	// t.Log("server inited")

	t.Run("testing panic over string ", func(t *testing.T) {
		resp, err := http.Get("http://127.0.0.1" + _testPort + "/api/testing/string")
		if err != nil {
			t.Errorf("error requesting the api : %s", err.Error())
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("reading body: " + err.Error())
		}

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.Contains(t, string(body), "some fatal")
		assert.Contains(t, string(body), "status\":500")
	})

	t.Run("testing panic over error hanlded", func(t *testing.T) {
		resp, err := http.Get("http://127.0.0.1" + _testPort + "/api/testing/error")
		if err != nil {
			t.Errorf("error requesting the api : %s", err.Error())
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("reading body: " + err.Error())
		}

		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
		assert.Contains(t, string(body), "some fatal error")
		assert.Contains(t, string(body), "status\":403")
	})
}
