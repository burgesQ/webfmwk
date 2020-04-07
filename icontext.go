package webfmwk

import (
	"context"
	"net/http"

	"github.com/burgesQ/webfmwk/v4/log"
)

type Header [2]string

// IContext Interface implement the context used in this project
type IContext interface {

	// GetRequest return the holded http.Request object
	GetRequest() *http.Request

	// SetRequest is used to save the request object
	SetRequest(rq *http.Request) IContext

	// SetWriter is used to save the ResponseWriter obj
	SetWriter(rw http.ResponseWriter) IContext

	// SetVars is used to save the url vars
	SetVars(vars map[string]string) IContext

	// GetVar return the url var parameters. Empty string for none
	GetVar(key string) (val string)

	// SetQuery save the query param object
	SetQuery(query map[string][]string) IContext

	// GetQueries return the queries object
	GetQueries() map[string][]string

	// GetQuery fetch the query object key
	GetQuery(key string) (val string, ok bool)

	// SetLogger set the logger of the ctx
	SetLogger(logger log.ILog) IContext

	// GetLogger return the logger of the ctx
	GetLogger() log.ILog

	// Save the given context object into the fmwk context
	SetContext(ctx context.Context) IContext

	// Fetch the previously saved context object
	GetContext() context.Context

	// GetRequest return the current request ID
	GetRequestID() string

	// SetRequest set the id of the current request
	SetRequestID(id string) IContext

	// FetchContent extract the content from the body
	FetchContent(content interface{})

	// Validate is used to validate a content of the content params
	Validate(content interface{})

	// Decode load the query param in the content object
	DecodeQP(content interface{})

	// IsPretty toggle the compact outptu mode
	IsPretty() bool

	// CheckHeader ensure the Content-Type of the request
	CheckHeader()

	// SetHeader set the header of the http response
	SetHeaders(headers ...Header)

	// OwnRecover is used to encapsulate the wanted panic
	OwnRecover()

	// SendResponse create & send a response according to the parameters
	SendResponse(op int, content []byte, headers ...Header)

	// JSONBlob answer the JSON content with the status code op
	JSONBlob(op int, content []byte)

	// JSON answer the JSON content with the status code op
	JSON(op int, content interface{})

	// JSONOk return the interface with an http.StatusOK (200)
	JSONOk(content interface{})

	// JSONCreated return the interface with an http.StatusCreated (201)
	JSONCreated(content interface{})

	// JSONAccepted return the interface with an http.StatusAccepted (202)
	JSONAccepted(content interface{})

	// JSONNoContent return an empty payload an http.StatusNoContent (204)
	JSONNoContent()

	// JSONBadRequest return the interface with an http.StatusBadRequest (400)
	JSONBadRequest(content interface{})

	// JSONUnauthorized return the interface with an http.StatusUnauthorized (401)
	JSONUnauthorized(content interface{})

	// JSONForbiden return the interface with an http.StatusForbidden (403)
	JSONForbiden(content interface{})

	// JSONNoContent return the interface with an http.StatusNotFound (404)
	JSONNotFound(content interface{})

	// JSONConflict return the interface with an http.StatusConflict (409)
	JSONConflict(content interface{})

	// JSONUnauthorized return the interface with an http.StatusUnprocessableEntity (422)
	JSONUnprocessable(content interface{})

	// JSONInternalError return the interface with an http.StatusInternalServerError (500)
	JSONInternalError(content interface{})

	// JSONNotImplemented return the interface with an http.StatusNotImplemented (501)
	JSONNotImplemented(content interface{})
}
