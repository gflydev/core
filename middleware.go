package core

// MiddlewareHandler must return RequestHandler for continuing or error to stop on it.
// It is a function type that takes a context object and processes it.
type MiddlewareHandler func(ctx *Ctx) error

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
	return func(handler IHandler) IHandler {
		return &middlewareEndpoint{
			handler:     handler,
			middlewares: middlewares,
		}
	}
}

// middlewareEndpoint Default handler
// Wraps middleware as an implementation of the IHandler interface
type middlewareEndpoint struct {
	Endpoint
	handler     IHandler
	middlewares []MiddlewareHandler
}

// Handle processes the middleware stack and executes the contained handler.
//
// Parameters:
//   - c: A pointer to the Ctx object representing the context of the current request.
//
// Returns:
//   - An error if any middleware or handler validation fails; otherwise, nil.
func (m *middlewareEndpoint) Handle(c *Ctx) error {
	// Run middleware functions
	for _, m := range m.middlewares {
		err := m(c)
		if err != nil {
			return err
		}
	}

	err := m.handler.Validate(c)
	if err != nil {
		return err
	}

	// Execute Handler's handling
	return m.handler.Handle(c)
}

// NewMiddleware creates a new instance of Middleware.
//
// Returns:
//   - An object implementing the IMiddleware interface.
func NewMiddleware() IMiddleware {
	return &Middleware{}
}
