package handler

import (
	golog "log"

	"github.com/burgesQ/gommon/log"
	"github.com/burgesQ/webfmwk/v4"
	"github.com/google/uuid"
)

type CustomLogger struct {
	lg      log.Log
	rid, ip string
	short   []interface{}
}

// Logging generate an request ID and log information about
// the newly receive request
// The logger is then overloaded to add the request ID to every futur log message
func Logging(next webfmwk.HandlerFunc) webfmwk.HandlerFunc {
	return webfmwk.HandlerFunc(func(c webfmwk.Context) error {
		r := c.GetRequest()

		IPAddress := r.Header.Get("X-Real-Ip")
		if IPAddress == "" {
			IPAddress = r.Header.Get("X-Forwarded-For")
		}
		if IPAddress == "" {
			IPAddress = r.RemoteAddr
		}

		oldLog := c.GetLogger()

		rid := uuid.New().String()

		// TODO: add rid to response header ?
		c.SetHeaders(webfmwk.Header{"X-REQUEST_ID", rid})

		c.SetLogger(newLogger(oldLog, rid))
		c.GetLogger().Infof("%q --> [%s]%s ", IPAddress, r.Method, r.RequestURI)

		return next(c)
	})
}

func newLogger(lg log.Log, rid string) log.Log {
	return CustomLogger{
		lg:    lg,
		short: []interface{}{rid},
	}
}

func (l CustomLogger) Errorf(format string, v ...interface{}) {
	l.lg.Errorf("[!] (%s): "+format, append(l.short, v...)...)
}

func (l CustomLogger) Warnf(format string, v ...interface{}) {
	l.lg.Warnf(" [-] (%s): "+format, append(l.short, v...)...)
}

func (l CustomLogger) Infof(format string, v ...interface{}) {
	l.lg.Infof(" [-] (%s): "+format, append(l.short, v...)...)
}

func (l CustomLogger) Debugf(format string, v ...interface{}) {
	l.lg.Debugf("[+] (%s): "+format, append(l.short, v...)...)
}

func (l CustomLogger) Fatalf(format string, v ...interface{}) {
	l.lg.Fatalf("[x] (%s): "+format, append(l.short, v...)...)
}

func (l CustomLogger) SetLogLevel(level log.Level) bool {
	return l.lg.SetLogLevel(level)
}

func (l CustomLogger) GetErrorLogger() *golog.Logger {
	return l.lg.GetErrorLogger()
}
