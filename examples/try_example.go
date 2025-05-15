package main

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/gflydev/core/errors"
	"github.com/gflydev/core/try"
)

// This example demonstrates how to use the try package in a real application context.
// It shows how to use the try-catch-finally pattern for error handling and recovery.

func main() {
	fmt.Println("gFly Core Try Package Examples")
	fmt.Println("============================")

	// Example 1: Basic try-catch
	fmt.Println("\nExample 1: Basic Try-Catch")
	fmt.Println("------------------------")
	basicTryCatchExample()

	// Example 2: Try-catch-finally
	fmt.Println("\nExample 2: Try-Catch-Finally")
	fmt.Println("--------------------------")
	tryCatchFinallyExample()

	// Example 3: Nested try-catch
	fmt.Println("\nExample 3: Nested Try-Catch")
	fmt.Println("-------------------------")
	nestedTryCatchExample()

	// Example 4: Custom error handling
	fmt.Println("\nExample 4: Custom Error Handling")
	fmt.Println("------------------------------")
	customErrorHandlingExample()

	// Example 5: File operations with try-catch
	fmt.Println("\nExample 5: File Operations with Try-Catch")
	fmt.Println("--------------------------------------")
	fileOperationsExample()

	// Example 6: Throwing errors
	fmt.Println("\nExample 6: Throwing Errors")
	fmt.Println("-----------------------")
	throwingErrorsExample()

	// Example 7: Error type checking
	fmt.Println("\nExample 7: Error Type Checking")
	fmt.Println("---------------------------")
	errorTypeCheckingExample()
}

func basicTryCatchExample() {
	// Basic try-catch example
	try.Perform(func() {
		// This is the "try" block
		fmt.Println("Trying to divide 10 by 0...")

		// This will cause a panic
		divideByZero(10, 0)

		// This line will not be executed
		fmt.Println("This line will not be executed")
	}).Catch(func(err try.E) {
		// This is the "catch" block
		fmt.Printf("Caught an error: %v\n", err)
	})

	fmt.Println("Program continues after the try-catch block")
}

func tryCatchFinallyExample() {
	// Try-catch-finally example
	try.Perform(func() {
		// This is the "try" block
		fmt.Println("Trying to parse a non-numeric string...")

		// This will cause an error
		_, err := strconv.Atoi("not-a-number")
		if err != nil {
			// Throw the error to be caught
			try.Throw(err)
		}

		// This line will not be executed
		fmt.Println("This line will not be executed")
	}).Catch(func(err try.E) {
		// This is the "catch" block
		fmt.Printf("Caught an error: %v\n", err)
	}).Finally(func() {
		// This is the "finally" block that always executes
		fmt.Println("Finally block executed")
	})

	fmt.Println("Program continues after the try-catch-finally block")
}

func nestedTryCatchExample() {
	// Nested try-catch example
	try.Perform(func() {
		// Outer try block
		fmt.Println("Outer try block")

		try.Perform(func() {
			// Inner try block
			fmt.Println("Inner try block")

			// This will cause an error
			_, err := strconv.Atoi("not-a-number")
			if err != nil {
				// Throw the error to be caught by the inner catch
				try.Throw(err)
			}

			// This line will not be executed
			fmt.Println("This line will not be executed")
		}).Catch(func(err try.E) {
			// Inner catch block
			fmt.Printf("Inner catch block caught an error: %v\n", err)

			// Rethrow the error to be caught by the outer catch
			try.Throw(fmt.Errorf("rethrown from inner catch: %w", err))
		})

		// This line will not be executed because we rethrew the error
		fmt.Println("This line will not be executed")
	}).Catch(func(err try.E) {
		// Outer catch block
		fmt.Printf("Outer catch block caught an error: %v\n", err)
	})

	fmt.Println("Program continues after the nested try-catch blocks")
}

