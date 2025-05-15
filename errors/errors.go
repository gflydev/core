// Package errors provides standardized error handling for the gFly framework.
// It defines error types, interfaces, and helper functions for creating,
// wrapping, and handling errors in a consistent way across all packages.
//
// The package implements a comprehensive error handling system with the following features:
//
// Error Categorization:
//   - Errors are categorized by error codes (e.g., not_found, invalid_input)
//   - Each error code maps to an appropriate HTTP status code for API responses
//
// Error Context and Metadata:
//   - Errors can carry context information for debugging purposes
//   - Errors can include metadata for additional information
//
// Error Wrapping:
//   - Preserves the error chain when wrapping errors
//   - Compatible with the standard Go errors package
//
// Stack Traces:
//   - Automatically captures stack traces for better debugging
//   - Provides detailed information about where errors occurred
//
// JSON Serialization:
//   - Errors can be serialized to JSON for API responses
//   - Includes all relevant error information in the JSON output
//
// Helper Functions:
//   - Provides helper functions for creating common error types
//   - Includes formatted versions of helper functions for convenience
//
// Usage Examples:
//
//	// Create a basic error
//	err := errors.New(errors.CodeNotFound, "User not found")
//
//	// Create an error with a formatted message
//	err := errors.Newf(errors.CodeInvalidInput, "Invalid value for field %s: %v", "email", value)
//
//	// Use helper functions for common error types
//	err := errors.NotFound("User not found")
//	err := errors.InvalidInput("Invalid email address")
//
//	// Wrap an existing error
//	if originalErr := db.QueryRow("SELECT * FROM users WHERE id = ?", id); originalErr != nil {
//	    return errors.Wrap(originalErr, errors.CodeNotFound, "User not found")
//	}
//
//	// Add context information
//	err := errors.InvalidInput("Invalid email address").
//	    WithContext("user_id", userId).
//	    WithContext("request_id", requestId)
//
// For more detailed information and best practices, see the package README.md.
package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"strings"
)

// Error codes for categorizing errors
const (
	// CodeUnknown represents an unknown or uncategorized error
	CodeUnknown = "unknown"
	// CodeNotFound represents a resource not found error
	CodeNotFound = "not_found"
	// CodeInvalidInput represents an invalid input error
	CodeInvalidInput = "invalid_input"
	// CodeUnauthorized represents an authorization error
	CodeUnauthorized = "unauthorized"
	// CodeUnauthenticated represents an authentication error
	CodeUnauthenticated = "unauthenticated"
	// CodeForbidden represents a forbidden access error
	CodeForbidden = "forbidden"
	// CodeConflict represents a conflict error (e.g., duplicate resource)
	CodeConflict = "conflict"
	// CodeInternal represents an internal server error
	CodeInternal = "internal"
	// CodeTimeout represents a timeout error
	CodeTimeout = "timeout"
	// CodeUnavailable represents a service unavailable error
	CodeUnavailable = "unavailable"
	// CodeNotImplemented represents a not implemented error
	CodeNotImplemented = "not_implemented"
)

// HTTP status codes mapped to error codes
var statusCodeMap = map[string]int{
	CodeUnknown:         500,
	CodeNotFound:        404,
	CodeInvalidInput:    400,
	CodeUnauthorized:    401,
	CodeUnauthenticated: 401,
	CodeForbidden:       403,
	CodeConflict:        409,
	CodeInternal:        500,
	CodeTimeout:         504,
	CodeUnavailable:     503,
	CodeNotImplemented:  501,
}

// Error is the base error interface that all gFly errors implement.
// It extends the standard error interface with additional methods for
// error categorization, context, and metadata.
type Error interface {
	error
	// Code returns the error code for categorization
	Code() string
	// Message returns a user-friendly error message
	Message() string
	// StatusCode returns the HTTP status code associated with this error
	StatusCode() int
	// Unwrap returns the underlying error if this is a wrapped error
	Unwrap() error
	// WithMetadata adds metadata to the error and returns a new error
	WithMetadata(key string, value interface{}) Error
	// GetMetadata retrieves metadata from the error
	GetMetadata(key string) (interface{}, bool)
	// WithContext adds context information to the error and returns a new error
	WithContext(key string, value interface{}) Error
	// GetContext retrieves context information from the error
	GetContext(key string) (interface{}, bool)
	// StackTrace returns the stack trace at the point the error was created
	StackTrace() string
}

// baseError is the implementation of the Error interface
type baseError struct {
	code       string                 // Error code for categorization
	message    string                 // User-friendly error message
	cause      error                  // Underlying cause of the error
	metadata   map[string]interface{} // Additional metadata about the error
	context    map[string]interface{} // Context information about where the error occurred
	stackTrace string                 // Stack trace at the point the error was created
}

// Error implements the standard error interface
func (e *baseError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %v", e.message, e.cause)
	}
	return e.message
}

// Code returns the error code
func (e *baseError) Code() string {
	return e.code
}

