package errors_test

import (
	"database/sql"
	"fmt"
	"github.com/gflydev/core/errors"
)

// This file contains examples of how to use the errors package.
// These examples are also used as documentation.

func Example_basic() {
	// Create a basic error
	err := errors.New(errors.CodeNotFound, "User not found")
	fmt.Println(err.Error())
	fmt.Println(err.Code())
	fmt.Println(err.StatusCode())

	// Output:
	// User not found
	// not_found
	// 404
}

func Example_formatted() {
	// Create an error with a formatted message
	err := errors.Newf(errors.CodeInvalidInput, "Invalid value for field %s: %v", "email", "not-an-email")
	fmt.Println(err.Error())
	fmt.Println(err.Code())
	fmt.Println(err.StatusCode())

	// Output:
	// Invalid value for field email: not-an-email
	// invalid_input
	// 400
}

func Example_helpers() {
	// Use helper functions for common error types
	notFoundErr := errors.NotFound("User not found")
	fmt.Println(notFoundErr.Error())
	fmt.Println(notFoundErr.Code())
	fmt.Println(notFoundErr.StatusCode())

	invalidInputErr := errors.InvalidInput("Invalid email address")
	fmt.Println(invalidInputErr.Error())
	fmt.Println(invalidInputErr.Code())
	fmt.Println(invalidInputErr.StatusCode())

	unauthorizedErr := errors.Unauthorized("Access denied")
	fmt.Println(unauthorizedErr.Error())
	fmt.Println(unauthorizedErr.Code())
	fmt.Println(unauthorizedErr.StatusCode())

	// Output:
	// User not found
	// not_found
	// 404
	// Invalid email address
	// invalid_input
	// 400
	// Access denied
	// unauthorized
	// 401
}

func Example_wrapping() {
	// Create an original error
	originalErr := sql.ErrNoRows

	// Wrap the error
	wrappedErr := errors.Wrap(originalErr, errors.CodeNotFound, "User not found")
	fmt.Println(wrappedErr.Error())
	fmt.Println(wrappedErr.Code())
	fmt.Println(wrappedErr.StatusCode())

	// Wrap with a formatted message
	formattedErr := errors.Wrapf(originalErr, errors.CodeInvalidInput, "Invalid input for field %s", "email")
	fmt.Println(formattedErr.Error())
	fmt.Println(formattedErr.Code())
	fmt.Println(formattedErr.StatusCode())

	// Output:
	// User not found: sql: no rows in result set
	// not_found
	// 404
	// Invalid input for field email: sql: no rows in result set
	// invalid_input
	// 400
}

func Example_contextAndMetadata() {
	// Add context information
	err := errors.InvalidInput("Invalid email address").
		WithContext("user_id", 123).
		WithContext("request_id", "req-456")

	// Add metadata
	err = err.WithMetadata("field", "email").
		WithMetadata("value", "not-an-email")

	// Retrieve context and metadata
	userId, ok := err.GetContext("user_id")
	fmt.Printf("User ID: %v (found: %v)\n", userId, ok)

	requestId, ok := err.GetContext("request_id")
	fmt.Printf("Request ID: %v (found: %v)\n", requestId, ok)

	field, ok := err.GetMetadata("field")
	fmt.Printf("Field: %v (found: %v)\n", field, ok)

	value, ok := err.GetMetadata("value")
	fmt.Printf("Value: %v (found: %v)\n", value, ok)

	// Output:
	// User ID: 123 (found: true)
	// Request ID: req-456 (found: true)
	// Field: email (found: true)
	// Value: not-an-email (found: true)
}

func Example_fileNotFound() {
	// Create a FileNotFound error
	err := errors.NewFileNotFound("config.json", "/etc/myapp")
	fmt.Println(err.Error())
	fmt.Println(err.Code())
	fmt.Println(err.StatusCode())

	// Add metadata
	err = err.WithMetadata("required", true)
	required, ok := err.GetMetadata("required")
	fmt.Printf("Required: %v (found: %v)\n", required, ok)

	// Output:
	// File config.json not found at location /etc/myapp
	// not_found
	// 404
	// Required: true (found: true)
}

func Example_legacySupport() {
	// Use legacy error constants
	fmt.Println(errors.ItemNotFoundError.Error())
	fmt.Println(errors.ItemNotFoundError.Code())
	fmt.Println(errors.ItemNotFoundError.StatusCode())

	fmt.Println(errors.InvalidRequestError.Error())
	fmt.Println(errors.InvalidRequestError.Code())
	fmt.Println(errors.InvalidRequestError.StatusCode())

	// Output:
	// Item not found
	// not_found
	// 404
	// Invalid request
	// invalid_input
	// 400
}
