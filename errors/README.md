# Errors Package

The errors package provides standardized error handling for the gFly framework. It defines error types, interfaces, and helper functions for creating, wrapping, and handling errors in a consistent way across all packages.

## Features

- **Error Categorization**: Errors are categorized by error codes for consistent handling
- **Error Context**: Add context information to errors for better debugging
- **Error Metadata**: Attach metadata to errors for additional information
- **Error Wrapping**: Preserve the error chain when wrapping errors
- **Stack Traces**: Automatically capture stack traces for better debugging
- **HTTP Status Codes**: Map error codes to HTTP status codes for API responses
- **JSON Serialization**: Serialize errors to JSON for API responses

## Usage

### Creating Errors

```go
import "github.com/gflydev/core/errors"

// Create a basic error
err := errors.New(errors.CodeNotFound, "User not found")

// Create an error with a formatted message
err := errors.Newf(errors.CodeInvalidInput, "Invalid value for field %s: %v", "email", value)

// Use helper functions for common error types
err := errors.NotFound("User not found")
err := errors.InvalidInput("Invalid email address")
err := errors.Unauthorized("Access denied")
```

### Wrapping Errors

```go
// Wrap an existing error
originalErr := db.QueryRow("SELECT * FROM users WHERE id = ?", id)
if originalErr != nil {
    return errors.Wrap(originalErr, errors.CodeNotFound, "User not found")
}

// Wrap with a formatted message
if err := validateInput(input); err != nil {
    return errors.Wrapf(err, errors.CodeInvalidInput, "Invalid input for field %s", fieldName)
}
```

### Adding Context and Metadata

```go
// Add context information
err := errors.InvalidInput("Invalid email address").
    WithContext("user_id", userId).
    WithContext("request_id", requestId)

// Add metadata
err := errors.NotFound("User not found").
    WithMetadata("query", query).
    WithMetadata("params", params)

// Retrieve context and metadata
if userId, ok := err.GetContext("user_id"); ok {
    // Use userId
}

if query, ok := err.GetMetadata("query"); ok {
    // Use query
}
```

### Error Codes

The package defines standard error codes for categorizing errors:

- `CodeUnknown`: Unknown or uncategorized error
- `CodeNotFound`: Resource not found error
- `CodeInvalidInput`: Invalid input error
- `CodeUnauthorized`: Authorization error
- `CodeUnauthenticated`: Authentication error
- `CodeForbidden`: Forbidden access error
- `CodeConflict`: Conflict error (e.g., duplicate resource)
- `CodeInternal`: Internal server error
- `CodeTimeout`: Timeout error
- `CodeUnavailable`: Service unavailable error
- `CodeNotImplemented`: Not implemented error

### Helper Functions

The package provides helper functions for creating common error types:

- `NotFound(message)`: Creates a not found error
- `InvalidInput(message)`: Creates an invalid input error
- `Unauthorized(message)`: Creates an unauthorized error
- `Internal(message)`: Creates an internal server error
- `Timeout(message)`: Creates a timeout error
- `Unavailable(message)`: Creates a service unavailable error
- `Conflict(message)`: Creates a conflict error
- `Unauthenticated(message)`: Creates an unauthenticated error
- `Forbidden(message)`: Creates a forbidden error
- `NotImplemented(message)`: Creates a not implemented error

Each helper function also has a formatted version (e.g., `NotFoundf`, `InvalidInputf`) that accepts a format string and arguments.

### Legacy Support

For backward compatibility, the package provides legacy error constants:

```go
// Legacy usage (still supported)
panic(errors.ItemNotFoundError)
```

## Best Practices

1. **Use Error Codes**: Always use appropriate error codes for categorization
2. **Add Context**: Include relevant context information for debugging
3. **Wrap Errors**: Preserve the error chain by wrapping errors
4. **Use Helper Functions**: Use the provided helper functions for common error types
5. **Check Error Types**: Use the `Is` and `As` functions to check error types