// Message returns the user-friendly error message
func (e *baseError) Message() string {
	return e.message
}

// StatusCode returns the HTTP status code associated with this error
func (e *baseError) StatusCode() int {
	if code, ok := statusCodeMap[e.code]; ok {
		return code
	}
	return 500 // Default to internal server error
}

// Unwrap returns the underlying error
func (e *baseError) Unwrap() error {
	return e.cause
}

// WithMetadata adds metadata to the error and returns a new error
func (e *baseError) WithMetadata(key string, value interface{}) Error {
	// Create a copy of the error to avoid modifying the original
	newErr := *e
	if newErr.metadata == nil {
		newErr.metadata = make(map[string]interface{})
	}
	newErr.metadata[key] = value
	return &newErr
}

// GetMetadata retrieves metadata from the error
func (e *baseError) GetMetadata(key string) (interface{}, bool) {
	if e.metadata == nil {
		return nil, false
	}
	value, ok := e.metadata[key]
	return value, ok
}

// WithContext adds context information to the error and returns a new error
func (e *baseError) WithContext(key string, value interface{}) Error {
	// Create a copy of the error to avoid modifying the original
	newErr := *e
	if newErr.context == nil {
		newErr.context = make(map[string]interface{})
	}
	newErr.context[key] = value
	return &newErr
}

// GetContext retrieves context information from the error
func (e *baseError) GetContext(key string) (interface{}, bool) {
	if e.context == nil {
		return nil, false
	}
	value, ok := e.context[key]
	return value, ok
}

// StackTrace returns the stack trace at the point the error was created
func (e *baseError) StackTrace() string {
	return e.stackTrace
}

// MarshalJSON implements json.Marshaler for custom JSON serialization
func (e *baseError) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"code":      e.code,
		"message":   e.message,
		"metadata":  e.metadata,
		"context":   e.context,
		"cause":     e.cause,
		"stack":     e.stackTrace,
		"http_code": e.StatusCode(),
	})
}

// New creates a new error with the given code and message
func New(code string, message string) Error {
	return &baseError{
		code:       code,
		message:    message,
		stackTrace: captureStackTrace(2),
	}
}

// Newf creates a new error with the given code and formatted message
func Newf(code string, format string, args ...interface{}) Error {
	return &baseError{
		code:       code,
		message:    fmt.Sprintf(format, args...),
		stackTrace: captureStackTrace(2),
	}
}

// Wrap wraps an existing error with additional context
func Wrap(err error, code string, message string) Error {
	if err == nil {
		return nil
	}

	// If the error is already a gFly error, preserve its code if none is provided
	if gflyErr, ok := err.(Error); ok && code == "" {
		code = gflyErr.Code()
	}

	// Default to unknown code if none is provided
	if code == "" {
		code = CodeUnknown
	}

	return &baseError{
		code:       code,
		message:    message,
		cause:      err,
		stackTrace: captureStackTrace(2),
	}
}

// Wrapf wraps an existing error with additional context and a formatted message
func Wrapf(err error, code string, format string, args ...interface{}) Error {
	if err == nil {
		return nil
	}

	// If the error is already a gFly error, preserve its code if none is provided
	if gflyErr, ok := err.(Error); ok && code == "" {
		code = gflyErr.Code()
	}

	// Default to unknown code if none is provided
	if code == "" {
		code = CodeUnknown
	}

	return &baseError{
		code:       code,
		message:    fmt.Sprintf(format, args...),
		cause:      err,
		stackTrace: captureStackTrace(2),
	}
}

// Is reports whether any error in err's chain matches target.
// This is compatible with the standard errors.Is function.
func Is(err, target error) bool {
	if err == target {
		return true
	}

	// Unwrap the error if it's a gFly error
	if gflyErr, ok := err.(Error); ok {
		if unwrapped := gflyErr.Unwrap(); unwrapped != nil {
			return Is(unwrapped, target)
		}
	}

	// Use the standard errors.Is for other error types
	return false
}

// As finds the first error in err's chain that matches the type of target,
// and if so, sets target to that error value and returns true.
// This is compatible with the standard errors.As function.
func As(err error, target interface{}) bool {
	if target == nil {
		panic("errors: target cannot be nil")
	}

	// Try to cast the error to the target type
	val := fmt.Sprintf("%T", target)
	if strings.Contains(val, "Error") {
		if gflyErr, ok := err.(Error); ok {
			if targetErr, ok := target.(*Error); ok {
				*targetErr = gflyErr
				return true
			}
		}
	}

	// Unwrap the error if it's a gFly error
	var gflyErr Error
	if errors.As(err, &gflyErr) {
		if unwrapped := gflyErr.Unwrap(); unwrapped != nil {
			return As(unwrapped, target)
		}
	}

	// Use the standard errors.As for other error types
	return false
}

