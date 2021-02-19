package handler

import (
	golog "log"

	"github.com/burgesQ/gommon/log"
)

type customLogger struct {
	lg     log.Log
	prefix string
	short  []interface{}
}

func newLoggerRID(lg log.Log, prefix, rid string) log.Log {
	return customLogger{
		lg:     lg,
		short:  []interface{}{rid},
		prefix: prefix,
	}
}

func (l customLogger) Errorf(format string, v ...interface{}) {
	l.lg.Errorf("[!] "+l.prefix+format, append(l.short, v...)...)
}

func (l customLogger) Warnf(format string, v ...interface{}) {
	l.lg.Warnf(" [-] "+l.prefix+format, append(l.short, v...)...)
}

func (l customLogger) Infof(format string, v ...interface{}) {
	l.lg.Infof(" [-] "+l.prefix+format, append(l.short, v...)...)
}

func (l customLogger) Debugf(format string, v ...interface{}) {
	l.lg.Debugf("[+] "+l.prefix+format, append(l.short, v...)...)
}

func (l customLogger) Fatalf(format string, v ...interface{}) {
	l.lg.Fatalf("[x] "+l.prefix+format, append(l.short, v...)...)
}

func (l customLogger) SetLogLevel(level log.Level) bool {
	return l.lg.SetLogLevel(level)
}

func (l customLogger) GetErrorLogger() *golog.Logger {
	return l.lg.GetErrorLogger()
}
