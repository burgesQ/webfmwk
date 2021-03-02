package security

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/burgesQ/webfmwk/v5"
	"github.com/stretchr/testify/assert"
)

const (
	_testPort = ":6668"
)

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
		// never reach
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

	for _, h := range [][2]string{
		{HeaderProtection, HeaderProtectionV},
		{HeaderSecu, HeaderSecuV},
		{HeaderOption, HeaderOptionV},
	} {
		key := h[0]
		val := h[1]
		assert.Contains(t, resp.Header, key, "asserting header %q is present", key)
		assert.Equal(t, val, resp.Header[key][0], "asserting security header %q value is %q", key, val)
	}

}
