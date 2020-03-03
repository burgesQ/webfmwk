package log

import (
	"fmt"
	"log"
)

const (
	LogERR   = 0
	LogWARN  = 1
	LogINFO  = 2
	LogDEBUG = 3
)

type logger struct {
	level int
}

var (
	_lg = logger{
		level: LogERR,
	}

	_out = map[int]string{
		LogERR:   "! ERR  : ",
		LogWARN:  "* WARN : ",
		LogINFO:  "+ INFO : ",
		LogDEBUG: "- DBG  : ",
	}
)

func SetLogLevel(level int) (ok bool) {
	if level >= LogERR && level <= LogDEBUG {
		_lg.level = level
		ok = true
	}
	return
}

func (l *logger) logContent(level int, format string, v ...interface{}) {
	if level <= l.level {
		fmt.Printf("%s"+format+"\n", append([]interface{}{
			_out[level],
		}, v...)...)
	}
}

func GetLogger() ILog {
	return _lg
}

var eLogger *log.Logger

func (l logger) GetErrorLogger() *log.Logger {
	return eLogger
}

func (l logger) SetLogLevel(level int) bool {
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
