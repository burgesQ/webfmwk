package webfmwk

import (
	"context"

	"github.com/fasthttp/router"
	"github.com/segmentio/encoding/json"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"github.com/valyala/fasthttp/pprofhandler"
)

const (
	// GET http verbe
	GET = "GET"
	// POST http verbe
	POST = "POST"
	// PATCH http verbe
	PATCH = "PATCH"
	// PUT http verbe
	PUT = "PUT"
	// DELETE http verbe
	DELETE = "DELETE"

	ANY = "ANY"

	_pingEndpoint = "/ping"
)

type (
	// HandlerFunc hold the signature of a Handler.
	// You may return an error implementing the ErrorHandled interface
	// to reduce the boilerplate. If the returned error doesn't implement
	// the interface, a error 500 is by default returned.
	//
	//   HandlerError(c webfmwk.Context) error {
	//     return webfmwk.NewUnauthorizedError("get me some credential !")
	//   }
	//
	// Will produce a http 500 json response.
	HandlerFunc func(c Context) error

	// Handler hold the function signature for webfmwk Handler chaning (middlware).
	//
	//  import (
	//    github.com/burgesQ/webfmwk/v5
	//    github.com/burgesQ/webfmwk/handler/logging
	//    github.com/burgesQ/webfmwk/handler/security
	//  )
	//
	//  s := webfmwk.InitServer(
	//    webfmwk.WithHandler(
	//      logging.Handler,
	//      security.Handler,
	//    ))
	Handler func(HandlerFunc) HandlerFunc

	// DocHandler hold the required data to expose a swagger documentation handlers.
	//
	// Example serving a redoc one:
	//  import (
	//    github.com/burgesQ/webfmwk/v5
	//    github.com/burgesQ/webfmwk/handler/redoc
	//  )
	//
	//  s := webfmwk.InitServer(
	//    webfmwk.WithDocHandler(redoc.GetHandler(
	//      redoc.DocURI("/swagger.json")
	//    ))
	//
	//  s.Get("/swagger.json", func(c webfmwk.Context) error{
	//    return c.JSONBlob(200, `{"title": "some swagger"`)
	//  })
	DocHandler struct {
		// H hold the doc Handler to expose.
		H HandlerFunc

		// Name is used in debug message.
		Name string

		// Path hold the URI one which the handler is reachable.
		// If a prefix is setup, the path will prefixed.
		Path string
	}

	// Route hold the data for one route.
	Route struct {
		// Handler hold the exposed Handler method.
		Handler HandlerFunc `json:"-"`

		// Verbe hold the verbe at which the handler is reachable.
		Verbe string `json:"verbe"`

		// Path hold the uri at which the handler is reachable.
		// If a prefix is setup, the path will be prefixed.
		Path string `json:"path"`

		// Name is used in message.
		Name string `json:"name"`
	}

	// Routes hold an array of route.
	Routes []Route

	// RoutesPerPrefix hold the routes and there respectiv prefix.
	RoutesPerPrefix map[string]Routes
)

var _pong = json.RawMessage(`{"ping": "pong"}`)

//
// Routes method
//

// AddRoutes add the endpoint to the server.
func (s *Server) AddRoutes(r ...Route) {
	s.meta.routes[s.meta.prefix] = append(s.meta.routes[s.meta.prefix], r...)
}

// GET expose a handler to the http verb GET.
func (s *Server) GET(path string, handler HandlerFunc) {
	s.AddRoutes(Route{
		Path:    path,
		Verbe:   GET,
		Handler: handler,
	})
}

// DELETE expose a handler to the http verb DELETE.
func (s *Server) DELETE(path string, handler HandlerFunc) {
	s.AddRoutes(Route{
		Path:    path,
		Verbe:   DELETE,
		Handler: handler,
	})
}

// POST expose a handler to the http verb POST.
func (s *Server) POST(path string, handler HandlerFunc) {
	s.AddRoutes(Route{
		Path:    path,
		Verbe:   POST,
		Handler: handler,
	})
}

// PUT expose a handler to the http verb PUT.
func (s *Server) PUT(path string, handler HandlerFunc) {
	s.AddRoutes(Route{
		Path:    path,
		Verbe:   PUT,
		Handler: handler,
	})
}

// PATCH expose a handler to the http verb PATCH.
func (s *Server) PATCH(path string, handler HandlerFunc) {
	s.AddRoutes(Route{
		Path:    path,
		Verbe:   PATCH,
		Handler: handler,
	})
}

