package logger

import "gurms/internal/infra/logging/core/model/loglevel"

type LoggerOptions struct {
	loggerName  string
	level       loglevel.LogLevel
	filePath    string
	shouldParse bool
}

func NewLoggerOptions(name string) *LoggerOptions {
	return &LoggerOptions{
		loggerName:  name,
		shouldParse: true,
	}
}

func (o *LoggerOptions) GetName() string {
	return o.loggerName
}

func (o *LoggerOptions) GetPath() string {
	return o.filePath
}

func (o *LoggerOptions) GetLevel() loglevel.LogLevel {
	return o.level
}

func (o *LoggerOptions) IsShouldParse() bool {
	return o.shouldParse
}
