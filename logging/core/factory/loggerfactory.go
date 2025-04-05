package factory

import (
	"loggergo/logging/core/appender"
	"loggergo/logging/core/appender/file"
	"loggergo/logging/core/layout"
	"loggergo/logging/core/logger"
	"loggergo/logging/core/model/logrecord"
	"loggergo/node/nodetype"
	"loggergo/property/env/common/logging"
	"loggergo/system"
	"strings"
	"sync"

	"loggergo/datastructures/copyonwriteslice"
	"loggergo/datastructures/linkedlist"
	"loggergo/mpscunboundedarrayqueue"
)

const (
	PROPERTY_NAME_TURMS_AI_SERVING_HOME = "GURMS_AI_SERVING_HOME"
	PROPERTY_NAME_TURMS_GATEWAY_HOME    = "GURMS_GATEWAY_HOME"
	PROPERTY_NAME_TURMS_SERVICE_HOME    = "GURMS_SERVICE_HOME"
	SERVER_TYPE_UNKNOWN                 = "unknown"
)

var once sync.Once

var loggerlayout *layout.GurmsTemplateLayout

var initialized bool

var ALL_APPENDERS = copyonwriteslice.NewCopyOnWriteSlice[appender.Appender]()
var DEFAULT_APPENDERS = make([]appender.Appender, 0, 2)
var Queue = mpscunboundedarrayqueue.NewMpscUnboundedQueue[logrecord.LogRecord](1024)
var UNINITIALIZED_LOGGERS linkedlist.LinkedList

var homeDir string
var serverTypeName string
var fileLoggingProperties *logging.FileLoggingProperties
var defaultConsoleAppender appender.Appender

var logprocessor LogProcessor

func Loggerfactory(runWithTests bool,
	nodeId string,
	nodeType nodetype.NodeType,
	properties *logging.LoggingProperties) {
	once.Do(func() {
		initialize(runWithTests, nodeId, nodeType, properties)
	})
}

func WaitClose(timeoutMillis int64) {
	logprocessor.waitClose(timeoutMillis)
}

func GetLogger(name string) logger.Logger {
	options := logger.NewLoggerOptions(name)
	return getLogger(options)
}

func initialize(
	runWithTests bool,
	nodeId string,
	nodeType nodetype.NodeType,
	properties *logging.LoggingProperties) {
	switch nodeType {
	case 0:
		homeDir = system.GetProperty("PROPERTY_NAME_GURMS_AI_SERVING_HOME")
	case 1:
		homeDir = system.GetProperty("PROPERTY_NAME_GURMS_GATEWAY_HOME")
	case 2:
		homeDir = system.GetProperty("PROPERTY_NAME_GURMS_SERVICE_HOME")
	}
	if homeDir == "" {
		homeDir = "."
	}
	serverTypeName = nodeType.GetId()
	consoleLoggingProperties := properties.GetConsole()
	fileLoggingProperties = properties.GetFile()
	if consoleLoggingProperties.IsEnabled() {
		var consoleAppender appender.Appender
		if runWithTests {
			consoleAppender = appender.NewSystemConsoleAppender(consoleLoggingProperties.Level())
		} else {
			consoleAppender = appender.NewChannelConsoleAppender(consoleLoggingProperties.Level())
		}
		defaultConsoleAppender = consoleAppender
		DEFAULT_APPENDERS = append(DEFAULT_APPENDERS, consoleAppender)
	}
	if fileLoggingProperties.IsEnabled() {
		fileAppender := file.NewRollingFileAppender(
			fileLoggingProperties.GetLevel(),
			getFilePath(fileLoggingProperties.GetFilePath()),
			fileLoggingProperties.GetMaxFiles(),
			int64(fileLoggingProperties.GetMaxFilesSizeMb()),
			fileLoggingProperties.GetCompression(),
		)
		DEFAULT_APPENDERS = append(DEFAULT_APPENDERS, fileAppender)
	}

	loggerlayout = layout.NewGurmsTemplateLayout(nodeType, nodeId)
	initialized = true

	processor := NewLogProcessor(Queue)
	processor.Start()
}

func getFilePath(path string) string {
	if path == "" {
		return "."
	}
	path = strings.Replace(path, "@HOME", homeDir, -1)
	path = strings.Replace(path, "@SERVICE_TYPE_NAME", serverTypeName, -1)
	return path
}

func getLogger(options *logger.LoggerOptions) logger.Logger {
	loggerName := options.GetName()
	filePath := options.GetPath()
	appenders := make([]appender.Appender, 2)
	if filePath != "" {
		filePath = getFilePath(filePath)
		level := options.GetLevel()
		if level == -1 {
			level = fileLoggingProperties.GetLevel()
		}
		appender := file.NewRollingFileAppender(
			level,
			filePath,
			fileLoggingProperties.GetMaxFiles(),
			int64(fileLoggingProperties.GetMaxFilesSizeMb()),
			fileLoggingProperties.GetCompression())
		appenders = append(appenders, appender)
		ALL_APPENDERS.Add(appender)
		if defaultConsoleAppender != nil {
			appenders = append(appenders, defaultConsoleAppender)
		}
	} else {
		appenders = DEFAULT_APPENDERS
	}
	return logger.NewAsyncLogger(loggerName, options.IsShouldParse(), appenders, loggerlayout, Queue)
}

func IsInitialized() bool {
	return initialized
}
