package webfmwk

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/burgesQ/log"
	"github.com/valyala/fasthttp"
)

type (
	// Option apply specific configuration to the server at init time
	// They are tu be used this way :
	//   s := w.InitServer(
	//     webfmwk.WithLogger(log.GetLogger()),
	//     webfmwk.WithCtrlC(),
	//     webfmwk.CheckIsUp(),
	//     webfmwk.WithCORS(),
	//     webfmwk.SetPrefix("/api"),
	//     webfmwk.WithDocHanlders(redoc.GetHandler()),
	//     webfmwk.SetIdleTimeout(1 * time.Second),
	//     webfmwk.SetReadTimeout(1 * time.Second),
	//     webfmwk.SetWriteTimeout(1 * time.Second),
	//     webfmwk.WithHanlders(
	//       recover.Handler,
	//       logging.Handler,
	//       security.Handler))
	Option func(s *Server)

	// Options is a list of options
	Options []Option

	serverMeta struct {
		socketIOHandler     http.Handler
		socketIOHandlerFunc http.HandlerFunc
		baseServer          *fasthttp.Server
		routes              RoutesPerPrefix
		prefix              string
		pprofPath           string
		socketIOPath        string
		docHandlers         []DocHandler
		handlers            []Handler
		cors                bool
		socketIOHF          bool
		socketIOH           bool
		pprof               bool
		enableKeepAlive     bool
		ctrlcStarted        bool
		checkIsUp           bool
		ctrlc               bool
	}
)

var once sync.Once

func initOnce() error { return initValidator() }

// UseOption apply the param o option to the params s server
func UseOption(s *Server, o Option) {
	o(s)
}

// useOptions apply the params opts option to the param s server
func useOptions(s *Server, opts ...Option) {
	for _, o := range opts {
		UseOption(s, o)
	}
}

// InitServer initialize a webfmwk.Server instance.
// It may take some server options as parameters.
// List of server options : WithLogger, WithCtrlC, CheckIsUp, WithCORS, SetPrefix,
// WithHandlers, WithDocHandler, SetReadTimeout, SetWriteTimeout, SetIdleTimeout,
// EnableKeepAlive.
// Any error returned by the method should be handled as a fatal one.
func InitServer(opts ...Option) (*Server, error) {
	var e error

	once.Do(func() { e = initOnce() })

	var (
		wg          sync.WaitGroup
		ctx, cancel = context.WithCancel(context.Background())
		s           = &Server{
			launcher: CreateWorkerLauncher(&wg, cancel),
			ctx:      ctx,
			cancel:   cancel,
			wg:       &wg,
			log:      fetchLogger(),
			isReady:  make(chan bool),
			meta:     getDefaultMeta(),
		}
	)

	useOptions(s, opts...)

	return s, e
}

// WithLogger set the server logger which implement the log.Log interface
// Try to set it the earliest posible.
func WithLogger(lg log.Log) Option {
	return func(s *Server) {
		s.registerLogger(lg)
		lg.Debugf("\t-- logger loaded")
	}
}

// WithCtrlC enable the internal ctrl+c support from the server.
func WithCtrlC() Option {
	return func(s *Server) {
		s.enableCtrlC()
		s.log.Debugf("\t-- crtl-c support enabled")
	}
}

// CheckIsUp expose a `/ping` endpoint and try to poll to check the server healt
// when it's started.
func CheckIsUp() Option {
	return func(s *Server) {
		s.enableCheckIsUp()
		s.log.Debugf("\t-- check is up support enabled")
	}
}

// WithCORS enable the CORS (Cross-Origin Resource Sharing) support.
func WithCORS() Option {
	return func(s *Server) {
		s.enableCORS()
		s.log.Debugf("\t-- CORS support enabled")
	}
}

// SetPrefix set the API root prefix.
func SetPrefix(prefix string) Option {
	return func(s *Server) {
		s.setPrefix(prefix)
		s.log.Debugf("\t-- api prefix loaded")
	}
}

// WithDocHandlers allow to register custom DocHandler struct (ex: swaggo, redoc).
// If use with SetPrefix, register WithDocHandler after the SetPrefix one.
// Example:
//
//	package main
//
//	import (
//		"github.com/burgesQ/webfmwk/v5"
//		"github.com/burgesQ/webfmwk/v5/handler/redoc"
//	)
//
//	func main() {
//		var s = webfmwk.InitServer(webfmwk.WithDocHandlers(redoc.GetHandler()))
//	}
func WithDocHandlers(handler ...DocHandler) Option {
	return func(s *Server) {
		s.addDocHandlers(handler...)
		s.log.Debugf("\t-- doc handlers loaded")
	}
}

