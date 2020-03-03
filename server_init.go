package webfmwk

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/burgesQ/webfmwk/v3/log"
	"github.com/gorilla/mux"
)

type (
	// Option is tu be used this way :
	// s := w.InitServer(
	//   webfmwk.EnableCheckIsUp()
	//   webfmwk.WithCORS(),
	//   webfmwk.WithLogger(log.GetLogger()),
	//   webfmwk.WithMiddlewars(
	// 	  middleware.Logging,
	// 	  middleware.Security))
	Option func(s *Server)

	ServerMeta struct {
		ctrlc        bool
		checkIsUp    bool
		ctrlcStarted bool
		cors         bool
		middlewares  []mux.MiddlewareFunc
		setter       Setter
		docHandler   http.Handler
		baseServer   *http.Server
		prefix       string
		routes       RoutesPerPrefix
	}
)

var (
	nakedSetter = func(c *Context) IContext {
		return c
	}

	defaultMeta = func() ServerMeta {
		return ServerMeta{
			baseServer: &http.Server{
				ReadTimeout:       20 * time.Second,
				ReadHeaderTimeout: 20 * time.Second,
				WriteTimeout:      20 * time.Second,
				MaxHeaderBytes:    1 << 20,
			},
			setter: nakedSetter,
			routes: make(RoutesPerPrefix),
		}
	}
)

func InitServer(opts ...Option) *Server {
	var (
		wg          sync.WaitGroup
		ctx, cancel = context.WithCancel(context.Background())
		s           = &Server{
			launcher: CreateWorkerLauncher(&wg, cancel),
			ctx:      ctx,
			wg:       &wg,
			log:      logger,
			isReady:  make(chan bool),
			meta:     defaultMeta(),
		}
	)

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (m *ServerMeta) toServer(addr string) http.Server {
	return http.Server{
		Addr:           addr,
		ReadTimeout:    m.baseServer.ReadTimeout,
		WriteTimeout:   m.baseServer.WriteTimeout,
		MaxHeaderBytes: m.baseServer.MaxHeaderBytes,
	}
}

func WithCtrlC() Option {
	return func(s *Server) {
		s.meta.ctrlc = true
		s.log.Debugf("\t-- crtl-c support enabled")
	}
}

func CheckIsUp() Option {
	return func(s *Server) {
		s.EnableCheckIsUp()
		s.log.Debugf("\t-- check is up support enabled")
	}
}

func WithCORS() Option {
	return func(s *Server) {
		s.EnableCORS()
		s.log.Debugf("\t-- CORS support enabled")
	}
}

func WithMiddlewares(mw ...mux.MiddlewareFunc) Option {
	return func(s *Server) {
		s.RegisterMiddlewares(mw...)
		s.log.Debugf("\t-- middlewares loaded")
	}
}

func WithDocHandler(handler http.Handler) Option {
	return func(s *Server) {
		s.RegisterDocHandler(handler)
		s.log.Debugf("\t-- doc handler loaded")
	}
}

func WithLogger(lg log.ILog) Option {
	return func(s *Server) {
		s.RegisterLogger(lg)
		lg.Debugf("\t-- logger loaded")
	}
}

func WithCustomContext(setter Setter) Option {
	return func(s *Server) {
		s.SetCustomContext(setter)
		s.log.Debugf("\t-- custom context loaded")
	}
}

func SetMaxHeaderBytes(val int) Option {
	return func(s *Server) {
		s.meta.baseServer.MaxHeaderBytes = val
	}
}

// ReadTimeout is a timing constraint on the client http request imposed by the server from the moment
// of initial connection up to the time the entire request body has been read.
// [Accept] --> [TLS Handshake] --> [Request Headers] --> [Request Body] --> [Response]
func SetReadTimeout(val time.Duration) Option {
	return func(s *Server) {
		s.meta.baseServer.ReadTimeout = val
	}
}

func SetReadHeaderTimeout(val time.Duration) Option {
	return func(s *Server) {
		s.meta.baseServer.ReadHeaderTimeout = val
	}
}

// WriteTimeout is a time limit imposed on client connecting to the server via http from the
// time the server has completed reading the request header up to the time it has finished writing the response.
// [Accept] --> [TLS Handshake] --> [Request Headers] --> [Request Body] --> [Response]
func SetWriteTimeout(val time.Duration) Option {
	return func(s *Server) {
		s.meta.baseServer.WriteTimeout = val
	}
}

func SetPrefix(prefix string) Option {
	return func(s *Server) {
		s.meta.prefix = prefix
		s.log.Debugf("\t-- api prefix loaded")
	}
}
