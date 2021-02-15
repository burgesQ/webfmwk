package handler

import (
	golog "log"

	"github.com/burgesQ/gommon/log"
	"github.com/burgesQ/webfmwk/v4"
)

type CustomLogger struct {
	lg      log.Log
	rid, ip string
	short   []interface{}
}

func newLogger(lg log.Log, rid, ip string) log.Log {
	return CustomLogger{
		lg:    lg,
		rid:   rid,
		ip:    ip,
		short: []interface{}{rid, ip},
	}
}

func (l CustomLogger) Errorf(format string, v ...interface{}) {
	l.lg.Errorf("[!] (%s) %s: "+format, append(l.short, v...)...)
}

func (l CustomLogger) Warnf(format string, v ...interface{}) {
	l.lg.Warnf("[-] (%s) %s: "+format, append(l.short, v...)...)
}

func (l CustomLogger) Infof(format string, v ...interface{}) {
	l.lg.Infof("[+] (%s) %s: "+format, append(l.short, v...)...)
}

func (l CustomLogger) Debugf(format string, v ...interface{}) {
	l.lg.Debugf("[x] (%s) %s: "+format, append(l.short, v...)...)
}

func (l CustomLogger) Fatalf(format string, v ...interface{}) {
	l.lg.Fatalf("[X] (%s) %s: "+format, append(l.short, v...)...)
}

func (l CustomLogger) SetLogLevel(level log.Level) bool {
	return l.lg.SetLogLevel(level)
}

func (l CustomLogger) GetErrorLogger() *golog.Logger {
	return l.lg.GetErrorLogger()
}

// Logging log information about the newly receive request
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
		c.SetLogger(newLogger(oldLog, c.GetRequestID(), IPAddress))
		c.GetLogger().Infof("[%s]%s", r.Method, r.RequestURI)

		return next(c)
	})
}
