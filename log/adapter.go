package log

import (
	"context"
	"io"
)

// Fatal calls the default logger's Fatal method and then exits the program.
//
// Parameters:
//   - v (...any): Variadic arguments that are passed to the logger.
//
// Returns:
//   - None (the program terminates with os.Exit(1)).
func Fatal(v ...any) {
	logger.Fatal(v...)
}

// Error logs an error message using the default logger's Error method.
//
// Parameters:
//   - v (...any): Variadic arguments that are passed to the logger.
//
// Returns:
//   - None.
func Error(v ...any) {
	logger.Error(v...)
}

// Warn logs a warning message using the default logger's Warn method.
//
// Parameters:
//   - v (...any): Variadic arguments that are passed to the logger.
//
// Returns:
//   - None.
func Warn(v ...any) {
	logger.Warn(v...)
}

// Info logs an informational message using the default logger's Info method.
//
// Parameters:
//   - v (...any): Variadic arguments that are passed to the logger.
//
// Returns:
//   - None.
func Info(v ...any) {
	logger.Info(v...)
}

// Debug logs a debug message using the default logger's Debug method.
//
// Parameters:
//   - v (...any): Variadic arguments that are passed to the logger.
//
// Returns:
//   - None.
func Debug(v ...any) {
	logger.Debug(v...)
}

// Trace logs a trace message using the default logger's Trace method.
//
// Parameters:
//   - v (...any): Variadic arguments that are passed to the logger.
//
// Returns:
//   - None.
func Trace(v ...any) {
	logger.Trace(v...)
}

// Panic logs a message at the Panic level using the default logger's Panic method.
// Note: Despite the name, this method does not actually panic.
//
// Parameters:
//   - v (...any): Variadic arguments that are passed to the logger.
//
// Returns:
//   - None.
func Panic(v ...any) {
	logger.Panic(v...)
}

// Fatalf logs a formatted fatal message and then exits the program.
//
// Parameters:
//   - format (string): The format string.
//   - v (...any): Values to format.
//
// Returns:
//   - None (the program terminates with os.Exit(1)).
func Fatalf(format string, v ...any) {
	logger.Fatalf(format, v...)
}

// Errorf logs a formatted error message using the default logger's Errorf method.
//
// Parameters:
//   - format (string): The format string.
//   - v (...any): Values to format.
//
// Returns:
//   - None.
func Errorf(format string, v ...any) {
	logger.Errorf(format, v...)
}

// Warnf logs a formatted warning message using the default logger's Warnf method.
//
// Parameters:
//   - format (string): The format string.
//   - v (...any): Values to format.
//
// Returns:
//   - None.
func Warnf(format string, v ...any) {
	logger.Warnf(format, v...)
}

// Infof logs a formatted informational message using the default logger's Infof method.
//
// Parameters:
//   - format (string): The format string.
//   - v (...any): Values to format.
//
// Returns:
//   - None.
func Infof(format string, v ...any) {
	logger.Infof(format, v...)
}

// Debugf logs a formatted debug message using the default logger's Debugf method.
//
// Parameters:
//   - format (string): The format string.
//   - v (...any): Values to format.
//
// Returns:
//   - None.
func Debugf(format string, v ...any) {
	logger.Debugf(format, v...)
}

// Tracef logs a formatted trace message using the default logger's Tracef method.
//
// Parameters:
//   - format (string): The format string.
//   - v (...any): Values to format.
//
// Returns:
//   - None.
func Tracef(format string, v ...any) {
	logger.Tracef(format, v...)
}

// Panicf logs a formatted message at the Panic level using the default logger's Panicf method.
// Note: Despite the name, this method does not actually panic.
//
// Parameters:
//   - format (string): The format string.
//   - v (...any): Values to format.
//
// Returns:
//   - None.
func Panicf(format string, v ...any) {
	logger.Panicf(format, v...)
}

// Tracew logs a trace message with context key-value pairs using the default logger's Tracew method.
//
// Parameters:
//   - msg (string): The message to log.
//   - keysAndValues (...any): Key-value pairs providing additional context.
//
// Returns:
//   - None.
func Tracew(msg string, keysAndValues ...any) {
	logger.Tracew(msg, keysAndValues...)
}

// Debugw logs a debug message with context key-value pairs using the default logger's Debugw method.
//
// Parameters:
//   - msg (string): The message to log.
//   - keysAndValues (...any): Key-value pairs providing additional context.
//
// Returns:
//   - None.
func Debugw(msg string, keysAndValues ...any) {
	logger.Debugw(msg, keysAndValues...)
}

// Infow logs an informational message with context key-value pairs using the default logger's Infow method.
//
// Parameters:
//   - msg (string): The message to log.
//   - keysAndValues (...any): Key-value pairs providing additional context.
//
// Returns:
//   - None.
func Infow(msg string, keysAndValues ...any) {
	logger.Infow(msg, keysAndValues...)
}

// Warnw logs a warning message with context key-value pairs using the default logger's Warnw method.
//
// Parameters:
//   - msg (string): The message to log.
//   - keysAndValues (...any): Key-value pairs providing additional context.
//
// Returns:
//   - None.
func Warnw(msg string, keysAndValues ...any) {
	logger.Warnw(msg, keysAndValues...)
}

// Errorw logs an error message with context key-value pairs using the default logger's Errorw method.
//
// Parameters:
//   - msg (string): The message to log.
//   - keysAndValues (...any): Key-value pairs providing additional context.
//
// Returns:
//   - None.
func Errorw(msg string, keysAndValues ...any) {
	logger.Errorw(msg, keysAndValues...)
}

// Fatalw logs a fatal message with context key-value pairs and then exits the program.
//
// Parameters:
//   - msg (string): The message to log.
//   - keysAndValues (...any): Key-value pairs providing additional context.
//
// Returns:
//   - None (the program terminates with os.Exit(1)).
func Fatalw(msg string, keysAndValues ...any) {
	logger.Fatalw(msg, keysAndValues...)
}

// Panicw logs a message at the Panic level with context key-value pairs using the default logger's Panicw method.
// Note: Despite the name, this method does not actually panic.
//
// Parameters:
//   - msg (string): The message to log.
//   - keysAndValues (...any): Key-value pairs providing additional context.
//
// Returns:
//   - None.
func Panicw(msg string, keysAndValues ...any) {
	logger.Panicw(msg, keysAndValues...)
}

// WithContext creates a new logger with the provided context.
//
// Parameters:
//   - ctx (context.Context): The context to associate with the logger.
//
// Returns:
//   - CommonLogger: A logger instance associated with the given context.
func WithContext(ctx context.Context) CommonLogger {
	return logger.WithContext(ctx)
}

// SetLogger sets the default logger to a new instance.
//
// Parameters:
//   - v (AllLogger): The new logger to be set.
//
// Returns:
//   - None.
func SetLogger(v AllLogger) {
	logger = v
}

// SetOutput changes the output destination of the logger.
//
// Parameters:
//   - w (io.Writer): The new output destination.
//
// Returns:
//   - None.
func SetOutput(w io.Writer) {
	logger.SetOutput(w)
}

// SetLevel adjusts the logging level to filter messages below a certain severity.
//
// Parameters:
//   - lv (Level): The logging level to set.
//
// Returns:
//   - None.
func SetLevel(lv Level) {
	logger.SetLevel(lv)
}
