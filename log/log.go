// Package log implement the Log interface used by the webfmwk
package log

import (
	"fmt"
)

type (
	// Log interface implement the logging system inside the API
	Log interface {
		Printf(format string, args ...interface{})

		Debugf(format string, v ...interface{})
		Infof(format string, v ...interface{})
		Warnf(format string, v ...interface{})
		Errorf(format string, v ...interface{})
		Fatalf(format string, v ...interface{})
	}

	logger struct {
		level Level
	}
)

var (
	_lg = logger{
		level: LogErr,
	}
)

func GetLogger() Log {
	return _lg
}

func (l *logger) logContentf(level Level, format string, v ...interface{}) {
	if level <= l.level || level == LogPrint {
		fmt.Printf("%s"+format+"\n", append([]interface{}{
			_out[level],
		}, v...)...)
	}
}

func (l logger) Printf(format string, v ...interface{}) {
	l.logContentf(LogPrint, format, v...)
}

func (l logger) Debugf(format string, v ...interface{}) {
	l.logContentf(LogDebug, format, v...)
}

func (l logger) Infof(format string, v ...interface{}) {
	l.logContentf(LogInfo, format, v...)
}

func (l logger) Warnf(format string, v ...interface{}) {
	l.logContentf(LogWarning, format, v...)
}

func (l logger) Errorf(format string, v ...interface{}) {
	l.logContentf(LogErr, format, v...)
}

func (l logger) Fatalf(format string, v ...interface{}) {
	l.logContentf(LogErr, format, v...)
	panic(fmt.Sprintf(format, v...))
}

func Debugf(format string, v ...interface{}) {
	_lg.logContentf(LogDebug, format, v...)
}

func Infof(format string, v ...interface{}) {
	_lg.logContentf(LogInfo, format, v...)
}

func Warnf(format string, v ...interface{}) {
	_lg.logContentf(LogWarning, format, v...)
}

func Errorf(format string, v ...interface{}) {
	_lg.logContentf(LogErr, format, v...)
}

func Fatalf(format string, v ...interface{}) {
	_lg.logContentf(LogErr, format, v...)
}
