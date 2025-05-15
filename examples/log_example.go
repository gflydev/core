package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gflydev/core/log"
)

// This example demonstrates how to use the log package in a real application context.
// It shows different logging levels, formatted logging, structured logging, and how to
// configure the logger.

func main() {
	fmt.Println("gFly Core Log Package Examples")
	fmt.Println("=============================")

	// Example 1: Basic logging
	fmt.Println("\nExample 1: Basic Logging")
	fmt.Println("----------------------")
	basicLoggingExample()

	// Example 2: Formatted logging
	fmt.Println("\nExample 2: Formatted Logging")
	fmt.Println("--------------------------")
	formattedLoggingExample()

	// Example 3: Structured logging
	fmt.Println("\nExample 3: Structured Logging")
	fmt.Println("---------------------------")
	structuredLoggingExample()

	// Example 4: Logging with context
	fmt.Println("\nExample 4: Logging with Context")
	fmt.Println("-----------------------------")
	contextLoggingExample()

	// Example 5: Configuring the logger
	fmt.Println("\nExample 5: Configuring the Logger")
	fmt.Println("------------------------------")
	configuringLoggerExample()

	// Example 6: Custom log output
	fmt.Println("\nExample 6: Custom Log Output")
	fmt.Println("-------------------------")
	customOutputExample()

	// Example 7: Error logging
	fmt.Println("\nExample 7: Error Logging")
	fmt.Println("----------------------")
	errorLoggingExample()

	// Example 8: JSON structured logging
	fmt.Println("\nExample 8: JSON Structured Logging")
	fmt.Println("--------------------------------")
	jsonStructuredLoggingExample()
}

func basicLoggingExample() {
	// Get the default logger
	logger := log.DefaultLogger()

	// Log messages at different levels
	logger.Trace("This is a trace message")
	logger.Debug("This is a debug message")
	logger.Info("This is an info message")
	logger.Warn("This is a warning message")
	logger.Error("This is an error message")

	// Note: Fatal and Panic are commented out as they would terminate the program
	// logger.Fatal("This is a fatal message")
	// logger.Panic("This is a panic message")

	fmt.Println("Basic logging demonstrated. Check the log output.")
}

func formattedLoggingExample() {
	// Get the default logger
	logger := log.DefaultLogger()

	// Log formatted messages at different levels
	logger.Tracef("User %s logged in from IP %s", "john_doe", "192.168.1.1")
	logger.Debugf("Processing request #%d from user %s", 12345, "john_doe")
	logger.Infof("Request completed in %.2f ms", 42.5)
	logger.Warnf("High latency detected: %.2f ms", 150.3)
	logger.Errorf("Failed to process request #%d: %s", 12345, "timeout")

	// Note: Fatal and Panic are commented out as they would terminate the program
	// logger.Fatalf("System failure: %s", "database connection lost")
	// logger.Panicf("Critical error: %s", "out of memory")

	fmt.Println("Formatted logging demonstrated. Check the log output.")
}

func structuredLoggingExample() {
	// Get the default logger
	logger := log.DefaultLogger()

	// Log structured messages at different levels
	logger.Tracew("User login",
		"user_id", 123,
		"username", "john_doe",
		"ip", "192.168.1.1",
		"timestamp", time.Now().Format(time.RFC3339),
	)

	logger.Debugw("Request processing",
		"request_id", "req-456",
		"user_id", 123,
		"endpoint", "/api/users",
		"method", "GET",
		"duration_ms", 42.5,
	)

	logger.Infow("Request completed",
		"request_id", "req-456",
		"status", 200,
		"duration_ms", 42.5,
	)

	logger.Warnw("High latency detected",
		"request_id", "req-456",
		"duration_ms", 150.3,
		"threshold_ms", 100,
	)

	logger.Errorw("Request failed",
		"request_id", "req-456",
		"status", 500,
		"error", "database timeout",
		"duration_ms", 3000.5,
	)

	// Note: Fatal and Panic are commented out as they would terminate the program
	// logger.Fatalw("System failure",
	//     "component", "database",
	//     "error", "connection lost",
	// )
	// logger.Panicw("Critical error",
	//     "component", "memory",
	//     "error", "out of memory",
	// )

	fmt.Println("Structured logging demonstrated. Check the log output.")
}

func contextLoggingExample() {
	// Create a context with values
	ctx := context.Background()
	ctx = context.WithValue(ctx, "request_id", "req-456")
	ctx = context.WithValue(ctx, "user_id", 123)

	// Get the default logger
	logger := log.DefaultLogger()

	// Use the context with the logger
	// Note: In a real application, the context would be used to extract
	// values and include them in the log messages automatically
	contextLogger := logger.WithContext(ctx)

	// Log messages using the context logger
	contextLogger.Info("Processing request with context")
	contextLogger.Infof("User %d made a request", 123)
	contextLogger.Infow("Request details",
		"request_id", ctx.Value("request_id"),
		"user_id", ctx.Value("user_id"),
	)

	fmt.Println("Context logging demonstrated. Check the log output.")
}

