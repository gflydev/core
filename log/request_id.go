package log

import (
	"context"
	"fmt"
)

// RequestIDKey is the key used to store the request ID in the context
const RequestIDKey = "request_id"

// WithRequestID adds a request ID to the logger context.
//
// Parameters:
//   - ctx (context.Context): The context to add the request ID to.
//   - requestID (string): The request ID to add.
//
// Returns:
//   - context.Context: A new context with the request ID added.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// GetRequestID retrieves the request ID from the context.
//
// Parameters:
//   - ctx (context.Context): The context to retrieve the request ID from.
//
// Returns:
//   - string: The request ID, or an empty string if not found.
func GetRequestID(ctx context.Context) string {
	requestID, ok := ctx.Value(RequestIDKey).(string)
	if !ok {
		return ""
	}
	return requestID
}

// RequestIDLogger is a middleware that adds the request ID to the logger context.
// This ensures that all log entries for the current request include the request ID.
type RequestIDLogger struct {
	logger CommonLogger
	ctx    context.Context
}

// NewRequestIDLogger creates a new logger that includes the request ID in all log entries.
//
// Parameters:
//   - logger (CommonLogger): The base logger to wrap.
//   - requestID (string): The request ID to include in log entries.
//
// Returns:
//   - CommonLogger: A new logger that includes the request ID in all log entries.
func NewRequestIDLogger(logger CommonLogger, requestID string) CommonLogger {
	ctx := WithRequestID(context.Background(), requestID)
	return &RequestIDLogger{
		logger: logger,
		ctx:    ctx,
	}
}

// Implement the CommonLogger interface methods to include the request ID in log entries

func (l *RequestIDLogger) Trace(v ...interface{}) {
	// Convert to a message with request ID as a field
	msg := fmt.Sprint(v...)
	l.logger.Tracew(msg, "request_id", GetRequestID(l.ctx))
}

func (l *RequestIDLogger) Debug(v ...interface{}) {
	// Convert to a message with request ID as a field
	msg := fmt.Sprint(v...)
	l.logger.Debugw(msg, "request_id", GetRequestID(l.ctx))
}

func (l *RequestIDLogger) Info(v ...interface{}) {
	// Convert to a message with request ID as a field
	msg := fmt.Sprint(v...)
	l.logger.Infow(msg, "request_id", GetRequestID(l.ctx))
}

func (l *RequestIDLogger) Warn(v ...interface{}) {
	// Convert to a message with request ID as a field
	msg := fmt.Sprint(v...)
	l.logger.Warnw(msg, "request_id", GetRequestID(l.ctx))
}

func (l *RequestIDLogger) Error(v ...interface{}) {
	// Convert to a message with request ID as a field
	msg := fmt.Sprint(v...)
	l.logger.Errorw(msg, "request_id", GetRequestID(l.ctx))
}

func (l *RequestIDLogger) Fatal(v ...interface{}) {
	// Convert to a message with request ID as a field
	msg := fmt.Sprint(v...)
	l.logger.Fatalw(msg, "request_id", GetRequestID(l.ctx))
}

func (l *RequestIDLogger) Panic(v ...interface{}) {
	// Convert to a message with request ID as a field
	msg := fmt.Sprint(v...)
	l.logger.Panicw(msg, "request_id", GetRequestID(l.ctx))
}

func (l *RequestIDLogger) Tracef(format string, v ...interface{}) {
	// Format the message and add request ID as a field
	msg := fmt.Sprintf(format, v...)
	l.logger.Tracew(msg, "request_id", GetRequestID(l.ctx))
}

func (l *RequestIDLogger) Debugf(format string, v ...interface{}) {
	// Format the message and add request ID as a field
	msg := fmt.Sprintf(format, v...)
	l.logger.Debugw(msg, "request_id", GetRequestID(l.ctx))
}

func (l *RequestIDLogger) Infof(format string, v ...interface{}) {
	// Format the message and add request ID as a field
	msg := fmt.Sprintf(format, v...)
	l.logger.Infow(msg, "request_id", GetRequestID(l.ctx))
}

func (l *RequestIDLogger) Warnf(format string, v ...interface{}) {
	// Format the message and add request ID as a field
	msg := fmt.Sprintf(format, v...)
	l.logger.Warnw(msg, "request_id", GetRequestID(l.ctx))
}

func (l *RequestIDLogger) Errorf(format string, v ...interface{}) {
	// Format the message and add request ID as a field
	msg := fmt.Sprintf(format, v...)
	l.logger.Errorw(msg, "request_id", GetRequestID(l.ctx))
}

func (l *RequestIDLogger) Fatalf(format string, v ...interface{}) {
	// Format the message and add request ID as a field
	msg := fmt.Sprintf(format, v...)
	l.logger.Fatalw(msg, "request_id", GetRequestID(l.ctx))
}

func (l *RequestIDLogger) Panicf(format string, v ...interface{}) {
	// Format the message and add request ID as a field
	msg := fmt.Sprintf(format, v...)
	l.logger.Panicw(msg, "request_id", GetRequestID(l.ctx))
}

func (l *RequestIDLogger) Tracew(msg string, keysAndValues ...interface{}) {
	// Add request ID to the key-value pairs
	keysAndValues = append(keysAndValues, "request_id", GetRequestID(l.ctx))
	l.logger.Tracew(msg, keysAndValues...)
}

func (l *RequestIDLogger) Debugw(msg string, keysAndValues ...interface{}) {
	// Add request ID to the key-value pairs
	keysAndValues = append(keysAndValues, "request_id", GetRequestID(l.ctx))
	l.logger.Debugw(msg, keysAndValues...)
}

func (l *RequestIDLogger) Infow(msg string, keysAndValues ...interface{}) {
	// Add request ID to the key-value pairs
	keysAndValues = append(keysAndValues, "request_id", GetRequestID(l.ctx))
	l.logger.Infow(msg, keysAndValues...)
}

func (l *RequestIDLogger) Warnw(msg string, keysAndValues ...interface{}) {
	// Add request ID to the key-value pairs
	keysAndValues = append(keysAndValues, "request_id", GetRequestID(l.ctx))
	l.logger.Warnw(msg, keysAndValues...)
}

func (l *RequestIDLogger) Errorw(msg string, keysAndValues ...interface{}) {
	// Add request ID to the key-value pairs
	keysAndValues = append(keysAndValues, "request_id", GetRequestID(l.ctx))
	l.logger.Errorw(msg, keysAndValues...)
}

func (l *RequestIDLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	// Add request ID to the key-value pairs
	keysAndValues = append(keysAndValues, "request_id", GetRequestID(l.ctx))
	l.logger.Fatalw(msg, keysAndValues...)
}

func (l *RequestIDLogger) Panicw(msg string, keysAndValues ...interface{}) {
	// Add request ID to the key-value pairs
	keysAndValues = append(keysAndValues, "request_id", GetRequestID(l.ctx))
	l.logger.Panicw(msg, keysAndValues...)
}
