package logger

import (
	"bytes"
	"gurms/internal/infra/logging/core/appender"
	"gurms/internal/infra/logging/core/layout"
	"gurms/internal/infra/logging/core/model/loglevel"
	"gurms/internal/infra/logging/core/model/logrecord"
	mpsc "gurms/internal/supportpkgs/datastructures/mpscunboundedarrayqueue"
	"gurms/internal/supportpkgs/mathsupport"
	"math"
	"time"
)

type AsyncLogger struct {
	name        string
	shouldParse bool
	appenders   []appender.Appender
	layout      *layout.GurmsTemplateLayout
	queue       *mpsc.MpscUnboundedArrayQueue[logrecord.LogRecord]
	nameForLog  []byte
	level       int
}

func NewAsyncLogger(
	name string,
	shouldParse bool,
	appenders []appender.Appender,
	layoutAL *layout.GurmsTemplateLayout,
	queue *mpsc.MpscUnboundedArrayQueue[logrecord.LogRecord]) *AsyncLogger {
	nameForLog := layout.FormatStructName(name)

	var level int
	if len(appenders) == 0 {
		level = math.MaxInt
	} else {
		level = -1
		for _, appender := range appenders {
			level = mathsupport.Max(level, int(appender.GetLevel()))
		}
	}

	return &AsyncLogger{
		name:        name,
		shouldParse: shouldParse,
		appenders:   appenders,
		layout:      layoutAL,
		queue:       queue,
		nameForLog:  nameForLog,
		level:       level,
	}
}

func (a *AsyncLogger) GetAppenders() []appender.Appender {
	return a.appenders
}

func (a *AsyncLogger) IsTraceEnabled() bool {
	return a.level <= int(loglevel.TRACE)
}

func (a *AsyncLogger) IsDebugEnabled() bool {
	return a.level <= int(loglevel.DEBUG)
}

func (a *AsyncLogger) IsInfoEnabled() bool {
	return a.level <= int(loglevel.INFO)
}

func (a *AsyncLogger) IsWarnEnabled() bool {
	return a.level <= int(loglevel.WARN)
}

func (a *AsyncLogger) IsErrorEnabled() bool {
	return a.level <= int(loglevel.ERROR)
}

func (a *AsyncLogger) IsFatalEnabled() bool {
	return a.level <= int(loglevel.FATAL)
}

func (a *AsyncLogger) IsEnabled(loglevel loglevel.LogLevel) bool {
	return a.level <= int(loglevel)
}

func (a *AsyncLogger) Log(level loglevel.LogLevel, message string) {
	if !a.IsEnabled(level) {
		return
	}
	a.doLog(level, message, nil, nil)
}

func (a *AsyncLogger) LogWithArguments(level loglevel.LogLevel, message string, arguments ...interface{}) {
	if !a.IsEnabled(level) {
		return
	}
	a.doLog(level, message, arguments, nil)
}

func (a *AsyncLogger) LogWithError(level loglevel.LogLevel, message string, err error) {
	if !a.IsEnabled(level) {
		return
	}
	a.doLog(level, message, nil, err)
}

func (a *AsyncLogger) Debug(message string, args ...interface{}) {
	a.LogWithArguments(loglevel.DEBUG, message, args)
}

func (a *AsyncLogger) InfoWithArgs(message string, args ...interface{}) {
	a.LogWithArguments(loglevel.INFO, message, args)
}

func (a *AsyncLogger) Info(message *bytes.Buffer) {
	a.doLogBasic(loglevel.INFO, message)
}

func (a *AsyncLogger) Warn(message string) {
	a.doLog(loglevel.WARN, message, nil, nil)
}

func (a *AsyncLogger) WarnWithArgs(message string, args ...interface{}) {
	a.LogWithArguments(loglevel.WARN, message, args)
}

func (a *AsyncLogger) Error(err error) {
	a.doLog(loglevel.ERROR, "", nil, err)
}

func (a *AsyncLogger) ErrorWithMessage(message string, err error) {
	a.doLog(loglevel.ERROR, message, nil, err)
}

func (a *AsyncLogger) ErrorWithArgs(message string, args ...interface{}) {
	a.LogWithArguments(loglevel.ERROR, message, args)
}

func (a *AsyncLogger) ErrorWithBuffer(message *bytes.Buffer) {
	a.doLogBasic(loglevel.ERROR, message)
}

func (a *AsyncLogger) Fatal(message string) {
	a.doLog(loglevel.FATAL, message, nil, nil)
}

func (a *AsyncLogger) FatalWithArgs(message string, args ...interface{}) {
	a.LogWithArguments(loglevel.FATAL, message, args)
}

func (a *AsyncLogger) FatalWithError(message string, err error) {
	a.doLog(loglevel.FATAL, message, nil, err)
}

func (a *AsyncLogger) doLog(level loglevel.LogLevel, message string, args []interface{}, err error) {
	buffer := layout.Format(a.layout, a.shouldParse, a.nameForLog, level, message, args, err)

	a.queue.Offer(logrecord.NewLogRecord(a, level, time.Now().UnixMilli(), buffer))
}

func (a *AsyncLogger) doLogBasic(level loglevel.LogLevel, message *bytes.Buffer) {
	buffer := layout.FormatBasic(a.layout, a.nameForLog, level, message)

	a.queue.Offer(logrecord.NewLogRecord(a, level, time.Now().UnixMilli(), buffer))
}
