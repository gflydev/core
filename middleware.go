package core

import (
	"sort"
)

// MiddlewarePhase defines the execution phase of a middleware
type MiddlewarePhase int

const (
	// PhasePreRequest middleware executes before the request is processed
	PhasePreRequest MiddlewarePhase = iota
	// PhasePostRequest middleware executes after the request is processed
	PhasePostRequest
)

// MiddlewareHandler must return RequestHandler for continuing or error to stop on it.
// It is a function type that takes a context object and processes it.
type MiddlewareHandler func(ctx *Ctx) error

// MiddlewareConfig defines configuration options for a middleware
type MiddlewareConfig struct {
	// Handler is the middleware function to execute
	Handler MiddlewareHandler
	// Phase determines when the middleware should be executed
	Phase MiddlewarePhase
	// Priority determines the order of execution (lower numbers execute first)
	Priority int
	// Condition is an optional function that determines if the middleware should be executed
	Condition func(ctx *Ctx) bool
	// Name is an optional identifier for the middleware
	Name string
}

// MiddlewareOption is a function that configures a MiddlewareConfig
type MiddlewareOption func(*MiddlewareConfig)

// WithPhase sets the execution phase of a middleware
func WithPhase(phase MiddlewarePhase) MiddlewareOption {
	return func(config *MiddlewareConfig) {
		config.Phase = phase
	}
}

// WithPriority sets the execution priority of a middleware
func WithPriority(priority int) MiddlewareOption {
	return func(config *MiddlewareConfig) {
		config.Priority = priority
	}
}

// WithCondition sets a condition function for a middleware
func WithCondition(condition func(ctx *Ctx) bool) MiddlewareOption {
	return func(config *MiddlewareConfig) {
		config.Condition = condition
	}
}

// WithName sets a name for a middleware
func WithName(name string) MiddlewareOption {
	return func(config *MiddlewareConfig) {
		config.Name = name
	}
}

// IMiddleware Middleware interface defines the contract for grouping middleware functions.
type IMiddleware interface {
	// Group groups multiple MiddlewareHandler functions together and returns a function
	// which wraps an IHandler with the middleware functionalities.
	//
	// Parameters:
	//   - middlewares: Variadic parameter accepting multiple MiddlewareHandler functions.
	//
	// Returns:
	//   - A function that takes an IHandler and returns a wrapped IHandler with the applied middlewares.
	Group(middlewares ...MiddlewareHandler) func(IHandler) IHandler

	// GroupWithOptions groups middleware functions with configuration options and returns a function
	// which wraps an IHandler with the middleware functionalities.
	//
	// Parameters:
	//   - configs: Variadic parameter accepting multiple middleware configurations.
	//
	// Returns:
	//   - A function that takes an IHandler and returns a wrapped IHandler with the applied middlewares.
	GroupWithOptions(configs ...MiddlewareConfig) func(IHandler) IHandler

	// Use adds a middleware handler with optional configuration.
	//
	// Parameters:
	//   - handler: The middleware handler function.
	//   - options: Optional configuration options for the middleware.
	//
	// Returns:
	//   - A middleware configuration that can be used with GroupWithOptions.
	Use(handler MiddlewareHandler, options ...MiddlewareOption) MiddlewareConfig
}

// Middleware represents a struct implementing the IMiddleware interface.
type Middleware struct{}

// Group Create a group of Middleware functions. Implements the Group method for the IMiddleware interface.
//
// Parameters:
//   - middlewares: Variadic parameter accepting multiple MiddlewareHandler functions.
//
// Returns:
//   - A function that takes an IHandler and returns a wrapped IHandler with the applied middlewares.
func (m *Middleware) Group(middlewares ...MiddlewareHandler) func(IHandler) IHandler {
	configs := make([]MiddlewareConfig, len(middlewares))
	for i, handler := range middlewares {
		configs[i] = MiddlewareConfig{
			Handler:  handler,
			Phase:    PhasePreRequest,
			Priority: i,
		}
	}
	return m.GroupWithOptions(configs...)
}

// GroupWithOptions groups middleware functions with configuration options and returns a function
// which wraps an IHandler with the middleware functionalities.
//
// Parameters:
//   - configs: Variadic parameter accepting multiple middleware configurations.
//
// Returns:
//   - A function that takes an IHandler and returns a wrapped IHandler with the applied middlewares.
func (m *Middleware) GroupWithOptions(configs ...MiddlewareConfig) func(IHandler) IHandler {
	return func(handler IHandler) IHandler {
		return &middlewareEndpoint{
			handler:     handler,
			middlewares: configs,
		}
	}
}

