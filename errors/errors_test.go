package errors

import (
	"fmt"
	"github.com/gflydev/core/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := map[string]struct {
		code     string
		message  string
		expected string
	}{
		"Test New Error with code": {
			code:     CodeInvalidInput,
			message:  "Invalid input provided",
			expected: "Invalid input provided",
		},
		"Test New Error with unknown code": {
			code:     CodeUnknown,
			message:  "Unknown error occurred",
			expected: "Unknown error occurred",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := New(tt.code, tt.message)
			require.Equal(t, tt.expected, err.Error())
			require.Equal(t, tt.code, err.Code())
			require.NotEmpty(t, err.StackTrace())
		})
	}
}

func TestNewf(t *testing.T) {
	tests := map[string]struct {
		code     string
		format   string
		args     []interface{}
		expected string
	}{
		"Test Newf with no args": {
			code:     CodeInvalidInput,
			format:   "Invalid input provided",
			args:     nil,
			expected: "Invalid input provided",
		},
		"Test Newf with one argument": {
			code:     CodeNotFound,
			format:   "File `%s` not found",
			args:     []interface{}{"my-file.pdf"},
			expected: "File `my-file.pdf` not found",
		},
		"Test Newf with multiple arguments": {
			code:     CodeInvalidInput,
			format:   "File `%s` is not `%s` type",
			args:     []interface{}{"my-file.pdf", "png"},
			expected: "File `my-file.pdf` is not `png` type",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := Newf(tt.code, tt.format, tt.args...)
			require.Equal(t, tt.expected, err.Error())
			require.Equal(t, tt.code, err.Code())
			require.NotEmpty(t, err.StackTrace())
		})
	}
}

func TestWrap(t *testing.T) {
	// Create an original error
	originalErr := New(CodeNotFound, "Original error")

	// Wrap the error
	wrappedErr := Wrap(originalErr, CodeInvalidInput, "Wrapped error")

	// Test the wrapped error
	require.Equal(t, "Wrapped error: Original error", wrappedErr.Error())
	require.Equal(t, CodeInvalidInput, wrappedErr.Code())
	require.Equal(t, originalErr, wrappedErr.Unwrap())

	// Test wrapping nil
	require.Nil(t, Wrap(nil, CodeInvalidInput, "Wrapped nil"))

	// Test wrapping with empty code (should use original code)
	preservedCodeErr := Wrap(originalErr, "", "Preserved code")
	require.Equal(t, CodeNotFound, preservedCodeErr.Code())
}

func TestWrapf(t *testing.T) {
	// Create an original error
	originalErr := New(CodeNotFound, "Original error")

	// Wrap the error with formatting
	wrappedErr := Wrapf(originalErr, CodeInvalidInput, "Wrapped error: %s", "formatted")

	// Test the wrapped error
	require.Equal(t, "Wrapped error: formatted: Original error", wrappedErr.Error())
	require.Equal(t, CodeInvalidInput, wrappedErr.Code())
	require.Equal(t, originalErr, wrappedErr.Unwrap())
}

func TestErrorHelpers(t *testing.T) {
	tests := map[string]struct {
		errorFunc func(string) Error
		code      string
		message   string
	}{
		"NotFound": {
			errorFunc: NotFound,
			code:      CodeNotFound,
			message:   "Item not found",
		},
		"InvalidInput": {
			errorFunc: InvalidInput,
			code:      CodeInvalidInput,
			message:   "Invalid input",
		},
		"Unauthorized": {
			errorFunc: Unauthorized,
			code:      CodeUnauthorized,
			message:   "Unauthorized",
		},
		"Internal": {
			errorFunc: Internal,
			code:      CodeInternal,
			message:   "Internal error",
		},
		"Timeout": {
			errorFunc: Timeout,
			code:      CodeTimeout,
			message:   "Timeout occurred",
		},
		"Unavailable": {
			errorFunc: Unavailable,
			code:      CodeUnavailable,
			message:   "Service unavailable",
		},
		"Conflict": {
			errorFunc: Conflict,
			code:      CodeConflict,
			message:   "Conflict occurred",
		},
		"Unauthenticated": {
			errorFunc: Unauthenticated,
			code:      CodeUnauthenticated,
			message:   "Not authenticated",
		},
		"Forbidden": {
			errorFunc: Forbidden,
			code:      CodeForbidden,
			message:   "Forbidden",
		},
		"NotImplemented": {
			errorFunc: NotImplemented,
			code:      CodeNotImplemented,
			message:   "Not implemented",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := tt.errorFunc(tt.message)
			require.Equal(t, tt.message, err.Error())
			require.Equal(t, tt.code, err.Code())
		})
	}
}

func TestMetadataAndContext(t *testing.T) {
	// Create a base error
	err := New(CodeInvalidInput, "Invalid input")

	// Add metadata
	errWithMetadata := err.WithMetadata("key1", "value1")

	// Verify metadata
	value, ok := errWithMetadata.GetMetadata("key1")
	require.True(t, ok)
	require.Equal(t, "value1", value)

	// Add context
	errWithContext := errWithMetadata.WithContext("request_id", "12345")

	// Verify context
	requestID, ok := errWithContext.GetContext("request_id")
	require.True(t, ok)
	require.Equal(t, "12345", requestID)

	// Verify original metadata is preserved
	value, ok = errWithContext.GetMetadata("key1")
	require.True(t, ok)
	require.Equal(t, "value1", value)

	// Verify non-existent keys
	_, ok = errWithContext.GetMetadata("non_existent")
	require.False(t, ok)
	_, ok = errWithContext.GetContext("non_existent")
	require.False(t, ok)
}

