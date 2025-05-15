package log

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/valyala/bytebufferpool"
)

// StructuredFormatter is an interface for formatting log entries in a structured format
type StructuredFormatter interface {
	// Format formats a log entry into a structured format
	Format(level Level, message string, fields map[string]interface{}) ([]byte, error)
}

// JSONFormatter implements StructuredFormatter for JSON output
type JSONFormatter struct {
	// TimeKey is the key for the time field in the JSON output
	TimeKey string
	// LevelKey is the key for the level field in the JSON output
	LevelKey string
	// MessageKey is the key for the message field in the JSON output
	MessageKey string
	// TimeFormat is the format for the time field in the JSON output
	TimeFormat string
}

// NewJSONFormatter creates a new JSONFormatter with default settings
func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{
		TimeKey:    "time",
		LevelKey:   "level",
		MessageKey: "message",
		TimeFormat: time.RFC3339,
	}
}

// Format formats a log entry into JSON
func (f *JSONFormatter) Format(level Level, message string, fields map[string]interface{}) ([]byte, error) {
	data := make(map[string]interface{}, len(fields)+3)

	// Add standard fields
	data[f.TimeKey] = time.Now().Format(f.TimeFormat)
	data[f.LevelKey] = level.toString()
	data[f.MessageKey] = message

	// Add custom fields
	for k, v := range fields {
		switch v := v.(type) {
		case error:
			// Handle errors specially
			data[k] = v.Error()
		default:
			data[k] = v
		}
	}

	return json.Marshal(data)
}

// StructuredLogger is a logger that supports structured logging formats
type StructuredLogger struct {
	logger    *defaultLogger
	formatter StructuredFormatter
	enabled   bool
}

// NewStructuredLogger creates a new structured logger with the specified formatter
func NewStructuredLogger(logger *defaultLogger, formatter StructuredFormatter) *StructuredLogger {
	return &StructuredLogger{
		logger:    logger,
		formatter: formatter,
		enabled:   true,
	}
}

// EnableStructuredLogging enables or disables structured logging
func (l *StructuredLogger) EnableStructuredLogging(enabled bool) {
	l.enabled = enabled
}

// IsStructuredLoggingEnabled returns whether structured logging is enabled
func (l *StructuredLogger) IsStructuredLoggingEnabled() bool {
	return l.enabled
}

// SetFormatter sets the formatter for structured logging
func (l *StructuredLogger) SetFormatter(formatter StructuredFormatter) {
	l.formatter = formatter
}

// logStructured logs a message with fields in a structured format
func (l *StructuredLogger) logStructured(level Level, message string, fields map[string]interface{}) {
	if l.logger.level > level || !l.enabled || l.formatter == nil {
		// If structured logging is disabled or the level is filtered, use the default logger
		if l.logger.level <= level {
			l.logger.privateLogw(level, message, mapToKeyValuePairs(fields))
		}
		return
	}

	// Format the log entry
	data, err := l.formatter.Format(level, message, fields)
	if err != nil {
		// If formatting fails, fall back to the default logger
		l.logger.privateLogw(level, message, mapToKeyValuePairs(fields))
		l.logger.privateLogw(LevelError, "Failed to format structured log entry", []interface{}{"error", err})
		return
	}

	// Write the formatted log entry
	buf := bytebufferpool.Get()
	defer func() {
		buf.Reset()
		bytebufferpool.Put(buf)
	}()

	buf.Write(data)
	buf.WriteString("\n")

	_ = l.logger.stdLog.Output(l.logger.depth, buf.String())

	if level == LevelFatal {
		// Exit if the level is fatal
		fmt.Fprintln(io.MultiWriter(l.logger.stdLog.Writer(), buf), "exiting...")
		buf.Reset()
		bytebufferpool.Put(buf)
		l.logger.stdLog.Writer().(io.Closer).Close()
		os.Exit(1)
	}
}

// mapToKeyValuePairs converts a map to a slice of key-value pairs
func mapToKeyValuePairs(m map[string]interface{}) []interface{} {
	result := make([]interface{}, 0, len(m)*2)
	for k, v := range m {
		result = append(result, k, v)
	}
	return result
}

// parseKeyValuePairs converts a slice of key-value pairs to a map
func parseKeyValuePairs(keyValuePairs []interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	// Check for unpaired key-value
	if (len(keyValuePairs) & 1) == 1 {
		keyValuePairs = append(keyValuePairs, "KEYVALS UNPAIRED")
	}

	for i := 0; i < len(keyValuePairs); i += 2 {
		key, ok := keyValuePairs[i].(string)
		if !ok {
			key = fmt.Sprintf("%v", keyValuePairs[i])
		}
		result[key] = keyValuePairs[i+1]
	}

	return result
}
