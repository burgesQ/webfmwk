package webfmwk

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/burgesQ/webfmwk/log"
	"github.com/burgesQ/webfmwk/util"
	"gopkg.in/go-playground/validator.v9"
)

// TODO: proper interface
// Context Interface iplement a ridiculous useless method
type Context interface {
	// all method used by CContext
	GetName() string

	SetReader()
	Getreader()

	SetWriter()
	GetWriter()

	SetRoutes()
	GetRoutes()

	SetVars()
	GetVars()

	SetQuery()
	GetQuery()

	GetPretty()
	SetPretty()

	SetCustomeContext()
	GetCustomContext()
}

// TODO: rename to Custom Context
// TODO: remove Custom Context field (up cast to user class)
// CContext hold the data used by the request
type CContext struct {
	// http request reqder
	R *http.Request
	// http request writer
	W *http.ResponseWriter
	// routes list
	Routes *Routes
	// form
	Vars map[string]string
	// query param
	Query         map[string][]string
	Pretty        bool
	CustomContext interface{}
}

// AnonymousError struct is used to output error code
type (
	AnonymousError struct {
		Error string `json:"error"`
	}
	// header [2]string
)

var (
	MissingContentType = func(c *CContext) {
		c.JSON(http.StatusNotAcceptable, AnonymousError{"Missing Content-Type header"})
	}
	MismatchContentType = func(c *CContext) {
		c.JSON(http.StatusNotAcceptable, AnonymousError{"Content-Type is not application/json"})
	}
	UnprocessableEntity = func(c *CContext) {
		c.JSONUnprocessable(AnonymousError{"Unprocessable Payload, wrong json ?"})
	}
	UnprocessableQueryParam = func(c *CContext) {
		c.JSONUnprocessable(AnonymousError{"Unprocessable query param"})
	}
	formChecker = validator.New()
)

func (c *CContext) GetName() string {
	return "base context"
}

// HandleError log the error and create a appropriate server response.
// The log format is defined by context. By default error.Error() is provided to the log output.
// The server response is defined by reason.
func (c *CContext) HandleError(e error, code int, reason, context string, data ...interface{}) error {
	log.Errorf(context, e.Error(), data[:])
	c.JSON(code, AnonymousError{reason})
	return e
}

// Param fetch a url param
func (c *CContext) Param(key string) string {
	return c.Vars[key]
}

// CheckFetchContent extract the content from the body.
// Then it check the struct via gorilla/schema.
func (c *CContext) CheckFetchContent(dest interface{}) (e error) {
	defer c.R.Body.Close()
	if e = json.NewDecoder(c.R.Body).Decode(&dest); e != nil {
		log.Errorf("error while decoding the payload : %s", e.Error())
		return
	}
	log.Debugf("successfully extracted from body : %#v", dest)
	return
}

// Validate validate the json payload destructured into a go struct
func (c *CContext) Validate(dest interface{}) (e error) {
	if e = formChecker.Struct(dest); e != nil {
		log.Errorf("error while validating the payload :\n%s", e.Error())
	}
	return
}

// CheckHeader from content request (POST, PUT, PATCH)
func (c *CContext) CheckHeader() bool {
	ctype := c.R.Header.Get("Content-Type")

	// if no application/json
	if len(ctype) == 0 {
		MissingContentType(c)
		return false
	} else if !strings.HasPrefix(ctype, "application/json") {
		MismatchContentType(c)
		return false
	}

	return true
}

// Use the pretty json utilitary to create well, pretty json ? :nerd_face:
func prettyJson(r io.Reader, pretty *bool) string {

	o := new(bytes.Buffer)
	pj := util.NewPrettyJson(r, o)

	if !(*pretty) {
		pj = pj.SetCompactMode()
	}

	pj.Start()

	if err := pj.Close(); err != nil {
		log.Errorf("JSON error: %s", err.Error())
	}

	return o.String()
}

// SetHeader set one header for the next response
func (c *CContext) SetHeader(key, val string) {
	(*c.W).Header().Set(key, val)
}

// SetHeader set a array of headers for the next response
func (c *CContext) SetHeaders(headers ...[2]string) {
	for _, h := range headers {
		if h[0] != "" && h[1] != "" {
			c.SetHeader(h[0], h[1])
		} else {
			log.Warnf("can't set header [%s] to [%s] (empty value)", h[0], h[1])
		}
	}
}

// SendResponse create & send a response according to the parameters
func (c *CContext) SendResponse(statusCode int, content []byte, headers ...[2]string) error {
	c.SetHeaders(headers...)
	return c.response(statusCode, content)
}

// JSONBlob sent a JSON response allready encoded.
func (c *CContext) JSONBlob(statusCode int, content []byte) error {

	c.SetHeader("Accept", "application/json; charset=UTF-8")
	if statusCode != http.StatusNoContent {
		c.SetHeader("Content-Type", "application/json; charset=UTF-8")
		c.SetHeader("Produce", "application/json; charset=UTF-8")
	}

	pcontent := prettyJson(bytes.NewReader(content), &c.Pretty)

	return c.response(statusCode, []byte(pcontent))
}

// JSON create a JSON response based on the param content.
func (c *CContext) JSON(statusCode int, content interface{}) error {

	if data, err := json.Marshal(content); err != nil {
		log.Errorf("%s", err.Error())
		return c.JSONInternalError(AnonymousError{"Error creating the JSON response."})
	} else {
		return c.JSONBlob(statusCode, data)
	}
}

func (c *CContext) JSONNotImplemented(content interface{}) error {
	return c.JSON(http.StatusNotImplemented, content)
}

func (c *CContext) JSONNoContent() error {
	return c.JSON(http.StatusNoContent, nil)
}

func (c *CContext) JSONBadRequest(content interface{}) error {
	return c.JSON(http.StatusBadRequest, content)
}

func (c *CContext) JSONCreated(content interface{}) error {
	return c.JSON(http.StatusCreated, content)
}

func (c *CContext) JSONUnprocessable(content interface{}) error {
	return c.JSON(http.StatusUnprocessableEntity, content)
}

func (c *CContext) JSONOk(content interface{}) error {
	return c.JSON(http.StatusOK, content)
}

func (c *CContext) JSONNotFound(content interface{}) error {
	return c.JSON(http.StatusNotFound, content)
}

func (c *CContext) JSONConflict(content interface{}) error {
	return c.JSON(http.StatusConflict, content)
}

// InternalError create a error 500 with the reason why
func (c *CContext) JSONInternalError(content interface{}) error {
	return c.JSON(http.StatusInternalServerError, content)
}

func (c *CContext) response(statusCode int, content []byte) error {

	(*c.W).WriteHeader(statusCode)
	(*c.W).Write(content)

	if utf8.Valid(content) {
		log.Infof("[%d](%d): >%s<", statusCode, len(content), content)
	} else {
		log.Infof("[%d](%d)", statusCode, len(content))
	}

	return nil
}

func (c *CContext) OwnRecover() {
	if r := recover(); r != nil {
		switch err := r.(type) {
		case ErrorHandled:
			c.JSON(err.GetOPCode(), err.GetContent())
		default:
			panic(err)
		}
	}
}
