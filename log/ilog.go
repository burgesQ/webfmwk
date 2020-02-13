package log

// ILog interface implement the logging system inside the API
type ILog interface {
	Errorf(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Debugf(format string, v ...interface{})
	Fatalf(format string, v ...interface{})
	SetLogLevel(level int) bool
}