func customErrorHandlingExample() {
	// Custom error handling example
	try.Perform(func() {
		// Try to validate user input
		username := ""
		email := "not-an-email"

		// Validate username
		if username == "" {
			// Throw a custom error
			try.Throw(errors.InvalidInput("Username cannot be empty"))
		}

		// Validate email
		if !isValidEmail(email) {
			// Throw a custom error
			try.Throw(errors.InvalidInput("Invalid email format"))
		}

		// This line will not be executed
		fmt.Println("User validation successful")
	}).Catch(func(err try.E) {
		// Check if it's a gFly error
		if gflyErr, ok := err.(errors.Error); ok {
			fmt.Printf("Validation error: %s (Code: %s)\n", gflyErr.Message(), gflyErr.Code())
		} else {
			fmt.Printf("Unknown error: %v\n", err)
		}
	})

	fmt.Println("Program continues after the custom error handling")
}

func fileOperationsExample() {
	// Create a temporary file for demonstration
	tempFile, err := os.CreateTemp("", "try-example-*.txt")
	if err != nil {
		fmt.Printf("Failed to create temp file: %v\n", err)
		return
	}
	tempFilePath := tempFile.Name()
	tempFile.Close()
	defer os.Remove(tempFilePath)

	// Write to the file using try-catch
	try.Perform(func() {
		fmt.Printf("Writing to file: %s\n", tempFilePath)

		// Open the file for writing
		file, err := os.OpenFile(tempFilePath, os.O_WRONLY, 0644)
		if err != nil {
			try.Throw(err)
		}
		defer file.Close()

		// Write some data
		_, err = file.WriteString("Hello, World!")
		if err != nil {
			try.Throw(err)
		}

		fmt.Println("Successfully wrote to the file")
	}).Catch(func(err try.E) {
		fmt.Printf("Failed to write to file: %v\n", err)
	})

	// Read from the file using try-catch-finally
	var fileContent string
	try.Perform(func() {
		fmt.Printf("Reading from file: %s\n", tempFilePath)

		// Open the file for reading
		file, err := os.Open(tempFilePath)
		if err != nil {
			try.Throw(err)
		}
		defer file.Close()

		// Read the content
		content, err := io.ReadAll(file)
		if err != nil {
			try.Throw(err)
		}

		fileContent = string(content)
		fmt.Printf("File content: %s\n", fileContent)
	}).Catch(func(err try.E) {
		fmt.Printf("Failed to read from file: %v\n", err)
	}).Finally(func() {
		fmt.Println("File operation completed")
	})
}

func throwingErrorsExample() {
	// Example of throwing different types of errors
	try.Perform(func() {
		// Decide which error to throw based on a condition
		errorType := "not_found"

		switch errorType {
		case "not_found":
			try.Throw(errors.NotFound("Resource not found"))
		case "invalid_input":
			try.Throw(errors.InvalidInput("Invalid input data"))
		case "unauthorized":
			try.Throw(errors.Unauthorized("Unauthorized access"))
		case "internal":
			try.Throw(errors.Internal("Internal server error"))
		default:
			try.Throw(fmt.Errorf("unknown error type: %s", errorType))
		}

		// This line will not be executed
		fmt.Println("This line will not be executed")
	}).Catch(func(err try.E) {
		fmt.Printf("Caught an error: %v\n", err)
	})
}

func errorTypeCheckingExample() {
	// Example of checking error types in the catch block
	try.Perform(func() {
		// Throw a not found error
		try.Throw(errors.NotFound("User not found"))
	}).Catch(func(err try.E) {
		// Check the error type
		if gflyErr, ok := err.(errors.Error); ok {
			switch gflyErr.Code() {
			case errors.CodeNotFound:
				fmt.Println("Handling not found error")
			case errors.CodeInvalidInput:
				fmt.Println("Handling invalid input error")
			case errors.CodeUnauthorized:
				fmt.Println("Handling unauthorized error")
			case errors.CodeInternal:
				fmt.Println("Handling internal server error")
			default:
				fmt.Printf("Handling unknown error code: %s\n", gflyErr.Code())
			}
		} else {
			fmt.Printf("Handling non-gFly error: %v\n", err)
		}
	})
}

// Helper functions

func divideByZero(a, b int) int {
	return a / b
}

func isValidEmail(email string) bool {
	// This is a simplified email validation
	// In a real application, you would use a more robust validation
	return len(email) > 0 && (email[len(email)-4:] == ".com" || email[len(email)-4:] == ".org" || email[len(email)-3:] == ".io")
}
