package webfmwk

import (
	"context"
	"net/http"

	"github.com/burgesQ/webfmwk/v3/log"
)

// IContext Interface implement the context used in this project
type IContext interface {
	GetRequestID() string
	SetRequestID(id string) IContext

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

	// Save the given context object into the fmwk context
	SetContext(ctx context.Context) IContext

	// Fetch the previously saved context object
	GetContext() context.Context

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

	// OwnRecover is used to encapsulate the wanted panic
	OwnRecover()

	// SendResponse create & send a response according to the parameters
	SendResponse(op int, content []byte, headers ...[2]string)

	// JSONBlob answer the JSON content with the status code op
	JSONBlob(op int, content []byte)

	// JSON answer the JSON content with the status code op
	JSON(op int, content interface{})

	// 200
	JSONOk(interface{})

	// 201
	JSONCreated(interface{})

	// 202
	JSONAccepted(interface{})

	//
	JSONNotImplemented(interface{})

	//
	JSONNoContent()

	//
	JSONBadRequest(interface{})

	//
	JSONUnprocessable(interface{})

	// 404
	JSONNotFound(interface{})

	//
	JSONConflict(interface{})

	// 500
	JSONInternalError(interface{})
}