// Use adds a middleware handler with optional configuration.
//
// Parameters:
//   - handler: The middleware handler function.
//   - options: Optional configuration options for the middleware.
//
// Returns:
//   - A middleware configuration that can be used with GroupWithOptions.
func (m *Middleware) Use(handler MiddlewareHandler, options ...MiddlewareOption) MiddlewareConfig {
	config := MiddlewareConfig{
		Handler:  handler,
		Phase:    PhasePreRequest,
		Priority: 0,
	}

	for _, option := range options {
		option(&config)
	}

	return config
}

// middlewareEndpoint Default handler
// Wraps middleware as an implementation of the IHandler interface
type middlewareEndpoint struct {
	Endpoint
	handler     IHandler
	middlewares []MiddlewareConfig
}

// Handle processes the middleware stack and executes the contained handler.
//
// Parameters:
//   - c: A pointer to the Ctx object representing the context of the current request.
//
// Returns:
//   - An error if any middleware or handler validation fails; otherwise, nil.
func (m *middlewareEndpoint) Handle(c *Ctx) error {
	// Sort middlewares by phase and priority
	preRequestMiddlewares := make([]MiddlewareConfig, 0)
	postRequestMiddlewares := make([]MiddlewareConfig, 0)

	for _, middleware := range m.middlewares {
		if middleware.Phase == PhasePreRequest {
			preRequestMiddlewares = append(preRequestMiddlewares, middleware)
		} else {
			postRequestMiddlewares = append(postRequestMiddlewares, middleware)
		}
	}

	// Sort pre-request middlewares by priority
	sort.Slice(preRequestMiddlewares, func(i, j int) bool {
		return preRequestMiddlewares[i].Priority < preRequestMiddlewares[j].Priority
	})

	// Sort post-request middlewares by priority
	sort.Slice(postRequestMiddlewares, func(i, j int) bool {
		return postRequestMiddlewares[i].Priority < postRequestMiddlewares[j].Priority
	})

	// Run pre-request middleware functions
	for _, middleware := range preRequestMiddlewares {
		// Skip middleware if condition is not met
		if middleware.Condition != nil && !middleware.Condition(c) {
			continue
		}

		err := middleware.Handler(c)
		if err != nil {
			return err
		}
	}

	// Validate handler
	err := m.handler.Validate(c)
	if err != nil {
		return err
	}

	// Execute Handler's handling
	err = m.handler.Handle(c)

	// Run post-request middleware functions even if handler returns an error
	for _, middleware := range postRequestMiddlewares {
		// Skip middleware if condition is not met
		if middleware.Condition != nil && !middleware.Condition(c) {
			continue
		}

		// Execute post-request middleware and capture any errors
		postErr := middleware.Handler(c)
		if postErr != nil && err == nil {
			// Only override the original error if there wasn't one
			err = postErr
		}
	}

	return err
}

// Skip allows a middleware to be skipped for the current request.
// This function should be called within a middleware handler to skip
// subsequent middleware in the chain.
//
// Parameters:
//   - c: A pointer to the Ctx object representing the context of the current request.
//   - middlewareName: The name of the middleware to skip.
//
// Returns:
//   - error: Always returns nil to continue the middleware chain.
func Skip(c *Ctx, middlewareName string) error {
	skippedMiddleware := c.GetData("skipped_middleware")
	var skipped []string

	if skippedMiddleware != nil {
		skipped = skippedMiddleware.([]string)
	} else {
		skipped = make([]string, 0)
	}

	skipped = append(skipped, middlewareName)
	c.SetData("skipped_middleware", skipped)

	return nil
}

// IsSkipped checks if a middleware has been marked to be skipped.
//
// Parameters:
//   - c: A pointer to the Ctx object representing the context of the current request.
//   - middlewareName: The name of the middleware to check.
//
// Returns:
//   - bool: True if the middleware should be skipped, false otherwise.
func IsSkipped(c *Ctx, middlewareName string) bool {
	skippedMiddleware := c.GetData("skipped_middleware")
	if skippedMiddleware == nil {
		return false
	}

	skipped := skippedMiddleware.([]string)
	for _, name := range skipped {
		if name == middlewareName {
			return true
		}
	}

	return false
}

// NewMiddleware creates a new instance of Middleware.
//
// Returns:
//   - An object implementing the IMiddleware interface.
func NewMiddleware() IMiddleware {
	return &Middleware{}
}
