package webfmwk

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/burgesQ/webfmwk/v3/log"
	"github.com/gorilla/mux"
)

type (
	// Option is tu be used this way :
	//   s := w.InitServer(
	//     webfmwk.WithLogger(log.GetLogger()),
	//     webfmwk.EnableCheckIsUp()
	//     webfmwk.WithCORS(),
	//     webfmwk.WithPrefix("/api"),
	//     webfmwk.WithMiddlewars(
	// 	    middleware.Logging,
	// 	    middleware.Security))
	Option func(s *Server)

	// ServerMeta hold the server meta information
	serverMeta struct {
		ctrlc        bool
		checkIsUp    bool
		ctrlcStarted bool
		cors         bool
		middlewares  []mux.MiddlewareFunc
		handlers     []Handler
		setter       Setter
		docHandler   http.Handler
		baseServer   *http.Server
		prefix       string
		routes       RoutesPerPrefix
	}
)

func getDefaultMeta() serverMeta {
	return serverMeta{
		baseServer: &http.Server{
			ReadTimeout:       20 * time.Second,
			ReadHeaderTimeout: 20 * time.Second,
			WriteTimeout:      20 * time.Second,
			MaxHeaderBytes:    1 << 20,
		},
		setter: func(c *Context) IContext {
			return c
		},
		routes:     make(RoutesPerPrefix),
		prefix:     "",
		docHandler: nil,
		handlers:   nil,
	}
}

func (m *serverMeta) toServer(addr string) http.Server {
	return http.Server{
		Addr:           addr,
		ReadTimeout:    m.baseServer.ReadTimeout,
		WriteTimeout:   m.baseServer.WriteTimeout,
		MaxHeaderBytes: m.baseServer.MaxHeaderBytes,
		ConnContext: func(ctx context.Context, c net.Conn) context.Context {
			log.Debugf("[+] new connection")

			return ctx
		},
	}
}

func applyOptions(s *Server, opts ...Option) {
	for _, opt := range opts {
		opt(s)
	}
}

// InitServer initialize a webfmwk.Server instance
// It take the server options as parameters.
// List of server options : WithLogger, WithCtrlC, CheckIsUp, WithCORS, SetPrefix,
// WithDocHandler, WithCustomContext, WithMiddlewares, WithHandlers,
// SetReadTimeout, SetWriteTimeout, SetMaxHeaderBytes, SetReadHeaderTimeout,
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
			meta:     getDefaultMeta(),
		}
	)

	applyOptions(s, opts...)

	return s
}

// WithLogger set the server logger which implement the log.ILog interface
// Try to set it the earliest posible.
func WithLogger(lg log.ILog) Option {
	return func(s *Server) {
		s.registerLogger(lg)
		lg.Debugf("\t-- logger loaded")
	}
}

// WithCtrlC enable the internal ctrl+c support from the server
func WithCtrlC() Option {
	return func(s *Server) {
		s.enableCtrlC()
		s.log.Debugf("\t-- crtl-c support enabled")
	}
}

// CheckIsUp expose a `/ping` endpoint and try to poll to check the server healt
// when it's started
func CheckIsUp() Option {
	return func(s *Server) {
		s.enableCheckIsUp()
		s.log.Debugf("\t-- check is up support enabled")
	}
}

// WithCORS enable the CORS (Cross-Origin Resource Sharing) support
func WithCORS() Option {
	return func(s *Server) {
		s.enableCORS()
		s.log.Debugf("\t-- CORS support enabled")
	}
}

// SetPrefix set the API root prefix
func SetPrefix(prefix string) Option {
	return func(s *Server) {
		s.setPrefix(prefix)
		s.log.Debugf("\t-- api prefix loaded")
	}
}

// WithDocHandler allow to register a http.Handler doc handler (ex: swaggo).
// If use with SetPrefix, register WithDocHandler after the SetPrefix one
func WithDocHandler(handler http.Handler) Option {
	return func(s *Server) {
		s.registerDocHandler(handler)
		s.log.Debugf("\t-- doc handler loaded")
	}
}

// WithCustomContext allow to register a custom context
//   package main
//
//   import "github.com/burgesQ/webfmwk/v3"
//
//   type CustomContext struct {
//     webfmwk.IContext
//     val string
//   }
//
//   func main() {
//     var s = webfmwk.InitServer(
//       webfmwk.WithCustomContext(func(c *webfmwk.Context) webfmwk.IContext {
//         return &CustomContext{*c, "custom"}
//       }))
//   }
//
func WithCustomContext(setter Setter) Option {
	return func(s *Server) {
		s.setCustomContext(setter)
		s.log.Debugf("\t-- custom context loaded")
	}
}

// WithMiddlewares allow to register a list of gorilla/mux.MiddlewareFunc.
// Middlwares signature is the http.Handler one (func(w http.ResponseWriterm r *http.Request))
//
//   package main
//
//   import (
//     "github.com/burgesQ/webfmwk/v3"
//     "github.com/burgesQ/webfmwk/v3/middleware"
//   )
//
//   func main() {
//     var s = webfmwk.InitServer(webfmwk.WithMiddlewares(middleware.Security))
//   }
func WithMiddlewares(mw ...mux.MiddlewareFunc) Option {
	return func(s *Server) {
		s.addMiddlewares(mw...)
		s.log.Debugf("\t-- middlewares loaded")
	}
}

// WithHandlers allow to register a list of webfmwk.Handler
// Handler signature is the webfmwk.HandlerFunc one (func(c IContext))
//
//   package main
//
//   import (
//     "github.com/burgesQ/webfmwk/v3"
//     "github.com/burgesQ/webfmwk/v3/handler"
//   )
//
//   func main() {
//     var s = webfmwk.InitServer(webfmwk.WithHandlers(handler.Logging, handler.RequestID))
//   }
func WithHandlers(h ...Handler) Option {
	return func(s *Server) {
		s.addHandlers(h...)
		s.log.Debugf("\t-- handlers loaded")
	}
}

// ReadTimeout is a timing constraint on the client http request imposed by the server from the moment
// of initial connection up to the time the entire request body has been read.
//
// [Accept] --> [TLS Handshake] --> [Request Headers] --> [Request Body] --> [Response]
func SetReadTimeout(val time.Duration) Option {
	return func(s *Server) {
		s.meta.baseServer.ReadTimeout = val
	}
}

// WriteTimeout is a time limit imposed on client connecting to the server via http from the
// time the server has completed reading the request header up to the time it has finished writing the response.
//
// [Accept] --> [TLS Handshake] --> [Request Headers] --> [Request Body] --> [Response]
func SetWriteTimeout(val time.Duration) Option {
	return func(s *Server) {
		s.meta.baseServer.WriteTimeout = val
	}
}

// SetMaxHeaderBytes set the max header bytes of both request and response
func SetMaxHeaderBytes(val int) Option {
	return func(s *Server) {
		s.meta.baseServer.MaxHeaderBytes = val
	}
}

// SetReadHeaderTimeout set the value of the timeout on the read header process
func SetReadHeaderTimeout(val time.Duration) Option {
	return func(s *Server) {
		s.meta.baseServer.ReadHeaderTimeout = val
	}
}
