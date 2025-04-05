package logrecord

import (
	"bytes"
	"loggergo/logging/core/model/loglevel"
)

type LogRecord struct {
	logger    AsyncLogger
	level     loglevel.LogLevel
	timestamp int64
	data      *bytes.Buffer
}

type AsyncLogger interface {
	IsTraceEnabled() bool
	IsDebugEnabled() bool
	IsInfoEnabled() bool
	IsWarnEnabled() bool
	IsErrorEnabled() bool
	IsFatalEnabled() bool
}

func NewLogRecord(logger AsyncLogger, level loglevel.LogLevel,
	timestamp int64, data *bytes.Buffer) LogRecord {

	if data == nil {
		panic("nil pointer in NewLogRecord")
	}
	return LogRecord{
		logger:    logger,
		level:     level,
		timestamp: timestamp,
		data:      data,
	}
}

func (l *LogRecord) Level() loglevel.LogLevel {
	return l.level
}

func (l *LogRecord) Timestamp() int64 {
	return l.timestamp
}

func (l *LogRecord) GetLogger() any {
	return l.logger
}

func (l *LogRecord) GetBuffer() *bytes.Buffer {
	return l.data
}

func (l *LogRecord) ClearData() {
	l.data = nil
}
