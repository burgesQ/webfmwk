package webfmwk

import (
	"bytes"

	"github.com/segmentio/encoding/json"
	"github.com/valyala/fasthttp"
)

var (
	// ErrMissingContentType is returned in case of missing content type header.
	ErrMissingContentType = NewNotAcceptable(NewError("Missing Content-Type header"))

	// ErrNotJSON is returned when the content type isn't json
	ErrNotJSON = NewNotAcceptable(NewError("Content-Type is not application/json"))

	_prefixContentType = []byte("application/json")
)

//
// helper method
//

// GetIPFromRequest try to extract the source IP from the
// request headers (X-Real-IP and X-Forwareded-For).
func GetIPFromRequest(fc *fasthttp.RequestCtx) string {
	if ip := fc.Request.Header.Peek("X-Real-IP"); len(ip) > 0 {
		return string(ip)
	} else if ip = fc.Request.Header.Peek("X-Forwarded-For"); len(ip) > 0 {
		return string(ip)
	}

	return fc.RemoteAddr().String()
}

//
// internal handler
//

func contentIsJSON(next HandlerFunc) HandlerFunc {
	return HandlerFunc(func(c Context) error {
		var (
			fc = c.GetFastContext()
			m  = fc.Method()
		)

		if string(m) == POST || string(m) == PUT || string(m) == PATCH {
			if ctype := fc.Request.Header.Peek("Content-Type"); len(ctype) == 0 {
				return ErrMissingContentType
			} else if !bytes.HasPrefix(ctype, _prefixContentType) {
				return ErrNotJSON
			}
		}

		return next(c)
	})
}

func handleNotFound(c Context) error {
	fc := c.GetFastContext()

	c.GetLogger().Infof("[!] 404 reached for [%s] %s %s",
		GetIPFromRequest(fc), fc.Method(), fc.RequestURI())

	return c.JSONNotFound(json.RawMessage(`{"status":404,"message":"not found"}`))
}

func handleNotAllowed(c Context) error {
	fc := c.GetFastContext()

	c.GetLogger().Infof("[!] 405 reached for [%s] %s %s",
		GetIPFromRequest(fc), fc.Method(), fc.RequestURI())

	return c.JSONMethodNotAllowed(json.RawMessage(`{"status":405,"message":"method not allowed"}`))
}
