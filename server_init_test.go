package webfmwk

import (
	"net/http"
	"testing"
	"time"

	"github.com/burgesQ/gommon/assert"
)

var (
	_emptyMiddleware = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		})
	}
	_emptyHandler = func(next HandlerFunc) HandlerFunc {
		return HandlerFunc(func(c Context) error {
			return next(c)
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

	assert.True(t, s.meta.ctrlc)
	assert.True(t, s.meta.checkIsUp)
	assert.True(t, s.meta.cors)
	assert.True(t, len(s.meta.middlewares) == 1)
	assert.True(t, len(s.meta.handlers) == 1)

	assert.Equal(t, s.meta.baseServer.ReadTimeout, testT)
	assert.Equal(t, s.meta.baseServer.WriteTimeout, testT)
	assert.Equal(t, s.meta.baseServer.ReadHeaderTimeout, testT)
	assert.Equal(t, s.meta.baseServer.MaxHeaderBytes, 42)

	assert.True(t, s.meta.docHandler != nil)

	assert.Equal(t, s.meta.prefix, "/api")
}
