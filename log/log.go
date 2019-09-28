package log

import (
	"errors"
	"fmt"
	"log"
	"log/syslog"
	"os"
)

// Level is a shortcut
type Level syslog.Priority

const (
	// LoggerSTDOUT is used to output to stdout
	LoggerSTDOUT = iota
	// LoggerSTDERR is used to putpit to stderr
	LoggerSTDERR

	loggerMask = 3

	// LogFormatShort is the flag for short logs
	LogFormatShort = 0

	// LogFormatLong is the flag for long logs
	LogFormatLong = 1 << 3

	// From /usr/include/sys/syslog.h.
	// These are the same on Linux, BSD, and OS X.

	// LogEmerg is the emergency log level
	LogEmerg LogLevel = iota
	// LogAlert is the alert log level
	LogAlert
	// LogCrit is the critical log level
	LogCrit
	// LogErr is the error log level
	LogErr
	// LogWarning is the warning log level
	LogWarning
	// LogNotice is the notice log level
	LogNotice
	// LogInfo is the info log level
	LogInfo
	// LogDebug is the debug log level
	LogDebug
)

var (
	logger         *log.Logger
	level2str      map[int]string
	fac2str        map[int]string
	loglevel       Level = LogErr
	errBadLogLevel error = errors.New("unknown log level")
)

// Init package
func init() {

	level2str = make(map[int]string)
	level2str[int(LogErr)] = "ERROR"
	level2str[int(LogWarning)] = "WARNING"
	level2str[int(LogInfo)] = "INFO"
	level2str[int(LogDebug)] = "DEBUG"
}

//
// Misc code
//

// Init initialize the logger package
func Init(flags int) {
	var logFlags int = log.Ldate | log.Ltime
	if (flags & LogFormatLong) != 0 {
		logFlags = logFlags | log.Lmicroseconds | log.Lshortfile
	}

	switch flags & loggerMask {
	case LoggerSTDOUT:
		logger = log.New(os.Stdout, "", logFlags)
	case LoggerSTDERR:
		logger = log.New(os.Stderr, "", logFlags)
	}
}

func (l *Level) String() string {
	return level2str[int(*l)]
}

// Set set the log level from a string input
func (l *Level) Set(val string) error {
	for i, str := range level2str {
		if val == str {
			*l = Level(i)
			return nil
		}
	}

	return errBadLogLevel
}

// Get return the log level as a string
func (l *Level) Get() interface{} {
	return l.String()
}

// SetLogLevel set the log level of the logger
func SetLogLevel(level Level) {
	loglevel = level
}

func logf(level Level, s string, v ...interface{}) {
	if logger != nil {
		logger.Output(3, fmt.Sprintf(level2str[int(level)]+": "+s, v...))
	}
}

// Debugf log as an debug
func Debugf(format string, v ...interface{}) {
	if loglevel < LogDebug {
		return
	}
	logf(LogDebug, format, v...)
}

// Infof log as an info
func Infof(format string, v ...interface{}) {
	if loglevel < LogInfo {
		return
	}
	logf(LogInfo, format, v...)
}

// Warnf log as an warning
func Warnf(format string, v ...interface{}) {
	if loglevel < LogWarning {
		return
	}
	logf(LogWarning, format, v...)
}

// Errorf log as an Error
func Errorf(format string, v ...interface{}) {
	if loglevel < LogErr {
		return
	}
	logf(LogErr, format, v...)
}

// Fatalf log as an Error and then run a panic
func Fatalf(format string, v ...interface{}) {
	// Copied code from Errorf() to have correct line numbers printed...
	if loglevel >= LogErr {
		logf(LogErr, format, v...)
	}
	os.Exit(1)
}
