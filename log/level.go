package log

// Level hold the logging level
type Level uint8

// Following contant hold the possible log levels
const (
	LogErr Level = iota
	LogWarning
	LogInfo
	LogDebug
	LogPrint
)

var (
	_l2s = [...]string{
		LogErr:     "error",
		LogWarning: "warning",
		LogInfo:    "info",
		LogDebug:   "debug",
		LogPrint:   "print",
	}
	_out = map[Level]string{
		LogErr:     "! ERR  : ",
		LogWarning: "* WARN : ",
		LogInfo:    "+ INFO : ",
		LogDebug:   "- DBG  : ",
		LogPrint:   "",
	}
)

func (l Level) String() string {
	return _l2s[l]
}

// SetLogLevel set the global logger log level
func SetLogLevel(level Level) (ok bool) {
	if level >= LogErr && level <= LogDebug {
		_lg.level = level
		ok = true
	}

	return
}
