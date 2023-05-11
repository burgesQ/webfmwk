package webfmwk

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
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
		s, e        = InitServer(
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

	requirer := require.New(t)

	requirer.Nil(e)

	requirer.True(s.meta.ctrlc)
	requirer.True(s.meta.checkIsUp)
	requirer.True(s.meta.cors)
	// require.True(t, len(s.meta.middlewares) == 1)
	requirer.True(len(s.meta.handlers) == 1)

	requirer.Equal(s.meta.baseServer.ReadTimeout, testT)
	requirer.Equal(s.meta.baseServer.WriteTimeout, testT)
	requirer.Equal(s.meta.baseServer.IdleTimeout, testT)

	requirer.True(len(s.meta.docHandlers) == 1)

	requirer.Equal(s.meta.prefix, "/api")

	ht := s.meta.toServer("testing")
	requirer.Equal("webfmwk testing", ht.Name)
	requirer.Equal(testT, ht.ReadTimeout)
	requirer.Equal(testT, ht.WriteTimeout)
	requirer.Equal(testT, ht.IdleTimeout)
}
