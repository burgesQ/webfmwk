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
	HandlerSign func(c IContext)

	// Route hold the data for one route
	Route struct {
		Pattern string      `json:"pattern"`
		Method  string      `json:"method"`
		Name    string      `json:"name"`
		Handler HandlerSign `json:"-"`
	}

	// Routes hold an array of route
	Routes []Route
)

// check if a routes is compilent
// TODO: all
// func (r *Route) check() bool { return true }

// SetRouter create a mux.Handler router and then :
// register the middlewares,
// register the user defined routes,
// and return the routes handler
func (s *Server) SetRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, mw := range s.middlewares {
		router.Use(mw)
	}

	// for path prefix - for route

	subRouter := router.PathPrefix(s.prefix).Subrouter()
	// register routes
	for _, route := range s.routes {
		subRouter.
			HandleFunc(route.Pattern, s.customHandler(route.Handler)).
			Methods(route.Method).Name(route.Name)
	}

	// register doc handler
	if s.docHandler != nil {
		s.log.Infof("load swagger doc")
		subRouter.PathPrefix("/doc/").Handler(s.docHandler)
	}

	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		var (
			pathTemplate, _ = route.GetPathTemplate()
			methods, _      = route.GetMethods()
		)
		s.log.Debugf("Methods: [%s] Path: (%s)", strings.Join(methods, ","), pathTemplate)
		return nil
	})

	return subRouter
}
