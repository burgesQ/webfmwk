package webfmwk

import "net/http"

// Context Interface iplement a ridiculous useless method
type IContext interface {
	SetRequest(*http.Request)
	// GetRequest() *http.Request

	SetWriter(*http.ResponseWriter)
	//	GetWriter() *http.ResponseWriter

	SetRoutes(*Routes)
	// GetRoutes() *Routes

	Validate(interface{}) error

	SetVars(map[string]string)
	GetVar(string) string
	SetQuery(map[string][]string)
	GetQuery(string) string
	GetQueries() map[string][]string
	IsPretty() bool

	CheckHeader() bool

	OwnRecover()

	FetchContent(interface{}) error

	JSONBlob(int, []byte) error
	JSON(int, interface{}) error
	JSONOk(interface{}) error
	JSONNotImplemented(interface{}) error
	JSONCreated(interface{}) error
	JSONConflict(interface{}) error
	JSONNotFound(interface{}) error
	JSONInternalError(interface{}) error
	JSONUnprocessable(interface{}) error
	//	JSON(interface{}) error
}
