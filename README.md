[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/burgesQ/webfmwk)
[![Build Status](http://img.shields.io/travis/burgesQ/webfmwk.svg?style=flat-square)](https://travis-ci.org/burgesQ/webfmwk)
[![Codecov](https://img.shields.io/codecov/c/github/burgesQ/webfmwk.svg?style=flat-square)](https://codecov.io/gh/burgesQ/webfmwk)
[![Go Report Card](https://goreportcard.com/badge/github.com/burgesQ/webfmwk?style=flat-square)](https://goreportcard.com/report/github.com/burgesQ/webfmwk)
[![CodeFactor](https://www.codefactor.io/repository/github/burgesq/webfmwk/badge)](https://www.codefactor.io/repository/github/burgesq/webfmwk)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/burgesQ/webfmwk/master/LICENSE)

# What

`webfmwk` is an internal framework build and own by Frafos GmbH. 

It was designed to be a minimalist go web framework supporting JSON API. 

The purpose of the framework is to use as few external library than possible.

The server handle ctrl+c on it's own.

## dep

| what                | for                                             |
| :-:                 | :-:                                             |
| [gorilla-mux][1]    | for a easy & robust routing logic               |
| [gorilla-hanler][2] | for some useful already coded middlewares       |
| json-iterator       | use by the custom implementation of the context |
| validator           | use by the custom implementation of the context |

# Test

Simply run `go test .`

# How to use it

Their is a few main under the `./test` directory

## Example

### Hello world !

<details><summary>hello world</summary>
<p>

```go
import (
    w "github.com/burgesQ/webfmwk"
)

func main() {
	// create server
	s := w.InitServer(true)

    s.GET("/hello", func(c w.IContext) error {
		return c.JSONBlob(http.StatusOK, []byte(`{ "message": "hello world" }`))
	})

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

### fetch query param

### fetch url params

### fetch body / validate

### Use tls

Simply use the method `Server.StartTLS(addr, certPath, keyPath string)`.

<details><summary>use tls</summary>
<p>

```go
// start tls asynchronously on :4242
go func() {
  s.StartTLS(":4242", TLSConfig{
    Cert:     "/path/to/cert",
    Key:      "/path/to/key",
    Insecure: true,
  })
}()
```

</p>
</details>

### Set a base url

<details><summary>base url</summary>
<p>

```go
package main

import (
	w "github.com/burgesQ/webfmwk/v2"
)

func main() {
	// init server w/ ctrl+c support
	s := w.InitServer(true)

	s.SetPrefix("/api")

	s.GET("/test", func(c w.IContext) error {
		return c.JSONOk("ok")
	})

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

Then reach `:4242/api/test`

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

	// ctrl+c is handled internaly
	defer s.WaitAndStop()
}
```

</p>
</details>

### Register a extended context

Create a struct that extend `webfmwk.Context`

<details><summary>extend context</summary>
<p>

```go
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

	s.GET("/test", func(c w.IContext) error {
		ctx := c.(*custom Context)
		return c.JSONOk(ctx.customVal)
	})
```

</p>
</details>

### Register middlewares

Import `github.com/burgesQ/webfmwk/v2/middleware`

<details><summary>extend middleware</summary>
<p>

```go
import (
    w "github.com/burgesQ/webfmwk/v2"
    m "github.com/burgesQ/webfmwk/v2/middleware"
)

func main() {
	// create server
	s := w.InitServer()

    s.AddMiddleware(m.WithLogging)
```

</p>
</details>

### Swagger doc compatibility

Import `github.com/swaggo/http-swagger`

<details><summary>swagger doc</summary>
<p>

```go
package main

import (
	w "github.com/burgesQ/webfmwk/v2"
	"github.com/burgesQ/webfmwk/v2/log"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Answer struct {
	message string `json:"message"`
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

	// init logging
	log.SetLogLevel(log.LogDebug)
	log.Init(log.LoggerSTDOUT | log.LogFormatLong)

	// init server w/ ctrl+c support
	s := w.InitServer(true)

	s.RegisterDocHandler(httpSwagger.WrapHandler)

	s.GET("/test", func(c w.IContext) error {
		return c.JSONOk("ok")
	})

	// start asynchronously on :4242
	go func() {
		s.Start(":4242")
	}()

	// ctrl+c is handled internaly
	defer s.WaitAndStop()
}
```

Then reach `:4242/api/doc/index.html`

</p>
</details>

[1]: https://github.com/gorilla/gorilla-mux
[2]: https://github.com/gorilla/gorilla-handler
