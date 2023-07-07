package webfmwk

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var _emptyHandler = func(next HandlerFunc) HandlerFunc {
	return HandlerFunc(func(c Context) error {
		return next(c)
	})
}

type ws struct{}

func (ws) ServeHTTP(http.ResponseWriter, *http.Request) {}

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
			EnableKeepAlive(),
			WithHTTP2(),
			EnablePprof("/some/path"),
			WithHandlers(_emptyHandler),
			MaxRequestBodySize(42),
			WithSocketHandler("/sock_1", ws{}),
			WithSocketHandlerFunc("/sock_1", func(http.ResponseWriter, *http.Request) {}),
		)
	)

	requirer := require.New(t)

	requirer.Nil(e)

	requirer.True(s.meta.pprof, "pprof should be enabled")
	requirer.True(s.meta.enableKeepAlive, "keep alive should be enabled")
	requirer.True(s.meta.http2, "http2 should be enabled")
	requirer.True(s.meta.ctrlc, "ctrl+c support enabled")
	requirer.True(s.meta.checkIsUp, "ping pong start enabled")
	requirer.True(s.meta.cors, "cors enabled")
	// require.True(t, len(s.meta.middlewares) == 1)
	requirer.True(len(s.meta.handlers) == 1, "one handler should be loaded")

	requirer.Equal(testT, s.meta.baseServer.ReadTimeout)
	requirer.Equal(testT, s.meta.baseServer.WriteTimeout)
	requirer.Equal(testT, s.meta.baseServer.IdleTimeout, "idle timeout should be custom")
	requirer.Equal(42, s.meta.baseServer.MaxRequestBodySize, "body max size should be 42")

	requirer.True(len(s.meta.docHandlers) == 1)

	requirer.Equal("/api", s.meta.prefix, "prefix should be custom")

	ht := s.meta.toServer("testing")
	requirer.Equal("webfmwk testing", ht.Name)
	requirer.Equal(testT, ht.ReadTimeout)
	requirer.Equal(testT, ht.WriteTimeout)
	requirer.Equal(testT, ht.IdleTimeout)
}
