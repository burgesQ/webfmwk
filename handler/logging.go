package handler

import (
	"github.com/burgesQ/webfmwk/v4"
	"github.com/google/uuid"
)

const (
	HeaderRequestID = "X-REQUEST_ID"
	_logRIDPrefix   = "(%s): "
)

// Logging generate an request ID and log information about
// the newly receive request
// The logger is then overloaded to add the request ID to every futur log message
func Logging(next webfmwk.HandlerFunc) webfmwk.HandlerFunc {
	return webfmwk.HandlerFunc(func(c webfmwk.Context) error {
		var (
			oldLog = c.GetLogger()
			rid    = uuid.New().String()
			r      = c.GetRequest()
		)

		c.SetHeader(HeaderRequestID, rid)
		c.SetLogger(newLoggerRID(oldLog, _logRIDPrefix, rid))
		c.GetLogger().Infof("%q --> [%s]%s ", webfmwk.GetIPFromRequest(r), r.Method, r.RequestURI)

		return next(c)
	})
}
