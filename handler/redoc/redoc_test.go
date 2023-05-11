package redoc

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/burgesQ/gommon/webtest"
	"github.com/burgesQ/webfmwk/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	_testPort = ":6667"
)

func TestRedocParam(t *testing.T) {
	def := GetHandler()

	assert.Equal(t, _defPath, def.Path)
	assert.Contains(t, string(genContent(defRedoc())),
		"<redoc spec-url="+_defURI+"></redoc>")
}

func TestGetHandler(t *testing.T) {
	s, e := webfmwk.InitServer(webfmwk.CheckIsUp(),
		webfmwk.SetPrefix("/api"),
		webfmwk.WithDocHandlers(
			GetHandler(Path("/another"), DocURI("/source")),
		),
	)

	require.Nil(t, e)

	t.Cleanup(func() { require.Nil(t, s.ShutAndWait()) })

	s.GET("/source", func(c webfmwk.Context) error {
		return c.JSONOk(json.RawMessage(`{}`))
	})

	go s.Start(_testPort)
	<-s.IsReady()

	webtest.RequestAndTestAPI(t, "http://127.0.0.1"+_testPort+"/api/another",
		func(t *testing.T, resp *http.Response) {
			t.Helper()
			webtest.BodyContains(t, "/source", resp)
		})
}
