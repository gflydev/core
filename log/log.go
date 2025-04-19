package log

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
)

var logger AllLogger = &defaultLogger{
	stdLog: log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile|log.Lmicroseconds),
	depth:  4,
}

// Logger is a logger interface that provides logging function with levels.
type Logger interface {
	// Trace logs a message at the Trace level.
	// Parameters:
	//   - v: The variables or objects to log.
	Trace(v ...interface{})
	// Debug logs a message at the Debug level.
	// Parameters:
	//   - v: The variables or objects to log.
	Debug(v ...interface{})
	// Info logs a message at the Info level.
	// Parameters:
	//   - v: The variables or objects to log.
	Info(v ...interface{})
	// Warn logs a message at the Warning level.
	// Parameters:
	//   - v: The variables or objects to log.
	Warn(v ...interface{})
	// Error logs a message at the Error level.
	// Parameters:
	//   - v: The variables or objects to log.
	Error(v ...interface{})
	// Fatal logs a message at the Fatal level and exits the application.
	// Parameters:
	//   - v: The variables or objects to log.
	Fatal(v ...interface{})
	// Panic logs a message at the Panic level.
	// Note: Despite the name, this method does not actually panic.
	// Parameters:
	//   - v: The variables or objects to log.
	Panic(v ...interface{})
}

// FormatLogger is a logger interface that outputs logs with a format.
type FormatLogger interface {
	// Tracef logs a formatted message at the Trace level.
	// Parameters:
	//   - format: The format string.
	//   - v: The variables or objects to format and log.
	Tracef(format string, v ...interface{})
	// Debugf logs a formatted message at the Debug level.
	// Parameters:
	//   - format: The format string.
	//   - v: The variables or objects to format and log.
	Debugf(format string, v ...interface{})
	// Infof logs a formatted message at the Info level.
	// Parameters:
	//   - format: The format string.
	//   - v: The variables or objects to format and log.
	Infof(format string, v ...interface{})
	// Warnf logs a formatted message at the Warning level.
	// Parameters:
	//   - format: The format string.
	//   - v: The variables or objects to format and log.
	Warnf(format string, v ...interface{})
	// Errorf logs a formatted message at the Error level.
	// Parameters:
	//   - format: The format string.
	//   - v: The variables or objects to format and log.
	Errorf(format string, v ...interface{})
	// Fatalf logs a formatted message at the Fatal level and exits the application.
	// Parameters:
	//   - format: The format string.
	//   - v: The variables or objects to format and log.
	Fatalf(format string, v ...interface{})
	// Panicf logs a formatted message at the Panic level.
	// Note: Despite the name, this method does not actually panic.
	// Parameters:
	//   - format: The format string.
	//   - v: The variables or objects to format and log.
	Panicf(format string, v ...interface{})
}

// WithLogger is a logger interface that outputs logs with a message and key-value pairs.
type WithLogger interface {
	// Tracew logs a message at the Trace level with key-value pairs.
	// Parameters:
	//   - msg: The message string.
	//   - keysAndValues: The key-value pairs to log.
	Tracew(msg string, keysAndValues ...interface{})
	// Debugw logs a message at the Debug level with key-value pairs.
	// Parameters:
	//   - msg: The message string.
	//   - keysAndValues: The key-value pairs to log.
	Debugw(msg string, keysAndValues ...interface{})
	// Infow logs a message at the Info level with key-value pairs.
	// Parameters:
	//   - msg: The message string.
	//   - keysAndValues: The key-value pairs to log.
	Infow(msg string, keysAndValues ...interface{})
	// Warnw logs a message at the Warning level with key-value pairs.
	// Parameters:
	//   - msg: The message string.
	//   - keysAndValues: The key-value pairs to log.
	Warnw(msg string, keysAndValues ...interface{})
	// Errorw logs a message at the Error level with key-value pairs.
	// Parameters:
	//   - msg: The message string.
	//   - keysAndValues: The key-value pairs to log.
	Errorw(msg string, keysAndValues ...interface{})
	// Fatalw logs a message at the Fatal level with key-value pairs and exits the application.
	// Parameters:
	//   - msg: The message string.
	//   - keysAndValues: The key-value pairs to log.
	Fatalw(msg string, keysAndValues ...interface{})
	// Panicw logs a message at the Panic level with key-value pairs.
	// Note: Despite the name, this method does not actually panic.
	// Parameters:
	//   - msg: The message string.
	//   - keysAndValues: The key-value pairs to log.
	Panicw(msg string, keysAndValues ...interface{})
}

type CommonLogger interface {
	Logger
	FormatLogger
	WithLogger
}

// ControlLogger provides methods to configure a logger.
type ControlLogger interface {
	// SetLevel sets the logging level for the logger.
	// Parameters:
	//   - level: The logging level to set.
	SetLevel(Level)
	// SetOutput sets the output destination for the logger.
	// Parameters:
	//   - writer: The io.Writer to write logs to.
	SetOutput(io.Writer)
}

// AllLogger is the combination of Logger, FormatLogger, CtxLogger, and ControlLogger.
// Custom extensions can be made through AllLogger.
type AllLogger interface {
	CommonLogger
	ControlLogger
	// WithContext returns a logger that is associated with a given context.
	// Parameters:
	//   - ctx: The context to associate with the logger.
	WithContext(ctx context.Context) CommonLogger
}

// Level defines the priority of a log message.
// When a logger is configured with a level, any log message with a lower
// log level (smaller by integer comparison) will not be output.
type Level int

// The levels of logs.
const (
	LevelTrace Level = iota // Trace represents very detailed debug information.
	LevelDebug              // Debug represents general debugging information.
	LevelInfo               // Info represents informational messages.
	LevelWarn               // Warn represents warning messages.
	LevelError              // Error represents error messages.
	LevelFatal              // Fatal represents critical error messages and causes the program to exit.
	LevelPanic              // Panic represents critical error messages (but does not actually cause a panic).
)

var levels = []string{
	"[TRACE] ", // Trace level string representation.
	"[DEBUG] ", // Debug level string representation.
	"[INFO] ",  // Info level string representation.
	"[WARN] ",  // Warn level string representation.
	"[ERROR] ", // Error level string representation.
	"[FATAL] ", // Fatal level string representation.
	"[PANIC] ", // Panic level string representation.
}

// toString converts a Level to its string representation.
// Returns:
//   - A string representation of the level. If the level is not valid, it returns a formatted unknown level string.
func (lv Level) toString() string {
	if lv >= LevelTrace && lv <= LevelPanic {
		return levels[lv]
	}
	return fmt.Sprintf("[?%d] ", lv)
}
