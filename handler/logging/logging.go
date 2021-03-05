package logging

import (
	"time"
	"unicode/utf8"

	"github.com/burgesQ/webfmwk/v5"
	"github.com/google/uuid"
)

const (
	// HeaderRequestID hold the header name to which the RIP is attached
	HeaderRequestID = "X-Request-Id"
	_limitOutput    = 2048
)

// TODO: check if RID in request header

// Handler generate an request ID and log information about
// the newly receive request
// The logger is then overloaded to add the request ID to every futur log message
func Handler(next webfmwk.HandlerFunc) webfmwk.HandlerFunc {
	return webfmwk.HandlerFunc(func(c webfmwk.Context) error {
		var (
			start = time.Now()
			lg    = c.GetLogger()
			fc    = c.GetFastContext()
			rid   = uuid.New().String()
		)

		c.SetHeader(HeaderRequestID, rid)

		lg.Infof("[%s] --> %q [%s]%s ", rid, webfmwk.GetIPFromRequest(fc), fc.Method(), fc.RequestURI())

		e := next(c)
		elapsed := time.Since(start)
		content := c.GetFastContext().Response.Body()
		l := len(content)

		if utf8.Valid(content) {
			if l > _limitOutput {
				lg.Debugf("[%s] >%s<", rid, content[:_limitOutput])
			} else {
				lg.Debugf("[%s] >%s<", rid, content)
			}
		}

		lg.Infof("[%s] <-- [%d]: took %s", rid, fc.Response.Header.StatusCode(), elapsed)

		return e
	})
}
