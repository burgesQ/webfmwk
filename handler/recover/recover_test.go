//nolint:predeclared
package recover

import (
	"net/http"
	"testing"

	"github.com/burgesQ/gommon/webtest"
	"github.com/burgesQ/webfmwk/v5"
	"github.com/stretchr/testify/require"
)

const _testPort = ":6671"

func TestHandler(t *testing.T) {
	s := webfmwk.InitServer(webfmwk.CheckIsUp(),
		webfmwk.SetPrefix("/api"),
		webfmwk.WithHandlers(Handler),
	)

	// t.Log("init server...")
	t.Cleanup(func() {
		ctx := s.GetContext()
		// t.Log("closing server ...")
		ctx.Done()
		s.Shutdown()
		s.WaitAndStop()
		webfmwk.Shutdown()
		// t.Log("server closed")
	})

	s.GET("/testing/string", func(c webfmwk.Context) error {
		panic("some fatal")
	})

	s.GET("/testing/error", func(c webfmwk.Context) error {
		panic(webfmwk.NewForbidden(webfmwk.NewError("some fatal error")))
	})

	go s.Start(_testPort)
	<-s.IsReady()
	// t.Log("server inited")

	t.Run("testing panic over string ", func(t *testing.T) {
		webtest.RequestAndTestAPI(t, "http://127.0.0.1"+_testPort+"/api/testing/string",
			func(t *testing.T, resp *http.Response) {
				t.Helper()
				webtest.StatusCode(t, http.StatusInternalServerError, resp)

				body := webtest.FetchBody(t, resp)

				require.Contains(t, body, "some fatal")
				require.Contains(t, body, "status\":500")
			})
	})

	t.Run("testing panic over error hanlded", func(t *testing.T) {
		webtest.RequestAndTestAPI(t, "http://127.0.0.1"+_testPort+"/api/testing/error",
			func(t *testing.T, resp *http.Response) {
				t.Helper()
				webtest.StatusCode(t, http.StatusForbidden, resp)

				body := webtest.FetchBody(t, resp)

				require.Contains(t, body, "some fatal")
				require.Contains(t, body, "status\":403")
			})
	})
}
