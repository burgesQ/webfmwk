# webfmwk

[![Build Status](https://github.com/burgesQ/webfmwk/workflows/GoBuild/badge.svg)](https://github.com/burgesQ/webfmwk/actions?query=workflow%3AGoBuild)
[![codecov](https://codecov.io/gh/burgesQ/webfmwk/branch/master/graph/badge.svg)](https://codecov.io/gh/burgesQ/webfmwk)
[![Go Report Card](https://goreportcard.com/badge/github.com/burgesQ/webfmwk?style=flat-square)](https://goreportcard.com/report/github.com/burgesQ/webfmwk)
[![CodeFactor](https://www.codefactor.io/repository/github/burgesq/webfmwk/badge)](https://www.codefactor.io/repository/github/burgesq/webfmwk)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/burgesQ/webfmwk/master/LICENSE)
[![Version compatibility with Go 1.13 onward using modules](https://img.shields.io/badge/compatible%20with-go1.13+-5272b4.svg)](https://github.com/burgesQ/webfmwk#run)
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/burgesQ/webfmwk)

## What

`webfmwk` is a go web API framework. Credits to Frafos GmbH.

It's a go web framework supporting JSON API.

## Use it

Import `github.com/burgesQ/webfmwk/v4`.

### Important 

Go 1.13 is required. Make sure the environment variable `GO111MODULE` is set to on when running the install command.

### example

#### Hello world

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
	s.GET("/hello", func(c webfmwk.Context) error {
		c.JSONBlob(http.StatusOK, []byte(`{ "message": "hello world" }`))
	})

	// start asynchronously on :4242
	s.Start(":4242")

	// ctrl+c is handled internaly
	defer s.WaitAndStop()
}
```

Reach the endpoint: 

```bash
$ curl -i 'http://localhost:4242/hello'
HTTP/1.1 200 OK
Accept: application/json; charset=UTF-8
Content-Type: application/json; charset=UTF-8
Produce: application/json; charset=UTF-8
Date: Mon, 18 May 2020 07:45:31 GMT
Content-Length: 25

{"message":"hello world"}%
```

</p>
</details>

#### code samples

Some samples are in the `./doc` directory. A main hande the sample running, use `go run . [sample file name]` to test it 

```bash
$ cd doc
$ go run . panic_to_error
. panic_to_error
running panic_to_error (use panic to handle some error case)
- DBG  : 	-- crtl-c support enabled
- DBG  : 	-- handlers loaded
- DBG  : exit handler: starting
- DBG  : http server :4242: starting
- DBG  : [+] server 1 (:4242) 
- DBG  : [+] new connection
+ INFO : [+] (f2124b89-414b-4361-96ec-5f227c0e3369) : [GET]/panic
+ INFO : [-] (f2124b89-414b-4361-96ec-5f227c0e3369) : [422](27)
- DBG  : [-] (f2124b89-414b-4361-96ec-5f227c0e3369) : >{"error":"user not logged"}<
```

## Test

Simply run `make`


[1]: https://github.com/gorilla/mux
[2]: https://github.com/gorilla/handlers
[3]: gopkg.in/go-playground/validator.v9
[4]: https://github.com/gorilla/schema
[5]: https://github.com/json-iterator/go
