package webfmwk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/burgesQ/gommon/pretty"
)

type JSONResponse interface {
	// JSONBlob answer the JSON content with the status code op
	JSONBlob(op int, content []byte) error

	// JSON answer the JSON content with the status code op
	JSON(op int, content interface{}) error

	// JSONOk return the interface with an http.StatusOK (200)
	JSONOk(content interface{}) error

	// JSONCreated return the interface with an http.StatusCreated (201)
	JSONCreated(content interface{}) error

	// JSONAccepted return the interface with an http.StatusAccepted (202)
	JSONAccepted(content interface{}) error

	// JSONNoContent return an empty payload an http.StatusNoContent (204)
	JSONNoContent() error

	// JSONBadRequest return the interface with an http.StatusBadRequest (400)
	JSONBadRequest(content interface{}) error

	// JSONUnauthorized return the interface with an http.StatusUnauthorized (401)
	JSONUnauthorized(content interface{}) error

	// JSONForbiden return the interface with an http.StatusForbidden (403)
	JSONForbiden(content interface{}) error

	// JSONNoContent return the interface with an http.StatusNotFound (404)
	JSONNotFound(content interface{}) error

	// JSONConflict return the interface with an http.StatusConflict (409)
	JSONConflict(content interface{}) error

	// JSONUnauthorized return the interface with an http.StatusUnprocessableEntity (422)
	JSONUnprocessable(content interface{}) error

	// JSONInternalError return the interface with an http.StatusInternalServerError (500)
	JSONInternalError(content interface{}) error

	// JSONNotImplemented return the interface with an http.StatusNotImplemented (501)
	JSONNotImplemented(content interface{}) error
}

// JSONBlob sent a JSON response already encoded
func (c *icontext) JSONBlob(statusCode int, content []byte) error {
	c.setHeaders(Header{"Accept", "application/json; charset=UTF-8"})

	if statusCode != http.StatusNoContent {
		c.setHeaders(Header{"Content-Type", "application/json; charset=UTF-8"},
			Header{"Produce", "application/json; charset=UTF-8"})
	}

	pcontent, e := pretty.SimplePrettyJSON(bytes.NewReader(content), c.IsPretty())
	if e != nil {
		return fmt.Errorf("canno't pretting the content : %w", e)
	}

	return c.response(statusCode, []byte(pcontent))
}

// JSON implement JSONResponse by returning a JSON encoded response.
func (c *icontext) JSON(statusCode int, content interface{}) error {
	data, e := json.Marshal(content)
	if e != nil {
		return fmt.Errorf("cannot json response : %w", e)
	}

	return c.JSONBlob(statusCode, data)
}

// JSONOk implement JSONResponse by returning a JSON encoded 200 response.
func (c *icontext) JSONOk(content interface{}) error {
	return c.JSON(http.StatusOK, content)
}

// JSONCreated implement JSONResponse by returning a JSON encoded 201 response.
func (c *icontext) JSONCreated(content interface{}) error {
	return c.JSON(http.StatusCreated, content)
}

// JSONAccepted implement JSONResponse by returning a JSON encoded 202 response.
func (c *icontext) JSONAccepted(content interface{}) error {
	return c.JSON(http.StatusAccepted, content)
}

// JSONNoContent implement JSONResponse by returning a JSON encoded 204 response.
func (c *icontext) JSONNoContent() error {
	return c.JSON(http.StatusNoContent, nil)
}

// JSONBadRequest implement JSONResponse by returning a JSON encoded 400 response.
func (c *icontext) JSONBadRequest(content interface{}) error {
	return c.JSON(http.StatusBadRequest, content)
}

// JSONUnauthorized implement JSONResponse by returning a JSON encoded 401 response.
func (c *icontext) JSONUnauthorized(content interface{}) error {
	return c.JSON(http.StatusUnauthorized, content)
}

// JSONForbiden implement JSONResponse by returning a JSON encoded 403 response.
func (c *icontext) JSONForbiden(content interface{}) error {
	return c.JSON(http.StatusForbidden, content)
}

// JSONNotFound implement JSONResponse by returning a JSON encoded 404 response.
func (c *icontext) JSONNotFound(content interface{}) error {
	return c.JSON(http.StatusNotFound, content)
}

// JSONConflict implement JSONResponse by returning a JSON encoded 429 response.
func (c *icontext) JSONConflict(content interface{}) error {
	return c.JSON(http.StatusConflict, content)
}

// JSONUnprocessable implement JSONResponse by returning a JSON encoded 422 response.
func (c *icontext) JSONUnprocessable(content interface{}) error {
	return c.JSON(http.StatusUnprocessableEntity, content)
}

// JSONInternalError implement JSONResponse by returning a JSON encoded 500 response.
func (c *icontext) JSONInternalError(content interface{}) error {
	return c.JSON(http.StatusInternalServerError, content)
}

// JSONNotImplemented implement JSONResponse by returning a JSON encoded 501 response.
func (c *icontext) JSONNotImplemented(content interface{}) error {
	return c.JSON(http.StatusNotImplemented, content)
}