func TestFileNotFound(t *testing.T) {
	// Create a FileNotFound error
	err := NewFileNotFound("my-file.pdf", "/tmp")

	// Test basic properties
	require.Equal(t, "File my-file.pdf not found at location /tmp", err.Error())
	require.Equal(t, CodeNotFound, err.Code())

	// Test as Error interface
	var gflyErr Error = err
	require.Equal(t, CodeNotFound, gflyErr.Code())

	// Test with metadata
	errWithMetadata := err.WithMetadata("size", 1024)
	size, ok := errWithMetadata.GetMetadata("size")
	require.True(t, ok)
	require.Equal(t, 1024, size)

	// Test with context
	errWithContext := err.WithContext("user", "admin")
	user, ok := errWithContext.GetContext("user")
	require.True(t, ok)
	require.Equal(t, "admin", user)
}

func TestLegacyConstants(t *testing.T) {
	// Test that legacy constants implement the Error interface
	assert.Implements(t, (*Error)(nil), UnknownError)
	assert.Implements(t, (*Error)(nil), NotImplementedError)
	assert.Implements(t, (*Error)(nil), UnauthenticatedError)
	assert.Implements(t, (*Error)(nil), UnauthorizedError)
	assert.Implements(t, (*Error)(nil), InvalidTokenError)
	assert.Implements(t, (*Error)(nil), InvalidRequestError)
	assert.Implements(t, (*Error)(nil), InvalidDataError)
	assert.Implements(t, (*Error)(nil), InvalidParameterError)
	assert.Implements(t, (*Error)(nil), TooManyRequestsError)
	assert.Implements(t, (*Error)(nil), InvalidHeaderError)
	assert.Implements(t, (*Error)(nil), ServiceUnavailableError)
	assert.Implements(t, (*Error)(nil), ItemNotFoundError)
	assert.Implements(t, (*Error)(nil), InternalErrorError)
	assert.Implements(t, (*Error)(nil), ActionTimeoutError)

	// Test specific error codes
	assert.Equal(t, CodeInternal, UnknownError.Code())
	assert.Equal(t, CodeNotImplemented, NotImplementedError.Code())
	assert.Equal(t, CodeUnauthenticated, UnauthenticatedError.Code())
	assert.Equal(t, CodeUnauthorized, UnauthorizedError.Code())
	assert.Equal(t, CodeInvalidInput, InvalidTokenError.Code())
	assert.Equal(t, CodeInvalidInput, InvalidRequestError.Code())
	assert.Equal(t, CodeInvalidInput, InvalidDataError.Code())
	assert.Equal(t, CodeInvalidInput, InvalidParameterError.Code())
	assert.Equal(t, "rate_limited", TooManyRequestsError.Code())
	assert.Equal(t, CodeInvalidInput, InvalidHeaderError.Code())
	assert.Equal(t, CodeUnavailable, ServiceUnavailableError.Code())
	assert.Equal(t, CodeNotFound, ItemNotFoundError.Code())
	assert.Equal(t, CodeInternal, InternalErrorError.Code())
	assert.Equal(t, CodeTimeout, ActionTimeoutError.Code())
}

// Legacy test for backward compatibility
func Test_LegacyNew(t *testing.T) {
	tests := map[string]struct {
		format   string
		args     any
		expected string
		isEqual  bool
	}{
		"Test New Error string": {
			format:   "empty file",
			args:     nil,
			expected: "not empty file",
			isEqual:  false,
		},
		"Test New Error one argument": {
			format:   "file `%s` not found",
			args:     "my-file.pdf",
			expected: "file `my-file.pdf` not found",
			isEqual:  true,
		},
		"Test New Error multi string arguments": {
			format:   "file `%s` is not `%s` type",
			args:     []string{"my-file.pdf", "png"},
			expected: "file `my-file.pdf` is not `png` type",
			isEqual:  true,
		},
		"Test New Error multi number arguments": {
			format:   "file1's size `%v` and file2's size `%v`",
			args:     []int{123, 324},
			expected: "file1's size `123` and file2's size `324`",
			isEqual:  true,
		},
		"Test file not found": {
			format:   "",
			args:     NewFileNotFound("my-file.pdf", "/tmp"),
			expected: "File my-file.pdf not found at location /tmp",
			isEqual:  true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var err error

			// How do I determine whether an array contains
			switch tt.args.(type) {
			case []string, []int:
				err = fmt.Errorf(tt.format, utils.UnpackArray(tt.args)...)
			case Error:
				err = tt.args.(Error)
			default:
				err = fmt.Errorf(tt.format, tt.args)
			}

			if tt.isEqual {
				require.Equal(t, tt.expected, err.Error())
			} else {
				require.NotEqual(t, tt.expected, err.Error())
			}
		})
	}
}
