package logging

import (
	"loggergo/infra/property/constants"
	"loggergo/logging/core/model/loglevel"

	"github.com/spf13/viper"
)

const FILE_DEFAULT_VALUE_ENABLED = true
const FILE_DEFAULT_VALUE_LEVEL int = 1

const FILE_DEFAULT_VALUE_FILE_PATH = "@HOME/log/.log"
const FILE_DEFAULT_VALUE_MAX_FILES = 320
const FILE_DEFAULT_VALUE_FILE_SIZE_MB = 32

type FileLoggingProperties struct {
	enabled       bool
	level         loglevel.LogLevel
	filePath      string
	maxFiles      int
	maxFileSizeMb int
	compression   *FileLoggingCompressionProperties
}

func NewFileLoggingProperties() *FileLoggingProperties {
	var enabled bool
	if viper.IsSet(constants.GURMS_LOGGING_FILE_ENABLED) {
		enabled = viper.GetBool(constants.GURMS_LOGGING_FILE_ENABLED)
	} else {
		enabled = FILE_DEFAULT_VALUE_ENABLED
	}
	var level int
	if viper.IsSet(constants.GURMS_LOGGING_FILE_LEVEL) {
		level = viper.GetInt(constants.GURMS_LOGGING_FILE_LEVEL)
	} else {
		level = FILE_DEFAULT_VALUE_LEVEL
	}
	var path string
	if viper.IsSet(constants.GURMS_LOGGING_FILE_PATH) {
		path = viper.GetString(constants.GURMS_LOGGING_FILE_PATH)
	} else {
		path = FILE_DEFAULT_VALUE_FILE_PATH
	}
	var maxFiles int
	if viper.IsSet(constants.GURMS_LOGGING_FILE_MAX_FILES) {
		maxFiles = viper.GetInt(constants.GURMS_LOGGING_FILE_MAX_FILES)
	} else {
		maxFiles = FILE_DEFAULT_VALUE_MAX_FILES
	}
	var fileSizeMb int
	if viper.IsSet(constants.GURMS_LOGGING_FILE_MAX_FILE_SIZE_MB) {
		fileSizeMb = viper.GetInt(constants.GURMS_LOGGING_FILE_MAX_FILE_SIZE_MB)
	} else {
		fileSizeMb = FILE_DEFAULT_VALUE_FILE_SIZE_MB
	}
	var compression bool
	if viper.IsSet(constants.GURMS_LOGGING_FILE_COMPRESSION_ENABLED) {
		compression = viper.GetBool(constants.GURMS_LOGGING_FILE_COMPRESSION_ENABLED)
	} else {
		compression = FILE_DEFAULT_VALUE_COMPRESSION_ENABLED
	}

	return &FileLoggingProperties{
		enabled:       enabled,
		level:         loglevel.LogLevel(level),
		filePath:      path,
		maxFiles:      maxFiles,
		maxFileSizeMb: fileSizeMb,
		compression:   NewFileLoggingCompressionProperties(compression),
	}
}

func (f *FileLoggingProperties) IsEnabled() bool {
	return f.enabled
}

func (f *FileLoggingProperties) GetLevel() loglevel.LogLevel {
	return f.level
}

func (f *FileLoggingProperties) GetFilePath() string {
	return f.filePath
}

func (f *FileLoggingProperties) GetMaxFiles() int {
	return f.maxFiles
}

func (f *FileLoggingProperties) GetMaxFilesSizeMb() int {
	return f.maxFileSizeMb
}

func (f *FileLoggingProperties) GetCompression() bool {
	return f.compression.enabled
}
