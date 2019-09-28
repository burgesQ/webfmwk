package webfmwk

import (
	scontext "context"
	"crypto/tls"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/burgesQ/webfmwk/log"
	"github.com/burgesQ/webfmwk/util"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// TODO: route restriction

// TLSConfig contain the tls config passed by the config file
type TLSConfig struct {
	Cert     string `json:"cert"`
	Key      string `json:"key"`
	Insecure bool   `json:"insecure"`
	// CaCert string `json:"ca-cert"`
}

// Server is a struct holding all the necessary data / struct
type Server struct {
	routes        Routes
	ctx           *scontext.Context
	wg            *sync.WaitGroup
	launcher      util.WorkerLauncher
	middlewares   []mux.MiddlewareFunc
	prefix        string
	customContext interface{}
	docHandler    http.Handler
	CORS          bool
	// contextual
}

var (
	// poolOfServers hold all the http(s) server to properly shut them down.
	poolOfServers []*http.Server
)

//
// Setter - Getter
//

func (s *Server) RegisterDocHandler(handler http.Handler) {
	s.docHandler = handler
}

// Save a custom context * so it can be fetched in the controller handler.
func (s *Server) SetCustomContext(setter func(c CContext) interface{}) {
	ctx, _ := s.customContext.(CContext)
	s.customContext = setter(ctx)
}

// SetPrefix set the url path to prefix.
func (s *Server) SetPrefix(prefix string) {
	s.prefix = prefix
}

// FetchLauncher return a pointer on the util.workerLauncher used.
func (s *Server) GetLauncher() *util.WorkerLauncher {
	return &s.launcher
}

// FetchLauncher return a pointer on the context.Context used.
func (s *Server) GetContext() *scontext.Context {
	return s.ctx
}

// AddMiddlware append a middleware to the list of middleware.
func (s *Server) AddMiddleware(mw mux.MiddlewareFunc) {
	s.middlewares = append(s.middlewares, mw)
}

// Add a extra route to expose.
func (s *Server) AddRoute(r Route) {
	s.routes = append(s.routes, r)
}

// Add extra routes to expose.
func (s *Server) AddRoutes(r []Route) {
	s.routes = append(s.routes, r...)
}

//
// Routes method
//

func (s *Server) GET(path string, handler HandlerSign) {
	s.AddRoute(Route{
		Pattern: path,
		Method:  "GET",
		Handler: handler,
	})
}

func (s *Server) DELETE(path string, handler HandlerSign) {
	s.AddRoute(Route{
		Pattern: path,
		Method:  "DELETE",
		Handler: handler,
	})
}

func (s *Server) POST(path string, handler HandlerSign) {
	s.AddRoute(Route{
		Pattern: path,
		Method:  "POST",
		Handler: handler,
	})
}

func (s *Server) PUT(path string, handler HandlerSign) {
	s.AddRoute(Route{
		Pattern: path,
		Method:  "PUT",
		Handler: handler,
	})
}

func (s *Server) PATCH(path string, handler HandlerSign) {
	s.AddRoute(Route{
		Pattern: path,
		Method:  "PATCH",
		Handler: handler,
	})
}

//
// Magic
//

// webfmwk main logic
// Return a http handler wrapped by webfmwk .
func (s *Server) customHandler(handler HandlerSign) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		// copy context & set data
		ctx, _ := (s.customContext).(CContext)

		ctx.R = r
		ctx.W = &w
		ctx.Routes = &s.routes

		// extract params
		s.HandleParam(&ctx, r)

		// check for pjson
		if len(ctx.Query["pjson"]) > 0 {
			ctx.Pretty = true
		} else {
			ctx.Pretty = false
		}

		// check for header if needed
		if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
			if !ctx.CheckHeader() {
				return
			}
		}

		// register the user custom context
		//TODO: so wrong in there ..
		ctx.CustomContext = s.customContext

		// run handler
		defer ctx.OwnRecover()

		if err := handler(ctx); err != nil {
			log.Errorf("%s", err.Error())
		}

	}
}

// Initialize a http.Server struct.
// Save the server in the pool of workers.
func (s *Server) setServer(addr string, tlsStuffs ...TLSConfig) *http.Server {

	// ! handlers.CORS() must be the first handler
	worker := http.Server{
		Addr:           addr,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if s.CORS {
		headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
		originsOk := handlers.AllowedOrigins([]string{"*"})
		methodsOk := handlers.AllowedMethods([]string{"POST", "PUT", "PATCH", "OPTIONS"})
		worker.Handler = handlers.CORS(originsOk, headersOk, methodsOk)(s.SetRouter())
	} else {
		worker.Handler = s.SetRouter()
	}

	// load tls for https
	if len(tlsStuffs) == 1 {
		tlsCfg := tlsStuffs[0]

		var err error

		worker.TLSConfig = &tls.Config{
			InsecureSkipVerify: tlsCfg.Insecure,
			Certificates:       make([]tls.Certificate, 1),
		}

		worker.TLSConfig.Certificates[0], err = tls.LoadX509KeyPair(tlsCfg.Cert, tlsCfg.Key)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}
	}

	// save the server
	poolOfServers = append(poolOfServers, &worker)
	log.Debugf("[+] server %d (%s) ", len(poolOfServers), addr)

	return &worker
}

// Run the web framework server on addr via https.
func (s *Server) StartTLS(addr string, tlsStuffs TLSConfig) error {
	s.launcher.Start("https server "+addr, func() error {
		return s.setServer(addr, tlsStuffs).ListenAndServeTLS(tlsStuffs.Cert, tlsStuffs.Key)
	})
	return nil
}

// Run the web framework server on addr.
func (s *Server) Start(addr string) error {
	s.launcher.Start("http server "+addr, func() error {
		return s.setServer(addr).ListenAndServe()
	})
	return nil
}

// Shutdown terminate all running servers.
// Call shutdown with a context.context on each http(s) server.
func (s *Server) Shutdown(ctx scontext.Context) error {

	for _, server := range poolOfServers {
		server.Shutdown(ctx)
		log.Debugf("-1")
	}

	poolOfServers = []*http.Server{}

	log.Infof("ctx bye")
	return nil
}

// WaitAndStop wait for all servers to terminate.
// Use of a sync.waitGroup to properly wait all group.
func (s *Server) WaitAndStop() {
	s.wg.Wait()
	log.Infof("wg bye")
}

// Handle ctrl+c.
func (s *Server) ExitHandler(ctx scontext.Context, sig ...os.Signal) {
	c := make(chan os.Signal)
	signal.Notify(c, sig...)

	select {
	case <-ctx.Done():
		s.Shutdown(ctx)
		return
	case si := <-c:
		log.Infof("captured %v, exiting...", si)
		s.Shutdown(ctx)
		return
	}
}

// InitServer set the server struct & pre-launch the exit handler.
// Init the worker internal launcher.
func InitServer(withCtrl bool) (s Server) {

	var wg sync.WaitGroup
	ctx, cancel := scontext.WithCancel(scontext.Background())

	s.launcher = util.CreateWorkerLauncher(&wg, cancel)

	// launch the ctrl+c job
	if withCtrl {
		s.launcher.Start("exit handler", func() error {
			s.ExitHandler(ctx, os.Interrupt)
			return nil
		})
	}
	// save the context & wait groupe
	s.ctx = &ctx
	s.wg = &wg

	return
}
