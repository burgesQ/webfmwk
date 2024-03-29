package slogging

import (
	"log/slog"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/burgesQ/webfmwk/v6"
)

const (
	// HeaderRequestID hold the header name to which the RIP is attached
	HeaderRequestID = "X-Request-Id"
	_limitOutput    = 2048
)

var LogLevel = ((*slog.Logger).Info)

// NewHandler generate an request ID and log information about
// the newly receive request
// The logger is then overloaded to add the request ID to every futur log message
func NewHandler(method ...func(*slog.Logger, string, ...any)) webfmwk.Handler {
	mh := ((*slog.Logger).Info)

	if len(method) > 0 {
		mh = method[0]
	}

	return func(next webfmwk.HandlerFunc) webfmwk.HandlerFunc {
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

			lg := c.GetStructuredLogger().With(
				slog.Group("request",
					slog.String("id", rid),
					slog.String("ip", webfmwk.GetIPFromRequest(fc)),
					slog.String("method", string(fc.Method())),
					slog.String("uri", string(fc.RequestURI()))))

			c.SetStructuredLogger(lg)

			(mh)(lg, "--> new request")

			// LogLevel

			e := next(c)
			elapsed := time.Since(start)
			content := c.GetFastContext().Response.Body()
			l := len(content)

			if utf8.Valid(content) {
				if l > _limitOutput {
					lg.Debug("trunkated response",
						slog.String("body", string(content[:_limitOutput])),
						slog.Int("lim", _limitOutput))
				} else {
					lg.Debug("full response", slog.String("body", string(content)))
				}
			}

			(mh)(lg, "<-- request done",
				slog.Int("code", fc.Response.Header.StatusCode()),
				slog.Any("elapsed", elapsed))

			return e
		})
	}
}
