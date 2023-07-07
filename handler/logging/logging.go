package logging

import (
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/burgesQ/webfmwk/v5"
)

const (
	// HeaderRequestID hold the header name to which the RIP is attached
	HeaderRequestID = "X-Request-Id"
	_limitOutput    = 2048
)

// Handler generate an request ID and log information about
// the newly receive request
// The logger is then overloaded to add the request ID to every futur log message
func Handler(next webfmwk.HandlerFunc) webfmwk.HandlerFunc {
	return webfmwk.HandlerFunc(func(c webfmwk.Context) error {
		var (
			start = time.Now()
			fc    = c.GetFastContext()
			rid   = string(fc.Request.Header.Peek(HeaderRequestID))
		)

		if rid == "" {
			rid = strconv.Itoa(int(c.GetFastContext().ID()))
		}
		c.SetHeader(HeaderRequestID, rid)

		lg := c.GetLogger().SetPrefix("[" + rid + "]: ")

		c.SetLogger(lg)
		lg.Infof("--> %q [%s]%s ", webfmwk.GetIPFromRequest(fc), fc.Method(), fc.RequestURI())

		e := next(c)
		elapsed := time.Since(start)
		content := c.GetFastContext().Response.Body()
		l := len(content)

		if utf8.Valid(content) {
			if l > _limitOutput {
				lg.Debugf(">%s<", content[:_limitOutput])
			} else {
				lg.Debugf(">%s<", content)
			}
		}

		lg.Infof("<-- [%d]: took %s", fc.Response.Header.StatusCode(), elapsed)

		return e
	})
}
