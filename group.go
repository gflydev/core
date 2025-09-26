package core

// Group is a sub-router to group paths
type Group struct {
	router *Router
	prefix string
	// Global middleware handler for router group.
	middlewares []MiddlewareHandler
}

// ===========================================================================================================
// 										Group Middleware
// ===========================================================================================================

// IGroupMiddleware Interface to declare all Middleware methods for gFly struct.
type IGroupMiddleware interface {
	// Use applies middleware for all router groups.
	// Important: This should be called at the top of the group router before defining routes.
	// Parameters:
	// - middlewares: Variadic parameter accepting one or more middleware handlers of type MiddlewareHandler.
	Use(middlewares ...MiddlewareHandler)

	// Apply creates an Interceptor chain handler for a specific handler.
	//
	// Note: Interceptors have access to response/request before and after the route handler is called
	//
	// Parameters:
	//   - middleware ([]MiddlewareHandler): Middleware handlers to be grouped.
	//
	// Returns:
	//   - func(IHandler) IHandler: A function that applies the interceptor handlers to an IHandler.
	Apply(middleware ...MiddlewareHandler) func(IHandler) IHandler
}

// Use applies middleware for all router groups.
// This function appends the provided middleware handlers to the group's middleware chain.
// Parameters:
// - middlewares: Variadic parameter accepting one or more middleware handlers of type MiddlewareHandler.
func (g *Group) Use(middlewares ...MiddlewareHandler) {
	g.middlewares = append(g.middlewares, middlewares...)
}

// Apply creates an Interceptor chain handler for a specific handler.
//
// Example usage:
//
//	group.POST("/one", group.Apply(middleware.RuleMiddlewareFunc)(api.NewDefaultApi()))
//
// Parameters:
//   - middlewares ([]MiddlewareHandler): Middleware handlers to be grouped.
//
// Returns:
//   - func(IHandler) IHandler: A function that applies the interceptor handlers to an IHandler.
func (g *Group) Apply(middlewares ...MiddlewareHandler) func(IHandler) IHandler {
	// Create a new middleware instance and use its Group method to apply the middlewares
	middleware := NewMiddleware()
	return middleware.Group(middlewares...)
}

// ===========================================================================================================
// 										Group Router
// ===========================================================================================================

// IGroupRouter Interface to declare all HTTP methods.
type IGroupRouter interface {
	// GET Http GET method
	// Parameters:
	// - path: The path for the route.
	// - handler: The handler implementing the IHandler interface to process the HTTP GET request.
	GET(path string, handler IHandler)

	// HEAD Http HEAD method
	// Parameters:
	// - path: The path for the route.
	// - handler: The handler implementing the IHandler interface to process the HTTP HEAD request.
	HEAD(path string, handler IHandler)

	// POST Http POST method
	// Parameters:
	// - path: The path for the route.
	// - handler: The handler implementing the IHandler interface to process the HTTP POST request.
	POST(path string, handler IHandler)

	// PUT Http PUT method
	// Parameters:
	// - path: The path for the route.
	// - handler: The handler implementing the IHandler interface to process the HTTP PUT request.
	PUT(path string, handler IHandler)

	// PATCH Http PATCH method
	// Parameters:
	// - path: The path for the route.
	// - handler: The handler implementing the IHandler interface to process the HTTP PATCH request.
	PATCH(path string, handler IHandler)

	// DELETE Http DELETE method
	// Parameters:
	// - path: The path for the route.
	// - handler: The handler implementing the IHandler interface to process the HTTP DELETE request.
	DELETE(path string, handler IHandler)

	// CONNECT Http CONNECT method
	// Parameters:
	// - path: The path for the route.
	// - handler: The handler implementing the IHandler interface to process the HTTP CONNECT request.
	CONNECT(path string, handler IHandler)

	// OPTIONS Http OPTIONS method
	// Parameters:
	// - path: The path for the route.
	// - handler: The handler implementing the IHandler interface to process the HTTP OPTIONS request.
	OPTIONS(path string, handler IHandler)

	// TRACE Http TRACE method
	// Parameters:
	// - path: The path for the route.
	// - handler: The handler implementing the IHandler interface to process the HTTP TRACE request.
	TRACE(path string, handler IHandler)

	// Group multi routers
	// Parameters:
	// - path: The prefix path for the group.
	// - groupFunc: A function that defines the routes and sub-groups in this group.
	Group(path string, groupFunc func(*Group))
}

