# Error Handling Guidelines for gFly Core

This document outlines the standard error handling patterns to be used across all packages in the gFly Core project. Following these guidelines ensures a consistent approach to error handling, making the codebase more maintainable and providing a better developer experience.

## Core Principles

1. **Use the gFly errors package**: Always use the `github.com/gflydev/core/errors` package for error handling, not the standard `errors` package.
2. **Categorize errors**: Use appropriate error codes to categorize errors (e.g., `errors.CodeNotFound`, `errors.CodeInvalidInput`).
3. **Provide context**: Include descriptive error messages and additional context information when creating errors.
4. **Preserve error chains**: Use `errors.Wrap` or `errors.Wrapf` to wrap errors while preserving the original error.
5. **Return errors, don't handle them**: Functions should generally return errors rather than handling them internally, allowing the caller to decide how to handle them.

## Error Creation

### Creating New Errors

Use the appropriate constructor function based on the error type:

Example:
- Create a basic error with a specific code: `errors.New(errors.CodeNotFound, "User not found")`
- Create an error with a formatted message: `errors.Newf(errors.CodeInvalidInput, "Invalid value for field %s: %v", "email", value)`
- Use helper functions for common error types:
  - `errors.NotFound("User not found")`
  - `errors.InvalidInput("Invalid email address")`
  - `errors.Unauthorized("Invalid API key")`
  - `errors.Internal("Database connection failed")`

### Wrapping Errors

When an error occurs in a lower-level function, wrap it with additional context:

Example:
- Wrap an error with additional context:
  ```
  if err := db.QueryRow("SELECT * FROM users WHERE id = ?", id); err != nil {
      return errors.Wrap(err, errors.CodeNotFound, "User not found")
  }
  ```
- Wrap with a formatted message:
  ```
  if err := validateInput(input); err != nil {
      return errors.Wrapf(err, errors.CodeInvalidInput, "Invalid input for field %s", fieldName)
  }
  ```

### Adding Context and Metadata

Enrich errors with additional context and metadata:

Example:
- Add context information:
  ```
  err := errors.InvalidInput("Invalid user data").
      WithContext("user_id", userId).
      WithContext("request_id", requestId)
  ```
- Add metadata:
  ```
  err := errors.NotFound("User not found").
      WithMetadata("field", "email").
      WithMetadata("value", value)
  ```

## Error Handling Patterns

### Basic Error Handling

Example:
```
result, err := someFunction()
if err != nil {
    // Check if it's a specific type of error
    if errors.Is(err, errors.ItemNotFoundError) {
        // Handle not found error
    } else {
        // Handle other errors
    }
    return errors.Wrap(err, errors.CodeInternal, "Operation failed")
}
```

### Structured Error Handling with try-catch

For more complex error handling scenarios, use the `try` package:

Example:
```
try.Perform(func() {
    // Code that might fail
    result, err := someFunction()
    if err != nil {
        try.Throw(err) // Throw the error to be caught
    }
    // Process result
}).Catch(func(err try.E) {
    // Handle the error
    if gflyErr, ok := err.(errors.Error); ok {
        // Handle gFly error with additional context
        log.Printf("Error code: %s, Message: %s", gflyErr.Code(), gflyErr.Message())
    } else {
        // Handle standard error
        log.Printf("Error: %v", err)
    }
}).Finally(func() {
    // Cleanup code that always runs
    closeResources()
})
```

## HTTP Error Handling

When handling errors in HTTP handlers, use the error's status code:

Example:
```
func handler(w http.ResponseWriter, r *http.Request) {
    result, err := someFunction()
    if err != nil {
        var apiErr errors.Error
        if errors.As(err, &apiErr) {
            // If it's already our custom error type
            w.WriteHeader(apiErr.StatusCode())
            json.NewEncoder(w).Encode(map[string]string{
                "error": apiErr.Error(),
                "code":  apiErr.Code(),
            })
        } else {
            // Wrap unknown errors as internal server errors
            wrappedErr := errors.Wrap(err, errors.CodeInternal, "Internal server error")
            w.WriteHeader(wrappedErr.StatusCode())
            json.NewEncoder(w).Encode(map[string]string{
                "error": wrappedErr.Error(),
                "code":  wrappedErr.Code(),
            })
        }
        return
    }
    // Handle success case
}
```

## Best Practices

1. **Be specific**: Use the most specific error type that applies to the situation.
2. **Be descriptive**: Error messages should be clear and descriptive, helping developers understand what went wrong.
3. **Add context**: Include relevant context information that can help with debugging.
4. **Don't expose sensitive information**: Ensure error messages don't contain sensitive information that could be exposed to users.
5. **Log appropriately**: Log errors at the appropriate level based on their severity.
6. **Test error paths**: Write tests that verify error handling works correctly.

## Common Error Codes

- `errors.CodeNotFound`: Resource not found
- `errors.CodeInvalidInput`: Invalid input data
- `errors.CodeUnauthorized`: Authorization error
- `errors.CodeUnauthenticated`: Authentication error
- `errors.CodeForbidden`: Forbidden access
- `errors.CodeConflict`: Conflict (e.g., duplicate resource)
- `errors.CodeInternal`: Internal server error
- `errors.CodeTimeout`: Timeout error
- `errors.CodeUnavailable`: Service unavailable
- `errors.CodeNotImplemented`: Not implemented

## Example

Here's a complete example of a function that follows these error handling guidelines:

```
func GetUser(id string) (*User, error) {
    if id == "" {
        return nil, errors.InvalidInput("User ID cannot be empty")
    }
    
    user, err := db.FindUserByID(id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, errors.NotFoundf("User with ID %s not found", id)
        }
        return nil, errors.Wrapf(err, errors.CodeInternal, "Failed to find user with ID %s", id)
    }
    
    return user, nil
}
```

By following these guidelines, we ensure a consistent approach to error handling across the entire gFly Core project, making the codebase more maintainable and providing a better developer experience.