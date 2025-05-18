package log

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/valyala/bytebufferpool"
)

var _ AllLogger = (*defaultLogger)(nil)

type defaultLogger struct {
	stdLog *log.Logger // The underlying standard logger used for output
	level  Level       // The current log level
	depth  int         // The call depth for logging
}

// privateLog logs a message at a given level using the default logger. It uses a buffer pool to optimize memory usage.
// Parameters:
//   - lv: The logging level of the message.
//   - fmtArgs: The arguments to be logged.
//
// Returns: None
func (l *defaultLogger) privateLog(lv Level, fmtArgs []any) {
	if l.level > lv {
		return
	}
	level := lv.toString()
	buf := bytebufferpool.Get()        // Borrow a buffer from the bytebufferpool
	_, _ = buf.WriteString(level)      // It is fine to ignore the error
	_, _ = fmt.Fprint(buf, fmtArgs...) // It is fine to ignore the error

	_ = l.stdLog.Output(l.depth, buf.String()) // Use the standard logger to output the constructed message
	buf.Reset()                                // Reset the buffer for reuse
	bytebufferpool.Put(buf)                    // Return the buffer to the pool
	if lv == LevelFatal {
		os.Exit(1) // Terminates the program if the level is fatal
	}
}

// privateLogf logs a formatted message at a given level using the default logger. It supports format specifiers.
// Parameters:
//   - lv: The logging level of the message.
//   - format: The format string for the message.
//   - fmtArgs: The arguments for the format string.
//
// Returns: None
func (l *defaultLogger) privateLogf(lv Level, format string, fmtArgs []any) {
	if l.level > lv {
		return
	}
	level := lv.toString()
	buf := bytebufferpool.Get()   // Borrow a buffer from the bytebufferpool
	_, _ = buf.WriteString(level) // It is fine to ignore the error

	if len(fmtArgs) > 0 {
		_, _ = fmt.Fprintf(buf, format, fmtArgs...)
	} else {
		_, _ = buf.WriteString(format) // Just write the format string if no args
	}
	_ = l.stdLog.Output(l.depth, buf.String()) // Use the standard logger to output the constructed message
	buf.Reset()                                // Reset the buffer for reuse
	bytebufferpool.Put(buf)                    // Return the buffer to the pool
	if lv == LevelFatal {
		os.Exit(1) // Terminates the program if the level is fatal
	}
}

// privateLogw logs a message with additional key-value pairs at a given level using the default logger.
// Parameters:
//   - lv: The logging level of the message.
//   - format: An optional format string for the message.
//   - keysAndValues: Key-value pairs to include in the log.
//
// Returns: None
func (l *defaultLogger) privateLogw(lv Level, format string, keysAndValues []any) {
	if l.level > lv {
		return
	}
	level := lv.toString()
	buf := bytebufferpool.Get()   // Borrow a buffer from the bytebufferpool
	_, _ = buf.WriteString(level) // It is fine to ignore the error

	// Write the format string to the buffer if provided
	if format != "" {
		_, _ = buf.WriteString(format) // It is fine to ignore the error
	}

	// Write the key-value pairs to the buffer
	if len(keysAndValues) > 0 {
		// Check for unpaired key-value
		if (len(keysAndValues) & 1) == 1 {
			keysAndValues = append(keysAndValues, "KEYVALS UNPAIRED") // Add placeholder for missing value
		}

		// Add a space before key-value pairs if there's a message
		if format != "" {
			_, _ = buf.WriteString(" ") // Add space after message
		}

		// Process key-value pairs
		for i := 0; i < len(keysAndValues); i += 2 {
			if i > 0 {
				_, _ = buf.WriteString(" ") // Add space between pairs
			}
			_, _ = fmt.Fprintf(buf, "%s=%v", keysAndValues[i], keysAndValues[i+1])
		}
	}

	_ = l.stdLog.Output(l.depth, buf.String()) // Use the standard logger to output the constructed message
	buf.Reset()                                // Reset the buffer for reuse
	bytebufferpool.Put(buf)                    // Return the buffer to the pool
	if lv == LevelFatal {
		os.Exit(1) // Terminates the program if the level is fatal
	}
}

