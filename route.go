package webfmwk

import (
	"strings"

	"github.com/gorilla/mux"
)

const (
	GET    = "GET"
	POST   = "POST"
	PATCH  = "PATCH"
	PUT    = "PUT"
	DELETE = "DELETE"
)

type (
	// HandlerSign hold the signature of the controller
	HandlerSign func(c IContext)

	// Handler Sign func(c IContext) error

	// Route hold the data for one route
	Route struct {
		Verbe   string      `json:"verbe"`
		Path    string      `json:"path"`
		Name    string      `json:"name"`
		Handler HandlerSign `json:"-"`
	}

	// Routes hold an array of route
	Routes []Route

	// RoutesPerPrefix hold the routes and there respectiv prefix
	RoutesPerPrefix map[string]Routes
)

func (rpp *RoutesPerPrefix) addRoute(p string, r Route) {
	(*rpp)[p] = append((*rpp)[p], r)
}

func (rpp *RoutesPerPrefix) addRoutes(p string, r Routes) {
	(*rpp)[p] = append((*rpp)[p], r...)
}

//
// Routes method
//

// SetPrefix set the url path to prefix
func (s *Server) SetPrefix(prefix string) {
	s.prefix = prefix
}

// AddRoute add a new route to expose
func (s *Server) AddRoute(r Route) {
	s.routes.addRoute(s.prefix, r)
}

// AddRoutes save all the routes to expose
func (s *Server) AddRoutes(r Routes) {
	s.routes.addRoutes(s.prefix, r)
}

// GET expose a route to the http verb GET
func (s *Server) GET(path string, handler HandlerSign) {
	s.AddRoute(Route{
		Path:    path,
		Verbe:   GET,
		Handler: handler,
	})
}

// DELETE expose a route to the http verb DELETE
func (s *Server) DELETE(path string, handler HandlerSign) {
	s.AddRoute(Route{
		Path:    path,
		Verbe:   DELETE,
		Handler: handler,
	})
}

// POST expose a route to the http verb POST
func (s *Server) POST(path string, handler HandlerSign) {
	s.AddRoute(Route{
		Path:    path,
		Verbe:   POST,
		Handler: handler,
	})
}

// PUT expose a route to the http verb PUT
func (s *Server) PUT(path string, handler HandlerSign) {
	s.AddRoute(Route{
		Path:    path,
		Verbe:   PUT,
		Handler: handler,
	})
}

// PATCH expose a route to the http verb PATCH
func (s *Server) PATCH(path string, handler HandlerSign) {
	s.AddRoute(Route{
		Path:    path,
		Verbe:   PATCH,
		Handler: handler,
	})
}

// RouteApplier apply the array of RoutePerPrefix
func (s *Server) RouteApplier(rpp RoutesPerPrefix) {
	for prefix, routes := range rpp {
		s.SetPrefix(prefix)
		for _, route := range routes {
			switch route.Verbe {
			case GET:
				s.GET(route.Path, route.Handler)
			case POST:
				s.POST(route.Path, route.Handler)
			case PUT:
				s.PUT(route.Path, route.Handler)
			case PATCH:
				s.PATCH(route.Path, route.Handler)
			case DELETE:
				s.DELETE(route.Path, route.Handler)
			default:
				s.log.Warnf("Cannot load route [%s](%s)", route.Path, route.Verbe)
			}
		}
	}
}

// SetRouter create a mux.Handler router and then :
// register the middlewares,
// register the user defined routes per prefix,
// and return the routes handler
func (s *Server) SetRouter() *mux.Router {
	var router = mux.NewRouter().StrictSlash(true)

	for _, mw := range s.middlewares {
		router.Use(mw)
	}

	for prefix, routes := range s.routes {
		subRouter := router.PathPrefix(prefix).Subrouter()
		// register routes
		for _, route := range routes {
			subRouter.
				HandleFunc(route.Path, s.customHandler(route.Handler)).
				Methods(route.Verbe).Name(route.Name)
		}

		// register doc handler
		if s.docHandler != nil {
			s.log.Infof("load swagger doc")
			subRouter.PathPrefix("/doc/").Handler(s.docHandler)
		}
	}

	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		var (
			pathTemplate, _ = route.GetPathTemplate()
			methods, _      = route.GetMethods()
		)
		s.log.Debugf("Methods: [%s] Path: (%s)", strings.Join(methods, ","), pathTemplate)
		return nil
	})

	return router
}
