package webfmwk

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/burgesQ/log"
	wlog "github.com/burgesQ/webfmwk/v6/log"
	"github.com/burgesQ/webfmwk/v6/tls"
	fasthttp2 "github.com/dgrr/http2"
	"github.com/lab259/cors"
	"github.com/valyala/fasthttp"
)

type (
	// Server is a struct holding all the necessary data / struct
	Server struct {
		ctx      context.Context //nolint:containedctx
		cancel   context.CancelFunc
		wg       *sync.WaitGroup
		launcher WorkerLauncher
		log      log.Log
		isReady  chan bool
		meta     serverMeta
	}
)

var (
	// TODO: use sync.Pool
	// poolOfServers hold all the http(s) server to properly shut them down
	poolOfServers []*fasthttp.Server
	poolMu        sync.Mutex
)

// Run allow to launch multiple server from a single call.
// It take an va arg list of Address as argument.
// The method wait for the server to end via a call to WaitAndStop.
func (s *Server) Run(addrs ...Address) {
	defer s.WaitForStop()

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

func fetchLogger() log.Log { return wlog.GetLogger() }

// Shututdown terminate all running servers.
func Shutdown() error {
	poolMu.Lock()
	defer poolMu.Unlock()

	var senti error

	for i, server := range poolOfServers {
		if e := server.Shutdown(); e != nil {
			senti = fmt.Errorf("shutdowning server %d : %w", i, e)
		}
	}

	poolOfServers = nil

	return senti
}

//
// Server implemtation
//

//
// Process methods
//

// Start expose an server to an HTTP endpoint.
func (s *Server) Start(addr string) {
	if s.meta.http2 {
		s.log.Warnf("https endpoints required with http2, skipping %q", addr)

		return
	}

	s.internalHandler()
	s.launcher.Start(func() {
		s.log.Debugf("http server %s: starting", addr)

		go s.pollPingEndpoint(addr)

		if e := s.internalInit(addr).ListenAndServe(addr); e != nil {
			s.log.Errorf("http server %s (%T): %s", addr, e, e)
		}
		s.log.Infof("http server %s: done", addr)
	})
}

// StartTLS expose an https server.
// The server may have mTLS and/or http2 capabilities.
func (s *Server) StartTLS(addr string, cfg tls.IConfig) {
	s.internalHandler()

	tlsCfg, err := tls.GetTLSCfg(cfg, s.meta.http2)
	if err != nil {
		s.log.Fatalf("loading tls config: %v", err)
	}

	listner, err := tls.LoadListner(addr, tlsCfg)
	if err != nil {
		s.log.Fatalf("loading tls listener: %v", err)
	}

	server := s.internalInit(addr)

	if s.meta.http2 {
		s.log.Infof("loading http2 support")
		fasthttp2.ConfigureServer(server, fasthttp2.ServerConfig{Debug: true})
	}

	so2 := sOr2(s.meta.http2)

	s.launcher.Start(func() {
		s.log.Debugf("%s server %s: starting", so2, addr)
		defer s.log.Infof("%s server %s: done", so2, addr)

		go s.pollPingEndpoint(addr)

		if e := server.Serve(listner); e != nil {
			s.log.Errorf("%s server %s (%T): %s", so2, addr, e, e)
		}
	})
}

func sOr2(http2 bool) string {
	if http2 {
		return "http2"
	}

	return "https"
}

// ShutdownAndWait call for Shutdown and wait for all server to terminate.
func (s *Server) ShutdownAndWait() error {
	defer s.WaitForStop()

	return s.Shutdown()
}

// Shutdown call the framework shutdown to stop all running server.
func (s *Server) Shutdown() error {
	s.cancel()

	return Shutdown()
}

// WaitForStop wait for all servers to terminate.
// Use of a sync.waitGroup to properly wait all running servers.
func (s *Server) WaitForStop() {
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

// launch the ctrl+c job if needed.
func (s *Server) internalHandler() {
	if s.meta.ctrlc && !s.meta.ctrlcStarted {
		s.launcher.Start(func() {
			s.log.Debugf("exit handler: starting")
			s.exitHandler(os.Interrupt, syscall.SIGHUP)
			s.log.Infof("exit handler: done")
		})

		s.meta.ctrlcStarted = true
	}
}

// handle ctrl+c internaly.
func (s *Server) exitHandler(sig ...os.Signal) {
	c := make(chan os.Signal, 1)

	signal.Notify(c, sig...)

	defer func() {
		if e := s.Shutdown(); e != nil {
			s.log.Errorf("cannot stop the server: %v", e)
		}
	}()

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
func (s *Server) GetLauncher() WorkerLauncher {
	return s.launcher
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
func (s *Server) addHandlers(h ...Handler) *Server { //nolint: unparam
	s.meta.handlers = append(s.meta.handlers, h...)

	return s
}

// RegisterDocHandler is used to register an swagger doc handler.
func (s *Server) addDocHandlers(h ...DocHandler) *Server {
	s.meta.docHandlers = append(s.meta.docHandlers, h...)

	return s
}

// SetPrefix save a custom context so it can be fetched in the controllers.
func (s *Server) setPrefix(prefix string) *Server {
	s.meta.prefix = prefix

	return s
}

// RegisterLogger register the Log used.
func (s *Server) registerLogger(lg log.Log) *Server {
	s.log = lg

	return s
}

// EnableCORS enable CORS verification.
func (s *Server) enableCORS() *Server {
	s.meta.cors = true

	return s
}

// enableCheckIsUp add an /ping endpoint.
// If used, once a server is started, the user can check weather the server is
// up or not by reading the isReady channel vie the IsReady() method.
func (s *Server) EnableCheckIsUp() *Server {
	s.meta.checkIsUp = true

	return s
}

// DisableHTTP2 allow to disable HTTP2 on the fly.
// It usage isn't recommanded.
// For testing purpore only.
func (s *Server) DisableHTTP2() *Server {
	s.meta.http2 = false

	return s
}

// EnableCtrlC let the server handle the SIGINT interuption.
// To add worker to the interuption pool, please use the `GetLauncher` method.
func (s *Server) enableCtrlC() *Server {
	s.meta.ctrlc = true

	return s
}
