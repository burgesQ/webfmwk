package log

import (
	"fmt"
)

const (
	logERR   = 0
	logWARN  = 1
	logINFO  = 2
	logDEBUG = 3
)

type logger struct {
	level int
}

var (
	lg = logger{
		level: logERR,
	}

	out = map[int]string{
		logERR:   "! ERR  : ",
		logWARN:  "* WARN : ",
		logINFO:  "+ INFO : ",
		logDEBUG: "- DBG  : ",
	}
)

func SetLogLevel(level int) {
	if level >= logERR && level <= logDEBUG {
		lg.level = level
	}
}

func (l *logger) logContent(level int, format string, v ...interface{}) {
	if level <= l.level {
		fmt.Printf("%s"+format+"\n", append([]interface{}{
			out[level],
		}, v...)...)
	}
}

func GetLogger() ILog {
	return lg
}

func (l logger) Debugf(format string, v ...interface{}) {
	l.logContent(logDEBUG, format, v...)
}

func (l logger) Infof(format string, v ...interface{}) {
	l.logContent(logINFO, format, v...)
}

func (l logger) Warnf(format string, v ...interface{}) {
	l.logContent(logWARN, format, v...)
}

func (l logger) Errorf(format string, v ...interface{}) {
	l.logContent(logERR, format, v...)
}

func (l logger) Fatalf(format string, v ...interface{}) {
	l.logContent(logERR, format, v...)
	panic(fmt.Sprintf(format, v...))
}

func Debugf(format string, v ...interface{}) {
	lg.logContent(logDEBUG, format, v...)
}

func Infof(format string, v ...interface{}) {
	lg.logContent(logINFO, format, v...)
}

func Warnf(format string, v ...interface{}) {
	lg.logContent(logWARN, format, v...)
}

func Errorf(format string, v ...interface{}) {
	lg.logContent(logERR, format, v...)
}

func Fatalf(format string, v ...interface{}) {
	lg.logContent(logERR, format, v...)
}
