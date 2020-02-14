package webfmwk

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/burgesQ/webfmwk/v3/log"
	"github.com/burgesQ/webfmwk/v3/pretty"
	"github.com/gorilla/schema"

	en_translator "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

// Context implement the IContext interface
// It hold the data used by the request
type (
	Context struct {
		r     *http.Request
		w     http.ResponseWriter
		vars  map[string]string
		query map[string][]string
		log   log.ILog
		ctx   *context.Context
		uid   int
	}

	// AnonymousError struct is used to answer error
	AnonymousError struct {
		Error string `json:"error"`
	}

	ValidationError struct {
		Error validator.ValidationErrorsTranslations `json:"error"`
	}
)

var (
	// validate annotation : `validate` : go-playground
	validate = validator.New()
	// decoder annotation : `schema` : gorilla
	decoder = schema.NewDecoder()
	// this is usually know or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	// TODO: extract from Accept-Language ?
	trans ut.Translator

	inited = false
)

func (c *Context) initOnce() {
	if inited {
		return
	}

	var (
		// tranlator
		en = en_translator.New()

		// universal translator
		uni *ut.UniversalTranslator
	)

	// universal translator
	uni = ut.New(en, en)
	// this is usually know or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	trans, _ = uni.GetTranslator("en")
	// TODO: check ret val
	if e := en_translations.RegisterDefaultTranslations(validate, trans); e != nil {
		c.log.Fatalf("cannot init translations : %v", e)
	}

	inited = true
}

// SetRequest implement IContext
func (c *Context) GetRequestID() int {
	return c.uid
}

// SetRequest implement IContext
func (c *Context) SetRequestID(id int) IContext {
	c.uid = id
	return c
}

// SetRequest implement IContext
func (c *Context) SetRequest(r *http.Request) IContext {
	c.r = r
	return c
}

// SetWriter implement IContext
func (c *Context) SetWriter(w http.ResponseWriter) IContext {
	c.w = w
	return c
}

// SetVars implement IContext
func (c *Context) SetVars(v map[string]string) IContext {
	c.vars = v
	return c
}

// GetVar implement IContext
func (c Context) GetVar(key string) string {
	return c.vars[key]
}

// SetQuery implement IContext
func (c *Context) SetQuery(q map[string][]string) IContext {
	c.query = q
	return c
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

// SetLogger implement IContext
func (c *Context) SetLogger(logger log.ILog) IContext {
	c.log = logger
	return c
}

// SetContext implement IContext
func (c *Context) SetContext(ctx *context.Context) IContext {
	c.ctx = ctx
	return c
}

// GetContent implement IContext
func (c *Context) GetContext() *context.Context {
	return c.ctx
}

// FetchContent implement IContext
// It load payload in the dest interface{} using the system json library
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
	c.initOnce()
	if e := validate.Struct(dest); e != nil {
		out := e.(validator.ValidationErrors).Translate(trans)
		c.log.Errorf("error while validating the payload :\n%s", out)
		panic(NewUnprocessable(ValidationError{out}))
	}
}

// DecodeQP implement IContext
func (c Context) DecodeQP(dest interface{}) {
	if e := decoder.Decode(dest, c.GetQueries()); e != nil {
		c.log.Errorf("error while validating the query params :\n%s", e.Error())
		c.log.Debugf("[%#v]", dest)
		panic(NewUnprocessable(AnonymousError{e.Error()}))
	}
}

// IsPretty implement IContext
func (c Context) IsPretty() bool {
	return len(c.query["pretty"]) > 0
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
		switch e := r.(type) {
		case IErrorHandled:
			c.JSON(e.GetOPCode(), e.GetContent())
		default:
			c.log.Errorf("catched %T %#v", e, e)
			panic(e)
		}
	}
}

func (c *Context) setHeader(key, val string) {
	c.w.Header().Set(key, val)
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
	if utf8.Valid(content) {
		c.log.Infof("[%d](%d): >%s<", statusCode, len(content), content)
	} else {
		c.log.Infof("[%d](%d)", statusCode, len(content))
	}

	c.w.WriteHeader(statusCode)

	l, e := c.w.Write(content)
	if e != nil {
		c.log.Errorf("while sending response (%d) : %s", l, e.Error())
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

	pcontent, err := pretty.SimplePrettyJSON(bytes.NewReader(content), c.IsPretty())
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

// JSONOk implement IContext
func (c *Context) JSONAccepted(content interface{}) {
	c.JSON(http.StatusAccepted, content)
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
