package webfmwk

import (
	"net/http"
	"testing"
	"time"

	z "github.com/burgesQ/webfmwk/v3/testing"
)

func TestServerNewInit(t *testing.T) {
	var (
		testT = time.Second * 42
		s     = InitServer(
			WithCtrlC(), CheckIsUp(), WithCORS(),
			SetMaxHeaderBytes(42), SetReadTimeout(testT),
			SetWriteTimeout(testT), SetReadHeaderTimeout(testT),
			WithMiddlewars(func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				})
			}),
		)
	)
	z.AssertTrue(t, s.meta.ctrlc)
	z.AssertTrue(t, s.meta.checkIsUp)
	z.AssertTrue(t, s.meta.cors)
	z.AssertTrue(t, len(s.meta.middlewares) == 1)

	z.AssertEqual(t, s.meta.baseServer.ReadTimeout, testT)
	z.AssertEqual(t, s.meta.baseServer.WriteTimeout, testT)
	z.AssertEqual(t, s.meta.baseServer.ReadHeaderTimeout, testT)
	z.AssertEqual(t, s.meta.baseServer.MaxHeaderBytes, 42)
}
