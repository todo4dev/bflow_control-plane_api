package interceptor

import (
	"src/application/adapter/logger"
	"src/core/di"
	"src/presentation/api/rest/core"
	"time"
)

func LoggingInterceptor() core.InterceptorFN {
	logger := di.Resolve[logger.ILoggerAdapter]()
	return func(ctx core.HttpContext, next core.HandlerFN) error {
		start := time.Now()
		err := next(ctx)
		optionalMap := map[string]any{"ms": time.Since(start).String()}
		if err != nil {
			optionalMap["err"] = err.Error()
		}
		logger.Info("Request completed", optionalMap)
		return err
	}
}
