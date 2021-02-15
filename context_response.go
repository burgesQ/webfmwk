package webfmwk

import (
	"fmt"
	"unicode/utf8"
)

type (
	Header [2]string

	SendResponse interface {
		JSONResponse
		XMLResponse

		// SendResponse create & send a response according to the parameters
		SendResponse(op int, content []byte, headers ...Header) error

		// SetHeader set the header of the http response
		SetHeaders(headers ...Header)
	}

	XMLResponse interface {
		// JSONBlob answer the JSON content with the status code op
		XMLBlob(op int, content []byte) error
	}
)

// SetHeaders implement Context
func (c *icontext) SetHeaders(headers ...Header) {
	c.setHeaders(headers...)
}

// setHeader set the header of the holded http.ResponseWriter
func (c *icontext) setHeaders(headers ...Header) {
	for _, h := range headers {
		key, val := h[0], h[1]
		if key == "" || val == "" {
			c.log.Warnf("can't set header [%s] to [%s] (empty value)", key, val)

			return
		}

		c.w.Header().Set(key, val)
	}
}

// XMLBlob sent a XML response already encoded
func (c *icontext) XMLBlob(statusCode int, content []byte) error {
	c.setHeaders(Header{"Content-Type", "application/xml; charset=UTF-8"},
		Header{"Produce", "application/xml; charset=UTF-8"})
	return c.response(statusCode, content)
}

// SendResponse implement Context
func (c *icontext) SendResponse(statusCode int, content []byte, headers ...Header) error {
	c.setHeaders(headers...)
	return c.response(statusCode, content)
}

// response generate the http.Response with the holded http.ResponseWriter
// IDEA: add toggler `logReponse` ?
func (c *icontext) response(statusCode int, content []byte) error {
	var l = len(content)

	c.log.Infof("[%d](%d)", statusCode, l)

	if utf8.Valid(content) {
		if l > _limitOutput {
			c.log.Debugf(">%s<", content[:_limitOutput])
		} else {
			c.log.Debugf(">%s<", content)
		}
	}

	c.w.WriteHeader(statusCode)

	if _, e := c.w.Write(content); e != nil {
		return fmt.Errorf("cannot write response : %w", e)
	}

	return nil
}
