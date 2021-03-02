package handler

import (
	"time"

	. "github.com/burgesQ/webfmwk/v4"
	"github.com/google/uuid"
)

const (
	HeaderRequestID = "X-REQUEST_ID"
	_logRIDPrefix   = "(%s): "
)

// Logging generate an request ID and log information about
// the newly receive request
// The logger is then overloaded to add the request ID to every futur log message
func Logging(next HandlerFunc) HandlerFunc {
	return HandlerFunc(func(c Context) error {
		var (
			start  = time.Now()
			oldLog = c.GetLogger()
			rid    = uuid.New().String()
			r      = c.GetRequest()
		)

		c.SetHeader(HeaderRequestID, rid)
		c.SetLogger(newLoggerRID(oldLog, _logRIDPrefix, rid))

		c.GetLogger().Infof("--> %q [%s]%s ", GetIPFromRequest(r), r.Method, r.RequestURI)
		e := next(c)
		c.GetLogger().Infof("<-- [STATUS_CODE]: took %s", time.Since(start))
		return e
	})
}
