// Package log implement the Log interface used by the webfmwk
package log

import (
	"fmt"

	qlog "github.com/burgesQ/log"
)

type (
	logger struct {
		prefix string
		level  Level
	}
)

// default internal logger
var _lg = logger{level: LogDebug, prefix: ""}

// GetLogger return an struct fullfilling the Log interface
func GetLogger() qlog.Log { return _lg }

// SetPrefix implement the LogPrefix interface
func (l logger) SetPrefix(prefix string) qlog.Log {
	return logger{level: l.level, prefix: prefix}
}

// GetPrefix implement the LogPrefix interface
func (l logger) GetPrefix() string {
	return l.prefix
}

func (l *logger) logContentf(level Level, format string, v ...interface{}) {
	if level <= l.level || level == LogPrint {
		//nolint: forbidigo
		fmt.Printf("%s"+l.prefix+format+"\n", append([]interface{}{
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
