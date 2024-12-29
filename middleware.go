package core

// MiddlewareHandler must return RequestHandler for continuing or error to stop on it.
type MiddlewareHandler func(ctx *Ctx) error

// IMiddleware Middleware interface
type IMiddleware interface {
	Group(middlewares ...MiddlewareHandler) func(IHandler) IHandler
}

// Middleware Middleware type
type Middleware struct{}

// Group Create a group Middleware functions. Implement for IMiddleware interface
func (m *Middleware) Group(middlewares ...MiddlewareHandler) func(IHandler) IHandler {
	return func(handler IHandler) IHandler {
		return &middlewareEndpoint{
			handler:     handler,
			middlewares: middlewares,
		}
	}
}

// middlewareEndpoint Default handler
// Need to wrap middleware as a implementation of IHandler interface
type middlewareEndpoint struct {
	Endpoint
	handler     IHandler
	middlewares []MiddlewareHandler
}

// Handle check handler validation
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

func NewMiddleware() IMiddleware {
	return &Middleware{}
}
