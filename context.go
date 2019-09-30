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
func (c *Context) FetchContent(dest interface{}) {
	defer c.r.Body.Close()
	if e := json.NewDecoder(c.r.Body).Decode(&dest); e != nil {
		c.log.Errorf("while decoding the payload : %s", e.Error())
		panic(NewUnprocessable(AnonymousError{"Unprocessable payload, wrong json ?"}))
	}
}

// Validate implement IContext
// this implemt use validator to anotate & check struct
func (c Context) Validate(dest interface{}) {
	if e := formChecker.Struct(dest); e != nil {
		c.log.Errorf("error while validating the payload :\n%s", e.Error())
		panic(NewUnprocessable(AnonymousError{e.Error()}))
	}
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
func (c Context) CheckHeader() {
	if ctype := c.r.Header.Get("Content-Type"); len(ctype) == 0 {
		panic(NewNotAcceptable(AnonymousError{"Missing Content-Type header"}))
	} else if !strings.HasPrefix(ctype, "application/json") {
		panic(NewNotAcceptable(AnonymousError{"Content-Type is not application/json"}))
	}
}

// OwnRecover implement IContext
func (c Context) OwnRecover() {
	if r := recover(); r != nil {
		switch err := r.(type) {
		case IErrorHandled:
			c.JSON(err.GetOPCode(), err.GetContent())
		default:
			log.Errorf("catched %T %#v", err, err)
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

func (c *Context) response(statusCode int, content []byte) {
	(*c.w).WriteHeader(statusCode)
	(*c.w).Write(content)

	if utf8.Valid(content) {
		c.log.Infof("[%d](%d): >%s<", statusCode, len(content), content)
	} else {
		c.log.Infof("[%d](%d)", statusCode, len(content))
	}
}

// Send Response implement IContext
func (c *Context) SendResponse(statusCode int, content []byte, headers ...[2]string) {
	c.setHeaders(headers...)
	c.response(statusCode, content)
}

// JSONBlob sent a JSON response already encoded
func (c *Context) JSONBlob(statusCode int, content []byte) {

	c.setHeader("Accept", "application/json; charset=UTF-8")
	if statusCode != http.StatusNoContent {
		c.setHeader("Content-Type", "application/json; charset=UTF-8")
		c.setHeader("Produce", "application/json; charset=UTF-8")
	}

	pcontent, err := util.SimplePrettyJSON(bytes.NewReader(content), c.IsPretty())
	if err != nil {
		c.log.Errorf("while prettier the content : %s", err.Error())
	}

	c.response(statusCode, []byte(pcontent))
}

// JSON create a JSON response based on the param content.
func (c *Context) JSON(statusCode int, content interface{}) {
	data, err := json.Marshal(content)
	if err != nil {
		c.log.Errorf("%s", err.Error())
		panic(NewInternal(AnonymousError{"Error creating the JSON response."}))
	}
	c.JSONBlob(statusCode, data)
}

// JSONOk implement IContext
func (c *Context) JSONOk(content interface{}) {
	c.JSON(http.StatusOK, content)
}

// JSONNoContent implement IContext
func (c *Context) JSONNoContent() {
	c.JSON(http.StatusNoContent, nil)
}

// JSONBadRequest implement IContext
func (c *Context) JSONBadRequest(content interface{}) {
	c.JSON(http.StatusBadRequest, content)
}

// JSONCreated implement IContext
func (c *Context) JSONCreated(content interface{}) {
	c.JSON(http.StatusCreated, content)
}

// JSONUnprocessable implement IContext
func (c *Context) JSONUnprocessable(content interface{}) {
	c.JSON(http.StatusUnprocessableEntity, content)
}

// JSONNotFound implement IContext
func (c *Context) JSONNotFound(content interface{}) {
	c.JSON(http.StatusNotFound, content)
}

// JSONConflict implement IContext
func (c *Context) JSONConflict(content interface{}) {
	c.JSON(http.StatusConflict, content)
}

// JSONNotImplemented implement IContext
func (c *Context) JSONNotImplemented(content interface{}) {
	c.JSON(http.StatusNotImplemented, content)
}

// JSONInternalError implement IContext
func (c *Context) JSONInternalError(content interface{}) {
	c.JSON(http.StatusInternalServerError, content)
}
