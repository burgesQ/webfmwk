package webfmwk

import (
	"net/http"
	"testing"
	"time"

	z "github.com/burgesQ/webfmwk/v3/testing"
)

var (
	_emptyMiddleware = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		})
	}
	_emptyHandler = func(next HandlerFunc) HandlerFunc {
		return HandlerFunc(func(c IContext) {
		})
	}
)

func TestServerNewInit(t *testing.T) {
	var (
		fakeswagger = http.HandlerFunc(func(q http.ResponseWriter, r *http.Request) {})
		testT       = time.Second * 42
		s           = InitServer(
			WithCtrlC(), CheckIsUp(), WithCORS(),
			WithDocHandler(fakeswagger),
			SetMaxHeaderBytes(42), SetReadTimeout(testT),
			SetWriteTimeout(testT), SetReadHeaderTimeout(testT),
			SetPrefix("/api"),
			WithMiddlewares(_emptyMiddleware),
			WithHandlers(_emptyHandler))
	)

	z.AssertTrue(t, s.meta.ctrlc)
	z.AssertTrue(t, s.meta.checkIsUp)
	z.AssertTrue(t, s.meta.cors)
	z.AssertTrue(t, len(s.meta.middlewares) == 1)
	z.AssertTrue(t, len(s.meta.handlers) == 1)

	z.AssertEqual(t, s.meta.baseServer.ReadTimeout, testT)
	z.AssertEqual(t, s.meta.baseServer.WriteTimeout, testT)
	z.AssertEqual(t, s.meta.baseServer.ReadHeaderTimeout, testT)
	z.AssertEqual(t, s.meta.baseServer.MaxHeaderBytes, 42)

	z.AssertTrue(t, s.meta.docHandler != nil)

	z.AssertEqual(t, s.meta.prefix, "/api")
}
