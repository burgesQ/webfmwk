package webfmwk

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/segmentio/encoding/json"
)

const _contentType = "application/json; charset=UTF-8"

type JSONHTTPResponse interface { //nolint: interfacebloat
	// JSONOk return the interface with an http.StatusOK (200).
	JSONOk(content interface{}) error

	// JSONCreated return the interface with an http.StatusCreated (201).
	JSONCreated(content interface{}) error

	// JSONAccepted return the interface with an http.StatusAccepted (202).
	JSONAccepted(content interface{}) error

	// JSONNoContent return an empty payload an http.StatusNoContent (204).
	JSONNoContent() error

	// JSONBadRequest return the interface with an http.StatusBadRequest (400).
	JSONBadRequest(content interface{}) error

	// JSONUnauthorized return the interface with an http.StatusUnauthorized (401).
	JSONUnauthorized(content interface{}) error

	// JSONForbiden return the interface with an http.StatusForbidden (403).
	JSONForbidden(content interface{}) error

	// JSONNoContent return the interface with an http.StatusNotFound (404).
	JSONNotFound(content interface{}) error

	// JSONMethodNotAllowed return the interface with an http.NotAllowed (405).
	JSONMethodNotAllowed(content interface{}) error

	// JSONConflict return the interface with an http.StatusConflict (409).
	JSONConflict(content interface{}) error

	// JSONUnauthorized return the interface with an http.StatusUnprocessableEntity (422).
	JSONUnprocessable(content interface{}) error

	// JSONInternalError return the interface with an http.StatusInternalServerError (500).
	JSONInternalError(content interface{}) error

	// JSONNotImplemented return the interface with an http.StatusNotImplemented (501).
	JSONNotImplemented(content interface{}) error
}

// JSONResponse interface is used to answer JSON content to the client.
type JSONResponse interface {
	JSONHTTPResponse

	// JSONBlob answer the JSON content with the status code op.
	JSONBlob(op int, content []byte) error

	// JSON answer the JSON content with the status code op.
	JSON(op int, content interface{}) error
}

// JSONBlob sent a JSON response already encoded
func (c *icontext) JSONBlob(statusCode int, content []byte) error {
	var (
		out bytes.Buffer
		run = func() error {
			if c.IsPretty() {
				return json.Indent(&out, content, "", "  ")
			}

			return json.Compact(&out, content)
		}
	)

	if e := run(); e != nil {
		c.slog.Error("cannot prettying the content", "error", e)
	} else {
		content = out.Bytes()
	}

	if statusCode != http.StatusNoContent {
		c.SetContentType(_contentType)
		c.SetHeader("Produce", _contentType)
	}

	c.SetHeader("Accept", _contentType)

	return c.response(statusCode, content)
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

// JSONForbidden implement JSONResponse by returning a JSON encoded 403 response.
func (c *icontext) JSONForbidden(content interface{}) error {
	return c.JSON(http.StatusForbidden, content)
}

// JSONNotFound implement JSONResponse by returning a JSON encoded 404 response.
func (c *icontext) JSONNotFound(content interface{}) error {
	return c.JSON(http.StatusNotFound, content)
}

// JSONMethodNotAllowed implement JSONResponse by returning a JSON encoded 405 response.
func (c *icontext) JSONMethodNotAllowed(content interface{}) error {
	return c.JSON(http.StatusMethodNotAllowed, content)
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
