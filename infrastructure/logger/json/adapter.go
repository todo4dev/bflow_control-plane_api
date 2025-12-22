package json

import (
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"os"
	"time"

	adapter "src/application/adapter/logger"
)

// internal level for comparison
type logLevelEnum int

const (
	logLevel_Debug logLevelEnum = iota
	logLevel_Info
	logLevel_Warn
	logLevel_Error
	logLevel_Fatal
)

// JSONLoggerAdapter implements the logger.Logger interface for JSON output.
type JSONLoggerAdapter struct {
	output io.Writer
	extra  map[string]any
	level  logLevelEnum
}

var _ adapter.ILoggerAdapter = (*JSONLoggerAdapter)(nil)

// NewJSONLoggerAdapter creates a new JSON logger adapter that writes to os.Stdout.
func NewJSONLoggerAdapter(config *adapter.LoggerConfig) *JSONLoggerAdapter {
	if config == nil {
		config = &adapter.LoggerConfig{}
	}

	var level logLevelEnum
	switch config.Level {
	case adapter.LoggerLevel_Debug:
		level = logLevel_Debug
	case adapter.LoggerLevel_Info:
		level = logLevel_Info
	case adapter.LoggerLevel_Warn:
		level = logLevel_Warn
	case adapter.LoggerLevel_Error:
		level = logLevel_Error
	default:
		level = logLevel_Info
	}

	return &JSONLoggerAdapter{
		output: os.Stdout,
		extra:  make(map[string]any),
		level:  level,
	}
}

// logEntry represents a single log entry structure for JSON marshaling.
type logEntry struct {
	Timestamp string         `json:"timestamp"`
	Level     string         `json:"level"`
	Message   string         `json:"message"`
	Extra     map[string]any `json:"extra,omitempty"`
}

// Set adds a key-value pair to the logger's default fields.
func (a *JSONLoggerAdapter) Set(key string, value any) {
	a.extra[key] = value
}

// log handles the common logic for logging messages at different levels.
func (a *JSONLoggerAdapter) log(level logLevelEnum, levelLabel, message string, optionalMap ...map[string]any) {
	// respects the configured log level
	if level < a.level {
		return
	}

	extraEntry := make(map[string]any)
	maps.Copy(extraEntry, a.extra)
	for _, m := range optionalMap {
		maps.Copy(extraEntry, m)
	}

	entry := logEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     levelLabel,
		Message:   message,
	}

	if len(extraEntry) > 0 {
		entry.Extra = extraEntry
	}

	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		// Fallback to plain text error if JSON marshaling fails
		fmt.Fprintf(a.output, `{"timestamp":"%s","level":"ERROR","message":"Failed to marshal log entry: %v - Original message: %s"}`+"\n",
			time.Now().Format(time.RFC3339), err, message)
		return
	}
	fmt.Fprintln(a.output, string(jsonBytes))
}

// Debug logs a message at Debug level.
func (a *JSONLoggerAdapter) Debug(msg string, optionalMap ...map[string]any) {
	a.log(logLevel_Debug, "DEBUG", msg, optionalMap...)
}

// Info logs a message at Info level.
func (a *JSONLoggerAdapter) Info(msg string, optionalMap ...map[string]any) {
	a.log(logLevel_Info, "INFO", msg, optionalMap...)
}

// Warn logs a message at Warn level.
func (a *JSONLoggerAdapter) Warn(msg string, optionalMap ...map[string]any) {
	a.log(logLevel_Warn, "WARN", msg, optionalMap...)
}

// Error logs a message at Error level.
func (a *JSONLoggerAdapter) Error(msg string, optionalMap ...map[string]any) {
	a.log(logLevel_Error, "ERROR", msg, optionalMap...)
}

// Fatal logs a message at Fatal level and then exits the application.
func (a *JSONLoggerAdapter) Fatal(msg string, optionalMap ...map[string]any) {
	a.log(logLevel_Fatal, "FATAL", msg, optionalMap...)
	os.Exit(1)
}
