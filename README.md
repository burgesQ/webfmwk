# webfmwk

[![Build Status](https://github.com/burgesQ/webfmwk/workflows/GoBuild/badge.svg)](https://github.com/burgesQ/webfmwk/actions?query=workflow%3AGoBuild)
[![codecov](https://codecov.io/gh/burgesQ/webfmwk/branch/master/graph/badge.svg)](https://codecov.io/gh/burgesQ/webfmwk)
[![Go Report Card](https://goreportcard.com/badge/github.com/burgesQ/webfmwk?style=flat-square)](https://goreportcard.com/report/github.com/burgesQ/webfmwk)
[![CodeFactor](https://www.codefactor.io/repository/github/burgesq/webfmwk/badge)](https://www.codefactor.io/repository/github/burgesq/webfmwk)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/burgesQ/webfmwk/master/LICENSE)
[![Version compatibility with Go 1.13 onward using modules](https://img.shields.io/badge/compatible%20with-go1.13+-5272b4.svg)](https://github.com/burgesQ/webfmwk#run)
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/burgesQ/webfmwk)

## What

`webfmwk` is a go web API framework

## Use it

Import `github.com/burgesQ/webfmwk/v4`.

### Important 

Go `1.13` is required. 

### Example

#### Hello world

<details><summary>go sample</summary>
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
	var s = webfmwk.InitServer()

	s.GET("/hello", func(c webfmwk.Context) error {
		c.JSONBlob(http.StatusOK, []byte(`{ "message": "hello world" }`))
	})

	s.Start(":4242")

	// ctrl+c is handled internaly
	defer s.WaitAndStop()
}
```
</p>
</details>

Reach the endpoint: 

<details><summary>curl sample</summary>
<p>

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

Some samples are in the `./doc` directory. The main (`doc.go`) hande the samples orchestration. Use `go run . [sample file name]` to run the example file.

<details><summary>see more sample</summary>
<p>

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


| what                                         | **filename**        |
| :-                                           | :-                  |
| return hello world                           | `hello_world.go`    |
| fetch value from url                         | `url_param.go`      |
| fetch query param value                      | `query_param.go`    |
| post content handling                        | `post_content.go`   |
| register mutliple endpoints                  | `routes.go`         |
| overload the framework context               | `custom_context.go` |
| register extra hanlders / middleware         | `handlers.go`       |
| generate and expose a swagger doc            | `swagger.go`        |
| start the server in https                    | `tls.go`            |
| attach worker to the server pool             | `custom_worker.go`  |
| add an ID per requrest (ease logging for ex) | `request_id.go`     |
| panic to return an http error                | `panic_to_error.go` |

</p>
</details>

## Test

Simply run `make`

## Contributing

First of all, **thank you** for contributing hearts

If you find any typo/misconfiguration/... please send me a PR or open an issue. 

Also, while creating your Pull Request on GitHub, please write a description which gives the context and/or explains why you are creating it.

## Credit

Frafos GmbH :tada: where I've writted most of that code


[1]: https://github.com/gorilla/mux
[2]: https://github.com/gorilla/handlers
[3]: gopkg.in/go-playground/validator.v9
[4]: https://github.com/gorilla/schema
[5]: https://github.com/json-iterator/go
