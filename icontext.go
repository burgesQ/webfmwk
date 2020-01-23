package webfmwk

import (
	"net/http"

	"github.com/burgesQ/webfmwk/v2/log"
)

// IContext Interface implement the context used in this project
type IContext interface {

	// SetRequest is used to save the request object
	SetRequest(rq *http.Request)

	// SetWriter is used to save the ResponseWriter obj
	// TODO: use io.Writer ?
	SetWriter(rw *http.ResponseWriter)

	// FetchContent extract the content from the body
	FetchContent(content interface{})

	// Validate is used to validate a content of the content params
	Validate(content interface{})

	// Decode load the query param in the content object
	DecodeQP(content interface{})

	// SetVars is used to save the url vars
	SetVars(vars map[string]string)

	// GetVar return the url var parameters. Empty string for none
	GetVar(key string) (val string)

	// SetQuery save the query param object
	SetQuery(query map[string][]string)

	// GetQueries return the queries object
	GetQueries() map[string][]string

	// GetQuery fetch the query object key
	GetQuery(key string) (val string, ok bool)

	// IsPretty toggle the compact outptu mode
	IsPretty() bool

	// SetLogger set the logger of the ctx
	SetLogger(logger log.ILog)

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

	//
	JSONNotImplemented(interface{})

	//
	JSONNoContent()

	//
	JSONBadRequest(interface{})

	// 201
	JSONCreated(interface{})

	//
	JSONUnprocessable(interface{})

	// 200
	JSONOk(interface{})

	// 404
	JSONNotFound(interface{})

	//
	JSONConflict(interface{})

	//
	JSONInternalError(interface{})

	//
	JSONAccepted(interface{})
}
