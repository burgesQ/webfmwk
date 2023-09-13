package webfmwk

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

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
		// log      log.Log
		slog    *slog.Logger
		isReady chan bool
		meta    serverMeta
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
			s.GetStructuredLogger().Error("invalid format", "address", addr)

			continue
		}

		if cfg := addr.GetTLS(); cfg != nil && !cfg.Empty() {
			s.GetStructuredLogger().Info("starting https server",
				"name", addr.GetName(), "address", "https://"+addr.GetAddr())
			s.StartTLS(addr.GetAddr(), cfg)

			continue
		}

		s.GetStructuredLogger().Info("starting http server",
			"name", addr.GetName(), "address", "http://"+addr.GetAddr())
		s.Start(addr.GetAddr())
	}
}

//
// Global methods
//

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
		s.slog.Warn("https endpoints required with http2, skipping", "address", addr)

		return
	}

	s.internalHandler()
	s.launcher.Start(func() {
		s.slog.Debug("http server: starting", "address", addr)

		go s.pollPingEndpoint(addr)

		if e := s.internalInit(addr).ListenAndServe(addr); e != nil {
			s.slog.Error("http server", "address", addr, "error", e)
		}

		s.slog.Info("http server: done", "address", addr)
	})
}

// StartTLS expose an https server.
// The server may have mTLS and/or http2 capabilities.
func (s *Server) StartTLS(addr string, cfg tls.IConfig) {
	s.internalHandler()

	tlsCfg, err := tls.GetTLSCfg(cfg, s.meta.http2)
	if err != nil {
		s.slog.Error("loading tls config", "error", err)
		os.Exit(2)
	}

	listner, err := tls.LoadListner(addr, tlsCfg)
	if err != nil {
		s.slog.Error("loading tls listener", "error", err)
		os.Exit(3)
	}

	server := s.internalInit(addr)

	if s.meta.http2 {
		s.slog.Info("loading http2 support")
		fasthttp2.ConfigureServer(server, fasthttp2.ServerConfig{Debug: true})
	}

	so2 := sOr2(s.meta.http2)

	s.launcher.Start(func() {
		s.slog.Debug(fmt.Sprintf("%s server: starting", so2), "address", addr)
		defer s.slog.Info(fmt.Sprintf("%s server: done", so2), "address", addr)

		go s.pollPingEndpoint(addr, cfg)

		if e := server.Serve(listner); e != nil {
			s.slog.Error(fmt.Sprintf("%s server", so2), "address", addr, "error", e)
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
			s.slog.Info("routes", "name", m, "route", p[i])
		}
	}

	return all
}

type FastLogger struct{ *slog.Logger }

func (flg *FastLogger) Printf(msg string, keys ...any) {
	flg.Info(fmt.Sprintf(msg, keys...))
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

	worker.Logger = &FastLogger{s.slog}

	// save the server
	poolMu.Lock()
	defer poolMu.Unlock()

	poolOfServers = append(poolOfServers, worker)

	s.slog.Debug("[+] server ", "address", addr, "total", len(poolOfServers))

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
			s.slog.Debug("exit handler: starting")
			s.exitHandler(os.Interrupt, syscall.SIGHUP)
			s.slog.Info("exit handler: done")
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
			s.slog.Error("cannot stop the server", "error", e)
		}
	}()

	for s.ctx.Err() == nil {
		select {
		case si := <-c:
			s.slog.Info("captured signal, exiting...", "signal", si)

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
func (s *Server) GetStructuredLogger() *slog.Logger {
	return s.slog
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
func (s *Server) registerStructuredLogger(slg *slog.Logger) *Server {
	s.slog = slg

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
