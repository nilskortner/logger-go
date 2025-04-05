package logger

import (
	"bytes"
	model "gurms/internal/infra/logging/core/model/loglevel"
)

type Logger interface {
	IsTraceEnabled() bool

	IsDebugEnabled() bool

	IsInfoEnabled() bool

	IsWarnEnabled() bool

	IsErrorEnabled() bool

	IsFatalEnabled() bool

	IsEnabled(model.LogLevel) bool

	Log(level model.LogLevel, message string)
	LogWithArguments(level model.LogLevel, format string, args ...interface{})
	LogWithError(level model.LogLevel, message string, err error)

	Debug(message string, args ...interface{})
	InfoWithArgs(message string, args ...interface{})
	Info(data *bytes.Buffer)
	Warn(message string)
	WarnWithArgs(message string, args ...interface{})
	Error(err error)
	ErrorWithMessage(message string, err error)
	ErrorWithArgs(message string, args ...interface{})
	ErrorWithBuffer(message *bytes.Buffer)
	Fatal(message string)
	FatalWithArgs(message string, args ...interface{})
	FatalWithError(message string, err error)
}
