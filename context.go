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
	validator "gopkg.in/go-playground/validator.v9"
)

// Context hold the data used by the request
type Context struct {
	r             *http.Request
	w             *http.ResponseWriter
	routes        *Routes
	vars          map[string]string
	query         map[string][]string
	CustomContext interface{}
}

// AnonymousError struct is used to output error code
type (
	AnonymousError struct {
		Error string `json:"error"`
	}
)

var (
	MissingContentType = func(c *Context) {
		c.JSON(http.StatusNotAcceptable, AnonymousError{"Missing Content-Type header"})
	}
	MismatchContentType = func(c *Context) {
		c.JSON(http.StatusNotAcceptable, AnonymousError{"Content-Type is not application/json"})
	}
	UnprocessableEntity = func(c *Context) {
		c.JSONUnprocessable(AnonymousError{"Unprocessable Payload, wrong json ?"})
	}
	UnprocessableQueryParam = func(c *Context) {
		c.JSONUnprocessable(AnonymousError{"Unprocessable query param"})
	}

	formChecker = validator.New()
)

func (c *Context) SetRequest(r *http.Request) {
	c.r = r
}

func (c Context) GetRequest() *http.Request {
	return c.r
}

func (c *Context) SetWriter(w *http.ResponseWriter) {
	c.w = w
}

func (c Context) GetWriter() *http.ResponseWriter {
	return c.w
}

func (c *Context) SetRoutes(r *Routes) {
	c.routes = r
}

func (c Context) GetRoutes() *Routes {
	return c.routes
}

func (c *Context) SetVars(v map[string]string) {
	c.vars = v
}

// Param fetch a url param
func (c Context) GetVar(key string) string {
	return c.vars[key]
}

func (c *Context) SetQuery(q map[string][]string) {
	c.query = q
}

func (c *Context) GetQueries() map[string][]string {
	return c.query
}

func (c *Context) GetQuery(key string) string {
	if len(c.query[key]) > 0 {
		return c.query[key][0]
	}
	return ""
}

func (c Context) IsPretty() bool {
	if len(c.query["pjson"]) > 0 {
		return true
	}
	return false
}

// FetchContent extract the content from the body.
func (c *Context) FetchContent(dest interface{}) (e error) {
	defer c.r.Body.Close()
	if e = json.NewDecoder(c.r.Body).Decode(&dest); e != nil {
		log.Errorf("while decoding the payload : %s", e.Error())
		return
	}
	return
}

// Validate validate the json payload destructured into a go struct
func (c Context) Validate(dest interface{}) (e error) {
	if e = formChecker.Struct(dest); e != nil {
		log.Errorf("error while validating the payload :\n%s", e.Error())
	}
	return
}

// CheckHeader from content request (POST, PUT, PATCH)
func (c Context) CheckHeader() (ret bool) {
	ctype := c.r.Header.Get("Content-Type")

	// if no application/json
	if len(ctype) == 0 {
		MissingContentType(&c)
		ret = false
	} else if !strings.HasPrefix(ctype, "application/json") {
		MismatchContentType(&c)
		ret = false
	} else {
		ret = true
	}

	return
}

// Use the pretty json utilitary to create well, pretty json ? :nerd_face:
func prettyJSON(r io.Reader, pretty bool) string {

	o := new(bytes.Buffer)
	pj := util.NewPrettyJson(r, o)

	if !pretty {
		pj = pj.SetCompactMode()
	}

	pj.Start()

	if err := pj.Close(); err != nil {
		log.Errorf("JSON error: %s", err.Error())
	}

	return o.String()
}

// SetHeader set one header for the next response
func (c *Context) SetHeader(key, val string) {
	(*c.w).Header().Set(key, val)
}

// SetHeader set a array of headers for the next response
func (c *Context) SetHeaders(headers ...[2]string) {
	for _, h := range headers {
		if h[0] != "" && h[1] != "" {
			c.SetHeader(h[0], h[1])
		} else {
			log.Warnf("can't set header [%s] to [%s] (empty value)", h[0], h[1])
		}
	}
}

// SendResponse create & send a response according to the parameters
func (c *Context) SendResponse(statusCode int, content []byte, headers ...[2]string) error {
	c.SetHeaders(headers...)
	return c.response(statusCode, content)
}

// JSONBlob sent a JSON response allready encoded.
func (c *Context) JSONBlob(statusCode int, content []byte) error {

	c.SetHeader("Accept", "application/json; charset=UTF-8")
	if statusCode != http.StatusNoContent {
		c.SetHeader("Content-Type", "application/json; charset=UTF-8")
		c.SetHeader("Produce", "application/json; charset=UTF-8")
	}

	pcontent := prettyJSON(bytes.NewReader(content), c.IsPretty())

	return c.response(statusCode, []byte(pcontent))
}

// JSON create a JSON response based on the param content.
func (c *Context) JSON(statusCode int, content interface{}) error {

	if data, err := json.Marshal(content); err != nil {
		log.Errorf("%s", err.Error())
		return c.JSONInternalError(AnonymousError{"Error creating the JSON response."})
	} else {
		return c.JSONBlob(statusCode, data)
	}
}

func (c *Context) JSONNotImplemented(content interface{}) error {
	return c.JSON(http.StatusNotImplemented, content)
}

func (c *Context) JSONNoContent() error {
	return c.JSON(http.StatusNoContent, nil)
}

func (c *Context) JSONBadRequest(content interface{}) error {
	return c.JSON(http.StatusBadRequest, content)
}

func (c *Context) JSONCreated(content interface{}) error {
	return c.JSON(http.StatusCreated, content)
}

func (c *Context) JSONUnprocessable(content interface{}) error {
	return c.JSON(http.StatusUnprocessableEntity, content)
}

func (c *Context) JSONOk(content interface{}) error {
	return c.JSON(http.StatusOK, content)
}

func (c *Context) JSONNotFound(content interface{}) error {
	return c.JSON(http.StatusNotFound, content)
}

func (c *Context) JSONConflict(content interface{}) error {
	return c.JSON(http.StatusConflict, content)
}

// InternalError create a error 500 with the reason why
func (c *Context) JSONInternalError(content interface{}) error {
	return c.JSON(http.StatusInternalServerError, content)
}

func (c *Context) response(statusCode int, content []byte) error {

	(*c.w).WriteHeader(statusCode)
	(*c.w).Write(content)

	if utf8.Valid(content) {
		log.Infof("[%d](%d): >%s<", statusCode, len(content), content)
	} else {
		log.Infof("[%d](%d)", statusCode, len(content))
	}

	return nil
}

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

// HandleError log the error and create a appropriate server response.
// The log format is defined by context. By default error.Error() is provided to the log output.
// The server response is defined by reason.
// func (c *Context) HandleError(e error, code int, reason, context string, data ...interface{}) error {
// 	log.Errorf(context, e.Error(), data[:])
// 	c.JSON(code, AnonymousError{reason})
// 	return e
// }
