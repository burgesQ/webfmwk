package redoc

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/burgesQ/webfmwk/v5"
	"github.com/stretchr/testify/assert"
)

const (
	_testPort = ":6667"
)

func TestRedocParam(t *testing.T) {
	def := GetHandler()

	assert.Equal(t, _defPath, def.Path)
	assert.Contains(t, string(genContent(_defRedoc)),
		"<redoc spec-url="+_defURI+"></redoc>")
}

func TestGetHandler(t *testing.T) {
	var (
		s = webfmwk.InitServer(webfmwk.CheckIsUp(),
			webfmwk.DisableKeepAlive(),
			webfmwk.SetPrefix("/api"),
			webfmwk.WithDocHandlers(
				GetHandler(Path("/another"), DocURI("/source")),
			),
		)
	)

	t.Cleanup(func() {
		var ctx = s.GetContext()
		ctx.Done()
		s.Shutdown()
		s.WaitAndStop()
		webfmwk.Shutdown()
	})

	s.GET("/source", func(c webfmwk.Context) error {
		return c.JSONOk(json.RawMessage(`{}`))
	})

	go s.Start(_testPort)
	<-s.IsReady()

	resp, err := http.Get("http://127.0.0.1" + _testPort + "/api/another")
	if err != nil {
		t.Fatalf("error requesting the api : %s", err.Error())
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("reading body: " + err.Error())
	}

	assert.Contains(t, string(bodyBytes), "/source")
}
