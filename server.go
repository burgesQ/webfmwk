package webfmwk

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/burgesQ/webfmwk/v5/log"
	"github.com/burgesQ/webfmwk/v5/tls"
	"github.com/lab259/cors"
	"github.com/valyala/fasthttp"
)

type (
	// Server is a struct holding all the necessary data / struct
	Server struct {
		ctx      context.Context
		cancel   context.CancelFunc
		wg       *sync.WaitGroup
		launcher WorkerLauncher
		log      log.Log
		isReady  chan bool
		meta     serverMeta
	}
)

var (
	// poolOfServers hold all the http(s) server to properly shut them down
	poolOfServers    []*fasthttp.Server
	logger           log.Log
	loggerMu, poolMu sync.Mutex
)

// Run allow to launch multiple server from a single call.
// It take an va arg list of Address as argument.
// The method wait for the server to end via a call to WaitAndStop.
func (s *Server) Run(addrs ...Address) {
	defer s.WaitAndStop()

	for i := range addrs {
		addr := addrs[i]
		if !addr.IsOk() {
			s.GetLogger().Errorf("invalid address format : %s", addr)

			continue
		}

		if cfg := addr.GetTLS(); cfg != nil && !cfg.Empty() {
			s.GetLogger().Infof("starting %s on https://%s", addr.GetName(), addr.GetAddr())
			s.StartTLS(addr.GetAddr(), cfg)

			continue
		}

		s.GetLogger().Infof("starting %s on http://%s", addr.GetName(), addr.GetAddr())
		s.Start(addr.GetAddr())
	}
}

//
// Global methods
//

func fetchLogger() {
	logger = log.GetLogger()
}

// GetLogger return an instance of the Log interface used.
func GetLogger() log.Log {
	// from init server - if the logger is fetched before
	// the server init (which happened pretty often)
	once.Do(initOnce)

	return logger
}

// Shutdown terminate all running servers.
func Shutdown() {
	poolMu.Lock()
	defer poolMu.Unlock()

	for _, server := range poolOfServers {
		logger.Infof("shutdowning server %s...", server.Name)

		if e := server.Shutdown(); e != nil {
			logger.Errorf("shutdowning server : %v", e)
		}

		logger.Infof("server %s down", server.Name)
	}

	poolOfServers = nil
}

//
// Server implemtation
//

//
// Process methods
//

// Start expose an server to an HTTP endpoint.
func (s *Server) Start(addr string) {
	s.internalHandler()
	s.launcher.Start("http server "+addr, func() error {
		go s.pollPingEndpoint(addr)

		return s.internalInit(addr).ListenAndServe(addr)
	})
}

// StartTLS expose an server to an HTTPS address..
func (s *Server) StartTLS(addr string, cfg tls.IConfig) {
	s.internalHandler()

	listener, err := tls.LoadListener(addr, cfg)
	if err != nil {
		s.GetLogger().Fatalf("loading tls: %s", err.Error())
	}

	s.launcher.Start("https server "+addr, func() error {
		return s.internalInit(addr).Serve(listener)
	})
}

// Shutdown call the framework shutdown to stop all running server.
func (s *Server) Shutdown() {
	s.cancel()
	Shutdown()
}

// WaitAndStop wait for all servers to terminate.
// Use of a sync.waitGroup to properly wait all running servers.
func (s *Server) WaitAndStop() {
	s.wg.Wait()
}

// DumpRoutes dump the API endpoints using the server logger.
func (s *Server) DumpRoutes() map[string][]string {
	all := s.GetRouter().List()

	for m, p := range all {
		for i := range p {
			s.log.Infof("routes: [%s]%s", m, p[i])
		}
	}

	return all
}

