package logger

// ILoggerAdapter is the interface that wraps the basic logging methods.
//
// Implementations of this interface should provide mechanisms to log
// messages at different severity levels (Debug, Info, Warn, Error, Fatal)
// and allow for structured logging by accepting a map of parameters.
type ILoggerAdapter interface {
	// Set a default value to logger
	Set(key string, value any)
	// Debug logs a message at Debug level.
	Debug(message string, optionalMap ...map[string]any)
	// Info logs a message at Info level.
	Info(message string, optionalMap ...map[string]any)
	// Warn logs a message at Warn level.
	Warn(message string, optionalMap ...map[string]any)
	// Error logs a message at Error level.
	Error(message string, optionalMap ...map[string]any)
	// Fatal logs a message at Fatal level, then exits the application.
	Fatal(message string, optionalMap ...map[string]any)
}
