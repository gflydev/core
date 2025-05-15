# Request ID Tracking

This document explains how to use the request ID tracking feature in the gFly Core framework.

## Overview

Request ID tracking is a feature that assigns a unique identifier to each incoming HTTP request. This identifier is then used to track the request throughout its lifecycle, making it easier to debug issues and trace request flows in logs.

## Features

- Automatic generation of unique request IDs for each incoming request
- Preservation of existing request IDs from client requests
- Inclusion of request IDs in response headers
- Integration with the logging system to include request IDs in log entries
- Helper functions to retrieve request IDs and request-specific loggers

## Usage

### Enabling Request ID Tracking

To enable request ID tracking, you need to register the request ID middleware in your application. This can be done in your application's bootstrap code:

```go
package main

import "github.com/gflydev/core"

func main() {
    app := core.New()

    // Register middleware
    app.RegisterMiddleware(func(fly core.IFlyMiddleware) {
        // Add request ID middleware
        fly.Use(core.RequestIDMiddleware)

        // Add request ID log middleware (must be after RequestIDMiddleware)
        fly.Use(core.RequestIDLogMiddleware)

        // Add other middleware...
    })

    // Register routes...

    app.Run()
}
```

### Using Request IDs in Your Code

Once the middleware is registered, you can retrieve the request ID from the context in your handlers:

```go
package handlers

import "github.com/gflydev/core"

func MyHandler(ctx *core.Ctx) error {
    // Get the request ID
    requestID := core.GetRequestID(ctx)

    // Use the request ID...

    return nil
}
```

### Using Request-Specific Loggers

The request ID log middleware adds a request-specific logger to the context. You can retrieve this logger and use it to include the request ID in your log entries:

```go
package handlers

import "github.com/gflydev/core"

func MyHandler(ctx *core.Ctx) error {
    // Get the request-specific logger
    logger := core.GetLogger(ctx)

    // Use the logger
    logger.Info("Processing request")

    // Example with data variable
    data := map[string]string{"key": "value"}
    logger.Infof("Request data: %v", data)

    // Example with duration variable
    duration := 100
    logger.Infow("Request processed", "status", "success", "duration", duration)

    return nil
}
```

All log entries created using this logger will automatically include the request ID.

## Configuration

The request ID middleware uses the following constants that can be customized if needed:

- `RequestIDHeader`: The HTTP header used for request IDs (default: "X-Request-ID")
- `RequestIDContextKey`: The key used to store the request ID in the context (default: "request_id")

## Best Practices

- Always use the request-specific logger obtained from `GetLogger(ctx)` instead of the global logger when logging within a request handler.
- Include the request ID in any external service calls to maintain traceability across services.
- Consider adding the request ID to error responses to help correlate client errors with server logs.

## Example

Here's a complete example of how to use request ID tracking in a handler:

```go
package handlers

import "github.com/gflydev/core"

func MyHandler(ctx *core.Ctx) error {
    // Get the request ID
    requestID := core.GetRequestID(ctx)

    // Get the request-specific logger
    logger := core.GetLogger(ctx)

    // Log with request ID
    logger.Infof("Processing request %s", requestID)

    // Do something...

    // Example variables
    duration := 100
    data := map[string]string{"key": "value"}

    // Log the result
    logger.Infow("Request processed", 
        "status", "success",
        "duration", duration,
        "data", data)

    return ctx.Success(map[string]interface{}{
        "status": "success",
        "request_id": requestID,
    })
}
```

In this example, all log entries will include the request ID, and the response will also include the request ID in the JSON payload.
