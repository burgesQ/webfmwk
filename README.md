# README.md

# What

`webfmwk` is a minimalist go web framework design for JSON API. His purpose is to use as few as possible external library for a lighter build package.

The server handle ctrl+c on it's own.

## use of

| what                | for                                       |
| :-:                 | :-:                                       |
| [gorilla-mux][1]    | for a easy & robust routing logic         |
| [gorilla-hanler][2] | for some useful already coded middlewares |

# How to use it

Sample `main.go`

## psjon

Set the `pjson` query param to anything to fetch a pretty json payload.

# Test

Simply run `go test .`.

Code coverage (`go test . -cover`) : 73%

# Exemple

## Pre-requisite : init frafos logging

```golib
// init frafos logging
log.SetLogLevel(log.LOG_DEBUG)
log.Init(log.LOGGER_STDOUT | log.LOGFORMAT_LONG)
```

## Hugh server instance

```golib
import (
    "frafos.com/golib/webfmwk"
    "frafos.com/golib/webfmwk/middleware"
)

func main() {
    // create server
    s := w.InitServer()

    // add middlewares
    s.AddMiddlware(m.WithLogging)
    s.AddMiddlware(m.Security)

    // declare routes
	routes := w.Routes{
		w.Route{
			Pattern: "/hello",
			Method:  "GET",
			Name:    "hello world",
			Handler: func(c w.CContext) error {
				// a basic string response
				return c.JSONBlob(http.StatusOK, []byte(`{ "message": "hello world" }`))
			},
		},
		w.Route{
			Pattern: "/routes",
			Method:  "GET",
			Name:    "list routes",
			Handler: func(c w.CContext) error {
				// you can serialize go struct
				return c.JSON(http.StatusOK, *c.Routes)
			},
		},
		w.Route{
			Pattern: "/hello/{who}",
			Method:  "GET",
			Name:    "hello who",
			Handler: func(c w.CContext) error {
				// url param are stored in Context.Vars
				var content = `{ "message": "hello ` + c.Vars["who"] + `" }`
				return c.JSONBlob(http.StatusOK, []byte(content))
			},
		},
		w.Route{
			Pattern: "/testQuery",
			Method:  "GET",
			Name:    "what's under",
			Handler: func(c w.CContext) error {
				// query param are stored in Context.Query
				return c.JSON(http.StatusOK, c.Query)
			},
		},
	}

    // set routes
	s.AddRoutes(routes)

    // you can also add routes
	s.AddRoute(w.Route{
		Pattern: "/world",
		Method:  "POST",
		Name:    "post content",
		Handler: func(c w.CContext) error {

			anonymous := struct {
				FirstName string `json:"first_name,omitempty"`
				LastName  string `json:"last_name,omitempty"`
			}{}

			// check body handle the error management, so no return needed
			if !c.CheckFetchContent(&anonymous) {
				return nil
			}

			return c.JSON(http.StatusCreated, anonymous)
		},
	})

    // start asynchronously on :4242
    go func() {
        s.Start(":4242")
    }()

	// ctrl+c is handled internaly
	defer s.WaitAndStop()


}
```

## Run tls

Simply use the method `webfmwk.Server.StartTLS(addr, certPath, keyPath string)`.

```golib
// start tls asynchronously on :4242
go func() {
    s.StartTLS(":4242", "server.crt", "server.key")
}()
```

## Listen on 2 port

Use the magic of go routine

```golib
// start on 2 different address
go func() {
	go func() {
		s.Start(":4444")
	}()
	s.Start(":4242")
}()
```

## Register a custom context

Create a struct that implement `webfmwk.CContext`

```golib
import (
    "frafos.com/golib/log"
    w "frafos.com/golib/webfmwk"
)

type CustomContext struct {
	w.CContext
	Value string
}

// implem ctx interface
func (c *CustomContext) GetName() string {
	return "custom context"
}

func main() {
	// init frafos logging
	log.SetLogLevel(log.LOG_DEBUG)
	log.Init(log.LOGGER_STDOUT | log.LOGFORMAT_LONG)

	// create server
	s := w.InitServer()

    s.Get("/testContext", func(c w.CContext) error {
	    // you can fetch your pre-setted custom context this way
		// log.Debugf("%#v", c)
		cc := c.CustomContext.(*CustomContext)
		var content = `{ "message": "hello ` + cc.Value + `" }`
		return c.JSONBlob(http.StatusOK, []byte(content))
	})


	s.SetCustomContext(func(c w.CContext) interface{} {
		cctx := &CustomContext{c, "turlu"}
		return cctx
	})

    // listen on :4242
	s.Start(":4242")

	// ctrl+c is handled internaly
	defer s.WaitAndStop()
}

```

## Register middlewares

Import `fraofs.com/golib/webfmwk/middleware`

```golib
import (
    "frafos.com/golib/log"
    w "frafos.com/golib/webfmwk"
    m "frafos.com/golib/webfmwk/middleware"
)

func main() {
	// init frafos logging
	log.SetLogLevel(log.LOG_DEBUG)
	log.Init(log.LOGGER_STDOUT | log.LOGFORMAT_LONG)

	// create server
	s := w.InitServer()

    s.AddMiddleware(m.WithLogging)

    s.GET("/hello", func(c w.CContext) error {
		return c.JSONBlob(http.StatusOK, []byte(`{ "message": "hello world" }`))
	})

    // listen on :4242
	s.Start(":4242")

	// ctrl+c is handled internaly
	defer s.WaitAndStop()
}

```
## Minimalistic hello world

```golib
import (
    "frafos.com/golib/log"
    w "frafos.com/golib/webfmwk"
)

func main() {
	// init frafos logging
	log.SetLogLevel(log.LOG_DEBUG)
	log.Init(log.LOGGER_STDOUT | log.LOGFORMAT_LONG)

	// create server
	s := w.InitServer()

    s.GET("/hello", func(c w.CContext) error {
		return c.JSONBlob(http.StatusOK, []byte(`{ "message": "hello world" }`))
	})

    // listen on :4242
	s.Start(":4242")

	// ctrl+c is handled internaly
	defer s.WaitAndStop()
}
```

# TODO

- [x] server
  - [x] Headers
  - [x] Middelware
    - [x] logging
    - [x] secu
    - [x] CORS
- [x] route
  - [x] GET/DELETE
  - [x] POST/PUT
  - [x] url params
  - [x] quert param
  - [x] routes prefix
- [x] test multiple listning address
- [x] pjson
- [ ] context
  - [x] register custom context
  - [ ] up cast
  - [x] down cast
- [ ] RFC's
- [ ] todo's
- [ ] json validation
- [ ] template
- [x] godoc compat
- [x] unit testing

[1]: https://github.com/gorilla/gorilla-mux
[2]: https://github.com/gorilla/gorilla-handler
