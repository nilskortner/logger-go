package logging

import (
	"loggergo/infra/property/constants"
	"loggergo/logging/core/model/loglevel"

	"github.com/spf13/viper"
)

const CONSOLE_DEFAULT_VALUE_ENABLED = false
const CONSOLE_DEFAULT_VALUE_LEVEL loglevel.LogLevel = 1

type ConsoleLoggingProperties struct {
	enabled bool
	level   loglevel.LogLevel
}

func NewConsoleLoggingProperties() *ConsoleLoggingProperties {
	var enabled bool
	if viper.IsSet(constants.GURMS_LOGGING_CONSOLE_ENABLED) {
		enabled = viper.GetBool(constants.GURMS_LOGGING_CONSOLE_ENABLED)
	} else {
		enabled = CONSOLE_DEFAULT_VALUE_ENABLED
	}
	var level loglevel.LogLevel
	if viper.IsSet(constants.GURMS_LOGGING_CONSOLE_LEVEL) {
		level = loglevel.LogLevel(viper.GetInt(constants.GURMS_LOGGING_CONSOLE_LEVEL))
	} else {
		level = CONSOLE_DEFAULT_VALUE_LEVEL
	}

	return &ConsoleLoggingProperties{
		enabled: enabled,
		level:   level,
	}
}

func (c *ConsoleLoggingProperties) IsEnabled() bool {
	return c.enabled
}

func (c *ConsoleLoggingProperties) Level() loglevel.LogLevel {
	return c.level
}
