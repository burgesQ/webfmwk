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

// GetLogger return an struct fullfilling the Log interface
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

// Printf implement the Log interface.
func (l logger) Printf(format string, v ...interface{}) {
	l.logContentf(LogPrint, format, v...)
}

// Debugf implement the Log interface.
func (l logger) Debugf(format string, v ...interface{}) {
	l.logContentf(LogDebug, format, v...)
}

// Infof implement the Log interface.
func (l logger) Infof(format string, v ...interface{}) {
	l.logContentf(LogInfo, format, v...)
}

// Warnf implement the Log interface.
func (l logger) Warnf(format string, v ...interface{}) {
	l.logContentf(LogWarning, format, v...)
}

// Errorf implement the Log interface.
func (l logger) Errorf(format string, v ...interface{}) {
	l.logContentf(LogErr, format, v...)
}

// Fatalf implement the Log interface.
func (l logger) Fatalf(format string, v ...interface{}) {
	l.logContentf(LogErr, format, v...)
	panic(fmt.Sprintf(format, v...))
}

// Debugf output a debug message.
func Debugf(format string, v ...interface{}) {
	_lg.logContentf(LogDebug, format, v...)
}

// Infof output an info message.
func Infof(format string, v ...interface{}) {
	_lg.logContentf(LogInfo, format, v...)
}

// Warnf output a warning message.
func Warnf(format string, v ...interface{}) {
	_lg.logContentf(LogWarning, format, v...)
}

// Errorf output an error message.
func Errorf(format string, v ...interface{}) {
	_lg.logContentf(LogErr, format, v...)
}

// Fatalf output an fatal message and then panic.
func Fatalf(format string, v ...interface{}) {
	_lg.logContentf(LogErr, format, v...)
}
