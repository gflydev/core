# gFly Framework Architecture

## Overview

gFly is a high-performance web framework for Go, designed to provide a robust foundation for building web applications and APIs. The framework is built on top of [fasthttp](https://github.com/valyala/fasthttp) and [fluentsql](https://github.com/jiveio/fluentsql) for high-performance HTTP and Database handling and provides a comprehensive set of features for routing, middleware, error handling, logging, and more.

This document provides a detailed explanation of the framework's architecture, components, and their interactions.

## Core Components

The gFly framework consists of several key components that work together to provide a complete web application framework:

### 1. Application Core (GFly)

The `GFly` struct is the central component of the framework, serving as the application container. It initializes and coordinates all other components.

**Key Responsibilities:**
- Application initialization and configuration
- HTTP server setup and management
- Router and middleware coordination
- Request handling and lifecycle management

**Location:** `gfly.go`

### 2. Router

The router is responsible for matching HTTP requests to the appropriate handlers based on the request method and path.

**Key Features:**
- High-performance routing using a radix tree implementation
- Support for all HTTP methods (GET, POST, PUT, etc.)
- Route grouping for organizing related routes
- Path parameters and optional path segments
- Automatic redirection for trailing slashes
- Case-insensitive path matching
- Custom error handling for 404 and 405 responses

**Location:** `router.go`, `redix_tree.go`

### 3. Context

The context provides access to the HTTP request and response, as well as utilities for handling the request and generating the response.

**Key Features:**
- HTTP header manipulation
- Response generation (JSON, HTML, files, etc.)
- Request data parsing (form data, query parameters, path parameters)
- File uploads and downloads
- Session management
- Data storage for the request lifecycle

**Location:** `context.go`

### 4. Middleware

Middleware provides a way to process requests before they reach the handler, allowing for cross-cutting concerns like authentication, logging, and error handling.

**Key Features:**
- Global middleware for all requests
- Group-specific middleware for route groups
- Middleware chaining
- Error handling within middleware

**Location:** `middleware.go`

### 5. Logging

The logging system provides structured logging with different levels of verbosity.

**Key Features:**
- Multiple log levels (Trace, Debug, Info, Warn, Error, Fatal, Panic)
- Formatted logging with printf-style formatting
- Structured logging with key-value pairs
- Context-aware logging
- Configurable output destinations
- Colored output for different log levels

**Location:** `log/log.go`, `log/default.go`, `log/adapter.go`

### 6. Error Handling

The error handling system provides a structured way to handle errors and panics.

**Key Features:**
- Custom error types for common error scenarios
- Try-catch-finally pattern for structured error handling
- Panic recovery
- Error propagation

**Location:** `errors/errors.go`, `try/perform_finally_catch.go`

### 7. Utilities

The utilities package provides a collection of helper functions for common tasks.

**Key Features:**
- String manipulation
- Byte manipulation
- File operations
- Environment variable handling
- Token generation
- Hashing
- Time utilities
- Reflection utilities
- Map and slice utilities

**Location:** `utils/*.go`

## Component Interactions

The components of the gFly framework interact in a structured way to handle HTTP requests and generate responses:

1. **Application Initialization:**
   - The application is initialized with `New()`, which creates a new `GFly` instance.
   - The application configuration is set up, either with default values or custom configuration.
   - Middleware and routes are registered with the application.

2. **Server Startup:**
   - The application is started with `Run()`, which initializes the HTTP server.
   - Global middleware is set up.
   - Routes are registered with the router.
   - The server starts listening for incoming requests.

3. **Request Handling:**
   - When a request is received, it is wrapped in a `Ctx` object.
   - The request is passed to the router for matching.
   - If a matching route is found, the associated handler is executed.
   - If no matching route is found, the appropriate error handler is executed.

4. **Middleware Processing:**
   - Before the handler is executed, all applicable middleware is executed in order.
   - Each middleware can modify the request or response, or terminate the request processing early.
   - If any middleware returns an error, the request processing is terminated and the error is returned.

5. **Handler Execution:**
   - The handler is executed with the `Ctx` object.
   - The handler can access request data, generate a response, and interact with other components.
   - If the handler returns an error, it is processed by the error handling system.

6. **Response Generation:**
   - The handler generates a response using the `Ctx` object.
   - The response is sent back to the client.

7. **Error Handling:**
   - If an error occurs during request processing, it is handled by the error handling system.
   - Depending on the error type and configuration, the error may be logged, returned to the client, or trigger a panic.

For a more detailed explanation of the request lifecycle, including visual diagrams, see [REQUEST_LIFECYCLE.md](REQUEST_LIFECYCLE.md).

## Design Patterns

The gFly framework uses several design patterns to provide a flexible and extensible architecture:

### 1. Middleware Pattern

The middleware pattern is used to process requests before they reach the handler. This allows for cross-cutting concerns like authentication, logging, and error handling to be separated from the main application logic.

### 2. Chain of Responsibility

The middleware system uses the chain of responsibility pattern to pass requests through a series of handlers, each of which can process the request or pass it to the next handler.

### 3. Dependency Injection

The framework uses dependency injection to provide components with access to other components they depend on. For example, the `Ctx` object is injected into handlers and middleware.

### 4. Factory Method

The framework uses factory methods like `New()` to create new instances of components with the appropriate configuration.

### 5. Singleton

Some components, like the logger, are implemented as singletons to ensure that there is only one instance of the component in the application.

### 6. Adapter

The logging system uses the adapter pattern to allow different logging implementations to be used with the same interface.

## Extension Points

The gFly framework provides several extension points for customizing and extending its functionality:

### 1. Custom Middleware

Custom middleware can be created by implementing the `MiddlewareHandler` function type and registered with the application using the `Use()` method.

### 2. Custom Handlers

Custom handlers can be created by implementing the `IHandler` interface and registered with the router using the appropriate HTTP method functions (`GET()`, `POST()`, etc.).

### 3. Custom Error Handlers

Custom error handlers can be registered with the router to handle specific error scenarios, such as 404 (Not Found) or 405 (Method Not Allowed) errors.

### 4. Custom Logging Adapters

Custom logging adapters can be created by implementing the `AllLogger` interface and registered with the logging system using the `SetLogger()` function.

## Conclusion

The gFly framework provides a comprehensive set of components for building web applications and APIs in Go. Its modular architecture allows for flexibility and extensibility, while its high-performance foundation ensures efficient request handling.

By understanding the framework's components and their interactions, developers can effectively use the framework to build robust and scalable web applications.
