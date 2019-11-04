[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/burgesQ/webfmwk)
[![Build Status](http://img.shields.io/travis/burgesQ/webfmwk.svg?style=flat-square)](https://travis-ci.org/burgesQ/webfmwk)
[![Codecov](https://img.shields.io/codecov/c/github/burgesQ/webfmwk.svg?style=flat-square)](https://codecov.io/gh/burgesQ/webfmwk) 
[![GolangCI](https://golangci.com/badges/github.com/burgesQ/webfmwk.svg)](https://golangci.com/r/github.com/burgesQ/webfmwk)
[![Go Report Card](https://goreportcard.com/badge/github.com/burgesQ/webfmwk?style=flat-square)](https://goreportcard.com/report/github.com/burgesQ/webfmwk)
[![CodeFactor](https://www.codefactor.io/repository/github/burgesq/webfmwk/badge)](https://www.codefactor.io/repository/github/burgesq/webfmwk)
[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/3372/badge)](https://bestpractices.coreinfrastructure.org/projects/3372)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/burgesQ/webfmwk/master/LICENSE)


# What

`webfmwk` is an internal framework build and own by Frafos GmbH. 

It was designed to be a minimalist go web framework supporting JSON API. 

The purpose of the framework is to use as few external library than possible.

The server handle ctrl+c on it's own.

## dep

| what                 | for                                             |
| :-:                  | :-:                                             |
| [gorilla/mux][1]     | for a easy & robust routing logic               |
| [gorilla/hanlers][2] | for some useful already coded middlewares       |
| [gorilla/schema][4]  | for some useful already coded middlewares       |
| [validator][3]       | use by the custom implementation of the context |
| [json-iterator][5]   | use by the custom implementation of the context |

# Test

Simply run `go test .`

# How to use it

Their are a few mains in the `./exmaple` directory. The content of the mains or used later in the `README.md`.

## Example

### Hello world !

Reach the endpoint with `curl -X GET 'http://localhost:4242/hello'`.

<details><summary>hello world</summary>
<p>

```go
package main

import (
	"net/http"

	w "github.com/burgesQ/webfmwk/v2"
)

func main() {
	// create server
	s := w.InitServer(true)

	s.GET("/hello", func(c w.IContext) {
		c.JSONBlob(http.StatusOK, []byte(`{ "message": "hello world" }`))
	})

	// start asynchronously on :4242
	go func() {
		s.Start(":4242")
	}()

	// ctrl+c is handled internally
	defer s.WaitAndStop()
}
```

</p>
</details>

### fetch query param

Reach the endpoint with `curl -X GET 'http://localhost:4242/hello?&pjson&turlu=tutu'`.

<details><summary>query param</summary>
<p>


```go
package main

import (
	"net/http"

	w "github.com/burgesQ/webfmwk/v2"
	"github.com/burgesQ/webfmwk/v2/log"
)

func main() {
	// create server
	s := w.InitServer(true)

	s.GET("/hello", func(c w.IContext) {
		var (
			queries   = c.GetQueries()
			pjson, ok = c.GetQuery("pjson")
		)
		if ok {
			log.Errorf("%#v", pjson)
		}
		c.JSON(http.StatusOK, queries)
	})

	// start asynchronously on :4242
	go func() {
		s.Start(":4242")
	}()

	// ctrl+c is handled internally
	defer s.WaitAndStop()
}
```

</p>
</details>

### fetch url params

Reach the endpoint with `curl -X GET 'http://localhost:4242/hello/you'`.

<details><summary>url param</summary>
<p>

```go
package main

import (
	"net/http"

	w "github.com/burgesQ/webfmwk/v2"
)

func main() {
	// create server
	s := w.InitServer(true)

	s.GET("/hello/{id}", func(c w.IContext) {
		c.JSONBlob(http.StatusOK, []byte(`{ "id": "`+c.GetVar("id")+`" }`))
	})

	// start asynchronously on :4242
	go func() {
		s.Start(":4242")
	}()

	// ctrl+c is handled internally
	defer s.WaitAndStop()
}
```

</p>
</details>

### deserialize body / query param / validate

Reach the endpoint with `curl -X POST -d '{"name": "test", "age": 12}' -H "Content-Type: application/json" "http://localhost:4242/hello"`.

Note that the `webfmwk` only accept `application/json` content (for the moment ?).

Don't hesitate to play with the payload to inspect the behavior of the Validate method.

The struct annotation are done via the `validator`  and `schema` keywords. Please refer to the [`validator` documentation][3] and the [`gorilla/schema`][4] one.

<details><summary>POST content</summary>
<p>

```go
package main

import (
	"net/http"

	w "github.com/burgesQ/webfmwk/v2"
)

type (
	// Content hold the body of the request
	Content struct {
		Name string `schema:"name" json:"name" validate:"omitempty"`
		Age  int    `schema:"age" json:"age" validate:"gte=1"`
	}

	// QueryParam hold the query params
	QueryParam struct {
		PJSON bool `schema:"pjson" json:"pjson"`
		Val   int  `schema:"val" json:"val" validate:"gte=1"`
	}

	// Payload hold the output of the endpoint
	Payload struct {
		Content Content    `json:"content"`
		QP      QueryParam `json:"query_param"`
	}
)

func main() {
	// create server
	s := w.InitServer(true)

	s.POST("/hello", func(c w.IContext) {

		out := Payload{}

		c.FetchContent(&out.content)
		c.Validate(out.content)

		c.DecodeQP(&out.qp)
		c.Validate(out.qp)

		c.JSON(http.StatusOK, out)
	})

	// start asynchronously on :4242
	go func() {
		s.Start(":4244")
	}()

	// ctrl+c is handled internally
	defer s.WaitAndStop()
}
```

</p>
</details>

### Set a base url

Reach the endpoint with `curl -X GET 'http://localhost:4242/api/v1/test` and `curl -X GET 'http://localhost:4242/api/v2/test`.

<details><summary>base url</summary>
<p>

```go
package main

import (
	"github.com/burgesQ/webfmwk/v2"
)

var (
	routes = webfmwk.RoutesPerPrefix{
		"/api/v1": {
			{
				Verbe: "GET",
				Path:  "/test",
				Name:  "test v1",
				Handler: func(c webfmwk.IContext) {
					c.JSONOk("v1 ok")
				},
			},
		},
		"/api/v2": {
			{
				Verbe: "GET",
				Path:  "/test",
				Name:  "test v2",
				Handler: func(c webfmwk.IContext) {
					c.JSONOk("v2 ok")
				},
			},
		},
	}
)

func main() {

	s := webfmwk.InitServer(true)

	s.RouteApplier(routes)

	// start asynchronously on :4242
	go func() {
		s.Start(":4242")
	}()

	// ctrl+c is handled internaly
	defer s.WaitAndStop()

}
```

</p>
</details>


### Use tls

Use the method `Server.StartTLS(addr, certPath, keyPath string)`.

<details><summary>use tls</summary>
<p>

```go
package main

import (
	w "github.com/burgesQ/webfmwk/v2"
)

func main() {
	// init server w/ ctrl+c support
	s := w.InitServer(true)

	s.GET("/test", func(c w.IContext) error {
		return c.JSONOk("ok")
	})

	// start asynchronously on :4242
	go func() {
		s.StartTLS(":4242", TLSConfig{
			Cert:     "/path/to/cert",
			Key:      "/path/to/key",
			Insecure: true,
		})
	}()

	// ctrl+c is handled internally
	defer s.WaitAndStop()
}
```

</p>
</details>

### Register a custom logger

The logger must implement the `webfmwk/log.ILog` interface.

<details><summary>custom logger</summary>
<p>

```go
package main

import (
	w "github.com/burgesQ/webfmwk/v2"
	"github.com/burgesQ/webfmwk/v2/log"
)

// GetLogger return a log.ILog interface
var logger = log.GetLogger()

func main() {
	// init server w/ ctrl+c support
	s := w.InitServer(true)

	s.SetLogger(logger)

	s.GET("/test", func(c w.IContext) error {
		return c.JSONOk("ok")
	})

	// start asynchronously on :4242
	go func() {
		s.StartTLS(":4242", TLSConfig{
			Cert:     "/path/to/cert",
			Key:      "/path/to/key",
			Insecure: true,
		})
	}()

	// ctrl+c is handled internally
	defer s.WaitAndStop()
}
```

</p>
</details>

### Register a extended context

Create a struct that extend `webfmwk.Context`. 

Then, add a middleware to extend the context using the `Server.SetCustomContext(func(c *Context) IContext)`

<details><summary>extend context</summary>
<p>

```go
package main

import (
	w "github.com/burgesQ/webfmwk/v2"
)

type customContext struct {
	w.Context
	customVal string
}

func main() {
	// init server w/ ctrl+c support
	s := w.InitServer(true)

	s.SetCustomContext(func(c *w.Context) w.IContext {
		ctx := &customContext{*c, "42"}
		return ctx
	})

	s.GET("/test", func(c w.IContext) {
		ctx := c.(*customContext)
		c.JSONOk(ctx.customVal)
	})

	// start asynchronously on :4242
	go func() {
		s.Start(":4244")
	}()

	// ctrl+c is handled internally
	defer s.WaitAndStop()
}
```

</p>
</details>

### Register middlewares

Import `github.com/burgesQ/webfmwk/v2/middleware`

<details><summary>extend middleware</summary>
<p>

```go
package main

import (
	w "github.com/burgesQ/webfmwk/v2"
	m "github.com/burgesQ/webfmwk/v2/middleware"
)

func main() {

	// init server w/ ctrl+c support
	s := w.InitServer(true)

	s.AddMiddleware(m.Logging)
	s.AddMiddleware(m.Security)

	s.GET("/test", func(c w.IContext) error {
		return c.JSONOk("ok")
	})

	// start asynchronously on :4242
	go func() {
		s.Start(":4242")
	}()

	// ctrl+c is handled internally
	defer s.WaitAndStop()
}
```

</p>
</details>

### Swagger doc compatibility

Import `github.com/swaggo/http-swagger`.

Then, from a browser reach `:4242/api/doc/index.html`. 

Or, run `curl -X GET 'http://localhost:4242/api/doc/swagger.json'`.

<details><summary>swagger doc</summary>
<p>

```go
package main

import (
	w "github.com/burgesQ/webfmwk/v2"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Answer struct {
	Message string `json:"message"`
}

// @Summary hello world
// @Description Return a simple greeting
// @Param pjson query bool false "return a pretty JSON"
// @Success 200 {object} db.Reply
// @Produce application/json
// @Router /hello [get]
func hello(c w.IContext) error {
	return c.JSONOk(Answer{"ok"})
}

// @title hello world API
// @version 1.0
// @description This is an simple API
// @termsOfService https://www.youtube.com/watch?v=DLzxrzFCyOs
// @contact.name Quentin Burgess
// @contact.url github.com/burgesQ
// @contact.email quentin@frafos.com
// @license.name GFO
// @host localhost:4242
func main() {
	// init server w/ ctrl+c support
	s := w.InitServer(true)

    s.SetPrefix("/api")

    s.RegisterDocHandler(httpSwagger.WrapHandler)

	s.GET("/test", func(c w.IContext) error {
		return c.JSONOk("ok")
	})

	// start asynchronously on :4242
	go func() {
		s.Start(":4242")
	}()

	// ctrl+c is handled internally
	defer s.WaitAndStop()
}
```
</p>
</details>

### Add worker 

Use the `Server.GetWorkerLauncher()` method.

Run the test main and wait 10 sec.

<details><summary>extra worker</summary>
<p>

```go
package main

import (
	"time"

	w "github.com/burgesQ/webfmwk/v2"
	"github.com/burgesQ/webfmwk/v2/log"
)

func main() {

	log.SetLogLevel(log.LogDEBUG)

	// init server w/ ctrl+c support
	s := w.InitServer(true)
	wl := s.GetLauncher()

	s.GET("/test", func(c w.IContext) {
		c.JSONOk("ok")
	})

	wl.Start("custom worker", func() error {
		time.Sleep(10 * time.Second)
		log.Debugf("done")
		return nil
	})

	// start asynchronously on :4242
	go func() {
		s.Start(":4242")
	}()

	// ctrl+c is handled internally
	defer s.WaitAndStop()
}
```

</p>
</details>

[1]: https://github.com/gorilla/mux
[2]: https://github.com/gorilla/handlers
[3]: gopkg.in/go-playground/validator.v9 
[4]: https://github.com/gorilla/schema
[5]: https://github.com/json-iterator/go
