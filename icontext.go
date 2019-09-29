package webfmwk

import (
	"net/http"

	"github.com/burgesQ/webfmwk/log"
)

// IContext Interface implement the context used in this project
type IContext interface {

	// SetRequest is used to save the request object
	SetRequest(rq *http.Request)

	// SetWriter is used to save the ResponseWriter obj
	// TODO: use io.Writer ?
	SetWriter(rw *http.ResponseWriter)

	// SetRoutes is used to save the route array
	SetRoutes(rt *Routes)

	// FetchContent extract the content from the body
	FetchContent(interface{}) error

	// Validate is used to validate a content of a request
	Validate(content interface{}) error

	// SetVars is used to save the url vars
	SetVars(vars map[string]string)

	// GetVar return the url var parameters. Empty string for none
	GetVar(key string) (val string)

	// SetQuery save the query param object
	SetQuery(query map[string][]string)

	// SetLogger set the logger of the ctx
	SetLogger(logger log.ILog)

	// GetQueries return the queries object
	GetQueries() map[string][]string

	// GetQuery fetch the query object key
	GetQuery(key string) (val string, ok bool)

	// IsPretty toggle the compact outptu mode
	IsPretty() bool

	// CheckHeader ensure the Content-Type in case of request
	CheckHeader() bool

	// OwnRecover is used to encapsulate the wanted panic
	OwnRecover()

	// JSONBlob answer the JSON content with the status code op
	JSONBlob(op int, content []byte) error
	JSON(int, interface{}) error

	//
	JSONNotImplemented(interface{}) error

	//
	JSONNoContent() error

	//
	JSONBadRequest(interface{}) error

	// 201
	JSONCreated(interface{}) error

	//
	JSONUnprocessable(interface{}) error

	// 200
	JSONOk(interface{}) error

	// 404
	JSONNotFound(interface{}) error

	//
	JSONConflict(interface{}) error

	//
	JSONInternalError(interface{}) error
}
