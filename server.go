package webfmwk

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/burgesQ/webfmwk/v3/log"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const noTime = 0

type (
	// TLSConfig contain the tls config passed by the config file
	TLSConfig struct {
		Cert     string `json:"cert"`
		Key      string `json:"key"`
		Insecure bool   `json:"insecure"`
		// CaCert string `json:"ca-cert"`
	}

	// Server is a struct holding all the necessary data / struct
	Server struct {
		routes      RoutesPerPrefix
		ctx         *context.Context
		wg          *sync.WaitGroup
		launcher    WorkerLauncher
		middlewares []mux.MiddlewareFunc
		prefix      string
		context     IContext
		docHandler  http.Handler
		CORS        bool
		log         log.ILog
	}

	// WorkerConfig hold the worker config per server instance
	WorkerConfig struct {
		// ReadTimeout is a timing constraint on the client http request imposed by the server from the moment
		// of initial connection up to the time the entire request body has been read.
		// [Accept] --> [TLS Handshake] --> [Request Headers] --> [Request Body] --> [Response]
		ReadTimeout       time.Duration
		ReadHeaderTimeout time.Duration
		// WriteTimeout is a time limit imposed on client connecting to the server via http from the
		// time the server has completed reading the request header up to the time it has finished writing the response.
		// [Accept] --> [TLS Handshake] --> [Request Headers] --> [Request Body] --> [Response]
		WriteTimeout   time.Duration
		MaxHeaderBytes int
	}
)

var (
	// poolOfServers hold all the http(s) server to properly shut them down
	poolOfServers []*http.Server
	logger        = log.GetLogger()
	workerConfig  = WorkerConfig{
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
)

//
// Setter - Getter
//

// GetLogger return an instance of the ILog interface used
func GetLogger() log.ILog {
	return logger
}

// InitServer set the server struct & pre-launch the exit handler.
// Init the worker internal launcher.
// If withCtrl is set to true, the server will handle ctrl+C internall.y
// Please add worker to the package's WorkerLauncher to sync them.
func InitServer(withCtrl bool) Server {
	var (
		wg          sync.WaitGroup
		ctx, cancel = context.WithCancel(context.Background())
		s           = Server{
			launcher: CreateWorkerLauncher(&wg, cancel),
			ctx:      &ctx,
			wg:       &wg,
			context:  &Context{},
			log:      logger,
			routes:   make(RoutesPerPrefix),
		}
	)

	// launch the ctrl+c job
	if withCtrl {
		s.launcher.Start("exit handler", func() error {
			s.ExitHandler(ctx, os.Interrupt)
			return nil
		})
	}

	return s
}

// GetLogger return an instance of the ILog interface used
func (s *Server) GetLogger() log.ILog {
	return s.log
}

// RegisterDocHandler is used to save an swagger doc handler
func (s *Server) RegisterDocHandler(handler http.Handler) {
	s.docHandler = handler
}

// SetCustomContext save a custom context so it can be fetched in the controller handler
func (s *Server) SetCustomContext(setter func(c *Context) IContext) bool {
	ctx, ok := s.context.(*Context)
	if ok {
		s.context = setter(ctx)
	}

	return ok
}

// GetLauncher return a pointer on the util.workerLauncher used
func (s *Server) GetLauncher() *WorkerLauncher {
	return &s.launcher
}

// GetContext return a pointer on the context.Context used
func (s *Server) GetContext() *context.Context {
	return s.ctx
}

// Enamelware append a middleware to the list of middleware
func (s *Server) AddMiddleware(mw mux.MiddlewareFunc) {
	s.middlewares = append(s.middlewares, mw)
}

func (s *Server) hasBody(r *http.Request) bool {
	return r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH"
}

// webfmwk main logic, return a http handler wrapped by webfmwk
func (s *Server) customHandler(handler HandlerSign) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var ctx = s.context

		// run handler
		ctx.SetRequest(r).SetWriter(&w).
			SetVars(mux.Vars(r)).SetQuery(r.URL.Query()).
			SetLogger(s.log).SetContext(s.ctx)

		defer ctx.OwnRecover()

		if s.hasBody(r) {
			ctx.CheckHeader()
		}
		handler(ctx)
	}
}

func (s *Server) loadTLS(worker *http.Server, tlsCfg TLSConfig) {
	worker.TLSConfig = &tls.Config{
		InsecureSkipVerify: tlsCfg.Insecure,
		Certificates:       make([]tls.Certificate, 1),
	}

	cert, err := tls.LoadX509KeyPair(tlsCfg.Cert, tlsCfg.Key)
	if err != nil {
		s.log.Fatalf("%s", err.Error())
	}

	worker.TLSConfig.Certificates[0] = cert
}