// Initialize a http.Server struct. Save the server in the pool of workers.
func (s *Server) internalInit(addr string) *fasthttp.Server {
	var (
		worker = s.meta.toServer(addr)
		router = s.GetRouter()
	)

	// register CORS handler - note that it should be the first one
	if s.meta.cors {
		worker.Handler = cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedHeaders:   []string{"X-Requested-With", "Content-Type"},
			AllowedMethods:   []string{"POST", "PUT", "PATCH", "OPTIONS"},
			AllowCredentials: true,
			// Debug: true,
		}).Handler(router.Handler)
	} else {
		worker.Handler = router.Handler
	}

	worker.Logger = s.log

	// save the server
	poolMu.Lock()
	defer poolMu.Unlock()

	poolOfServers = append(poolOfServers, worker)

	s.log.Debugf("[+] server %d (%s) ", len(poolOfServers), addr)

	return worker
}

func concatAddr(addr, prefix string) string {
	if len(addr) > 1 && addr[0] == ':' {
		return "http://127.0.0.1" + addr + prefix + _pingEndpoint
	} else if strings.HasPrefix(addr, "127.0.0.1") {
		return "http://" + addr + prefix + _pingEndpoint
	}

	return addr + prefix + _pingEndpoint
}

// launch the ctrl+c job if needed
func (s *Server) internalHandler() {
	if s.meta.ctrlc && !s.meta.ctrlcStarted {
		s.launcher.Start("exit handler", func() error {
			s.exitHandler(os.Interrupt, syscall.SIGHUP)

			return nil
		})

		s.meta.ctrlcStarted = true
	}
}

// handle ctrl+c internaly
func (s *Server) exitHandler(sig ...os.Signal) {
	c := make(chan os.Signal, 1)

	signal.Notify(c, sig...)

	defer s.Shutdown()

	for s.ctx.Err() == nil {
		select {
		case si := <-c:
			s.log.Infof("captured %v, exiting...", si)

			return
		case <-s.ctx.Done():
			return
		}
	}
}

//
// Setter/Getter
//

// GetLogger return the used Log instance.
func (s *Server) GetLogger() log.Log {
	return s.log
}

// GetLauncher return a pointer to the internal workerLauncher.
func (s *Server) GetLauncher() *WorkerLauncher {
	return &s.launcher
}

// GetContext return the context.Context used.
func (s *Server) GetContext() context.Context {
	return s.ctx
}

// GetContext return the server' context cancel func.
func (s *Server) GetCancel() context.CancelFunc {
	return s.cancel
}

// IsReady return the channel on which `true` is send once the server is up.
func (s *Server) IsReady() chan bool {
	return s.isReady
}

// AddHandlers register the Handler handlers. Handler are executed from the top most.
// The followig examle run the RequestID handler BEFORE the Logging one, to produce a
// log which look like :
// + INFO : [+] (bc339ac1-a62a-48df-8e97-adf9dec32c42) : [GET]/test
//
//	s.AddHandlers(handler.Logging, handler.RequestID)
//
//nolint:unparam
func (s *Server) addHandlers(h ...Handler) *Server {
	s.meta.handlers = append(s.meta.handlers, h...)

	return s
}

// RegisterDocHandler is used to register an swagger doc handler
func (s *Server) addDocHandlers(h ...DocHandler) *Server {
	s.meta.docHandlers = append(s.meta.docHandlers, h...)

	return s
}

// SetPrefix save a custom context so it can be fetched in the controllers
func (s *Server) setPrefix(prefix string) *Server {
	s.meta.prefix = prefix

	return s
}

// RegisterLogger register the Log used
func (s *Server) registerLogger(lg log.Log) *Server {
	loggerMu.Lock()
	defer loggerMu.Unlock()

	logger, s.log = lg, lg

	return s
}

// EnableCORS enable CORS verification
func (s *Server) enableCORS() *Server {
	s.meta.cors = true

	return s
}

// enableCheckIsUp add an /ping endpoint. Is used, cnce a server is started,
// the user can check weather the server is up or not by reading the isReady channel
// vie the IsReady() method.
func (s *Server) enableCheckIsUp() *Server {
	s.meta.checkIsUp = true

	return s
}

// EnableCtrlC let the server handle the SIGINT interuption. To add
// worker to the interuption pool, please use the `GetLauncher` method
func (s *Server) enableCtrlC() *Server {
	s.meta.ctrlc = true

	return s
}
