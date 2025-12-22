package logger

import "src/core/validator"

type LoggerLevelEnum string

const (
	LoggerLevel_Debug LoggerLevelEnum = "debug"
	LoggerLevel_Info  LoggerLevelEnum = "info"
	LoggerLevel_Warn  LoggerLevelEnum = "warn"
	LoggerLevel_Error LoggerLevelEnum = "error"
)

type LoggerConfig struct {
	Level LoggerLevelEnum
}

var _ validator.IValidable = (*LoggerConfig)(nil)

func (c *LoggerConfig) Validate() error {
	return validator.Object(c,
		validator.String(&c.Level).
			Required().
			Allow(string(LoggerLevel_Debug), string(LoggerLevel_Info), string(LoggerLevel_Warn), string(LoggerLevel_Error)).
			Default(string(LoggerLevel_Info)),
	).Validate()
}
