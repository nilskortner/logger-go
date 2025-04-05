package loglevel

type LogLevel int

const (
	TRACE LogLevel = iota - 1
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
	PANIC
)

func (level LogLevel) IsLoggable(enabledLevel LogLevel) bool {
	return enabledLevel <= level
}

func (level LogLevel) IsErrorOrFatal() bool {
	return level == ERROR || level == FATAL || level == PANIC
}

func (level LogLevel) String() string {
	switch level {
	case TRACE:
		return "TRACE"
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	case PANIC:
		return "PANIC"
	default:
		return "UNKNOWN"
	}
}
