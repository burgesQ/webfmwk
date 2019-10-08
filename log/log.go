package log

import (
	"fmt"
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
	lg = logger{
		level: LogERR,
	}

	out = map[int]string{
		LogERR:   "! ERR  : ",
		LogWARN:  "* WARN : ",
		LogINFO:  "+ INFO : ",
		LogDEBUG: "- DBG  : ",
	}
)

func init() {

}

func SetLogLevel(level int) (ok bool) {
	if level >= LogERR && level <= LogDEBUG {
		lg.level = level
		ok = true
	}
	return
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
	lg.logContent(LogDEBUG, format, v...)
}

func Infof(format string, v ...interface{}) {
	lg.logContent(LogINFO, format, v...)
}

func Warnf(format string, v ...interface{}) {
	lg.logContent(LogWARN, format, v...)
}

func Errorf(format string, v ...interface{}) {
	lg.logContent(LogERR, format, v...)
}

func Fatalf(format string, v ...interface{}) {
	lg.logContent(LogERR, format, v...)
}