// Trace logs a message at the Trace level.
// Parameters:
//   - v: Variadic arguments to be logged.
//
// Returns: None
func (l *defaultLogger) Trace(v ...any) {
	l.privateLog(LevelTrace, v)
}

// Debug logs a message at the Debug level.
// Parameters:
//   - v: Variadic arguments to be logged.
//
// Returns: None
func (l *defaultLogger) Debug(v ...any) {
	l.privateLog(LevelDebug, v)
}

// Info logs a message at the Info level.
// Parameters:
//   - v: Variadic arguments to be logged.
//
// Returns: None
func (l *defaultLogger) Info(v ...any) {
	l.privateLog(LevelInfo, v)
}

// Warn logs a message at the Warn level.
// Parameters:
//   - v: Variadic arguments to be logged.
//
// Returns: None
func (l *defaultLogger) Warn(v ...any) {
	l.privateLog(LevelWarn, v)
}

// Error logs a message at the Error level.
// Parameters:
//   - v: Variadic arguments to be logged.
//
// Returns: None
func (l *defaultLogger) Error(v ...any) {
	l.privateLog(LevelError, v)
}

// Fatal logs a message at the Fatal level and terminates the program.
// Parameters:
//   - v: Variadic arguments to be logged.
//
// Returns: None
func (l *defaultLogger) Fatal(v ...any) {
	l.privateLog(LevelFatal, v)
}

// Panic logs a message at the Panic level.
// Note: Despite the name, this method does not actually panic.
// Parameters:
//   - v: Variadic arguments to be logged.
//
// Returns: None
func (l *defaultLogger) Panic(v ...any) {
	l.privateLog(LevelPanic, v)
}

// Tracef logs a formatted message at the Trace level.
// Parameters:
//   - format: The format string for the message.
//   - v: Variadic arguments for the format string.
//
// Returns: None
func (l *defaultLogger) Tracef(format string, v ...any) {
	l.privateLogf(LevelTrace, format, v)
}

// Debugf logs a formatted message at the Debug level.
// Parameters:
//   - format: The format string for the message.
//   - v: Variadic arguments for the format string.
//
// Returns: None
func (l *defaultLogger) Debugf(format string, v ...any) {
	l.privateLogf(LevelDebug, format, v)
}

// Infof logs a formatted message at the Info level.
// Parameters:
//   - format: The format string for the message.
//   - v: Variadic arguments for the format string.
//
// Returns: None
func (l *defaultLogger) Infof(format string, v ...any) {
	l.privateLogf(LevelInfo, format, v)
}

// Warnf logs a formatted message at the Warn level.
// Parameters:
//   - format: The format string for the message.
//   - v: Variadic arguments for the format string.
//
// Returns: None
func (l *defaultLogger) Warnf(format string, v ...any) {
	l.privateLogf(LevelWarn, format, v)
}

// Errorf logs a formatted message at the Error level.
// Parameters:
//   - format: The format string for the message.
//   - v: Variadic arguments for the format string.
//
// Returns: None
func (l *defaultLogger) Errorf(format string, v ...any) {
	l.privateLogf(LevelError, format, v)
}

// Fatalf logs a formatted message at the Fatal level and terminates the program.
// Parameters:
//   - format: The format string for the message.
//   - v: Variadic arguments for the format string.
//
// Returns: None
func (l *defaultLogger) Fatalf(format string, v ...any) {
	l.privateLogf(LevelFatal, format, v)
}

// Panicf logs a formatted message at the Panic level.
// Note: Despite the name, this method does not actually panic.
// Parameters:
//   - format: The format string for the message.
//   - v: Variadic arguments for the format string.
//
// Returns: None
func (l *defaultLogger) Panicf(format string, v ...any) {
	l.privateLogf(LevelPanic, format, v)
}

