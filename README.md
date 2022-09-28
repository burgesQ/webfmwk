# webfmwk

[![Build Status][11]][1]
[![GoDoc][12]][2]
[![codecov][13]][3]
[![Go Report Card][14]][4]
[![Version compatibility with Go 1.15 onward using modules][15]][5]
[![License][16]][6]


[![Internal Build Status][17]][7]
[![Internal Coverage Report][18]][8]

[1]:https://github.com/burgesQ/webfmwk/actions?query=workflow%3AGoTest
[2]:http://godoc.org/github.com/burgesQ/webfmwk/
[3]:https://codecov.io/gh/burgesQ/webfmwk
[4]:https://goreportcard.com/report/github.com/burgesQ/webfmwk
[5]:https://github.com/burgesQ/webfmwk#run
[6]:https://raw.githubusercontent.com/burgesQ/webfmwk/master/LICENSE
[7]:https://gitlab.frafos.net/gommon/webfmwk/-/commits/master
[8]:https://gitlab.frafos.net/gommon/webfmwk/-/commits/master

[11]:https://github.com/burgesQ/webfmwk/workflows/GoTest/badge.svg
[12]:https://godoc.org/github.com/burgesQ/webfmwk/v5?status.svg
[13]:https://codecov.io/gh/burgesQ/webfmwk/branch/master/graph/badge.svg
[14]:https://goreportcard.com/badge/github.com/burgesQ/webfmwk
[15]:https://img.shields.io/badge/compatible%20with-go1.15+-5272b4.svg
[16]:http://img.shields.io/badge/license-mit-blue.svg
[17]:https://gitlab.frafos.net/gommon/webfmwk/badges/master/pipeline.svg
[18]:https://gitlab.frafos.net/gommon/webfmwk/badges/master/coverage.svg

## What

`webfmwk` is a go web API framework build on top of several packages :
- [`valya/fasthttp`][21] instead of `net/http`
- [`fasthttp/router`][22] to route API endpoints
-  [`gorilla/schema`][23], [`go-playground/validator/v10`][24],
[`go-playground/local`][25] and [`go-playground/local`][26] for
structure validation and error translation
- [`segmito/encoding`][27] instead of `encoding/json`


[21]:github.com/valyala/fasthttp
[22]:github.com/fasthttp/router
[23]:github.com/gorilla/schema
[24]:gopkg.in/go-playground/validator.v10
[25]:gopkg.in/go-playground/universal-translator
[26]:gopkg.in/go-playground/local
[27]:github.com/segmentio/encoding

## Use it

Import `github.com/burgesQ/webfmwk/v5`.

### Important

Go `1.15` is required.

### Example

#### Hello world

<details><summary>go sample</summary>
<p>

```go
package main

import (
    "net/http"

    "github.com/burgesQ/webfmwk/v5"
)

// curl -X GET 127.0.0.1:4242/hello
// { "message": "hello world" }
func main() {
    var s = webfmwk.InitServer()

    s.GET("/hello", func(c webfmwk.Context) error {
        c.JSONBlob(http.StatusOK, []byte(`{ "message": "hello world" }`))
    })

    // ctrl+c is handled internaly
    defer s.WaitAndStop()

    s.Start(":4242")
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

Some samples are in the `./doc` directory. The main (`doc.go`) hande the samples
orchestration. Use `go run . [sample file name]` to run the example file.

<details><summary>see more sample</summary>
<p>

```bash
$ cd doc
$ go run . panic_to_error
. panic_to_error
running panic_to_error (use panic to handle some error case)
- DBG  :    -- crtl-c support enabled
- DBG  :    -- handlers loaded
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

Simply run `make`, lint is also available via `make lint`

## Contributing

First of all, **thank you** for contributing hearts

If you find any typo/misconfiguration/... please send me a PR or open an issue.

Also, while creating your Pull Request on GitHub, please write a description
which gives the context and/or explains why you are creating it.

## Credit

Frafos GmbH :tada: where I've writted most of that code
