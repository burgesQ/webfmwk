package webfmwk

import (
	"net/http"
	"testing"

	validator "github.com/go-playground/validator/v10"
)

type customContext struct {
	Context
	Value string
}

type testSerial struct {
	A string `json:"test"`
}

type userForm struct {
	Firstname string `json:"first_name" validate:"required,alpha"`
	Lastname  string `json:"last_name" validate:"required,custom"`
}

type queryParam struct {
	Pretty bool `json:"pretty" schema:"pretty"`
	Some   *int `json:"some,omitempty" schema:"some" validate:"omitempty,min=-1"`
}

func TestUseCase(t *testing.T) {
	var (
		s = InitServer(
			CheckIsUp(), SetPrefix("/api"), DisableKeepAlive(),
			WithHandlers(func(next HandlerFunc) HandlerFunc {
				return HandlerFunc(func(c Context) error {
					cc := customContext{c, "turlu"}
					return next(cc)
				})
			}),
		)
		routes = RoutesPerPrefix{
			"/hello": {
				{
					Verbe: "GET", Path: "", Handler: func(c Context) error {
						return c.JSONBlob(http.StatusOK, []byte(`{ "message": "hello world" }`))
					}},
				{
					Verbe: "GET", Path: "/{who}", Handler: func(c Context) error {
						var content = `{ "message": "hello ` + c.GetVar("who") + `" }`
						return c.JSONBlob(http.StatusOK, []byte(content))
					}},
			},
			"/test": {
				{
					Verbe: "GET", Path: "query", Handler: func(c Context) error {
						return c.JSONOk(c.GetQuery())
					}},
				{
					Verbe: "GET", Path: "Context", Handler: func(c Context) error {
						return c.JSONBlob(http.StatusOK, []byte(`{ "message": "hello `+
							c.(customContext).Value+`" }`))
					}},

				{
					Verbe: "GET", Path: "/queryToStruct", Handler: func(c Context) error {
						qp := queryParam{}
						if e := c.DecodeAndValidateQP(&qp); e != nil {
							return e
						}

						return c.JSONOk(qp)
					}},
			},
			"": {
				{
					Verbe: "GET", Path: "/routes", Handler: func(c Context) error {
						return c.JSON(http.StatusOK, &testSerial{"hello"})
					}},
				{
					Verbe: "POST", Path: "/world", Handler: func(c Context) error {
						anonymous := userForm{}

						if e := c.FetchAndValidateContent(&anonymous); e != nil {
							return e
						}

						return c.JSONCreated(anonymous)
					}},
			},
		}
	)

	s.RouteApplier(routes)

	RegisterValidatorRule("custom", func(fi validator.FieldLevel) bool {
		return fi.Field().String() != "fail"
	})
	RegisterValidatorTrans("custom", "'{0} is invalid :)")
	//RegisterValidatorAlias("alpha", "letters")

	defer stopServer(t, s)
	go s.Start(_testPort)
	<-s.isReady

	const (
		_reqNTest = iota
		_pushNTest
		_pushNTestContain
		_deleteNTest
		// _patchNTest
		// _putNTest
	)

	tests := map[string]struct {
		action      int
		header      bool
		bodyDiffer  bool
		url         string
		body        string
		code        int
		headers     []Header
		pushContent []byte
	}{
		"hello world": {
			action: _reqNTest, header: true, url: "/api/hello",
			body: `{"message":"hello world"}`, code: http.StatusOK,
		},

		"not found": {
			action: _reqNTest, header: true, url: "/api/undef",
			body: `{"status":404,"message":"not found"}`, code: http.StatusNotFound,
		},

		"not allowed": {
			action: _deleteNTest, header: true, url: "/api/hello",
			body: `{"status":405,"message":"method not allowed"}`, code: http.StatusMethodNotAllowed,
		},

		"simple fetch": {
			action: _reqNTest, url: "/api/routes",
			body: `{"test":"hello"}`, code: http.StatusOK,
		},
		"url params": {
			action: _reqNTest, url: "/api/hello/you",
			body: `{"message":"hello you"}`, code: http.StatusOK,
		},
		"query params": {
			action: _reqNTest, bodyDiffer: true, url: "/api/testquery?pretty=1",
			body: `{"pretty":["1"]}`, code: http.StatusOK,
		},

		"query to struct": {
			action: _reqNTest, url: "/api/test/queryToStruct",
			body: `{"pretty":false}`, code: http.StatusOK,
		},
		"query to struct invalide value": {
			action: _reqNTest, url: "/api/test/queryToStruct?some=-5",
			body: `{"status":422,"message":{"queryParam.Some":"Some must be -1 or greater"}}`, code: http.StatusUnprocessableEntity,
		},
		"query to struct invalide field": {
			action: _reqNTest, url: "/api/test/queryToStruct?else=true",
			body: `{"status":422,"message":"schema: invalid path \"else\""}`, code: http.StatusUnprocessableEntity,
		},
		"query to struct filled": {
			action: _reqNTest, url: "/api/test/queryToStruct?some=10&pretty",
			body: "{\n  \"pretty\": false,\n  \"some\": 10\n}", code: http.StatusOK,
		},

		"query params pretty": {
			action: _reqNTest, url: "/api/testquery?pretty",
			code: http.StatusOK, body: `{}`,
		},
		"context": {
			action: _reqNTest, bodyDiffer: true, url: "/api/testContext",
			body: `{"message":"hello tutu"}`, code: http.StatusOK,
		},

		"push": {
			action: _pushNTest, url: "/api/world", pushContent: []byte(`{"first_name":"jean", "last_name":"claude"}`),
			body: `{"first_name":"jean","last_name":"claude"}`, code: http.StatusCreated,
		},

		"push_wrong_header": {
			action: _pushNTest, url: "/api/world", pushContent: []byte(`{"first_name":"jean", "last_name":"claude"}`),
			headers: []Header{{"Content-Type", "plain-text"}},
			body:    `{"status":406,"message":"Content-Type is not application/json"}`, code: http.StatusNotAcceptable,
		},

		"push_form_miss_field": {
			action: _pushNTestContain, url: "/api/world", pushContent: []byte(`{"first_name":"jean"}`),
			body: `Lastname is a required field`, code: http.StatusUnprocessableEntity,
		},
		"push_invalid_empty": {
			action: _pushNTest, url: "/api/world", pushContent: []byte(`{}`), bodyDiffer: true,
			body: `{}`, code: http.StatusUnprocessableEntity,
		},

		"push_invalid_wrong": {
			action: _pushNTest, url: "/api/world", pushContent: []byte(`{`),
			body: `{"status":422,"message":"Unprocessable payload"}`, code: http.StatusUnprocessableEntity,
		},

		"push_custom": {
			action: _pushNTest, url: "/api/world", pushContent: []byte(`{"first_name":"uno", "last_name":"fail"}`),
			body: `{"status":422,"message":{"userForm.Lastname":"'Lastname is invalid :)"}}`, code: http.StatusUnprocessableEntity,
		},

		// TODO: test GET/DELETE/PATCH/PUT ?
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			//			t.Helper()
			switch test.action {

			case _reqNTest:
				requestAndTestAPI(t, _testAddr+test.url, func(t *testing.T, resp *http.Response) {

					if test.header {
						for _, testVal := range []string{"Content-Type", "Accept", "Produce"} {
							assertHeader(t, testVal, jsonEncode, resp)
						}
					}

					if test.bodyDiffer {
						assertBodyDiffere(t, test.body, resp)
					} else {
						assertBody(t, test.body, resp)
					}
					assertStatusCode(t, test.code, resp)
				})

			case _deleteNTest:
				deleteAndTestAPI(t, _testAddr+test.url, func(t *testing.T, resp *http.Response) {
					assertBody(t, test.body, resp)
					assertStatusCode(t, test.code, resp)
				})

			case _pushNTest:
				pushAndTestAPI(t, _testAddr+test.url, test.pushContent,
					func(t *testing.T, resp *http.Response) {
						if test.bodyDiffer {
							assertBodyDiffere(t, test.body, resp)
						} else {
							assertBody(t, test.body, resp)
						}

						assertStatusCode(t, test.code, resp)
					}, test.headers...)

			case _pushNTestContain:
				pushAndTestAPI(t, _testAddr+test.url, test.pushContent,
					func(t *testing.T, resp *http.Response) {
						assertBodyContain(t, test.body, resp)
						assertStatusCode(t, test.code, resp)
					}, test.headers...)

			}

		})
	}
}
