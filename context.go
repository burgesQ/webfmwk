package webfmwk

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/burgesQ/webfmwk/log"
	"github.com/burgesQ/webfmwk/util"
	validator "gopkg.in/go-playground/validator.v9"
)

// Context implement the IContext interface
// It hold the data used by the request
type (
	Context struct {
		r      *http.Request
		w      *http.ResponseWriter
		routes *Routes
		vars   map[string]string
		query  map[string][]string
		custom IContext
		log    log.ILog
	}

	// AnonymousError struct is used to answer error
	AnonymousError struct {
		Error string `json:"error"`
	}
)

var (
	missingContentType = func(c *Context) {
		c.JSON(http.StatusNotAcceptable, AnonymousError{"Missing Content-Type header"})
	}
	mismatchContentType = func(c *Context) {
		c.JSON(http.StatusNotAcceptable, AnonymousError{"Content-Type is not application/json"})
	}
	unprocessableEntity = func(c *Context) {
		c.JSONUnprocessable(AnonymousError{"Unprocessable Payload, wrong json ?"})
	}
	unprocessableQueryParam = func(c *Context) {
		c.JSONUnprocessable(AnonymousError{"Unprocessable query param"})
	}
	validationFailed = func(c *Context, e error) {
		c.JSONUnprocessable(AnonymousError{e.Error()})
	}

	formChecker = validator.New()
)

// SetRequest implement IContext
func (c *Context) SetRequest(r *http.Request) {
	c.r = r
}

// SetWriter implement IContext
func (c *Context) SetWriter(w *http.ResponseWriter) {
	c.w = w
}

// SetRoutes implement IContext
func (c *Context) SetRoutes(r *Routes) {
	c.routes = r
}

// FetchContent implement IContext
func (c *Context) FetchContent(dest interface{}) (e error) {
	defer c.r.Body.Close()
	if e = json.NewDecoder(c.r.Body).Decode(&dest); e != nil {
		c.log.Errorf("while decoding the payload : %s", e.Error())
		//		panic(New422(AnonymousError{"Unprocessable Payload, wrong json ?"}))
		unprocessableEntity(c)
		return e
	}
	return
}

// Validate implement IContext
// this implemt use validator to anotate & check struct
func (c Context) Validate(dest interface{}) (e error) {
	if e = formChecker.Struct(dest); e != nil {
		c.log.Errorf("error while validating the payload :\n%s", e.Error())
		// panic(New422(AnonymousError{e.Error()}))
		validationFailed(&c, e)
		return e
	}
	return
}

// SetVars implement IContext
func (c *Context) SetVars(v map[string]string) {
	c.vars = v
}

// GetVar implement IContext
func (c Context) GetVar(key string) string {
	return c.vars[key]
}

// SetQuery implement IContext
func (c *Context) SetQuery(q map[string][]string) {
	c.query = q
}

// GetQueries implement IContext
func (c *Context) GetQueries() map[string][]string {
	return c.query
}

// GetQuery implement IContext
func (c *Context) GetQuery(key string) (string, bool) {
	if len(c.query[key]) > 0 {
		return c.query[key][0], true
	}
	return "", false
}

// IsPretty implement IContext
func (c Context) IsPretty() bool {
	if len(c.query["pjson"]) > 0 {
		return true
	}
	return false
}

func (c *Context) SetLogger(logger log.ILog) {
	c.log = logger

}

// CheckHeader implement IContext
func (c Context) CheckHeader() (ret bool) {
	ctype := c.r.Header.Get("Content-Type")

	// if no application/json
	if len(ctype) == 0 {
		missingContentType(&c)
		// panic(New406(AnonymousError{"Missing Content-Type header"}))
	} else if !strings.HasPrefix(ctype, "application/json") {
		mismatchContentType(&c)
		// panic(New406(AnonymousError{"Content-Type is not application/json"}))
	} else {
		ret = true
	}

	return
}

// OwnRecover implement IContext
func (c Context) OwnRecover() {
	if r := recover(); r != nil {
		switch err := r.(type) {
		case ErrorHandled:
			c.JSON(err.GetOPCode(), err.GetContent())
		default:
			panic(err)
		}
	}
}

func (c *Context) setHeader(key, val string) {
	(*c.w).Header().Set(key, val)
}

func (c *Context) setHeaders(headers ...[2]string) {
	for _, h := range headers {
		if h[0] != "" && h[1] != "" {
			c.setHeader(h[0], h[1])
		} else {
			c.log.Warnf("can't set header [%s] to [%s] (empty value)", h[0], h[1])
		}
	}
}

func (c *Context) response(statusCode int, content []byte) error {

	(*c.w).WriteHeader(statusCode)
	(*c.w).Write(content)

	if utf8.Valid(content) {
		c.log.Infof("[%d](%d): >%s<", statusCode, len(content), content)
	} else {
		c.log.Infof("[%d](%d)", statusCode, len(content))
	}

	return nil
}

// sendResponse create & send a response according to the parameters
func (c *Context) sendResponse(statusCode int, content []byte, headers ...[2]string) error {
	c.setHeaders(headers...)
	return c.response(statusCode, content)
}

// JSONBlob sent a JSON response already encoded
func (c *Context) JSONBlob(statusCode int, content []byte) error {

	c.setHeader("Accept", "application/json; charset=UTF-8")
	if statusCode != http.StatusNoContent {
		c.setHeader("Content-Type", "application/json; charset=UTF-8")
		c.setHeader("Produce", "application/json; charset=UTF-8")
	}

	pcontent, err := util.SimplePrettyJSON(bytes.NewReader(content), c.IsPretty())
	if err != nil {
		c.log.Errorf("while prettier the content : %s", err.Error())
	}

	return c.response(statusCode, []byte(pcontent))
}

// JSON create a JSON response based on the param content.
func (c *Context) JSON(statusCode int, content interface{}) error {
	data, err := json.Marshal(content)
	if err != nil {
		c.log.Errorf("%s", err.Error())
		return c.JSONInternalError(AnonymousError{"Error creating the JSON response."})
	}
	return c.JSONBlob(statusCode, data)
}

// JSONOk implement IContext
func (c *Context) JSONOk(content interface{}) error {
	return c.JSON(http.StatusOK, content)
}

// JSONNoContent implement IContext
func (c *Context) JSONNoContent() error {
	return c.JSON(http.StatusNoContent, nil)
}

// JSONBadRequest implement IContext
func (c *Context) JSONBadRequest(content interface{}) error {
	return c.JSON(http.StatusBadRequest, content)
}

// JSONCreated implement IContext
func (c *Context) JSONCreated(content interface{}) error {
	return c.JSON(http.StatusCreated, content)
}

// JSONUnprocessable implement IContext
func (c *Context) JSONUnprocessable(content interface{}) error {
	return c.JSON(http.StatusUnprocessableEntity, content)
}

// JSONNotFound implement IContext
func (c *Context) JSONNotFound(content interface{}) error {
	return c.JSON(http.StatusNotFound, content)
}

// JSONConflict implement IContext
func (c *Context) JSONConflict(content interface{}) error {
	return c.JSON(http.StatusConflict, content)
}

// JSONNotImplemented implement IContext
func (c *Context) JSONNotImplemented(content interface{}) error {
	return c.JSON(http.StatusNotImplemented, content)
}

// JSONInternalError implement IContext
func (c *Context) JSONInternalError(content interface{}) error {
	return c.JSON(http.StatusInternalServerError, content)
}