// WithHandlers allow to register a list of webfmwk.Handler
// Handler signature is the webfmwk.HandlerFunc one (func(c Context)).
// To register a custom context, simply do it in the toppest handler.
//
//	package main
//
//	import (
//		"github.com/burgesQ/webfmwk/v5"
//		"github.com/burgesQ/webfmwk/v5/handler/security"
//	)
//
//	type CustomContext struct {
//		webfmwk.Context
//		val String
//	}
//
//	func main() {
//		var s = webfmwk.InitServer(webfmwk.WithHandlers(security.Handler,
//			func(next Habdler) Handler {
//				return func(c webfmwk.Context) error {
//					cc := Context{c, "val"}
//					return next(cc)
//		}}))
func WithHandlers(h ...Handler) Option {
	return func(s *Server) {
		s.addHandlers(h...)
		s.log.Debugf("\t-- handlers loaded")
	}
}

// SetReadTimeout is a timing constraint on the client http request imposed by
// the server from the moment
// of initial connection up to the time the entire request body has been read.
//
// [Accept] --> [TLS Handshake] --> [Request Headers] --> [Request Body] --> [Response]
func SetReadTimeout(val time.Duration) Option {
	return func(s *Server) {
		s.meta.baseServer.ReadTimeout = val
		s.log.Debugf("\t-- read timeout loaded")
	}
}

// SetWriteTimeout is a time limit imposed on client connecting to the server
// via http from the
// time the server has completed reading the request header up to the time
// it has finished writing the response.
//
// [Accept] --> [TLS Handshake] --> [Request Headers] --> [Request Body] --> [Response]
func SetWriteTimeout(val time.Duration) Option {
	return func(s *Server) {
		s.meta.baseServer.WriteTimeout = val
		s.log.Debugf("\t-- write timeout loaded")
	}
}

// SetIDLETimeout the server IDLE timeout AKA keepalive timeout.
func SetIDLETimeout(val time.Duration) Option {
	return func(s *Server) {
		s.meta.baseServer.IdleTimeout = val
		s.log.Debugf("\t-- idle aka keepalive timeout loaded")
	}
}

// EnableKeepAlive disable the server keep alive functions.
func EnableKeepAlive() Option {
	return func(s *Server) {
		s.meta.enableKeepAlive = true
		s.log.Debugf("\t-- keepalive disabled")
	}
}

// EnablePprof enable the pprof endpoints.
func EnablePprof(path ...string) Option {
	return func(s *Server) {
		s.meta.pprof = true

		if len(path) > 0 {
			s.meta.pprofPath = path[0]
		}

		s.log.Debugf("\t-- pprof endpoint enabled")
	}
}

func MaxRequestBodySize(size int) Option {
	return func(s *Server) {
		s.meta.baseServer.MaxRequestBodySize = size
		s.log.Debugf("\t-- request max body size set to %d", size)
	}
}

const (
	ReadTimeout  = 20
	WriteTimeout = 20
	IdleTimeout  = 1
)

func WithSocketHandlerFunc(path string, hf http.HandlerFunc) Option {
	return func(s *Server) {
		s.meta.socketIOHandlerFunc, s.meta.socketIOPath, s.meta.socketIOHF = hf, path, true
		s.log.Debugf("\t-- socket io handler func loaded")
	}
}

func WithSocketHandler(path string, h http.Handler) Option {
	return func(s *Server) {
		s.meta.socketIOHandler, s.meta.socketIOPath, s.meta.socketIOH = h, path, true
		s.log.Debugf("\t-- socket io handlers loaded")
	}
}

// return default cfg
func getDefaultMeta() serverMeta {
	return serverMeta{
		baseServer: &fasthttp.Server{
			ReadTimeout:        ReadTimeout * time.Second,
			WriteTimeout:       WriteTimeout * time.Second,
			IdleTimeout:        IdleTimeout * time.Minute,
			MaxRequestBodySize: fasthttp.DefaultMaxRequestBodySize,
		},
		routes:    make(RoutesPerPrefix),
		pprofPath: "/debug/pprof/{profile:*}",
	}
}

func (m *serverMeta) toServer(addr string) *fasthttp.Server {
	return &fasthttp.Server{
		ReadTimeout:                   m.baseServer.ReadTimeout,
		WriteTimeout:                  m.baseServer.WriteTimeout,
		IdleTimeout:                   m.baseServer.IdleTimeout,
		MaxRequestBodySize:            m.baseServer.MaxRequestBodySize,
		Name:                          "webfmwk " + addr,
		DisableKeepalive:              !m.enableKeepAlive,
		DisableHeaderNamesNormalizing: true,
		ReduceMemoryUsage:             false,
		LogAllErrors:                  true,
		CloseOnShutdown:               true,
	}
}