// PATCH expose a handler to the http verb PATCH.
func (s *Server) ANY(path string, handler HandlerFunc) {
	s.AddRoutes(Route{
		Path:    path,
		Verbe:   ANY,
		Handler: handler,
	})
}

// RouteApplier apply the array of RoutePerPrefix.
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
				case ANY:
					s.ANY(prefix+route.Path, route.Handler)
				default:
					s.slog.Warn("Cannot load route [%s](%s)", "route", prefix+route.Path, "verbe", route.Verbe)
				}
			}
		}
	}
}

// GetRouter create a fasthttp/router.Router whit:
// - registered handlers (webfmwk/v5/handler)
// - doc handler is registered
// - test handler (/ping) is registered
// - registered fmwk routes
func (s *Server) GetRouter() *router.Router {
	r := router.New()

	r.HandleMethodNotAllowed, r.HandleOPTIONS = true, true
	r.RedirectTrailingSlash, r.RedirectFixedPath = false, false

	// IDEA: router.PanicHandler
	r.NotFound, r.MethodNotAllowed = s.CustomHandler(handleNotFound), s.CustomHandler(handleNotAllowed)

	// register doc handler
	if len(s.meta.docHandlers) > 0 {
		for i := range s.meta.docHandlers {
			h := s.meta.docHandlers[i]
			s.slog.Info("load doc handler", "name", h.Name)
			r.ANY(s.meta.prefix+h.Path, s.CustomHandler(h.H))
		}
	}

	// register test handler
	if s.meta.checkIsUp {
		r.GET(s.meta.prefix+_pingEndpoint, s.CustomHandler(func(c Context) error {
			return c.JSONOk(_pong)
		}))
	}

	// register socket.io (goplog) handlers
	switch {
	case s.meta.socketIOHF:
		s.slog.Info("loading socket io handler func", "path", s.meta.socketIOPath)
		r.ANY(s.meta.socketIOPath,
			fasthttpadaptor.NewFastHTTPHandlerFunc(s.meta.socketIOHandlerFunc))
	case s.meta.socketIOH:
		s.slog.Info("loading socket io handler", "path", s.meta.socketIOPath)
		r.ANY(s.meta.socketIOPath,
			fasthttpadaptor.NewFastHTTPHandler(s.meta.socketIOHandler))
	}

	if s.meta.pprof {
		s.slog.Info("loading pprof handler", "path", "/debug/pprof/{profile:*}'")
		r.GET(s.meta.prefix+s.meta.pprofPath, pprofhandler.PprofHandler)
	}

	// register routes
	for p, rs := range s.meta.routes {
		prefix, routes := p, rs // never sure if I should copy

		var group *router.Group
		if len(prefix) != 0 {
			group = r.Group(prefix)
		}

		for _, r1 := range routes {
			route := r1
			handler := route.Handler

			// register internal Handlers
			handler = contentIsJSON(handleHandlerError(handler))

			// register user server wise custom Handlers
			if s.meta.handlers != nil {
				for _, h := range s.meta.handlers {
					handler = h(handleHandlerError(handler))
				}
			}

			// TODO: register group wise / route wise custom Handlers
			// if route.handlers != nil {
			// 	for _, h := range s.meta.handlers {
			// 		handler = h(handleHandlerError(handler))
			// 	}
			// }
			// if group.handlers != nil {
			// 	for _, h := range s.meta.handlers {
			// 		handler = h(handleHandlerError(handler))
			// 	}
			// }

			if len(prefix) == 0 {
				r.Handle(route.Verbe, route.Path, s.CustomHandler(handler))
			} else {
				group.Handle(route.Verbe, route.Path, s.CustomHandler(handler))
			}
		}
	}

	return r
}

// CustomHandler return the webfmwk Handler main logic,
// which return a HandlerFunc wrapper in an fasthttp.Handler.
func (s *Server) CustomHandler(handler HandlerFunc) fasthttp.RequestHandler {
	return func(c *fasthttp.RequestCtx) {
		ctx, cancel := s.genContext(c)
		defer cancel()

		// we skip verification as it's done in the useHandler
		_ = handler(ctx)
	}
}

func (s *Server) genContext(c *fasthttp.RequestCtx) (Context, context.CancelFunc) {
	ctx, fn := context.WithCancel(s.ctx)

	return &icontext{c, s.slog, ctx}, fn
}