// Tracew logs a message at the Trace level with additional key-value pairs.
// Parameters:
//   - msg: The message to be logged.
//   - keysAndValues: Variadic key-value pairs to include in the log.
//
// Returns: None
func (l *defaultLogger) Tracew(msg string, keysAndValues ...any) {
	l.privateLogw(LevelTrace, msg, keysAndValues)
}

// Debugw logs a message at the Debug level with additional key-value pairs.
// Parameters:
//   - msg: The message to be logged.
//   - keysAndValues: Variadic key-value pairs to include in the log.
//
// Returns: None
func (l *defaultLogger) Debugw(msg string, keysAndValues ...any) {
	l.privateLogw(LevelDebug, msg, keysAndValues)
}

// Infow logs a message at the Info level with additional key-value pairs.
// Parameters:
//   - msg: The message to be logged.
//   - keysAndValues: Variadic key-value pairs to include in the log.
//
// Returns: None
func (l *defaultLogger) Infow(msg string, keysAndValues ...any) {
	l.privateLogw(LevelInfo, msg, keysAndValues)
}

// Warnw logs a message at the Warn level with additional key-value pairs.
// Parameters:
//   - msg: The message to be logged.
//   - keysAndValues: Variadic key-value pairs to include in the log.
//
// Returns: None
func (l *defaultLogger) Warnw(msg string, keysAndValues ...any) {
	l.privateLogw(LevelWarn, msg, keysAndValues)
}

// Errorw logs a message at the Error level with additional key-value pairs.
// Parameters:
//   - msg: The message to be logged.
//   - keysAndValues: Variadic key-value pairs to include in the log.
//
// Returns: None
func (l *defaultLogger) Errorw(msg string, keysAndValues ...any) {
	l.privateLogw(LevelError, msg, keysAndValues)
}

// Fatalw logs a message at the Fatal level with additional key-value pairs and terminates the program.
// Parameters:
//   - msg: The message to be logged.
//   - keysAndValues: Variadic key-value pairs to include in the log.
//
// Returns: None
func (l *defaultLogger) Fatalw(msg string, keysAndValues ...any) {
	l.privateLogw(LevelFatal, msg, keysAndValues)
}

// Panicw logs a message at the Panic level with additional key-value pairs.
// Note: Despite the name, this method does not actually panic.
// Parameters:
//   - msg: The message to be logged.
//   - keysAndValues: Variadic key-value pairs to include in the log.
//
// Returns: None
func (l *defaultLogger) Panicw(msg string, keysAndValues ...any) {
	l.privateLogw(LevelPanic, msg, keysAndValues)
}

// WithContext creates a new logger with a modified call depth. This is useful for logging with context.
// Parameters:
//   - _ (context.Context): The context to associate with the logger (ignored in this implementation).
//
// Returns: (CommonLogger) A new logger instance with adjusted depth.
func (l *defaultLogger) WithContext(_ context.Context) CommonLogger {
	return &defaultLogger{
		stdLog: l.stdLog,    // The underlying standard logger for output.
		level:  l.level,     // The current log level to preserve settings.
		depth:  l.depth - 1, // Adjusted call depth for logging.
	}
}

// SetLevel sets the logging level for the logger.
// Parameters:
//   - level (Level): The desired logging level to set.
//
// Returns: None
func (l *defaultLogger) SetLevel(level Level) {
	l.level = level
}

// SetOutput sets the output destination for the logger.
// Parameters:
//   - writer (io.Writer): The destination to write log output to.
//
// Returns: None
func (l *defaultLogger) SetOutput(writer io.Writer) {
	l.stdLog.SetOutput(writer)
}

// DefaultLogger returns the default logger instance.
// Parameters: None
// Returns: (AllLogger) The default logger instance.
func DefaultLogger() AllLogger {
	return logger
}
