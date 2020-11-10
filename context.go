package webfmwk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/burgesQ/gommon/log"
	"github.com/burgesQ/gommon/pretty"
	"github.com/gorilla/schema"

	validator "github.com/go-playground/validator/v10"
)

const (
	_prettyTag   = "pretty"
	_limitOutput = 2014
)

type (
	Header [2]string

	// Context Interface implement the context used in this project
	Context interface {

		// GetRequest return the holded http.Request object
		GetRequest() *http.Request

		// SetRequest is used to save the request object
		SetRequest(rq *http.Request) Context

		// SetWriter is used to save the ResponseWriter obj
		SetWriter(rw http.ResponseWriter) Context

		// SetVars is used to save the url vars
		SetVars(vars map[string]string) Context

		// GetVar return the url var parameters. Empty string for none
		GetVar(key string) (val string)

		// SetQuery save the query param object
		SetQuery(query map[string][]string) Context

		// GetQueries return the queries object
		GetQueries() map[string][]string

		// GetQuery fetch the query object key
		GetQuery(key string) (val string, ok bool)

		// SetLogger set the logger of the ctx
		SetLogger(logger log.Log) Context

		// GetLogger return the logger of the ctx
		GetLogger() log.Log

		// Save the given context object into the fmwk context
		SetContext(ctx context.Context) Context

		// GetContext fetch the previously saved context object
		GetContext() context.Context

		// GetRequest return the current request ID
		GetRequestID() string

		// SetRequest set the id of the current request
		SetRequestID(id string) Context

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

		// IsPretty toggle the compact outptu mode
		IsPretty() bool

		// SetHeader set the header of the http response
		SetHeaders(headers ...Header)

		// SendResponse create & send a response according to the parameters
		SendResponse(op int, content []byte, headers ...Header) error

		// JSONBlob answer the JSON content with the status code op
		XMLBlob(op int, content []byte) error

		// JSONBlob answer the JSON content with the status code op
		JSONBlob(op int, content []byte) error

		// JSON answer the JSON content with the status code op
		JSON(op int, content interface{}) error

		// JSONOk return the interface with an http.StatusOK (200)
		JSONOk(content interface{}) error

		// JSONCreated return the interface with an http.StatusCreated (201)
		JSONCreated(content interface{}) error

		// JSONAccepted return the interface with an http.StatusAccepted (202)
		JSONAccepted(content interface{}) error

		// JSONNoContent return an empty payload an http.StatusNoContent (204)
		JSONNoContent() error

		// JSONBadRequest return the interface with an http.StatusBadRequest (400)
		JSONBadRequest(content interface{}) error

		// JSONUnauthorized return the interface with an http.StatusUnauthorized (401)
		JSONUnauthorized(content interface{}) error

		// JSONForbiden return the interface with an http.StatusForbidden (403)
		JSONForbiden(content interface{}) error

		// JSONNoContent return the interface with an http.StatusNotFound (404)
		JSONNotFound(content interface{}) error

		// JSONConflict return the interface with an http.StatusConflict (409)
		JSONConflict(content interface{}) error

		// JSONUnauthorized return the interface with an http.StatusUnprocessableEntity (422)
		JSONUnprocessable(content interface{}) error

		// JSONInternalError return the interface with an http.StatusInternalServerError (500)
		JSONInternalError(content interface{}) error

		// JSONNotImplemented return the interface with an http.StatusNotImplemented (501)
		JSONNotImplemented(content interface{}) error
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
	decoder = schema.NewDecoder()

	errMissingContentType   = NewNotAcceptable(NewError("Missing Content-Type header"))
	errNotJSON              = NewNotAcceptable(NewError("Content-Type is not application/json"))
	errUnprocessablePayload = NewUnprocessable(NewError("Unprocessable payload"))
)

// GetRequest implement Context
func (c *icontext) GetRequest() *http.Request {
	return c.r
}

// SetRequest implement Context
func (c *icontext) SetRequest(r *http.Request) Context {
	c.r = r
	return c
}

// SetWriter implement Context
func (c *icontext) SetWriter(w http.ResponseWriter) Context {
	c.w = w
	return c
}

// SetVars implement Context
func (c *icontext) SetVars(v map[string]string) Context {
	c.vars = v
	return c
}

// GetVar implement Context
func (c *icontext) GetVar(key string) string {
	return c.vars[key]
}

// IsPretty implement Context
func (c *icontext) IsPretty() bool {
	return len(c.query[_prettyTag]) > 0
}

// SetQuery implement Context
func (c *icontext) SetQuery(q map[string][]string) Context {
	c.query = q
	return c
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

// SetContext implement Context
func (c *icontext) SetContext(ctx context.Context) Context {
	c.ctx = ctx
	return c
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
		c.log.Errorf("[!] (%s) fetching payload: %s", c.GetRequestID(), e.Error())
		return errUnprocessablePayload
	}

	return nil
}

// Validate implement Context
// this implemt use validator to anotate & check struct
func (c *icontext) Validate(dest interface{}) ErrorHandled {
	if e := validate.Struct(dest); e != nil {
		c.log.Errorf("[!] (%s) validating : %s", c.GetRequestID(), e.Error())

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
		c.log.Errorf("[!] (%s) validating qp : %s", c.GetRequestID(), e.Error())
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
		if key != "" && val != "" {
			c.w.Header().Set(key, val)
		} else {
			c.log.Warnf("can't set header [%s] to [%s] (empty value)", key, val)
		}
	}
}

// response generate the http.Response with the holded http.ResponseWriter
// IDEA: add toggler `logReponse` ?
func (c *icontext) response(statusCode int, content []byte) error {
	var l = len(content)

	c.log.Infof("[-] (%s) : [%d](%d)", c.GetRequestID(), statusCode, l)

	if utf8.Valid(content) {
		if l > _limitOutput {
			c.log.Debugf("[-] (%s) : >%s<", c.GetRequestID(), content[:_limitOutput])
		} else {
			c.log.Debugf("[-] (%s) : >%s<", c.GetRequestID(), content)
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

// JSONBlob sent a JSON response already encoded
func (c *icontext) JSONBlob(statusCode int, content []byte) error {
	c.setHeaders(Header{"Accept", "application/json; charset=UTF-8"})

	if statusCode != http.StatusNoContent {
		c.setHeaders(Header{"Content-Type", "application/json; charset=UTF-8"},
			Header{"Produce", "application/json; charset=UTF-8"})
	}

	pcontent, e := pretty.SimplePrettyJSON(bytes.NewReader(content), c.IsPretty())
	if e != nil {
		return fmt.Errorf("canno't pretting the content : %w", e)
	}

	return c.response(statusCode, []byte(pcontent))
}

// JSON create a JSON response based on the param content.
func (c *icontext) JSON(statusCode int, content interface{}) error {
	data, e := json.Marshal(content)
	if e != nil {
		return fmt.Errorf("cannot json response : %w", e)
	}

	return c.JSONBlob(statusCode, data)
}

// JSONOk implement Context
func (c *icontext) JSONOk(content interface{}) error {
	return c.JSON(http.StatusOK, content)
}

// JSONCreated implement Context
func (c *icontext) JSONCreated(content interface{}) error {
	return c.JSON(http.StatusCreated, content)
}

// JSONAccepted implement Context
func (c *icontext) JSONAccepted(content interface{}) error {
	return c.JSON(http.StatusAccepted, content)
}

// JSONNoContent implement Context
func (c *icontext) JSONNoContent() error {
	return c.JSON(http.StatusNoContent, nil)
}

// JSONBadRequest implement Context
func (c *icontext) JSONBadRequest(content interface{}) error {
	return c.JSON(http.StatusBadRequest, content)
}

// JSONUnauthorized implement Context
func (c *icontext) JSONUnauthorized(content interface{}) error {
	return c.JSON(http.StatusUnauthorized, content)
}

// JSONForbiden implement Context
func (c *icontext) JSONForbiden(content interface{}) error {
	return c.JSON(http.StatusForbidden, content)
}

// JSONNotFound implement Context
func (c *icontext) JSONNotFound(content interface{}) error {
	return c.JSON(http.StatusNotFound, content)
}

// JSONConflict implement Context
func (c *icontext) JSONConflict(content interface{}) error {
	return c.JSON(http.StatusConflict, content)
}

// JSONUnprocessable implement Context
func (c *icontext) JSONUnprocessable(content interface{}) error {
	return c.JSON(http.StatusUnprocessableEntity, content)
}

// JSONInternalError implement Context
func (c *icontext) JSONInternalError(content interface{}) error {
	return c.JSON(http.StatusInternalServerError, content)
}

// JSONNotImplemented implement Context
func (c *icontext) JSONNotImplemented(content interface{}) error {
	return c.JSON(http.StatusNotImplemented, content)
}
