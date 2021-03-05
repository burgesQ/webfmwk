package webfmwk

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/burgesQ/webfmwk/v5/log"
	validator "github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"github.com/valyala/fasthttp"
)

const _prettyTag = "pretty"

type (
	// InputHandling interface introduce I/O actions.
	InputHandling interface {
		// FetchContent extract the json content from the body into the content interface.
		FetchContent(content interface{}) ErrorHandled

		// Validate is used to validate a content of the content params.
		// See https://go-playground/validator.v10 for more.
		Validate(content interface{}) ErrorHandled

		// FetchAndValidateContent fetch the content then validate it.
		FetchAndValidateContent(content interface{}) ErrorHandled

		// DecodeQP load the query param in the content object.
		// Seee https://github.com/gorilla/query for more.
		DecodeQP(content interface{}) ErrorHandled

		// DecodeAndValidateQP load the query param in the content object and then validate it.
		DecodeAndValidateQP(content interface{}) ErrorHandled
	}

	// ContextLogger interface implement the context Logger needs.
	ContextLogger interface {
		// SetLogger set the logger of the ctx.
		SetLogger(logger log.Log) Context

		// GetLogger return the logger of the ctx.
		GetLogger() log.Log
	}

	// Context interface implement the context used in this project.
	Context interface {
		SendResponse
		InputHandling
		ContextLogger

		// GetFastContext return a pointer to the internal fasthttp.RequestCtx.
		GetFastContext() *fasthttp.RequestCtx

		// GetContext return the request context.Context.
		GetContext() context.Context

		// GetVar return the url var parameters. An empty string for missing case.
		GetVar(key string) (val string)

		// GetQueries return the queries into a fasthttp.Args object.
		GetQuery() *fasthttp.Args

		// GetQuery fetch the query object key
		// GetQuery(key string) (val string, ok bool)
	}

	// icontext implement the Context interface
	// It hold the data used by the request
	icontext struct {
		*fasthttp.RequestCtx
		// vars  map[string]string
		// query map[string][]string
		log log.Log
		ctx context.Context
	}
)

var (
	// decoder annotation : `schema` : gorilla
	decoder                 = schema.NewDecoder()
	errUnprocessablePayload = NewUnprocessable(NewError("Unprocessable payload"))
)

// GetVar implement Context
func (c *icontext) GetVar(key string) string {
	v, ok := c.UserValue(key).(string)
	if !ok {
		return ""
	}

	return v
}

// GetQuery implement Context
func (c *icontext) GetQuery() *fasthttp.Args {
	return c.QueryArgs()
}

// GetFastContext implement Context
func (c *icontext) GetFastContext() *fasthttp.RequestCtx {
	return c.RequestCtx
}

// SetLogger implement Context
func (c *icontext) SetLogger(logger log.Log) Context {
	c.log = logger

	return c
}

// GetLogger implement Context
func (c *icontext) GetLogger() log.Log {
	return c.log
}

// GetContext implement Context
func (c *icontext) GetContext() context.Context {
	return c.ctx
}

// FetchContent implement Context.
// It load payload in the dest interface{} using the system json library.
func (c *icontext) FetchContent(dest interface{}) ErrorHandled {
	b := c.PostBody()

	if e := json.Unmarshal(b, &dest); e != nil {
		c.log.Errorf("fetching payload: %s", e.Error())
		return errUnprocessablePayload
	}

	return nil
}

// Validate implement Context
// this implemtation use validator to anotate & check struct
func (c *icontext) Validate(dest interface{}) ErrorHandled {
	if e := validate.Struct(dest); e != nil {
		c.log.Errorf("validating : %s", e.Error())

		return NewUnprocessable(ValidationError{
			Status: http.StatusUnprocessableEntity,
			Error:  e.(validator.ValidationErrors).Translate(trans),
		})
	}

	return nil
}

// FetchAndValidateContent implemt Context.
// It sucesively call FetchContent then Validate on the dest param
func (c *icontext) FetchAndValidateContent(dest interface{}) ErrorHandled {
	if e := c.FetchContent(&dest); e != nil {
		return e
	}

	return c.Validate(dest)
}

// DecodeQP implement Context
func (c *icontext) DecodeQP(dest interface{}) ErrorHandled {
	m := map[string][]string{}

	c.QueryArgs().VisitAll(func(k, v []byte) {
		key := string(k)
		m[key] = make([]string, 1)
		m[key][0] = string(v)
	})

	if e := decoder.Decode(dest, m); e != nil {
		c.log.Errorf("validating qp : %s", e.Error())
		return NewUnprocessable(NewErrorFromError(e))
	}

	return nil
}

// DecodeAndValidateQP implement Context
func (c *icontext) DecodeAndValidateQP(qp interface{}) ErrorHandled {
	e := c.DecodeQP(qp)

	if e != nil {
		return e
	}

	return c.Validate(qp)
}
