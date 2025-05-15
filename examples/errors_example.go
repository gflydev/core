package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gflydev/core/errors"
)

// This example demonstrates how to use the errors package in a real application context.
// It shows how to create, wrap, and handle errors with different error codes and metadata.

func main() {
	fmt.Println("gFly Core Errors Package Examples")
	fmt.Println("=================================")

	// Example 1: Basic error creation and handling
	fmt.Println("\nExample 1: Basic Error Creation and Handling")
	fmt.Println("------------------------------------------")
	basicErrorExample()

	// Example 2: Error wrapping
	fmt.Println("\nExample 2: Error Wrapping")
	fmt.Println("------------------------")
	errorWrappingExample()

	// Example 3: Error context and metadata
	fmt.Println("\nExample 3: Error Context and Metadata")
	fmt.Println("-----------------------------------")
	errorContextExample()

	// Example 4: HTTP error handling
	fmt.Println("\nExample 4: HTTP Error Handling")
	fmt.Println("----------------------------")
	httpErrorHandlingExample()

	// Example 5: File not found error
	fmt.Println("\nExample 5: File Not Found Error")
	fmt.Println("-----------------------------")
	fileNotFoundExample()
}

func basicErrorExample() {
	// Create a basic error
	err := errors.New(errors.CodeNotFound, "User not found")
	printErrorDetails(err)

	// Create an error with a formatted message
	err = errors.Newf(errors.CodeInvalidInput, "Invalid value for field %s: %v", "email", "not-an-email")
	printErrorDetails(err)

	// Use helper functions for common error types
	err = errors.NotFound("User with ID 123 not found")
	printErrorDetails(err)

	err = errors.InvalidInput("Invalid email address format")
	printErrorDetails(err)

	err = errors.Unauthorized("Invalid API key")
	printErrorDetails(err)

	err = errors.Forbidden("User does not have permission to access this resource")
	printErrorDetails(err)

	err = errors.Internal("Database connection failed")
	printErrorDetails(err)
}

func errorWrappingExample() {
	// Simulate a database error
	originalErr := sql.ErrNoRows

	// Wrap the error with domain-specific context
	wrappedErr := errors.Wrap(originalErr, errors.CodeNotFound, "User not found in database")
	printErrorDetails(wrappedErr)

	// Wrap with a formatted message
	formattedErr := errors.Wrapf(originalErr, errors.CodeInvalidInput, "Invalid input for user ID %d", 123)
	printErrorDetails(formattedErr)

	// Demonstrate error unwrapping
	fmt.Println("Original error through unwrapping:", wrappedErr.Unwrap())

	// Check if the wrapped error is of a specific type
	if errors.Is(wrappedErr, sql.ErrNoRows) {
		fmt.Println("The wrapped error is sql.ErrNoRows")
	}
}

func errorContextExample() {
	// Create an error with context information
	err := errors.InvalidInput("Invalid user data").
		WithContext("user_id", 123).
		WithContext("request_id", "req-456").
		WithContext("timestamp", "2023-04-01T12:34:56Z")

	// Add metadata
	err = err.WithMetadata("field", "email").
		WithMetadata("value", "not-an-email").
		WithMetadata("validation_rule", "email_format")

	printErrorDetails(err)

	// Retrieve and display context and metadata
	fmt.Println("Context values:")
	if userId, ok := err.GetContext("user_id"); ok {
		fmt.Printf("  - user_id: %v\n", userId)
	}
	if requestId, ok := err.GetContext("request_id"); ok {
		fmt.Printf("  - request_id: %v\n", requestId)
	}
	if timestamp, ok := err.GetContext("timestamp"); ok {
		fmt.Printf("  - timestamp: %v\n", timestamp)
	}

	fmt.Println("Metadata values:")
	if field, ok := err.GetMetadata("field"); ok {
		fmt.Printf("  - field: %v\n", field)
	}
	if value, ok := err.GetMetadata("value"); ok {
		fmt.Printf("  - value: %v\n", value)
	}
	if rule, ok := err.GetMetadata("validation_rule"); ok {
		fmt.Printf("  - validation_rule: %v\n", rule)
	}
}

