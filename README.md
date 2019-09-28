[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/burgesQ/webfmwk)
[![Go Report Card](https://goreportcard.com/badge/github.com/burgesQ/webfmwk?style=flat-square)](https://goreportcard.com/report/github.com/burgesQ/webfmwk)
[![Build Status](http://img.shields.io/travis/burgesQ/webfmwk.svg?style=flat-square)](https://travis-ci.org/burgesQ/webfmwk)
[![Codecov](https://img.shields.io/codecov/c/github/burgesQ/webfmwk.svg?style=flat-square)](https://codecov.io/gh/burgesQ/webfmwk)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/burgesQ/webfmwk/master/LICENSE)

<!-- markdown-toc start - Don't edit this section. Run M-x markdown-toc-refresh-toc -->
**Table of Contents**

- [What](#what)
    - [dep](#dep)
- [Test](#test)
- [How to use it](#how-to-use-it)
    - [Example](#example)
        - [Hello world !](#hello-world-)
        - [fetch query param](#fetch-query-param)
        - [fetch url params](#fetch-url-params)
        - [fetch body / validate](#fetch-body--validate)
        - [Use tls](#use-tls)
        - [Register a extended context](#register-a-extended-context)
        - [Register middlewares](#register-middlewares)
        - [swagger doc compat](#swagger-doc-compat)

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

### Register a extended context

Create a struct that extend `webfmwk.Context`

<details><summary>extend context</summary>
<p>

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

</p>
</details>

### Register middlewares

Import `github.com/burgesQ/webfmwk/middleware`

<details><summary>extend middleware</summary>
<p>

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

</p>
</details>

### swagger doc compat



[1]: https://github.com/gorilla/gorilla-mux
[2]: https://github.com/gorilla/gorilla-handler
