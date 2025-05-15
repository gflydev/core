package core

import (
	"github.com/gflydev/core/log"
)

// RequestIDLogMiddleware is a middleware that adds the request ID to the logger context.
// This ensures that all log entries for the current request include the request ID.
//
// Parameters:
//   - ctx (*Ctx): The request context.
//
// Returns:
//   - error: Always returns nil to continue the middleware chain.
func RequestIDLogMiddleware(ctx *Ctx) error {
	// Get the request ID from the context
	requestID := GetRequestID(ctx)
	if requestID == "" {
		// If no request ID is found, this middleware should be used after RequestIDMiddleware
		return nil
	}

	// Get the default logger
	defaultLogger := log.DefaultLogger()

	// Create a new logger with the request ID
	requestLogger := log.NewRequestIDLogger(defaultLogger, requestID)

	// Store the logger in the context for use by other middleware and handlers
	ctx.SetData("logger", requestLogger)

	return nil
}

// GetLogger retrieves the request-specific logger from the context.
// If no logger is found, it returns the default logger.
//
// Parameters:
//   - ctx (*Ctx): The request context.
//
// Returns:
//   - log.CommonLogger: The request-specific logger, or the default logger if not found.
func GetLogger(ctx *Ctx) log.CommonLogger {
	logger := ctx.GetData("logger")
	if logger == nil {
		return log.DefaultLogger()
	}
	return logger.(log.CommonLogger)
}
