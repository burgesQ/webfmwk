package webfmwk

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	GET    = "GET"
	POST   = "POST"
	PATCH  = "PATCH"
	PUT    = "PUT"
	DELETE = "DELETE"

	_pingEndpoint = "/ping"
)

type (
	// HandlerSign hold the signature of a webfmwk handler (chain of middlware)
	HandlerFunc func(c Context) error

	// Handler hold the function signature of a webfmwk handler chaning (middlware)
	Handler func(HandlerFunc) HandlerFunc

	// DocHandler hold the required data to expose a documentation handlers
	DocHandler struct {
		// H represent the doc handler to expose
		H http.HandlerFunc
		// Name is used in debug message
		Name string
		// Path hold the URI at  which the handler is exposed.
		// If one used, will be prefixed by the api base.
		Path string
	}

	// Route hold the data for one route
	Route struct {
		Verbe   string      `json:"verbe"`
		Path    string      `json:"path"`
		Name    string      `json:"name"`
		Handler HandlerFunc `json:"-"`
	}

	// Routes hold an array of route
	Routes []Route

	// RoutesPerPrefix hold the routes and there respectiv prefix
	RoutesPerPrefix map[string]Routes
)

//
// Routes method
//

// AddRoute add the endpoint to the server
func (s *Server) AddRoutes(r ...Route) {
	s.meta.routes[s.meta.prefix] = append(s.meta.routes[s.meta.prefix], r...)
}

// GET expose a route to the http verb GET
func (s *Server) GET(path string, handler HandlerFunc) {
	s.AddRoutes(Route{
		Path:    path,
		Verbe:   GET,
		Handler: handler,
	})
}

// DELETE expose a route to the http verb DELETE
func (s *Server) DELETE(path string, handler HandlerFunc) {
	s.AddRoutes(Route{
		Path:    path,
		Verbe:   DELETE,
		Handler: handler,
	})
}

// POST expose a route to the http verb POST
func (s *Server) POST(path string, handler HandlerFunc) {
	s.AddRoutes(Route{
		Path:    path,
		Verbe:   POST,
		Handler: handler,
	})
}

// PUT expose a route to the http verb PUT
func (s *Server) PUT(path string, handler HandlerFunc) {
	s.AddRoutes(Route{
		Path:    path,
		Verbe:   PUT,
		Handler: handler,
	})
}

// PATCH expose a route to the http verb PATCH
func (s *Server) PATCH(path string, handler HandlerFunc) {
	s.AddRoutes(Route{
		Path:    path,
		Verbe:   PATCH,
		Handler: handler,
	})
}

// RouteApplier apply the array of RoutePerPrefix
func (s *Server) RouteApplier(rpps ...RoutesPerPrefix) {
	for _, rpp := range rpps {
		for prefix, routes := range rpp {
			for _, route := range routes {
				switch route.Verbe {
				case GET:
					s.GET(prefix+route.Path, route.Handler)
				case POST:
					s.POST(prefix+route.Path, route.Handler)
				case PUT:
					s.PUT(prefix+route.Path, route.Handler)
				case PATCH:
					s.PATCH(prefix+route.Path, route.Handler)
				case DELETE:
					s.DELETE(prefix+route.Path, route.Handler)
				default:
					s.log.Warnf("Cannot load route [%s](%s)", prefix+route.Path, route.Verbe)
				}
			}
		}
	}
}

// UseHanlder apply the HandlerFunc method
func UseHanlder(next HandlerFunc) HandlerFunc {
	return HandlerFunc(func(c Context) error {
		return next(c)
	})
}

// SetRouter create a mux.Handler router and then :
// register the middle wares,
// register the user defined routes per prefix,
// and return the routes handler
func (s *Server) SetRouter() *mux.Router {
	var router = mux.NewRouter().StrictSlash(true)

	router.NotFoundHandler, router.MethodNotAllowedHandler =
		http.HandlerFunc(s.notFoundHandler), http.HandlerFunc(s.notAllowedFunc)

	// register http handler / mux.Middleware
	for _, mw := range s.meta.middlewares {
		router.Use(mw)
	}

	// register doc handler
	if len(s.meta.docHandlers) > 0 {
		for i := range s.meta.docHandlers {
			h := s.meta.docHandlers[i]
			s.log.Infof("load %q doc handler", h.Name)
			router.HandleFunc(s.meta.prefix+h.Path, h.H)
		}
	}

	// register test handler
	if s.meta.checkIsUp {
		router.HandleFunc(s.meta.prefix+_pingEndpoint, func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "pong")
		}).Methods("GET").Name("ping endpoint")
	}

	// register routes
	for prefix, routes := range s.meta.routes {
		subRouter := router.PathPrefix(prefix).Subrouter()

		for _, route := range routes {
			var handler = route.Handler

			// register webfmwk.Handlers
			if s.meta.handlers != nil {
				for _, h := range s.meta.handlers {
					handler = h(UseHanlder(handler))
				}
			}

			subRouter.HandleFunc(route.Path, s.CustomHandler(handler)).
				Methods(route.Verbe).Name(route.Name)
		}
	}

	return router
}

// webfmwk main logic, return a http handler wrapped by webfmwk
func (s *Server) CustomHandler(handler HandlerFunc) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var ctx, cancel = s.genContext(w, r)
		defer cancel()

		if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
			if e := ctx.CheckHeader(); e != nil {
				s.handleError(ctx, e)
				return
			}
		}

		if e := handler(ctx); e != nil {
			s.log.Errorf("catch an e: %s", e.Error())
			s.handleError(ctx, e)
		}
	}
}

func getIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-Ip")
	if ip == "" {
		if ip = r.Header.Get("X-Forwarded-For"); ip == "" {
			ip = r.RemoteAddr
		}
	}

	return ip
}

func (s *Server) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	s.log.Infof("[!] 404 reached for [%s] %s %s", getIP(r), r.Method, r.RequestURI)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	if _, e := w.Write([]byte(`{"status":404,"message":"not found"}`)); e != nil {
		s.log.Errorf("[!] cannot write 404 ! %s", e.Error())
	}
}

func (s *Server) notAllowedFunc(w http.ResponseWriter, r *http.Request) {
	s.log.Infof("[!] 405 reached for [%s] %s %s", getIP(r), r.Method, r.RequestURI)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)

	if _, e := w.Write([]byte(`{"status":405,"message":"method not allowed"}`)); e != nil {
		s.log.Errorf("cannot write 405 ! %s", e.Error())
	}
}

func (s *Server) handleError(ctx Context, e error) {
	var eh ErrorHandled
	if errors.As(e, &eh) {
		_ = ctx.JSON(eh.GetOPCode(), eh.GetContent())
		return
	}

	s.log.Errorf("catched from controller (%T) : %s", e, e.Error())
	_ = ctx.JSONInternalError(NewErrorFromError(e))
}

func (s *Server) genContext(w http.ResponseWriter, r *http.Request) (Context, context.CancelFunc) {
	var ctx, fn = context.WithCancel(s.ctx)

	return &icontext{
		r:     r,
		w:     w,
		vars:  mux.Vars(r),
		query: r.URL.Query(),
		log:   s.log,
		ctx:   ctx,
	}, fn
}
