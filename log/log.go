package log

import (
	//"io"
    "os"
	"log"
    "fmt"
    "log/syslog"
	"errors"
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

    // From /usr/include/sys/syslog.h.
    // These are the same up to LOG_FTP on Linux, BSD, and OS X.
    LOG_KERN SyslogFacility = iota << 3
    LOG_USER
    LOG_MAIL
    LOG_DAEMON
    LOG_AUTH
    LOG_SYSLOG
    LOG_LPR
    LOG_NEWS
    LOG_UUCP
    LOG_CRON
    LOG_AUTHPRIV
    LOG_FTP

    LOG_LOCAL0
    LOG_LOCAL1
    LOG_LOCAL2
    LOG_LOCAL3
    LOG_LOCAL4
    LOG_LOCAL5
    LOG_LOCAL6
    LOG_LOCAL7
)

//
// Global var
//

type syslogMethodFunc func(*syslog.Writer,string) error

type LogLevel       syslog.Priority
type SyslogFacility syslog.Priority

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

func (f *SyslogFacility) String() string {
    return fac2str[int(*f)]
}

func (f *SyslogFacility) Set(val string) error {
    for i, str := range fac2str {
        if val == str {
            *f = SyslogFacility(i)
            return nil
        }
    }

    return BadFacility
}

var (
	logger        *log.Logger
    syslogger     *syslog.Writer

    level2str     map[int]string
    fac2str       map[int]string
    syslog_method map[int]syslogMethodFunc

    syslogTag     string
    loglevel      LogLevel       = LOG_ERR
    facility      SyslogFacility = LOG_USER

    BadLogLevel   error = errors.New("unknown log level")
    BadFacility   error = errors.New("unknown syslog facility")
)

// Init package
func init() {

    level2str = make(map[int]string)
    level2str[int(LOG_ERR)]     = "ERROR"
    level2str[int(LOG_WARNING)] = "WARNING"
    level2str[int(LOG_INFO)]    = "INFO"
    level2str[int(LOG_DEBUG)]   = "DEBUG"

    syslog_method = make(map[int]syslogMethodFunc)
    syslog_method[int(LOG_ERR)]     = (*syslog.Writer).Err
    syslog_method[int(LOG_WARNING)] = (*syslog.Writer).Warning
    syslog_method[int(LOG_INFO)]    = (*syslog.Writer).Info
    syslog_method[int(LOG_DEBUG)]   = (*syslog.Writer).Debug

    fac2str = make(map[int]string)

    fac2str[int(LOG_USER)]   = "USER"
    fac2str[int(LOG_DAEMON)] = "DAEMON"
    fac2str[int(LOG_SYSLOG)] = "SYSLOG"
    fac2str[int(LOG_LOCAL0)] = "LOCAL0"
    fac2str[int(LOG_LOCAL1)] = "LOCAL1"
    fac2str[int(LOG_LOCAL2)] = "LOCAL2"
    fac2str[int(LOG_LOCAL3)] = "LOCAL3"
    fac2str[int(LOG_LOCAL4)] = "LOCAL4"
    fac2str[int(LOG_LOCAL5)] = "LOCAL5"
    fac2str[int(LOG_LOCAL6)] = "LOCAL6"
    fac2str[int(LOG_LOCAL7)] = "LOCAL7"
}


//
// Misc code
//

func Init(flags int) error {

    var err error

	var logFlags int = log.Ldate | log.Ltime
    if (flags & LOGFORMAT_LONG) != 0 {
        logFlags = logFlags | log.Lmicroseconds | log.Lshortfile
    }

    switch flags & loggerMask {
    case LOGGER_STDOUT:
        logger = log.New(os.Stdout, "", logFlags)
    case LOGGER_STDERR:
        logger = log.New(os.Stderr, "", logFlags)
    case LOGGER_SYSLOG:
        syslogger, err = syslog.New(syslog.Priority(facility), syslogTag)
        if err != nil {
            return err
        }
    }

    return nil
}

// Must to be called before Init
func SetSyslogTag(tag string) {
    syslogTag = tag
}

func SetLogLevel(level LogLevel) {
    loglevel = level
}

func SetSyslogFacility(f SyslogFacility) {
    facility = f
}

func logf(level LogLevel, s string, v ...interface{}) {

	if logger != nil {
        logger.Output(3, fmt.Sprintf(level2str[int(level)] + ": " + s, v...))
	} else if syslogger != nil {
        syslog_method[int(level)](syslogger,
            fmt.Sprintf(level2str[int(level)] + ": " + s, v...))
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
