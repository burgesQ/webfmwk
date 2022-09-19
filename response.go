package webfmwk

type (
	// Header represent a header in a string key:value form.
	Header [2]string

	// SendResponse interface is used to reponde content to the client.
	SendResponse interface {
		JSONResponse

		// SendResponse create & send a response according to the parameters.
		SendResponse(op int, content []byte, headers ...Header) error

		// SetHeader set the header of the http response.
		SetHeaders(headers ...Header)

		// SetHeader set the k header to value v.
		SetHeader(k, v string)

		// SetContentType set the Content-Type to v.
		SetContentType(v string)

		// SetContentType set the Content-Type to v.
		SetStatusCode(code int)

		// IsPretty toggle the compact output mode.
		IsPretty() bool
	}
)

// SetHeaders implement Context
func (c *icontext) SetHeaders(headers ...Header) {
	c.setHeaders(headers...)
}

// SetHeaders implement Context
func (c *icontext) SetHeader(k, v string) {
	c.setHeaders(Header{k, v})
}

// IsPretty implement Context
func (c *icontext) IsPretty() bool {
	return c.QueryArgs().Has(_prettyTag)
}

// SendResponse implement Context
func (c *icontext) SendResponse(statusCode int, content []byte, headers ...Header) error {
	c.setHeaders(headers...)

	return c.response(statusCode, content)
}

// setHeader set the header of the holded http.ResponseWriter
func (c *icontext) setHeaders(headers ...Header) {
	for _, h := range headers {
		key, val := h[0], h[1]
		if key == "" || val == "" {
			c.log.Warnf("can't set header [%s] to [%s] (empty value)", key, val)

			return
		}

		c.Response.Header.Set(key, val)
	}
}

// response answer the client
// IDEA: add toggler `logReponse` ?
func (c *icontext) response(statusCode int, content []byte) error {
	c.SetStatusCode(statusCode)
	c.SetBody(content)

	return nil
}