// captureStackTrace captures the current stack trace
func captureStackTrace(skip int) string {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(skip, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])

	var builder strings.Builder
	for {
		frame, more := frames.Next()
		if !more {
			break
		}

		// Skip runtime and standard library frames
		if strings.Contains(frame.File, "runtime/") {
			continue
		}

		fmt.Fprintf(&builder, "%s:%d - %s\n", frame.File, frame.Line, frame.Function)
	}

	return builder.String()
}

// Common error constructors for frequently used error types

// NotFound creates a new not found error
func NotFound(message string) Error {
	return New(CodeNotFound, message)
}

// NotFoundf creates a new not found error with a formatted message
func NotFoundf(format string, args ...interface{}) Error {
	return Newf(CodeNotFound, format, args...)
}

// InvalidInput creates a new invalid input error
func InvalidInput(message string) Error {
	return New(CodeInvalidInput, message)
}

// InvalidInputf creates a new invalid input error with a formatted message
func InvalidInputf(format string, args ...interface{}) Error {
	return Newf(CodeInvalidInput, format, args...)
}

// Unauthorized creates a new unauthorized error
func Unauthorized(message string) Error {
	return New(CodeUnauthorized, message)
}

// Unauthorizedf creates a new unauthorized error with a formatted message
func Unauthorizedf(format string, args ...interface{}) Error {
	return Newf(CodeUnauthorized, format, args...)
}

// Internal creates a new internal server error
func Internal(message string) Error {
	return New(CodeInternal, message)
}

// Internalf creates a new internal server error with a formatted message
func Internalf(format string, args ...interface{}) Error {
	return Newf(CodeInternal, format, args...)
}

// Timeout creates a new timeout error
func Timeout(message string) Error {
	return New(CodeTimeout, message)
}

// Timeoutf creates a new timeout error with a formatted message
func Timeoutf(format string, args ...interface{}) Error {
	return Newf(CodeTimeout, format, args...)
}

// Unavailable creates a new service unavailable error
func Unavailable(message string) Error {
	return New(CodeUnavailable, message)
}

// Unavailablef creates a new service unavailable error with a formatted message
func Unavailablef(format string, args ...interface{}) Error {
	return Newf(CodeUnavailable, format, args...)
}

// Conflict creates a new conflict error
func Conflict(message string) Error {
	return New(CodeConflict, message)
}

// Conflictf creates a new conflict error with a formatted message
func Conflictf(format string, args ...interface{}) Error {
	return Newf(CodeConflict, format, args...)
}

// Unauthenticated creates a new unauthenticated error
func Unauthenticated(message string) Error {
	return New(CodeUnauthenticated, message)
}

// Unauthenticatedf creates a new unauthenticated error with a formatted message
func Unauthenticatedf(format string, args ...interface{}) Error {
	return Newf(CodeUnauthenticated, format, args...)
}

// Forbidden creates a new forbidden error
func Forbidden(message string) Error {
	return New(CodeForbidden, message)
}

// Forbiddenf creates a new forbidden error with a formatted message
func Forbiddenf(format string, args ...interface{}) Error {
	return Newf(CodeForbidden, format, args...)
}

// NotImplemented creates a new not implemented error
func NotImplemented(message string) Error {
	return New(CodeNotImplemented, message)
}

// NotImplementedf creates a new not implemented error with a formatted message
func NotImplementedf(format string, args ...interface{}) Error {
	return Newf(CodeNotImplemented, format, args...)
}

// Legacy error constants for backward compatibility
var (
	// UnknownError represents a general undefined error.
	UnknownError = Internal("Unknown error")

	// NotImplementedError indicates that a feature is not implemented.
	NotImplementedError = NotImplemented("Not implemented")

	// UnauthenticatedError indicates an authentication failure.
	UnauthenticatedError = Unauthenticated("Failed auth")

	// UnauthorizedError indicates a failed authorization attempt.
	UnauthorizedError = Unauthorized("Access Denied")

	// InvalidTokenError indicates an invalid authentication token.
	InvalidTokenError = InvalidInput("Invalid token")

	// InvalidRequestError indicates a malformed or invalid request.
	InvalidRequestError = InvalidInput("Invalid request")

	// InvalidDataError indicates a malformed or invalid data.
	InvalidDataError = InvalidInput("Invalid data")

	// InvalidParameterError indicates an invalid parameter in the request.
	InvalidParameterError = InvalidInput("Invalid parameter")

	// TooManyRequestsError indicates that the request rate limit has been exceeded.
	TooManyRequestsError = New("rate_limited", "Too many requests")

	// InvalidHeaderError indicates an issue with an HTTP header.
	InvalidHeaderError = InvalidInput("Invalid header")

	// ServiceUnavailableError indicates that the service is unavailable.
	ServiceUnavailableError = Unavailable("Service unavailable")

	// ItemNotFoundError indicates that a requested item could not be found.
	ItemNotFoundError = NotFound("Item not found")

	// InternalErrorError indicates an internal server error.
	InternalErrorError = Internal("Internal server error")

	// ActionTimeoutError indicates that the action has timed out.
	ActionTimeoutError = Timeout("Action timeout")
)
