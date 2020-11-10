package webfmwk

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/burgesQ/gommon/log"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type (
	// Server is a struct holding all the necessary data / struct
	Server struct {
		ctx      context.Context
		wg       *sync.WaitGroup
		launcher WorkerLauncher
		log      log.Log
		isReady  chan bool
		meta     serverMeta
	}
)

var (
	// poolOfServers hold all the http(s) server to properly shut them down
	poolOfServers []*http.Server
	logger        log.Log
)

//
// Global methods
//

func fetchLogger() {
	logger = log.GetLogger()
}

// GetLogger return an instance of the Log interface used
func GetLogger() log.Log {
	// from init server - if the logger is fetched before
	// the server init (which happened pretty often)
	once.Do(initOnce)
	return logger
}

// Shutdown terminate all running servers.
// Call shutdown with a context.context on each http(s) server.
func Shutdown(ctx context.Context) {
	for _, server := range poolOfServers {
		if e := server.Shutdown(ctx); e != nil {
			logger.Errorf("shutdowning server : %v", e)
		}
		logger.Infof("server %s down", server.Addr)
	}

	poolOfServers = []*http.Server{}
}

//
// Server implemtation
//

//
// Process methods
//

// Start expose an server to an HTTP endpoint
func (s *Server) Start(addr string) {
	s.internalHandler()
	s.launcher.Start("http server "+addr, func() error {
		go s.pollPingEndpoint(addr)
		return s.internalInit(addr).ListenAndServe()
	})
}

// Shutdown call the framework shutdown to stop all running server
func (s *Server) Shutdown(ctx context.Context) {
	Shutdown(ctx)
}

// WaitAndStop wait for all servers to terminate.
// Use of a sync.waitGroup to properly wait all running servers.
func (s *Server) WaitAndStop() {
	s.wg.Wait()
}

// DumpRoutes dump the API endpoints using the server logger
func (s *Server) DumpRoutes() {
	var router = s.SetRouter()
	if e := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		var (
			pathTemplate, _ = route.GetPathTemplate()
			methods, _      = route.GetMethods()
		)
		s.log.Infof("Methods: [%s] Path: (%s)", strings.Join(methods, ","), pathTemplate)
		return nil
	}); e != nil {
		log.Errorf("can't walk trough routing : %s", e.Error())
	}
}

// Initialize a http.Server struct. Save the server in the pool of workers.
func (s *Server) internalInit(addr string, tlsStuffs ...ITLSConfig) *http.Server {
	var (
		worker = s.meta.toServer(addr)
		h      = http.TimeoutHandler(s.SetRouter(),
			worker.WriteTimeout-(50*time.Millisecond),
			`{"error": "timeout reached"}`)
	)

	// register mox.CORS handler - note that it should be the first one
	if s.meta.cors {
		h = handlers.CORS(
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"}),
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedMethods([]string{"POST", "PUT", "PATCH", "OPTIONS"}))(
			h)
	}

	worker.Handler = h
	worker.ErrorLog = s.log.GetErrorLogger()

	// load tls for https
	if len(tlsStuffs) == 1 {
		s.loadTLS(&worker, tlsStuffs[0])
	}

	// save the server
	poolOfServers = append(poolOfServers, &worker)
	s.log.Debugf("[+] server %d (%s) ", len(poolOfServers), addr)

	return &worker
}

func concatAddr(addr, prefix string) (new string) {
	new = addr
	if len(addr) > 1 && addr[0] == ':' {
		new = "http://127.0.0.1" + addr
	} else if strings.HasPrefix(addr, "127.0.0.1") {
		new = "http://" + addr
	}

	new += prefix + _pingEndpoint

	return
}

// launch the ctrl+c job if needed
func (s *Server) internalHandler() {
	if s.meta.ctrlc && !s.meta.ctrlcStarted {
		s.launcher.Start("exit handler", func() error {
			s.exitHandler(s.ctx, os.Interrupt)
			return nil
		})
		time.Sleep(1 * time.Millisecond)
		s.meta.ctrlcStarted = true
	}
}

// handle ctrl+c internaly
func (s *Server) exitHandler(ctx context.Context, sig ...os.Signal) {
	var c = make(chan os.Signal, 1)
	signal.Notify(c, sig...)

	defer Shutdown(ctx)

	for ctx.Err() == nil {
		select {
		case si := <-c:
			s.log.Infof("captured %v, exiting...", si)
			return
		case <-ctx.Done():
			return
		}
	}
}

//
// Setter/Getter
//

// GetLogger return the used Log instance
func (s *Server) GetLogger() log.Log {
	return s.log
}

// GetLauncher return a pointer to the internal workerLauncher
func (s *Server) GetLauncher() *WorkerLauncher {
	return &s.launcher
}

// GetContext return the context.Context used
func (s *Server) GetContext() context.Context {
	return s.ctx
}

// IsReady return the channel on which `true` is send once the server is up
func (s *Server) IsReady() chan bool {
	return s.isReady
}

// RegisterDocHandler is used to register an swagger doc handler
func (s *Server) registerDocHandler(handler http.Handler) *Server {
	s.meta.docHandler = handler
	return s
}

// AddMiddlewares register the mux.MiddlewaresFunc middlewares
func (s *Server) addMiddlewares(mw ...mux.MiddlewareFunc) *Server {
	s.meta.middlewares = append(s.meta.middlewares, mw...)
	return s
}

// AddHandlers register the Handler handlers. Handler are executed from the top most.
// The followig examle run the RequestID handler BEFORE the Logging one, to produce a
// log which look like :
// + INFO : [+] (bc339ac1-a62a-48df-8e97-adf9dec32c42) : [GET]/test
//
//   s.AddHandlers(handler.Logging, handler.RequestID)
func (s *Server) addHandlers(h ...Handler) *Server {
	s.meta.handlers = append(s.meta.handlers, h...)
	return s
}

// SetPrefix save a custom context so it can be fetched in the controllers
func (s *Server) setPrefix(prefix string) *Server {
	s.meta.prefix = prefix
	return s
}

// RegisterLogger register the Log used
func (s *Server) registerLogger(lg log.Log) *Server {
	logger = lg
	s.log = lg
	return s
}

// EnableCORS enable CORS verification
func (s *Server) enableCORS() *Server {
	s.meta.cors = true
	return s
}

// EnableCheckIsUp add an /ping endpoint. Is used, cnce a server is started,
// the user can check weather the server is up or not by reading the isReady channel
// vie the IsReady() method
func (s *Server) EnableCheckIsUp() *Server {
	s.meta.checkIsUp = true
	return s
}

// EnableCtrlC let the server handle the SIGINT interuption. To add
// worker to the interuption pool, please use the `GetLauncher` method
func (s *Server) enableCtrlC() *Server {
	s.meta.ctrlc = true
	return s
}
