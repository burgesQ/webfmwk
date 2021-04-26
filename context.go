package webfmwk

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/burgesQ/webfmwk/v4/log"
	validator "github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
)

const (
	_prettyTag   = "pretty"
	_limitOutput = 2048
)

type (
	InputHandling interface {
		// FetchContent extract the content from the body
		FetchContent(content interface{}) ErrorHandled

		// Validate is used to validate a content of the content params
		Validate(content interface{}) ErrorHandled

		// FetchAndValidateContent fetch the content then validate it
		FetchAndValidateContent(content interface{}) ErrorHandled

		// Decode load the query param in the content object
		DecodeQP(content interface{}) ErrorHandled
	}

	ContextLogger interface {
		// SetLogger set the logger of the ctx
		SetLogger(logger log.Log) Context

		// GetLogger return the logger of the ctx
		GetLogger() log.Log
	}

	// Context Interface implement the context used in this project
	Context interface {
		SendResponse
		InputHandling
		ContextLogger

		// GetRequest return the holded http.Request object
		GetRequest() *http.Request

		// GetVar return the url var parameters. Empty string for none
		GetVar(key string) (val string)

		// GetQueries return the queries object
		GetQueries() map[string][]string

		// GetQuery fetch the query object key
		GetQuery(key string) (val string, ok bool)

		// GetContext fetch the previously saved context object
		GetContext() context.Context
	}

	// icontext implement the Context interface
	// It hold the data used by the request
	icontext struct {
		r     *http.Request
		w     http.ResponseWriter
		vars  map[string]string
		query map[string][]string
		log   log.Log
		ctx   context.Context
	}
)

var (
	// decoder annotation : `schema` : gorilla
	decoder                 = schema.NewDecoder()
	errUnprocessablePayload = NewUnprocessable(NewError("Unprocessable payload"))
)

// GetRequest implement Context
func (c *icontext) GetRequest() *http.Request {
	return c.r
}

// GetVar implement Context
func (c *icontext) GetVar(key string) string {
	return c.vars[key]
}

// GetQueries implement Context
func (c *icontext) GetQueries() map[string][]string {
	return c.query
}

// GetQuery implement Context
func (c *icontext) GetQuery(key string) (string, bool) {
	if len(c.query[key]) > 0 {
		return c.query[key][0], true
	}

	return "", false
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

// FetchContent implement Context
// It load payload in the dest interface{} using the system json library
func (c *icontext) FetchContent(dest interface{}) ErrorHandled {
	defer c.r.Body.Close()

	if e := json.NewDecoder(c.r.Body).Decode(&dest); e != nil {
		c.log.Errorf("fetching payload: %s", e.Error())
		return errUnprocessablePayload
	}

	return nil
}

// Validate implement Context
// this implemt use validator to anotate & check struct
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

func (c *icontext) FetchAndValidateContent(dest interface{}) ErrorHandled {
	if e := c.FetchContent(&dest); e != nil {
		return e
	}

	return c.Validate(dest)
}

// DecodeQP implement Context
func (c *icontext) DecodeQP(dest interface{}) (e ErrorHandled) {
	if e := decoder.Decode(dest, c.GetQueries()); e != nil {
		c.log.Errorf("validating qp : %s", e.Error())
		return NewUnprocessable(NewErrorFromError(e))
	}

	return nil
}
