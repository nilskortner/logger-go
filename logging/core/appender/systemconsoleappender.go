package appender

import (
	"fmt"
	"loggergo/logging/core/model/loglevel"
	"loggergo/logging/core/model/logrecord"
	"os"
)

type SystemConsoleAppender struct {
	appender *BaseAppender
}

func NewSystemConsoleAppender(level loglevel.LogLevel) *SystemConsoleAppender {
	return &SystemConsoleAppender{
		appender: NewAppender(level),
	}
}

func (s *SystemConsoleAppender) GetLevel() loglevel.LogLevel {
	return s.appender.GetLevel()
}

func (s *SystemConsoleAppender) Append(record logrecord.LogRecord) int {
	if !record.Level().IsLoggable(s.appender.level) {
		return 0
	}
	str := record.GetBuffer().String()

	if record.Level().IsErrorOrFatal() {
		fmt.Fprintln(os.Stderr, str)
	} else {
		fmt.Println(str)
	}

	return record.GetBuffer().Len()

}
