# webfmwk

[![Build Status](https://github.com/burgesQ/webfmwk/workflows/GoBuild/badge.svg)](https://github.com/burgesQ/webfmwk/actions?query=workflow%3AGoBuild)
[![codecov](https://codecov.io/gh/burgesQ/webfmwk/branch/master/graph/badge.svg)](https://codecov.io/gh/burgesQ/webfmwk)
[![Go Report Card](https://goreportcard.com/badge/github.com/burgesQ/webfmwk?style=flat-square)](https://goreportcard.com/report/github.com/burgesQ/webfmwk)
[![CodeFactor](https://www.codefactor.io/repository/github/burgesq/webfmwk/badge)](https://www.codefactor.io/repository/github/burgesq/webfmwk)
[![DeepSource](https://static.deepsource.io/deepsource-badge-light.svg)](https://deepsource.io/gh/burgesQ/webfmwk/?ref=repository-badge)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/burgesQ/webfmwk/master/LICENSE)
[![Version compatibility with Go 1.13 onward using modules](https://img.shields.io/badge/compatible%20with-go1.13+-5272b4.svg)](https://github.com/burgesQ/webfmwk#run)
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/burgesQ/webfmwk)

## What

`webfmwk` is a go web API framework. Proprietary of Frafos GmbH.

It was designed to be a minimalist go web framework supporting JSON API.

The purpose of the framework is to use as few external library than possible.

**TODO: explain that perf is not the purpose. that's why panic are used - more user friendly.**

## Test

Simply run `make`

## Contribute

## Run

### Important 

For Go 1.13, make sure the environment variable GO111MODULE is set as on when running the install command.

### Example

go >= 1.11. Simply import `github.com/burgesQ/webfmwk`. 

`Run go get `

Their are a few mains in the `./exmaple` directory. The content of the mains or used later in the `README.md`.

#### Hello world !

Reach the endpoint with `curl -X GET 'http://localhost:4242/hello'`.

<details><summary>hello world</summary>
<p>

```go
package main

import (
	"net/http"

	"github.com/burgesQ/webfmwk/v4"
)

// curl -X GET 127.0.0.1:4242/hello
// { "message": "hello world" }
func main() {
	// create server
	s := webfmwk.InitServer()

	// expose /hello
	s.GET("/hello", func(c webfmwk.IContext) {
		c.JSONBlob(http.StatusOK, []byte(`{ "message": "hello world" }`))
	})

	// start asynchronously on :4242
	s.Start(":4242")

	// ctrl+c is handled internaly
	defer s.WaitAndStop()
}
```

</p>
</details>

#### fetch query param

Reach the endpoint with `curl -X GET 'http://localhost:4242/hello?&pjson&turlu=tutu'`.

<details><summary>query param</summary>
<p>

```go
package main

import (
	"net/http"

	"github.com/burgesQ/webfmwk/v4"
)

// curl -i -X GET "127.0.0.1:4242/hello?pretty"
// {
//   "pretty": [
//     ""
// 		]
// }
// curl -i -X GET "127.0.0.1:4242/hello?prete"
// {"prete":[""]}%
func main() {
	var s = webfmwk.InitServer()

	// expose /hello
	s.GET("/hello", func(c webfmwk.IContext) {
		c.JSON(http.StatusOK, c.GetQueries())
	})

	// start asynchronously on :4242
	s.Start(":4242")

	// ctrl+c is handled internaly
	defer s.WaitAndStop()
}
```

</p>
</details>

#### fetch url params

Reach the endpoint with `curl -X GET 'http://localhost:4242/hello/you'`.

<details><summary>url param</summary>
<p>

```go
package main

import (
	"net/http"

	"github.com/burgesQ/webfmwk/v4"
)

// curl -X GET 127.0.0.1:4242/hello/world
// {"content":"hello world"}
func main() {
	// init server
	var s = webfmwk.InitServer()

	// expose /hello/name
	s.GET("/hello/{name}", func(c webfmwk.IContext) {
		c.JSONBlob(http.StatusOK, []byte(`{ "content": "hello `+c.GetVar("name")+`" }`))
	})

	// start asynchronously on :4242
	s.Start(":4242")

	// ctrl+c is handled internaly
	defer s.WaitAndStop()
}
```

</p>
</details>

#### deserialize body / query param / validate

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

	"github.com/burgesQ/webfmwk/v4"
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
	var s = webfmwk.InitServer()

	s.POST("/hello", func(c webfmwk.IContext) {
		var out = Payload{}

		// process query params
		c.DecodeQP(&out.qp)
		c.Validate(out.qp)

		// process payload
		c.FetchContent(&out.content)
		c.Validate(out.content)

		c.JSON(http.StatusOK, out)
	})

	// start asynchronously on :4242
	s.Start(":4244")

	// ctrl+c is handled internaly
	defer s.WaitAndStop()
}
```

</p>
</details>

#### Set a base url

Reach the endpoint with `curl -X GET 'http://localhost:4242/api/v1/test` and `curl -X GET 'http://localhost:4242/api/v2/test`.

<details><summary>base url</summary>
<p>

```go
package main

import (
    "github.com/burgesQ/webfmwk/v4"
)

var (
    routes = webfmwk.RoutesPerPrefix{
        "/v1": {
            {
                Verbe: "GET",
                Path:  "/test",
                Name:  "test v1",
                Handler: func(c webfmwk.IContext) {
                    c.JSONOk("v1 ok")
                },
            },
        },
        "/v2": {
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

    s := webfmwk.InitServer(webfmwk.SetPrefix("/api"))

    s.RouteApplier(routes)

    // start asynchronously on :4242
    s.Start(":4242")

    // ctrl+c is handled internaly
    defer s.WaitAndStop()
}
```

</p>
</details>


#### Use tls

Use the method `Server.StartTLS(addr, certPath, keyPath string)`.

<details><summary>use tls</summary>
<p>

```go
package main

import (
    w "github.com/burgesQ/webfmwk/v4"
)

func main() {
    // init server w/ ctrl+c support
    s := w.InitServer(WithCtrlC())

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

#### Register a custom logger

The logger must implement the `webfmwk/log.ILog` interface.

<details><summary>custom logger</summary>
<p>

```go
package main

import (
    w "github.com/burgesQ/webfmwk/v4"
    "github.com/burgesQ/webfmwk/v4/log"
)

// GetLogger return a log.ILog interface
var logger = log.GetLogger()

func main() {
    s := w.InitServer(WithLogger(logger))

    s.GET("/test", func(c w.IContext) error {
        return c.JSONOk("ok")
    })

    // start asynchronously on :4242
    s.StartTLS(":4242", TLSConfig{
    Cert:     "/path/to/cert",
    Key:      "/path/to/key",
    Insecure: true,
    })

    // ctrl+c is handled internally
    defer s.WaitAndStop()
}
```

</p>
</details>

#### Register a extended context

Create a struct that extend `webfmwk.Context`.

Then, add a middleware to extend the context using the `Server.SetCustomContext(func(c *Context) IContext)`

<details><summary>extend context</summary>
<p>

```go
package main

import "github.com/burgesQ/webfmwk/v4"

// customContext extend the webfmwk.Context
type customContext struct {
	webfmwk.Context
	val string
}

// curl -X GET 127.0.0.1:4242/test
// {"content":"42"}
func main() {
	// init server w/ ctrl+c support and custom context options
	var s = webfmwk.InitServer(
		webfmwk.WithCustomContext(func(c *webfmwk.Context) webfmwk.IContext {
			return &customContext{*c, "42"}
		}))

	// expose /test
	s.GET("/test", func(c webfmwk.IContext) {
		c.JSONOk(webfmwk.NewResponse(c.(*customContext).val))
	})

	// start asynchronously on :4242
	s.Start(":4242")

	// ctrl+c is handled internaly
	defer s.WaitAndStop()
}
```

</p>
</details>

#### Register middlewares

Import `github.com/burgesQ/webfmwk/v4/middleware`

<details><summary>extend middleware</summary>
<p>

```go
package main

import (
	"github.com/burgesQ/webfmwk/v4"
	"github.com/burgesQ/webfmwk/v4/middleware"
)

// Middleware implement http.Handler methods
// Check the server logs
//
// curl -i -X GET 127.0.0.1:4242/test
// Accept: application/json; charset=UTF-8
// Content-Type: application/json; charset=UTF-8
// Produce: application/json; charset=UTF-8
// Strict-Transport-Security: max-age=3600; includeSubDomains
// X-Content-Type-Options: nosniff
// X-Xss-Protection: 1; mode=block
// Date: Mon, 06 Apr 2020 14:58:44 GMT
// Content-Length: 4
func main() {
	// init server w/ ctrl+c support and middlewares
	s := webfmwk.InitServer(
		webfmwk.WithCtrlC(),
        webfmwk.WithMiddlewares(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				log.Infof("[%s] %s", r.Method, r.RequestURI)
				next.ServeHTTP(w, r)
			})
        }))
        
	// expose /test
	s.GET("/test", func(c webfmwk.IContext) {
		c.JSONOk("ok")
	})

	// start asynchronously on :4242
	s.Start(":4242")

	// ctrl+c is handled internaly
	defer s.WaitAndStop()
}
```

</p>
</details>

#### Register handlers

Import `github.com/burgesQ/webfmwk/v4/handler`

<details><summary>extend middleware</summary>
<p>

```go
package main

import (
	"github.com/burgesQ/webfmwk/v4"
	"github.com/burgesQ/webfmwk/v4/handler"
)

// Handlers implement webfmwk.Handler methods
// Check the server logs
//
// curl -i -X GET 127.0.0.1:4242/test
// Accept: application/json; charset=UTF-8
// Content-Type: application/json; charset=UTF-8
// Produce: application/json; charset=UTF-8
// Strict-Transport-Security: max-age=3600; includeSubDomains
// X-Content-Type-Options: nosniff
// X-Xss-Protection: 1; mode=block
// Date: Mon, 06 Apr 2020 14:58:44 GMT
// Content-Length: 4
func main() {
	// init server w/ ctrl+c support and middlewares
	s := webfmwk.InitServer(
		webfmwk.WithCtrlC(),
		webfmwk.WithHandlers(handler.Logging))

	// expose /test
	s.GET("/test", handler.Security(func(c webfmwk.IContext) {
		c.JSONOk("ok")
	}))

	// start asynchronously on :4242
	s.Start(":4242")

	// ctrl+c is handled internaly
	s.WaitAndStop()
}
```

</p>
</details>


#### Swagger doc compatibility

Import `github.com/swaggo/http-swagger`.

Then, from a browser reach `:4242/api/doc/index.html`.

Or, run `curl -X GET 'http://localhost:4242/api/doc/swagger.json'`.

<details><summary>swagger doc</summary>
<p>

```go
package main

import (
    w "github.com/burgesQ/webfmwk/v4"
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
    s := w.InitServer(WithDocHandler(httpSwagger.WrapHandler))

    s.SetPrefix("/api")

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

#### Add worker

Use the `Server.GetWorkerLauncher()` method.

Run the test main and wait 10 sec.

<details><summary>extra worker</summary>
<p>

```go
package main

import (
    "time"

    w "github.com/burgesQ/webfmwk/v4"
    "github.com/burgesQ/webfmwk/v4/log"
)

func main() {
    log.SetLogLevel(log.LogDEBUG)
    var (
      s  = w.InitServer()
      wl = s.GetLauncher()
   )


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
