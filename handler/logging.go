package handler

import (
	"time"

	. "github.com/burgesQ/webfmwk/v4"
	"github.com/google/uuid"
)

const (
	HeaderRequestID = "X-REQUEST_ID"
)

// Logging generate an request ID and log information about
// the newly receive request
// The logger is then overloaded to add the request ID to every futur log message
func Logging(next HandlerFunc) HandlerFunc {
	return HandlerFunc(func(c Context) error {
		var (
			start = time.Now()
			rid   = uuid.New().String()
			r     = c.GetRequest()
		)

		c.SetHeader(HeaderRequestID, rid)

		c.GetLogger().Infof(" >%s< --> %q [%s]%s ", rid, GetIPFromRequest(r), r.Method, r.RequestURI)
		e := next(c)
		c.GetLogger().Infof(" >%s< <-- [STATUS_CODE]: took %s", rid, time.Since(start))

		return e
	})
}
