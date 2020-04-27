// Package log implement the Log interface used by the webfmwk
package log

import (
	"fmt"
	"log"
)

// LogLevel
type Level int

const (
	LogERR Level = iota
	LogWARN
	LogINFO
	LogDEBUG
)

type (
	// Log interface implement the logging system inside the API
	Log interface {
		Errorf(format string, v ...interface{})
		Warnf(format string, v ...interface{})
		Infof(format string, v ...interface{})
		Debugf(format string, v ...interface{})
		Fatalf(format string, v ...interface{})
		SetLogLevel(level Level) bool
		GetErrorLogger() *log.Logger
	}

	logger struct {
		level Level
	}
)

var (
	_lg = logger{
		level: LogERR,
	}
	eLogger *log.Logger

	_out = map[Level]string{
		LogERR:   "! ERR  : ",
		LogWARN:  "* WARN : ",
		LogINFO:  "+ INFO : ",
		LogDEBUG: "- DBG  : ",
	}
)

func SetLogLevel(level Level) (ok bool) {
	if level >= LogERR && level <= LogDEBUG {
		_lg.level = level
		ok = true
	}
	return
}

func GetLogger() Log {
	return _lg
}

func (l *logger) logContent(level Level, format string, v ...interface{}) {
	if level <= l.level {
		fmt.Printf("%s"+format+"\n", append([]interface{}{
			_out[level],
		}, v...)...)
	}
}

func (l logger) GetErrorLogger() *log.Logger {
	return eLogger
}

func (l logger) SetLogLevel(level Level) bool {
	return SetLogLevel(level)
}

func (l logger) Debugf(format string, v ...interface{}) {
	l.logContent(LogDEBUG, format, v...)
}

func (l logger) Infof(format string, v ...interface{}) {
	l.logContent(LogINFO, format, v...)
}

func (l logger) Warnf(format string, v ...interface{}) {
	l.logContent(LogWARN, format, v...)
}

func (l logger) Errorf(format string, v ...interface{}) {
	l.logContent(LogERR, format, v...)
}

func (l logger) Fatalf(format string, v ...interface{}) {
	l.logContent(LogERR, format, v...)
	panic(fmt.Sprintf(format, v...))
}

func Debugf(format string, v ...interface{}) {
	_lg.logContent(LogDEBUG, format, v...)
}

func Infof(format string, v ...interface{}) {
	_lg.logContent(LogINFO, format, v...)
}

func Warnf(format string, v ...interface{}) {
	_lg.logContent(LogWARN, format, v...)
}

func Errorf(format string, v ...interface{}) {
	_lg.logContent(LogERR, format, v...)
}

func Fatalf(format string, v ...interface{}) {
	_lg.logContent(LogERR, format, v...)
}