func configuringLoggerExample() {
	// Get the default logger
	logger := log.DefaultLogger()

	// Show the current log level
	fmt.Println("Default log level is set to show all messages")

	// Set the log level to only show warnings and above
	logger.SetLevel(log.LevelWarn)
	fmt.Println("Log level set to WARN")

	// These messages should not appear in the log
	logger.Trace("This trace message should not appear")
	logger.Debug("This debug message should not appear")
	logger.Info("This info message should not appear")

	// These messages should appear in the log
	logger.Warn("This warning message should appear")
	logger.Error("This error message should appear")

	// Reset the log level to show all messages
	logger.SetLevel(log.LevelTrace)
	fmt.Println("Log level reset to TRACE")

	// Now all messages should appear
	logger.Trace("This trace message should now appear")
	logger.Debug("This debug message should now appear")
	logger.Info("This info message should now appear")
	logger.Warn("This warning message should still appear")
	logger.Error("This error message should still appear")

	fmt.Println("Logger configuration demonstrated. Check the log output.")
}

func customOutputExample() {
	// Get the default logger
	logger := log.DefaultLogger()

	// Create a temporary file for log output
	tempFile, err := os.CreateTemp("", "gfly-log-example-*.log")
	if err != nil {
		fmt.Printf("Failed to create temp file: %v\n", err)
		return
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	fmt.Printf("Logging to temporary file: %s\n", tempFile.Name())

	// Save the original output
	originalOutput := os.Stdout

	// Set the logger to output to the temporary file
	logger.SetOutput(tempFile)

	// Log some messages
	logger.Info("This message goes to the temporary file")
	logger.Infof("Formatted message with value: %d", 42)
	logger.Infow("Structured message", "key", "value")

	// Reset the logger to output to stdout
	logger.SetOutput(originalOutput)

	// Read and display the contents of the temporary file
	tempFile.Seek(0, 0)
	fileContents, err := io.ReadAll(tempFile)
	if err != nil {
		fmt.Printf("Failed to read temp file: %v\n", err)
		return
	}

	fmt.Println("\nContents of the log file:")
	fmt.Println("------------------------")
	fmt.Println(string(fileContents))

	fmt.Println("Custom log output demonstrated.")
}

func errorLoggingExample() {
	// Get the default logger
	logger := log.DefaultLogger()

	// Simulate some errors
	err1 := fmt.Errorf("simple error: something went wrong")
	logger.Error(err1)

	// Log an error with additional context
	logger.Errorw("Failed to process request",
		"error", err1,
		"request_id", "req-456",
		"user_id", 123,
	)

	// Log a formatted error message
	logger.Errorf("Error occurred: %v", err1)

	// Create a more complex error scenario with a custom error type
	err2 := fmt.Errorf("custom error: code=%d, message=%s", 500, "Internal server error")
	logger.Error(err2)

	fmt.Println("Error logging demonstrated. Check the log output.")
}

func jsonStructuredLoggingExample() {
	// Get the default logger
	logger := log.DefaultLogger()

	// Enable structured logging with JSON format
	logger.EnableStructuredLogging(true)

	fmt.Println("Structured logging enabled with JSON format")

	// Log structured messages at different levels
	logger.Tracew("User login",
		"user_id", 123,
		"username", "john_doe",
		"ip", "192.168.1.1",
		"timestamp", time.Now().Format(time.RFC3339),
	)

	logger.Debugw("Request processing",
		"request_id", "req-456",
		"user_id", 123,
		"endpoint", "/api/users",
		"method", "GET",
		"duration_ms", 42.5,
	)

	logger.Infow("Request completed",
		"request_id", "req-456",
		"status", 200,
		"duration_ms", 42.5,
	)

	logger.Warnw("High latency detected",
		"request_id", "req-456",
		"duration_ms", 150.3,
		"threshold_ms", 100,
	)

	logger.Errorw("Request failed",
		"request_id", "req-456",
		"status", 500,
		"error", "database timeout",
		"duration_ms", 3000.5,
	)

	// Log with nested data structures
	metadata := map[string]interface{}{
		"server":      "api-server-01",
		"environment": "production",
		"version":     "1.2.3",
	}

	logger.Infow("Application metadata",
		"metadata", metadata,
		"uptime_hours", 72.5,
	)

	// Disable structured logging
	logger.EnableStructuredLogging(false)
	fmt.Println("Structured logging disabled, reverting to text format")

	// This will be logged in text format
	logger.Infow("Back to text logging",
		"note", "This is in text format now",
	)

	fmt.Println("JSON structured logging demonstrated. Check the log output.")
}