// SetWorkerParams merge the WorkerConfig param with the
// package variable workerConfig. The workerConfig is then used
// to spawn an http.Server
func (s *Server) SetWorkerParams(w WorkerConfig) {
	if workerConfig.ReadTimeout != w.ReadTimeout && w.ReadTimeout != noTime {
		workerConfig.ReadTimeout = w.ReadTimeout
		s.log.Debugf("read timeout setted to %d", w.ReadTimeout)
	}

	if workerConfig.ReadHeaderTimeout != w.ReadHeaderTimeout && w.ReadHeaderTimeout != noTime {
		workerConfig.ReadHeaderTimeout = w.ReadHeaderTimeout
		s.log.Debugf("read header timeout setted to %d", w.ReadHeaderTimeout)
	}

	if workerConfig.WriteTimeout != w.WriteTimeout && w.WriteTimeout != noTime {
		workerConfig.WriteTimeout = w.WriteTimeout
		s.log.Debugf("write timeout setted to %d", w.WriteTimeout)
	}

	if workerConfig.MaxHeaderBytes != w.MaxHeaderBytes && w.MaxHeaderBytes != 0 {
		workerConfig.MaxHeaderBytes = w.MaxHeaderBytes
		s.log.Debugf("max header bytes setted to %d", w.MaxHeaderBytes)
	}
}

func toWorker(addr string) http.Server {
	return http.Server{
		Addr:           addr,
		ReadTimeout:    workerConfig.ReadTimeout,
		WriteTimeout:   workerConfig.WriteTimeout,
		MaxHeaderBytes: workerConfig.MaxHeaderBytes,
	}
}

// Initialize a http.Server struct. Save the server in the pool of workers.
func (s *Server) setServer(addr string, tlsStuffs ...TLSConfig) *http.Server {
	worker := toWorker(addr)
	// ! handlers.CORS() must be the first handler
	if s.CORS {
		worker.Handler = handlers.CORS(
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"}),
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedMethods([]string{"POST", "PUT", "PATCH", "OPTIONS"}))(
			s.SetRouter())
	} else {
		worker.Handler = http.TimeoutHandler(
			s.SetRouter(),
			worker.WriteTimeout-(50*time.Millisecond),
			`{"error": "timeout reached"}`)
	}

	// load tls for https
	if len(tlsStuffs) == 1 {
		s.loadTLS(&worker, tlsStuffs[0])
	}

	// save the server
	poolOfServers = append(poolOfServers, &worker)
	s.log.Debugf("[+] server %d (%s) ", len(poolOfServers), addr)

	return &worker
}

// StartTLS expose an server to an HTTPS endpoint
func (s *Server) StartTLS(addr string, tlsStuffs TLSConfig) {
	s.launcher.Start("https server "+addr, func() error {
		return s.setServer(addr, tlsStuffs).ListenAndServeTLS(tlsStuffs.Cert, tlsStuffs.Key)
	})
}

// Start expose an server to an HTTP endpoint
func (s *Server) Start(addr string) {
	s.launcher.Start("http server "+addr, func() error {
		worker := s.setServer(addr)
		return worker.ListenAndServe()
	})
}

// Shutdown terminate all running servers.
// Call shutdown with a context.context on each http(s) server.
func Shutdown(ctx context.Context) {
	for _, server := range poolOfServers {
		server.Shutdown(ctx)
	}

	poolOfServers = []*http.Server{}
}

// Shutdown call the framework shutdown to stop all running server
func (s *Server) Shutdown(ctx context.Context) {
	s.log.Debugf("-1")
	Shutdown(ctx)
}

// WaitAndStop wait for all servers to terminate.
// Use of a sync.waitGroup to properly wait all group.
func (s *Server) WaitAndStop() {
	s.wg.Wait()
}

// ExitHandler handle ctrl+c in intern
func (s *Server) ExitHandler(ctx context.Context, sig ...os.Signal) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, sig...)

	defer s.Shutdown(ctx)

	select {
	case si := <-c:
		s.log.Infof("captured %v, exiting...", si)
		return
	case <-ctx.Done():
		return
	}
}

// SetLogger set the logger of the server
func (s *Server) SetLogger(lg log.ILog) {
	logger = lg
	s.log = logger
}
