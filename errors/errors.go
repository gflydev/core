package errors

import "fmt"

// Error type represents a string-based error.
type Error string

// Error returns the error message for the Error type.
func (e Error) Error() string {
	return string(e)
}

// New returns a formatted error message based on the given format and arguments.
func New(format string, args ...any) error {
	return fmt.Errorf(format, args...)
}

const (
	// UnknownError represents a general undefined error.
	UnknownError = Error("Unknown error")

	// NotImplemented indicates that a feature is not implemented.
	NotImplemented = Error("Not implemented")

	// Unauthenticated indicates an authentication failure.
	// Authentication - Who are you?
	Unauthenticated = Error("Failed auth")

	// Unauthorized (Forbidden) indicates a failed authorization attempt.
	// Authorization - Are you allowed to do that?
	Unauthorized = Error("Access Denied")

	// InvalidToken indicates an invalid authentication token.
	InvalidToken = Error("Invalid token")

	// InvalidRequest indicates a malformed or invalid request.
	InvalidRequest = Error("Invalid request")

	// InvalidData indicates a malformed or invalid data.
	InvalidData = Error("Invalid data")

	// InvalidParameter indicates an invalid parameter in the request.
	InvalidParameter = Error("Invalid parameter")

	// TooManyRequests indicates that the request rate limit has been exceeded.
	TooManyRequests = Error("Too many requests")

	// InvalidHeader indicates an issue with an HTTP header.
	InvalidHeader = Error("Invalid header")

	// ServiceUnavailable indicates that the service is unavailable.
	ServiceUnavailable = Error("Service unavailable")

	// ItemNotFound indicates that a requested item could not be found.
	ItemNotFound = Error("Item not found")

	// InternalError indicates an internal server error.
	InternalError = Error("Internal server error")

	// ActionTimeout indicates that the action has timed out.
	ActionTimeout = Error("Action timeout")
)
