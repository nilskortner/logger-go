package logging

type LoggingProperties struct {
	console *ConsoleLoggingProperties
	file    *FileLoggingProperties
}

func NewLoggingProperties() *LoggingProperties {
	return &LoggingProperties{
		console: NewConsoleLoggingProperties(),
		file:    NewFileLoggingProperties(),
	}
}

func (properties *LoggingProperties) GetConsole() *ConsoleLoggingProperties {
	return properties.console
}

func (properties *LoggingProperties) GetFile() *FileLoggingProperties {
	return properties.file
}