func httpErrorHandlingExample() {
	// Simulate HTTP handlers that return different errors
	errors := []errors.Error{
		errors.NotFound("User not found"),
		errors.InvalidInput("Invalid email format"),
		errors.Unauthorized("Invalid API key"),
		errors.Forbidden("Insufficient permissions"),
		errors.Internal("Database connection failed"),
	}

	// Handle each error as if it were returned from an HTTP handler
	for _, err := range errors {
		fmt.Printf("HTTP %d response for error: %s\n", err.StatusCode(), err.Error())

		// Example of how you might use this in an HTTP handler
		fmt.Println("Example response body:")
		fmt.Printf("  {\n    \"error\": \"%s\",\n    \"code\": \"%s\",\n    \"status\": %d\n  }\n",
			err.Error(), err.Code(), err.StatusCode())
	}
}

func fileNotFoundExample() {
	// Create a temporary file for demonstration
	tempFile, err := os.CreateTemp("", "example")
	if err != nil {
		log.Fatalf("Failed to create temp file: %v", err)
	}
	tempFilePath := tempFile.Name()
	tempFile.Close()

	// Remove the file to simulate a "not found" scenario
	os.Remove(tempFilePath)

	// Try to open the file that doesn't exist
	_, err = os.Open(tempFilePath)
	if err != nil {
		// Create a domain-specific error
		notFoundErr := errors.NewFileNotFound(tempFilePath, "")
		printErrorDetails(notFoundErr)

		// Add metadata
		notFoundErr = notFoundErr.WithMetadata("required", true)
		notFoundErr = notFoundErr.WithMetadata("file_type", "configuration")

		// Retrieve and display metadata
		fmt.Println("Metadata values:")
		if required, ok := notFoundErr.GetMetadata("required"); ok {
			fmt.Printf("  - required: %v\n", required)
		}
		if fileType, ok := notFoundErr.GetMetadata("file_type"); ok {
			fmt.Printf("  - file_type: %v\n", fileType)
		}
	}
}

// Helper function to print error details
func printErrorDetails(err errors.Error) {
	fmt.Printf("Error: %s\n", err.Error())
	fmt.Printf("Code: %s\n", err.Code())
	fmt.Printf("Status Code: %d\n", err.StatusCode())
	fmt.Println()
}

// Example of how to use errors in an HTTP handler
func exampleHTTPHandler(w http.ResponseWriter, r *http.Request) {
	// Simulate a user lookup that fails
	userID := r.URL.Query().Get("id")
	if userID == "" {
		err := errors.InvalidInput("User ID is required")
		w.WriteHeader(err.StatusCode())
		fmt.Fprintf(w, "{\"error\":\"%s\",\"code\":\"%s\"}", err.Error(), err.Code())
		return
	}

	// Simulate a database lookup
	user, err := findUser(userID)
	if err != nil {
		var apiErr errors.Error
		if errors.As(err, &apiErr) {
			// If it's already our custom error type
			w.WriteHeader(apiErr.StatusCode())
			fmt.Fprintf(w, "{\"error\":\"%s\",\"code\":\"%s\"}", apiErr.Error(), apiErr.Code())
		} else {
			// Wrap unknown errors as internal server errors
			wrappedErr := errors.Wrap(err, errors.CodeInternal, "Failed to find user")
			w.WriteHeader(wrappedErr.StatusCode())
			fmt.Fprintf(w, "{\"error\":\"%s\",\"code\":\"%s\"}", wrappedErr.Error(), wrappedErr.Code())
		}
		return
	}

	// Success case
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "{\"user\":%s}", user)
}

// Simulate a user lookup function that might return errors
func findUser(id string) (string, error) {
	// Simulate different error scenarios based on the ID
	switch id {
	case "1":
		return "{\"id\":1,\"name\":\"John Doe\"}", nil
	case "404":
		return "", errors.NotFound("User not found")
	case "401":
		return "", errors.Unauthorized("Authentication required")
	case "403":
		return "", errors.Forbidden("Access denied")
	case "500":
		return "", errors.Internal("Database error")
	default:
		// Simulate a database error that we'll wrap
		return "", fmt.Errorf("database query failed: %w", sql.ErrNoRows)
	}
}
