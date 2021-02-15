package webfmwk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/burgesQ/gommon/log"
	validator "github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
)

const (
	_prettyTag   = "pretty"
	_limitOutput = 2048
)

type (
	Header [2]string

	RequestID interface {
		// GetRequest return the current request ID
		GetRequestID() string

		// SetRequest set the id of the current request
		SetRequestID(id string) Context
	}

	SendResponse interface {
		JSONResponse
		XMLResponse

		// SendResponse create & send a response according to the parameters
		SendResponse(op int, content []byte, headers ...Header) error
	}

	XMLResponse interface {
		// JSONBlob answer the JSON content with the status code op
		XMLBlob(op int, content []byte) error
	}

	InputHandling interface {
		// FetchContent extract the content from the body
		FetchContent(content interface{}) ErrorHandled

		// Validate is used to validate a content of the content params
		Validate(content interface{}) ErrorHandled

		// FetchAndValidateContent fetch the content then validate it
		FetchAndValidateContent(content interface{}) ErrorHandled

		// Decode load the query param in the content object
		DecodeQP(content interface{}) ErrorHandled

		// CheckHeader ensure the Content-Type of the request
		CheckHeader() ErrorHandled
	}

	// Context Interface implement the context used in this project
	Context interface {
		SendResponse
		RequestID
		InputHandling

		// GetRequest return the holded http.Request object
		GetRequest() *http.Request

		// GetVar return the url var parameters. Empty string for none
		GetVar(key string) (val string)

		// GetQueries return the queries object
		GetQueries() map[string][]string

		// GetQuery fetch the query object key
		GetQuery(key string) (val string, ok bool)

		// SetLogger set the logger of the ctx
		SetLogger(logger log.Log) Context

		// GetLogger return the logger of the ctx
		GetLogger() log.Log

		// GetContext fetch the previously saved context object
		GetContext() context.Context

		// IsPretty toggle the compact outptu mode
		IsPretty() bool

		// SetHeader set the header of the http response
		SetHeaders(headers ...Header)
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
		uid   string
	}
)

var (
	// decoder annotation : `schema` : gorilla
	decoder                 = schema.NewDecoder()
	errMissingContentType   = NewNotAcceptable(NewError("Missing Content-Type header"))
	errNotJSON              = NewNotAcceptable(NewError("Content-Type is not application/json"))
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

// IsPretty implement Context
func (c *icontext) IsPretty() bool {
	return len(c.query[_prettyTag]) > 0
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

// GetRequestID implement Context
func (c *icontext) GetRequestID() string {
	return c.uid
}

// SetRequestID implement Context
func (c *icontext) SetRequestID(id string) Context {
	c.uid = id
	return c
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

// CheckHeader implement Context
func (c *icontext) CheckHeader() ErrorHandled {
	if ctype := c.r.Header.Get("Content-Type"); ctype == "" {
		return errMissingContentType
	} else if !strings.HasPrefix(ctype, "application/json") {
		c.log.Errorf("%q != application/json", ctype)
		return errNotJSON
	}

	return nil
}

// SetHeaders implement Context
func (c *icontext) SetHeaders(headers ...Header) {
	c.setHeaders(headers...)
}

// setHeader set the header of the holded http.ResponseWriter
func (c *icontext) setHeaders(headers ...Header) {
	for _, h := range headers {
		key, val := h[0], h[1]
		if key == "" || val == "" {
			c.log.Warnf("can't set header [%s] to [%s] (empty value)", key, val)

			return
		}

		c.w.Header().Set(key, val)
	}
}

// response generate the http.Response with the holded http.ResponseWriter
// IDEA: add toggler `logReponse` ?
func (c *icontext) response(statusCode int, content []byte) error {
	var l = len(content)

	c.log.Infof("[%d](%d)", statusCode, l)

	if utf8.Valid(content) {
		if l > _limitOutput {
			c.log.Debugf(">%s<", content[:_limitOutput])
		} else {
			c.log.Debugf(">%s<", content)
		}
	}

	c.w.WriteHeader(statusCode)

	if _, e := c.w.Write(content); e != nil {
		return fmt.Errorf("cannot write response : %w", e)
	}

	return nil
}

// SendResponse implement Context
func (c *icontext) SendResponse(statusCode int, content []byte, headers ...Header) error {
	c.setHeaders(headers...)
	return c.response(statusCode, content)
}

// XMLBlob sent a XML response already encoded
func (c *icontext) XMLBlob(statusCode int, content []byte) error {
	c.setHeaders(Header{"Content-Type", "application/xml; charset=UTF-8"},
		Header{"Produce", "application/xml; charset=UTF-8"})
	return c.response(statusCode, content)
}
