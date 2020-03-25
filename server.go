package webfmwk

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/burgesQ/webfmwk/v3/log"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type (
	Setter func(c *Context) IContext

	// Server is a struct holding all the necessary data / struct
	Server struct {
		ctx      context.Context
		wg       *sync.WaitGroup
		launcher WorkerLauncher
		log      log.ILog
		isReady  chan bool
		meta     ServerMeta
	}
)

var (
	// poolOfServers hold all the http(s) server to properly shut them down
	poolOfServers []*http.Server
	logger        = log.GetLogger()
)

// GetLogger return an instance of the ILog interface used
func GetLogger() log.ILog {
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
// Getter - Setter
//

// GetLogger return the used ILog instance
func (s *Server) GetLogger() log.ILog {
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
func (s *Server) RegisterDocHandler(handler http.Handler) *Server {
	s.meta.docHandler = handler
	return s
}

// RegisterMiddlewares register the middlewares arguments
func (s *Server) RegisterMiddlewares(mw ...mux.MiddlewareFunc) *Server {
	s.meta.middlewares = append(s.meta.middlewares, mw...)
	return s
}

// SetCustomContext save a custom context so it can be fetched in the controllers
func (s *Server) SetCustomContext(setter Setter) *Server {
	s.meta.setter = setter
	return s
}

// RegisterLogger register the ILog used
func (s *Server) RegisterLogger(lg log.ILog) *Server {
	logger = lg
	s.log = lg
	return s
}

// EnableCORS enable CORS verification
func (s *Server) EnableCORS() *Server {
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
func (s *Server) EnableCtrlC() *Server {
	s.meta.ctrlc = true
	return s
}

//
// Process method
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
		s.log.Debugf("Methods: [%s] Path: (%s)", strings.Join(methods, ","), pathTemplate)
		return nil
	}); e != nil {
		log.Errorf("can't walk trough routing : %v", e)
	}
}

// private method

func hasBody(r *http.Request) bool {
	return r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH"
}

// webfmwk main logic, return a http handler wrapped by webfmwk
func (s *Server) customHandler(handler HandlerSign) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var ctx = s.meta.setter(&Context{})

		ctx.SetRequest(r).SetWriter(w).
			SetVars(mux.Vars(r)).SetQuery(r.URL.Query()).
			SetLogger(s.log).SetContext(s.ctx).SetRequestID(GetRequestID(r.Context()))

		defer ctx.OwnRecover()

		if hasBody(r) {
			ctx.CheckHeader()
		}

		handler(ctx)
	}
}

// Initialize a http.Server struct. Save the server in the pool of workers.
func (s *Server) internalInit(addr string, tlsStuffs ...TLSConfig) *http.Server {
	var (
		worker = s.meta.toServer(addr)
		h      = http.TimeoutHandler(s.SetRouter(),
			worker.WriteTimeout-(50*time.Millisecond),
			`{"error": "timeout reached"}`)
	)

	// CORS should be one of the first handler
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

// pollPingEndpoint try to reach the /ping endpoint of the server
// to then infrome that the server is up via the isReady channel
func (s *Server) pollPingEndpoint(addr string) {
	if !s.meta.checkIsUp {
		return
	}

	if len(addr) > 1 && addr[0] == ':' {
		addr = "http://127.0.0.1" + addr
	}

	addr += _pingEndpoint

	for {
		time.Sleep(time.Millisecond * 10)

		/* #nosec  */
		if resp, e := http.Get(addr); e != nil {
			s.log.Infof("server not up ... %s", e.Error())
			continue
		} else if e = resp.Body.Close(); e != nil || resp.StatusCode != http.StatusOK {
			s.log.Infof("unexpected status code, %s : %v", resp.StatusCode, e)
			continue
		}

		s.log.Infof("server is up")
		s.isReady <- true
		break
	}
}

// launch the ctrl+c job if needed
func (s *Server) internalHandler() {
	if s.meta.ctrlc && !s.meta.ctrlcStarted {
		s.launcher.Start("exit handler", func() error {
			s.exitHandler(s.ctx, os.Interrupt)
			return nil
		})
		time.Sleep(5 * time.Millisecond)
		s.meta.ctrlcStarted = true
	}
}

// handle ctrl+c in intern
func (s *Server) exitHandler(ctx context.Context, sig ...os.Signal) {
	var c = make(chan os.Signal, 1)
	signal.Notify(c, sig...)

	defer Shutdown(ctx)

	for ctx.Err() == nil {
		select {
		case si := <-c:
			logger.Infof("captured %v, exiting...", si)
			return
		case <-ctx.Done():
			return
		}
	}
}
