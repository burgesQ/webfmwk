package webfmwk

import (
	"strings"

	"github.com/burgesQ/webfmwk/log"
	"github.com/gorilla/mux"
)

type HandlerSign func(c IContext) error

type Route struct {
	Pattern string      `json:"pattern"`
	Method  string      `json:"method"`
	Name    string      `json:"name"`
	Handler HandlerSign `json:"-"`
}

// readability
type Routes []Route

// check if a routes is compilent
// TODO: all
func (r *Route) check() bool { return true }

// SetRouter create a mux.Handler router and then :
// register the middlewares,
// register the user defined routes,
// and return the routes handler
func (s *Server) SetRouter() *mux.Router {
	var (
		router    = mux.NewRouter().StrictSlash(true)
		subRouter = router.PathPrefix(s.prefix).Subrouter()
	)

	// regster middlewares
	for _, mw := range s.middlewares {
		subRouter.Use(mw)
	}

	// register routes
	for _, route := range s.routes {
		subRouter.
			HandleFunc(route.Pattern, s.customHandler(route.Handler)).
			Methods(route.Method).
			Name(route.Name)
	}

	// register doc handler
	if s.docHandler != nil {
		log.Infof("load swagger doc")
		subRouter.PathPrefix("/doc/").Handler(s.docHandler)
	}

	subRouter.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		var (
			pathTemplate, _ = route.GetPathTemplate()
			methods, _      = route.GetMethods()
		)
		log.Debugf("Methods: [%s] Path: (%s)", strings.Join(methods, ","), pathTemplate)
		return nil
	})

	return subRouter
}
