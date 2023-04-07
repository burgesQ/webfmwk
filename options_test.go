package webfmwk

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var _emptyHandler = func(next HandlerFunc) HandlerFunc {
	return HandlerFunc(func(c Context) error {
		return next(c)
	})
}

func TestServerNewInit(t *testing.T) {
	var (
		fakeswagger = DocHandler{}
		testT       = time.Second * 42
		s           = InitServer(
			WithCtrlC(), CheckIsUp(), WithCORS(),
			WithDocHandlers(fakeswagger),
			// SetMaxHeaderBytes(42),
			SetReadTimeout(testT),
			SetWriteTimeout(testT),
			SetIDLETimeout(testT),
			// SetReadHeaderTimeout(testT),
			SetPrefix("/api"),
			//			WithMiddlewares(_emptyMiddleware),
			WithHandlers(_emptyHandler))
	)

	asserter := assert.New(t)

	asserter.True(s.meta.ctrlc)
	asserter.True(s.meta.checkIsUp)
	asserter.True(s.meta.cors)
	// assert.True(t, len(s.meta.middlewares) == 1)
	asserter.True(len(s.meta.handlers) == 1)

	asserter.Equal(s.meta.baseServer.ReadTimeout, testT)
	asserter.Equal(s.meta.baseServer.WriteTimeout, testT)
	asserter.Equal(s.meta.baseServer.IdleTimeout, testT)

	asserter.True(len(s.meta.docHandlers) == 1)

	asserter.Equal(s.meta.prefix, "/api")

	ht := s.meta.toServer("testing")
	asserter.Equal("webfmwk testing", ht.Name)
	asserter.Equal(testT, ht.ReadTimeout)
	asserter.Equal(testT, ht.WriteTimeout)
	asserter.Equal(testT, ht.IdleTimeout)
}
