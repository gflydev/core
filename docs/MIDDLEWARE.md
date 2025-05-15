# Middleware in gFly Core

This document provides an overview of the middleware system in gFly Core and demonstrates how to use its flexible features.

## Table of Contents

1. [Introduction](#introduction)
2. [Basic Middleware Usage](#basic-middleware-usage)
3. [Advanced Middleware Features](#advanced-middleware-features)
   - [Middleware Phases](#middleware-phases)
   - [Middleware Priorities](#middleware-priorities)
   - [Conditional Middleware](#conditional-middleware)
   - [Named Middleware and Skipping](#named-middleware-and-skipping)
4. [Best Practices](#best-practices)
5. [Examples](#examples)

## Introduction

Middleware in gFly Core provides a way to process requests before they reach the handler and after the handler has processed the request. Middleware can be used for various purposes such as authentication, logging, error handling, and more.

The middleware system in gFly Core is designed to be flexible and powerful, allowing for:

- Global middleware for all requests
- Group-specific middleware for route groups
- Middleware ordering and prioritization
- Conditional middleware execution
- Post-processing middleware
- Middleware skipping

## Basic Middleware Usage

### Creating a Simple Middleware

A middleware in gFly Core is a function that takes a context object and returns an error:

```go
func LoggerMiddleware(ctx *core.Ctx) error {
    // Log the request
    log.Infof("Request: %s %s", ctx.Method(), ctx.Path())

    // Continue to the next middleware or handler
    return nil
}
```

### Registering Global Middleware

Global middleware is applied to all routes in the application:

```go
app := core.New()
app.Use(LoggerMiddleware)
```

### Registering Group-Specific Middleware

Group-specific middleware is applied only to routes in a specific group:

```go
app.Group("/api", func(g *core.Group) {
    g.Use(AuthMiddleware)

    g.GET("/users", GetUsersHandler)
})
```

## Advanced Middleware Features

### Middleware Phases

gFly Core supports two middleware phases:

1. **Pre-Request Phase (PhasePreRequest)**: Middleware executes before the handler processes the request (default)
2. **Post-Request Phase (PhasePostRequest)**: Middleware executes after the handler has processed the request

Example of a post-request middleware:

```go
// Response time middleware that measures how long it takes to process a request
func ResponseTimeMiddleware(ctx *core.Ctx) error {
    // Store the start time in the context
    startTime := time.Now()
    ctx.SetData("start_time", startTime)

    // Continue to the next middleware or handler
    return nil
}

// Response time logger middleware that logs the time taken to process a request
func ResponseTimeLoggerMiddleware(ctx *core.Ctx) error {
    // Get the start time from the context
    startTime := ctx.GetData("start_time").(time.Time)

    // Calculate the time taken
    duration := time.Since(startTime)

    // Log the time taken
    log.Infof("Request processed in %v", duration)

    // Continue to the next middleware
    return nil
}

// Register the middlewares with appropriate phases
app := core.New()
app.UseWithOptions(ResponseTimeMiddleware, core.WithName("response_time"))
app.UseWithOptions(ResponseTimeLoggerMiddleware, 
    core.WithName("response_time_logger"), 
    core.WithPhase(core.PhasePostRequest))
```

### Middleware Priorities

Middleware execution order can be controlled using priorities. Lower priority values execute first:

```go
app := core.New()

// This middleware will execute first (priority 1)
app.UseWithOptions(SecurityHeadersMiddleware, core.WithPriority(1))

// This middleware will execute second (priority 2)
app.UseWithOptions(LoggerMiddleware, core.WithPriority(2))

// This middleware will execute third (priority 3)
app.UseWithOptions(AuthMiddleware, core.WithPriority(3))
```

### Conditional Middleware

Middleware can be conditionally executed based on request properties:

```go
// Only apply CORS middleware to API routes
app.UseWithOptions(CorsMiddleware, core.WithCondition(func(ctx *core.Ctx) bool {
    return strings.HasPrefix(ctx.Path(), "/api")
}))

// Only apply rate limiting to non-admin users
app.UseWithOptions(RateLimitMiddleware, core.WithCondition(func(ctx *core.Ctx) bool {
    user := ctx.GetData("user")
    if user == nil {
        return true // Apply rate limiting to unauthenticated users
    }
    return user.(User).Role != "admin" // Don't apply rate limiting to admin users
}))
```

### Named Middleware and Skipping

Middleware can be named and skipped for specific requests:

```go
// Define a named middleware
app.UseWithOptions(AuthMiddleware, core.WithName("auth"))

// Define a middleware that can skip other middleware
func SkipAuthForPublicApisMiddleware(ctx *core.Ctx) error {
    if strings.HasPrefix(ctx.Path(), "/api/public") {
        // Skip the auth middleware for public APIs
        return core.Skip(ctx, "auth")
    }
    return nil
}

// Register the middleware that can skip other middleware
app.UseWithOptions(SkipAuthForPublicApisMiddleware, core.WithPriority(0))
```

## Best Practices

1. **Order Matters**: Be mindful of the order in which middleware is registered. Use priorities to ensure correct execution order.

2. **Keep Middleware Focused**: Each middleware should have a single responsibility.

3. **Use Phases Appropriately**: 
   - Use PhasePreRequest for middleware that needs to run before the handler (authentication, validation, etc.)
   - Use PhasePostRequest for middleware that needs to run after the handler (logging, metrics, etc.)

4. **Name Your Middleware**: Always give your middleware a name, especially if it might need to be skipped.

5. **Use Conditions Wisely**: Conditional middleware can improve performance by avoiding unnecessary middleware execution.

6. **Error Handling**: Middleware should handle errors appropriately and decide whether to pass them along or handle them internally.

## Examples

### Complete Authentication Example

```go
func AuthMiddleware(ctx *core.Ctx) error {
    // Get the token from the Authorization header
    token := ctx.GetReqHeaders()["Authorization"]
    if len(token) == 0 || !strings.HasPrefix(token[0], "Bearer ") {
        return errors.Unauthorized("Missing or invalid authorization token")
    }

    // Validate the token
    tokenString := strings.TrimPrefix(token[0], "Bearer ")
    user, err := validateToken(tokenString)
    if err != nil {
        return errors.Unauthorized("Invalid token").WithContext("error", err.Error())
    }

    // Store the user in the context
    ctx.SetData("user", user)

    // Continue to the next middleware or handler
    return nil
}

// Skip authentication for public routes
func SkipAuthForPublicRoutesMiddleware(ctx *core.Ctx) error {
    publicPaths := []string{"/login", "/register", "/public"}
    for _, path := range publicPaths {
        if strings.HasPrefix(ctx.Path(), path) {
            return core.Skip(ctx, "auth")
        }
    }
    return nil
}

// Register the middlewares
app := core.New()
app.UseWithOptions(SkipAuthForPublicRoutesMiddleware, core.WithPriority(1), core.WithName("skip_auth"))
app.UseWithOptions(AuthMiddleware, core.WithPriority(2), core.WithName("auth"))
```

### Logging and Metrics Example

```go
func RequestLoggerMiddleware(ctx *core.Ctx) error {
    // Generate a request ID
    requestID := utils.GenerateUUID()
    ctx.SetData("request_id", requestID)

    // Log the request
    log.WithFields(log.Fields{
        "request_id": requestID,
        "method":     ctx.Method(),
        "path":       ctx.Path(),
        "ip":         ctx.IP(),
    }).Info("Request received")

    // Store the start time
    ctx.SetData("request_start_time", time.Now())

    return nil
}

func ResponseLoggerMiddleware(ctx *core.Ctx) error {
    // Get the request ID and start time
    requestID := ctx.GetData("request_id").(string)
    startTime := ctx.GetData("request_start_time").(time.Time)

    // Calculate the duration
    duration := time.Since(startTime)

    // Log the response
    log.WithFields(log.Fields{
        "request_id": requestID,
        "status":     ctx.Response.StatusCode(),
        "duration":   duration.String(),
    }).Info("Response sent")

    return nil
}

// Register the middlewares
app := core.New()
app.UseWithOptions(RequestLoggerMiddleware, core.WithPriority(1), core.WithName("request_logger"))
app.UseWithOptions(ResponseLoggerMiddleware, core.WithPriority(1), core.WithPhase(core.PhasePostRequest), core.WithName("response_logger"))
```

These examples demonstrate how to use the flexible middleware system in gFly Core to implement common patterns like authentication, logging, and metrics.

### Rate Limiting Example

```go
// Apply global rate limiting middleware (60 requests per minute)
app.UseWithOptions(
    core.RateLimitMiddleware(
        core.WithRequestsPerMinute(60),
    ),
    core.WithName("global-rate-limit"),
)

// Apply a more restrictive rate limit to the API group (10 requests per minute)
apiGroup.UseWithOptions(
    core.RateLimitMiddleware(
        core.WithRequestsPerMinute(10),
        core.WithRateLimitMessage("API rate limit exceeded. Please slow down your requests."),
    ),
    core.WithName("api-rate-limit"),
)

// Apply user-specific rate limiting (5 requests per minute per user)
userGroup.UseWithOptions(
    core.RateLimitMiddleware(
        core.WithRequestsPerMinute(5),
        // Use a custom key function that uses the user ID from the query parameter
        core.WithRateLimitKeyFunc(func(ctx *core.Ctx) string {
            userID := ctx.QueryStr("user_id")
            if userID == "" {
                // If no user ID is provided, use the request path as a fallback
                return "anonymous-" + ctx.Path()
            }
            return userID
        }),
        core.WithRateLimitMessage("User-specific rate limit exceeded."),
    ),
    core.WithName("user-rate-limit"),
)
```

## Rate Limiting Middleware

The rate limiting middleware provides a way to limit the number of requests a client can make to your API within a specified time period. This helps protect your API from abuse, prevents resource exhaustion, and ensures fair usage.

### Features

- Token bucket algorithm for rate limiting
- Configurable rate limits (requests per second or minute)
- Configurable burst size for handling traffic spikes
- Rate limit headers in the response
- Customizable key function to identify clients
- Automatic cleanup of expired buckets to prevent memory leaks

### Configuration Options

The rate limiting middleware can be configured with the following options:

- `WithRequestsPerMinute(n)`: Sets the rate limiter to allow n requests per minute
- `WithRequestsPerSecond(n)`: Sets the rate limiter to allow n requests per second
- `WithTokenBucketLimiter(rate, burst, cleanupInterval, expiration)`: Creates a custom token bucket limiter
- `WithRateLimitKeyFunc(func)`: Sets a custom key function to identify clients
- `WithRateLimitMessage(message)`: Sets a custom message for rate limit exceeded errors
- `WithRateLimitHeaders(bool)`: Enables or disables rate limit headers in the response

### Rate Limit Headers

When enabled, the middleware adds the following headers to the response:

- `X-RateLimit-Limit`: The maximum number of requests allowed in the current time window
- `X-RateLimit-Remaining`: The number of requests remaining in the current time window
- `X-RateLimit-Reset`: The time when the rate limit will reset, in Unix timestamp format

### Example Usage

See the complete example in `examples/rate_limit_example.go`.
