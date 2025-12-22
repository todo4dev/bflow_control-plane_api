package logger

import (
	adapter "src/application/adapter/logger"
	"src/core/di"
	"src/core/env"
	json_impl "src/infrastructure/logger/json"
)

func init() {
	di.RegisterAs[adapter.ILoggerAdapter](func() adapter.ILoggerAdapter {
		config := adapter.LoggerConfig{
			Level: adapter.LoggerLevelEnum(env.Get("LOGGER_LEVEL", "info")),
		}
		if err := config.Validate(); err != nil {
			panic(err)
		}

		return json_impl.NewJSONLoggerAdapter(&config)
	})
}
