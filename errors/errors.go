package errors

import (
	sysError "errors"
	"fmt"
)

// New returns a formatted error message based on the given format and arguments.
//
// Parameters:
//   - format: String that contains the error message with optional format verbs
//   - args: Optional arguments to be formatted according to the format string
//
// Returns:
//   - error: Formatted error object containing the message
func New(format string, args ...any) error {
	return fmt.Errorf(format, args...)
}

// Is reports whether any error in err's tree matches target.
//
// The tree consists of err itself, followed by the errors obtained by repeatedly
// calling its Unwrap() error or Unwrap() []error method. When err wraps multiple
// errors, Is examines err followed by a depth-first traversal of its children.
//
// An error is considered to match a target if it is equal to that target or if
// it implements a method Is(error) bool such that Is(target) returns true.
//
// An error type might provide an Is method so it can be treated as equivalent
// to an existing error. For example, if MyError defines
//
//	func (m MyError) Is(target error) bool { return target == fs.ErrExist }
//
// then Is(MyError{}, fs.ErrExist) returns true. See [syscall.Errno.Is] for
// an example in the standard library. An Is method should only shallowly
// compare err and the target and not call [Unwrap] on either.
func Is(err, target error) bool {
	return sysError.Is(err, target)
}

// ToStr formats and returns an error string using the provided error format string and arguments.
//
// Note: The method should be used by custom error is not a string type.
//
// Parameters:
//   - format: Format string that may contain format specifiers
//   - args: Optional arguments to be formatted according to the format string
//
// Returns:
//   - string: Formatted string with arguments substituted into format
func ToStr(format string, args ...any) string {
	return fmt.Sprintf(format, args...)
}

var (
	// UnknownError represents a general undefined error.
	UnknownError = New("Unknown error")

	// NotImplemented indicates that a feature is not implemented.
	NotImplemented = New("Not implemented")

	// Unauthenticated indicates an authentication failure.
	// Authentication - Who are you?
	Unauthenticated = New("Failed auth")

	// Unauthorized (Forbidden) indicates a failed authorization attempt.
	// Authorization - Are you allowed to do that?
	Unauthorized = New("Access Denied")

	// InvalidToken indicates an invalid authentication token.
	InvalidToken = New("Invalid token")

	// InvalidRequest indicates a malformed or invalid request.
	InvalidRequest = New("Invalid request")

	// InvalidData indicates a malformed or invalid data.
	InvalidData = New("Invalid data")

	// InvalidParameter indicates an invalid parameter in the request.
	InvalidParameter = New("Invalid parameter")

	// TooManyRequests indicates that the request rate limit has been exceeded.
	TooManyRequests = New("Too many requests")

	// InvalidHeader indicates an issue with an HTTP header.
	InvalidHeader = New("Invalid header")

	// ServiceUnavailable indicates that the service is unavailable.
	ServiceUnavailable = New("Service unavailable")

	// ItemNotFound indicates that a requested item could not be found.
	ItemNotFound = New("Item not found")

	// InternalError indicates an internal server error.
	InternalError = New("Internal server error")

	// ActionTimeout indicates that the action has timed out.
	ActionTimeout = New("Action timeout")
)
