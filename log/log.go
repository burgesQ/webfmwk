package log

import (
	"errors"
	"fmt"
	"log"
	"log/syslog"
	"os"
)

const (
	LOGGER_STDOUT = iota
	LOGGER_STDERR
	LOGGER_SYSLOG

	loggerMask = 3

	LOGFORMAT_SHORT = 0
	LOGFORMAT_LONG  = 1 << 3

	// From /usr/include/sys/syslog.h.
	// These are the same on Linux, BSD, and OS X.
	LOG_EMERG LogLevel = iota
	LOG_ALERT
	LOG_CRIT
	LOG_ERR
	LOG_WARNING
	LOG_NOTICE
	LOG_INFO
	LOG_DEBUG
)

//
// Global var
//

type LogLevel syslog.Priority

var (
	logger      *log.Logger
	level2str   map[int]string
	fac2str     map[int]string
	loglevel    LogLevel = LOG_ERR
	BadLogLevel error    = errors.New("unknown log level")
)

// Init package
func init() {

	level2str = make(map[int]string)
	level2str[int(LOG_ERR)] = "ERROR"
	level2str[int(LOG_WARNING)] = "WARNING"
	level2str[int(LOG_INFO)] = "INFO"
	level2str[int(LOG_DEBUG)] = "DEBUG"
}

//
// Misc code
//

func Init(flags int) {
	var logFlags int = log.Ldate | log.Ltime
	if (flags & LOGFORMAT_LONG) != 0 {
		logFlags = logFlags | log.Lmicroseconds | log.Lshortfile
	}

	switch flags & loggerMask {
	case LOGGER_STDOUT:
		logger = log.New(os.Stdout, "", logFlags)
	case LOGGER_STDERR:
		logger = log.New(os.Stderr, "", logFlags)
	}
}

func (l *LogLevel) String() string {
	return level2str[int(*l)]
}

func (l *LogLevel) Set(val string) error {
	for i, str := range level2str {
		if val == str {
			*l = LogLevel(i)
			return nil
		}
	}

	return BadLogLevel
}

func (l *LogLevel) Get() interface{} {
	return l.String()
}

func SetLogLevel(level LogLevel) {
	loglevel = level
}

func logf(level LogLevel, s string, v ...interface{}) {
	if logger != nil {
		logger.Output(3, fmt.Sprintf(level2str[int(level)]+": "+s, v...))
	}
}

func Debugf(format string, v ...interface{}) {
	if loglevel < LOG_DEBUG {
		return
	}
	logf(LOG_DEBUG, format, v...)
}

func Infof(format string, v ...interface{}) {
	if loglevel < LOG_INFO {
		return
	}
	logf(LOG_INFO, format, v...)
}

func Warnf(format string, v ...interface{}) {
	if loglevel < LOG_WARNING {
		return
	}
	logf(LOG_WARNING, format, v...)
}

func Errorf(format string, v ...interface{}) {
	if loglevel < LOG_ERR {
		return
	}
	logf(LOG_ERR, format, v...)
}

func Fatalf(format string, v ...interface{}) {
	// Copied code from Errorf() to have correct line numbers printed...
	if loglevel >= LOG_ERR {
		logf(LOG_ERR, format, v...)
	}
	os.Exit(1)
}
