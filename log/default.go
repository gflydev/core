package log

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/valyala/bytebufferpool"
)

var _ AllLogger = (*defaultLogger)(nil)

// logMessage represents a message to be logged
type logMessage struct {
	level   Level
	message string
	fatal   bool
}

type defaultLogger struct {
	stdLog    *log.Logger     // The underlying standard logger used for output
	level     Level           // The current log level
	depth     int             // The call depth for logging
	msgChan   chan logMessage // Channel for log messages
	wg        sync.WaitGroup  // WaitGroup for graceful shutdown
	closeChan chan struct{}   // Channel to signal logger shutdown
	closed    bool            // Flag to indicate if logger is closed
	mu        sync.Mutex      // Mutex to protect closed flag
}

// privateLog logs a message at a given level using the default logger. It uses a buffer pool to optimize memory usage.
// Parameters:
//   - lv: The logging level of the message.
//   - fmtArgs: The arguments to be logged.
//
// Returns: None
func (l *defaultLogger) privateLog(lv Level, fmtArgs []interface{}) {
	if l.level > lv {
		return
	}

	// Check if logger is closed
	l.mu.Lock()
	if l.closed {
		l.mu.Unlock()
		return
	}
	l.mu.Unlock()

	level := lv.toString()
	buf := bytebufferpool.Get()     // Borrow a buffer from the bytebufferpool
	_, _ = buf.WriteString(level)   // It is fine to ignore the error
	_, _ = fmt.Fprint(buf, fmtArgs) // It is fine to ignore the error

	// Send the message to the channel for asynchronous logging
	msg := logMessage{
		level:   lv,
		message: buf.String(),
		fatal:   lv == LevelFatal,
	}

	// Reset and return the buffer to the pool
	buf.Reset()
	bytebufferpool.Put(buf)

	// Send the message to the channel
	select {
	case l.msgChan <- msg:
		// Message sent successfully
	case <-l.closeChan:
		// Logger is shutting down
		return
	}

	// If it's a fatal message, wait for it to be logged before exiting
	if lv == LevelFatal {
		// Wait a moment for the message to be processed
		l.wg.Wait()
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
func (l *defaultLogger) privateLogf(lv Level, format string, fmtArgs []interface{}) {
	if l.level > lv {
		return
	}

	// Check if logger is closed
	l.mu.Lock()
	if l.closed {
		l.mu.Unlock()
		return
	}
	l.mu.Unlock()

	level := lv.toString()
	buf := bytebufferpool.Get()   // Borrow a buffer from the bytebufferpool
	_, _ = buf.WriteString(level) // It is fine to ignore the error

	if len(fmtArgs) > 0 {
		_, _ = fmt.Fprintf(buf, format, fmtArgs...)
	} else {
		_, _ = buf.WriteString(format) // Just write the format string if no args
	}

	// Send the message to the channel for asynchronous logging
	msg := logMessage{
		level:   lv,
		message: buf.String(),
		fatal:   lv == LevelFatal,
	}

	// Reset and return the buffer to the pool
	buf.Reset()
	bytebufferpool.Put(buf)

	// Send the message to the channel
	select {
	case l.msgChan <- msg:
		// Message sent successfully
	case <-l.closeChan:
		// Logger is shutting down
		return
	}

	// If it's a fatal message, wait for it to be logged before exiting
	if lv == LevelFatal {
		// Wait a moment for the message to be processed
		l.wg.Wait()
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
func (l *defaultLogger) privateLogw(lv Level, format string, keysAndValues []interface{}) {
	if l.level > lv {
		return
	}

	// Check if logger is closed
	l.mu.Lock()
	if l.closed {
		l.mu.Unlock()
		return
	}
	l.mu.Unlock()

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

	// Send the message to the channel for asynchronous logging
	msg := logMessage{
		level:   lv,
		message: buf.String(),
		fatal:   lv == LevelFatal,
	}

	// Reset and return the buffer to the pool
	buf.Reset()
	bytebufferpool.Put(buf)

	// Send the message to the channel
	select {
	case l.msgChan <- msg:
		// Message sent successfully
	case <-l.closeChan:
		// Logger is shutting down
		return
	}

	// If it's a fatal message, wait for it to be logged before exiting
	if lv == LevelFatal {
		// Wait a moment for the message to be processed
		l.wg.Wait()
		os.Exit(1) // Terminates the program if the level is fatal
	}
}

// Trace logs a message at the Trace level.
// Parameters:
//   - v: Variadic arguments to be logged.
//
// Returns: None
func (l *defaultLogger) Trace(v ...interface{}) {
	l.privateLog(LevelTrace, v)
}

// Debug logs a message at the Debug level.
// Parameters:
//   - v: Variadic arguments to be logged.
//
// Returns: None
func (l *defaultLogger) Debug(v ...interface{}) {
	l.privateLog(LevelDebug, v)
}

// Info logs a message at the Info level.
// Parameters:
//   - v: Variadic arguments to be logged.
//
// Returns: None
func (l *defaultLogger) Info(v ...interface{}) {
	l.privateLog(LevelInfo, v)
}

// Warn logs a message at the Warn level.
// Parameters:
//   - v: Variadic arguments to be logged.
//
// Returns: None
func (l *defaultLogger) Warn(v ...interface{}) {
	l.privateLog(LevelWarn, v)
}

// Error logs a message at the Error level.
// Parameters:
//   - v: Variadic arguments to be logged.
//
// Returns: None
func (l *defaultLogger) Error(v ...interface{}) {
	l.privateLog(LevelError, v)
}

// Fatal logs a message at the Fatal level and terminates the program.
// Parameters:
//   - v: Variadic arguments to be logged.
//
// Returns: None
func (l *defaultLogger) Fatal(v ...interface{}) {
	l.privateLog(LevelFatal, v)
}

// Panic logs a message at the Panic level.
// Note: Despite the name, this method does not actually panic.
// Parameters:
//   - v: Variadic arguments to be logged.
//
// Returns: None
func (l *defaultLogger) Panic(v ...interface{}) {
	l.privateLog(LevelPanic, v)
}

// Tracef logs a formatted message at the Trace level.
// Parameters:
//   - format: The format string for the message.
//   - v: Variadic arguments for the format string.
//
// Returns: None
func (l *defaultLogger) Tracef(format string, v ...interface{}) {
	l.privateLogf(LevelTrace, format, v)
}

// Debugf logs a formatted message at the Debug level.
// Parameters:
//   - format: The format string for the message.
//   - v: Variadic arguments for the format string.
//
// Returns: None
func (l *defaultLogger) Debugf(format string, v ...interface{}) {
	l.privateLogf(LevelDebug, format, v)
}

// Infof logs a formatted message at the Info level.
// Parameters:
//   - format: The format string for the message.
//   - v: Variadic arguments for the format string.
//
// Returns: None
func (l *defaultLogger) Infof(format string, v ...interface{}) {
	l.privateLogf(LevelInfo, format, v)
}

// Warnf logs a formatted message at the Warn level.
// Parameters:
//   - format: The format string for the message.
//   - v: Variadic arguments for the format string.
//
// Returns: None
func (l *defaultLogger) Warnf(format string, v ...interface{}) {
	l.privateLogf(LevelWarn, format, v)
}

// Errorf logs a formatted message at the Error level.
// Parameters:
//   - format: The format string for the message.
//   - v: Variadic arguments for the format string.
//
// Returns: None
func (l *defaultLogger) Errorf(format string, v ...interface{}) {
	l.privateLogf(LevelError, format, v)
}

// Fatalf logs a formatted message at the Fatal level and terminates the program.
// Parameters:
//   - format: The format string for the message.
//   - v: Variadic arguments for the format string.
//
// Returns: None
func (l *defaultLogger) Fatalf(format string, v ...interface{}) {
	l.privateLogf(LevelFatal, format, v)
}

// Panicf logs a formatted message at the Panic level.
// Note: Despite the name, this method does not actually panic.
// Parameters:
//   - format: The format string for the message.
//   - v: Variadic arguments for the format string.
//
// Returns: None
func (l *defaultLogger) Panicf(format string, v ...interface{}) {
	l.privateLogf(LevelPanic, format, v)
}

// Tracew logs a message at the Trace level with additional key-value pairs.
// Parameters:
//   - msg: The message to be logged.
//   - keysAndValues: Variadic key-value pairs to include in the log.
//
// Returns: None
func (l *defaultLogger) Tracew(msg string, keysAndValues ...interface{}) {
	l.privateLogw(LevelTrace, msg, keysAndValues)
}

// Debugw logs a message at the Debug level with additional key-value pairs.
// Parameters:
//   - msg: The message to be logged.
//   - keysAndValues: Variadic key-value pairs to include in the log.
//
// Returns: None
func (l *defaultLogger) Debugw(msg string, keysAndValues ...interface{}) {
	l.privateLogw(LevelDebug, msg, keysAndValues)
}

// Infow logs a message at the Info level with additional key-value pairs.
// Parameters:
//   - msg: The message to be logged.
//   - keysAndValues: Variadic key-value pairs to include in the log.
//
// Returns: None
func (l *defaultLogger) Infow(msg string, keysAndValues ...interface{}) {
	l.privateLogw(LevelInfo, msg, keysAndValues)
}

// Warnw logs a message at the Warn level with additional key-value pairs.
// Parameters:
//   - msg: The message to be logged.
//   - keysAndValues: Variadic key-value pairs to include in the log.
//
// Returns: None
func (l *defaultLogger) Warnw(msg string, keysAndValues ...interface{}) {
	l.privateLogw(LevelWarn, msg, keysAndValues)
}

// Errorw logs a message at the Error level with additional key-value pairs.
// Parameters:
//   - msg: The message to be logged.
//   - keysAndValues: Variadic key-value pairs to include in the log.
//
// Returns: None
func (l *defaultLogger) Errorw(msg string, keysAndValues ...interface{}) {
	l.privateLogw(LevelError, msg, keysAndValues)
}

// Fatalw logs a message at the Fatal level with additional key-value pairs and terminates the program.
// Parameters:
//   - msg: The message to be logged.
//   - keysAndValues: Variadic key-value pairs to include in the log.
//
// Returns: None
func (l *defaultLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.privateLogw(LevelFatal, msg, keysAndValues)
}

// Panicw logs a message at the Panic level with additional key-value pairs.
// Note: Despite the name, this method does not actually panic.
// Parameters:
//   - msg: The message to be logged.
//   - keysAndValues: Variadic key-value pairs to include in the log.
//
// Returns: None
func (l *defaultLogger) Panicw(msg string, keysAndValues ...interface{}) {
	l.privateLogw(LevelPanic, msg, keysAndValues)
}

// WithContext creates a new logger with a modified call depth. This is useful for logging with context.
// Parameters:
//   - _ (context.Context): The context to associate with the logger (ignored in this implementation).
//
// Returns: (CommonLogger) A new logger instance with adjusted depth.
func (l *defaultLogger) WithContext(_ context.Context) CommonLogger {
	return &defaultLogger{
		stdLog:    l.stdLog,    // The underlying standard logger for output.
		level:     l.level,     // The current log level to preserve settings.
		depth:     l.depth - 1, // Adjusted call depth for logging.
		msgChan:   l.msgChan,   // Share the same message channel
		closeChan: l.closeChan, // Share the same close channel
		wg:        l.wg,        // Share the same wait group
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

// startLoggerWorker starts a goroutine that processes log messages from the channel
// and writes them to the output.
func (l *defaultLogger) startLoggerWorker() {
	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		for {
			select {
			case msg := <-l.msgChan:
				// Write the message to the output
				_ = l.stdLog.Output(l.depth, msg.message)
			case <-l.closeChan:
				// Process any remaining messages in the channel
				for {
					select {
					case msg := <-l.msgChan:
						_ = l.stdLog.Output(l.depth, msg.message)
					default:
						return
					}
				}
			}
		}
	}()
}

// Close gracefully shuts down the logger, ensuring all pending messages are processed.
func (l *defaultLogger) Close() {
	l.mu.Lock()
	if l.closed {
		l.mu.Unlock()
		return
	}
	l.closed = true
	close(l.closeChan)
	l.mu.Unlock()

	// Wait for the logger goroutine to finish
	l.wg.Wait()
}

// newDefaultLogger creates and initializes a new defaultLogger instance.
// It sets up the channels and starts the logger goroutine.
// Returns: (AllLogger) A new defaultLogger instance.
func newDefaultLogger() AllLogger {
	l := &defaultLogger{
		stdLog:    log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile|log.Lmicroseconds),
		level:     LevelInfo, // Default level
		depth:     4,
		msgChan:   make(chan logMessage, 1000), // Buffer size of 1000 messages
		closeChan: make(chan struct{}),
	}

	// Start the logger goroutine
	l.startLoggerWorker()

	return l
}

// DefaultLogger returns the default logger instance.
// Parameters: None
// Returns: (AllLogger) The default logger instance.
func DefaultLogger() AllLogger {
	return logger
}
