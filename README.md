[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/burgesQ/webfmwk)
[![Go Report Card](https://goreportcard.com/badge/github.com/burgesQ/webfmwk?style=flat-square)](https://goreportcard.com/report/github.com/burgesQ/webfmwk)
[![Build Status](http://img.shields.io/travis/burgesQ/webfmwk.svg?style=flat-square)](https://travis-ci.org/burgesQ/webfmwk)
[![Codecov](https://img.shields.io/codecov/c/github/burgesQ/webfmwk.svg?style=flat-square)](https://codecov.io/gh/burgesQ/webfmwk)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/burgesQ/webfmwk/master/LICENSE)

<!-- markdown-toc start - Don't edit this section. Run M-x markdown-toc-refresh-toc -->
**Table of Contents**

- [README.md](#readmemd)
- [What](#what)
    - [dep](#dep)
- [Test](#test)
- [How to use it](#how-to-use-it)
    - [Example](#example)
        - [Basic server](#basic-server)
        - [fetch query param](#fetch-query-param)
        - [fetch url params](#fetch-url-params)
        - [fetch body / validate](#fetch-body--validate)
        - [Run tls](#run-tls)
        - [Register a custom context](#register-a-custom-context)
        - [Register middlewares](#register-middlewares)
        - [Minimalistic hello world](#minimalistic-hello-world)

<!-- markdown-toc end -->

# What

`webfmwk` is a minimalist go web framework design for JSON API. 
His purpose is to use as few as possible external library for a lighter build package.

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

***Pre-requisite : init logging***

For the moment webfmwk use a static logger 

```go
import "github.com/burgesQ/webfmwk/log"

// init logging
log.SetLogLevel(log.LOG_DEBUG)
log.Init(log.LOGGER_STDOUT | log.LOGFORMAT_LONG)
```

### Basic server 

```go
import (
	w "github.com/burgesQ/webfmwk"
)

func main() {
	// init server w/ ctrl+c support
	s := w.InitServer(true)

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

### fetch query param

### fetch url params

### fetch body / validate

### Run tls

Simply use the method `Server.StartTLS(addr, certPath, keyPath string)`.

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

### Register a custom context

Create a struct that extend `webfmwk.Context`

```go
import (
    w "github.com/burgesQ/webfmwk"
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

### Register middlewares

Import `github.com/burgesQ/webfmwk/middleware`

```go
import (
    w "github.com/burgesQ/webfmwk"
    m "github.com/burgesQ/webfmwk/middleware"
)

func main() {
	// create server
	s := w.InitServer()

    s.AddMiddleware(m.WithLogging)
```

### swagger doc compat

### Minimalistic hello world

```golib
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

[1]: https://github.com/gorilla/gorilla-mux
[2]: https://github.com/gorilla/gorilla-handler