// GET is a shortcut for Router.GET(path, handler).
// Parameters:
// - path: The path for the GET route.
// - handler: The handler implementing the IHandler interface to process the HTTP GET request.
func (g *Group) GET(path string, handler IHandler) {
	if g.prefix == "" {
		validatePath(path)
	}

	g.router.GET(g.prefix+path, g.wrapMiddlewares(handler))
}

// HEAD is a shortcut for Router.HEAD(path, handler).
// Parameters:
// - path: The path for the HEAD route.
// - handler: The handler implementing the IHandler interface to process the HTTP HEAD request.
func (g *Group) HEAD(path string, handler IHandler) {
	if g.prefix == "" {
		validatePath(path)
	}

	g.router.HEAD(g.prefix+path, g.wrapMiddlewares(handler))
}

// POST is a shortcut for Router.POST(path, handler).
// Parameters:
// - path: The path for the POST route.
// - handler: The handler implementing the IHandler interface to process the HTTP POST request.
func (g *Group) POST(path string, handler IHandler) {
	if g.prefix == "" {
		validatePath(path)
	}

	g.router.POST(g.prefix+path, g.wrapMiddlewares(handler))
}

// PUT is a shortcut for Router.PUT(path, handler).
// Parameters:
// - path: The path for the PUT route.
// - handler: The handler implementing the IHandler interface to process the HTTP PUT request.
func (g *Group) PUT(path string, handler IHandler) {
	if g.prefix == "" {
		validatePath(path)
	}

	g.router.PUT(g.prefix+path, g.wrapMiddlewares(handler))
}

// PATCH is a shortcut for Router.PATCH(path, handler).
// Parameters:
// - path: The path for the PATCH route.
// - handler: The handler implementing the IHandler interface to process the HTTP PATCH request.
func (g *Group) PATCH(path string, handler IHandler) {
	if g.prefix == "" {
		validatePath(path)
	}

	g.router.PATCH(g.prefix+path, g.wrapMiddlewares(handler))
}

// DELETE is a shortcut for Router.DELETE(path, handler).
// Parameters:
// - path: The path for the DELETE route.
// - handler: The handler implementing the IHandler interface to process the HTTP DELETE request.
func (g *Group) DELETE(path string, handler IHandler) {
	if g.prefix == "" {
		validatePath(path)
	}

	g.router.DELETE(g.prefix+path, g.wrapMiddlewares(handler))
}

// CONNECT is a shortcut for Router.CONNECT(path, handler).
// Parameters:
// - path: The path for the CONNECT route.
// - handler: The handler implementing the IHandler interface to process the HTTP CONNECT request.
func (g *Group) CONNECT(path string, handler IHandler) {
	if g.prefix == "" {
		validatePath(path)
	}

	g.router.CONNECT(g.prefix+path, g.wrapMiddlewares(handler))
}

// OPTIONS is a shortcut for Router.OPTIONS(path, handler).
// Parameters:
// - path: The path for the OPTIONS route.
// - handler: The handler implementing the IHandler interface to process the HTTP OPTIONS request.
func (g *Group) OPTIONS(path string, handler IHandler) {
	if g.prefix == "" {
		validatePath(path)
	}

	g.router.OPTIONS(g.prefix+path, g.wrapMiddlewares(handler))
}

// TRACE is a shortcut for Router.TRACE(path, handler).
// Parameters:
// - path: The path for the TRACE route.
// - handler: The handler implementing the IHandler interface to process the HTTP TRACE request.
func (g *Group) TRACE(path string, handler IHandler) {
	if g.prefix == "" {
		validatePath(path)
	}

	g.router.TRACE(g.prefix+path, g.wrapMiddlewares(handler))
}

// Group creates a new sub-group of routes with a common path prefix.
// Parameters:
// - path: The prefix path for the group.
// - groupFunc: A function that defines the routes and sub-groups within this group.
func (g *Group) Group(path string, groupFunc func(*Group)) {
	group := g.router.Group(g.prefix + path)
	// Automatically append middleware from the parent to child routes.
	// For example, if a parent group has prefix "/user" with middlewares A and B,
	// all handlers in a subgroup with prefix "/info" inherit those middlewares (A, B).
	group.middlewares = g.middlewares

	groupFunc(group)
}

// wrapMiddlewares applies all middlewares in the current group to the provided handler.
//
// Parameters:
// - handler: The handler implementing the IHandler interface to process the HTTP request.
//
// Returns:
// - IHandler: The handler wrapped with all middlewares applied.
func (g *Group) wrapMiddlewares(handler IHandler) IHandler {
	if len(g.middlewares) > 0 {
		middlewareGroup := NewMiddleware()

		return middlewareGroup.Group(g.middlewares...)(handler)
	}

	return handler
}
