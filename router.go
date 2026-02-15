package core

import (
	"fmt"
	"github.com/gflydev/core/errors"
	"github.com/gflydev/core/log"
	"github.com/gflydev/core/utils"
	"net/http"
	"strings"

	"github.com/valyala/bytebufferpool"
	"github.com/valyala/fasthttp"
)

// MethodWild wild HTTP method
const MethodWild = "*"

var (
	questionMark = byte('?')

	// matchedRoutePathParam is the param name under which the path of the matched
	// route is stored, if Router.SaveMatchedRoutePath is set.
	matchedRoutePathParam = fmt.Sprintf("__matchedRoutePath::%s__", utils.RandByte(make([]byte, 15)))
)

// Router is a RequestHandler which can be used to dispatch requests to different
// handler functions via configurable routes

// Router is a high-performance HTTP request router.
// It dispatches requests to handlers based on the HTTP method and the requested path.
// Provides options for path matching, request handling, redirection, and error handling.
type Router struct {
	// trees contains the routing tree for each HTTP method.
	trees []*Tree

	// treeMutable toggles the mutability of the route handlers.
	// If true, route handlers can be updated after being set.
	treeMutable bool

	// methodsIndex maps HTTP methods to their respective indices in the trees.
	// This includes both standard and custom methods.
	methodsIndex map[string]int

	// registeredPaths stores all registered routes grouped by HTTP method.
	registeredPaths map[string][]string

	// SaveMatchedRoutePath, if enabled, adds the matched route path onto
	// the ctx.UserValue context before invoking the handler.
	// The matched route path is only added to handlers of routes registered
	// with this option enabled.
	SaveMatchedRoutePath bool

	// RedirectTrailingSlash Enables automatic redirection if the current route can't be matched but a
	// handler for the path with (without) the trailing slash exists.
	// For example if /foo/ is requested but a route only exists for /foo, the
	// client is redirected to /foo with http status code 301 for GET requests
	// and 308 for all other request methods.
	RedirectTrailingSlash bool

	// RedirectFixedPath enables auto-correction of paths if no handler exists for the requested path.
	// If enabled, the router tries to fix the current request path, if no
	// handle is registered for it.
	// First superfluous path elements like ../ or // are removed.
	// Afterward the router does a case-insensitive lookup of the cleaned path.
	// If a handle can be found for this route, the router makes a redirection
	// to the corrected path with status code 301 for GET requests and 308 for
	// all other request methods.
	// For example /FOO and /..//Foo could be redirected to /foo.
	// RedirectTrailingSlash is independent of this option.
	RedirectFixedPath bool

	// HandleMethodNotAllowed, when enabled, checks if a route is reachable
	// with a different HTTP method when the current request method fails.
	// If enabled, the router checks if another method is allowed for the
	// current route, if the current request can not be routed.
	// If this is the case, the request is answered with 'Method Not Allowed'
	// and HTTP status code 405.
	// If no other Method is allowed, the request is delegated to the NotFound
	// handler.
	HandleMethodNotAllowed bool

	// HandleOPTIONS automatically handles OPTIONS requests when enabled.
	// Custom OPTIONS handlers take precedence over automatic responses.
	HandleOPTIONS bool

	// GlobalOPTIONS is an optional handler invoked for automatic OPTIONS requests
	// when HandleOPTIONS is true and no specific OPTIONS handler is set.
	// The handler is only called if HandleOPTIONS is true and no OPTIONS
	// handler for the specific path was set.
	// The "Allowed" header is set before calling the handler.
	GlobalOPTIONS RequestHandler

	// NotFound defines a custom handler for unmatched requests.
	// Configurable RequestHandler which is called when no matching route is
	// found. If it is not set, default NotFound is used.
	NotFound RequestHandler

	// MethodNotAllowed defines a custom handler for requests with unsupported methods
	// Configurable RequestHandler which is called when a request
	// cannot be routed and HandleMethodNotAllowed is true.
	// If it is not set, ctx.Error with fasthttp.StatusMethodNotAllowed is used.
	// The "Allow" header with allowed request methods is set before the handler is called.
	MethodNotAllowed RequestHandler

	// PanicHandler is a function to recover from panics triggered in HTTP handlers.
	// It should be used to generate a error page and return the http error code
	// 500 (Internal Server Error).
	// The handler can be used to keep your server from crashing because of
	// unrecoverable panics.
	PanicHandler func(*Ctx, any)

	// globalAllowed caches allowed methods for server-wide requests (e.g., for OPTIONS requests on "*").
	globalAllowed string
}

// NewRouter returns a new instance of the Router struct.
//
// This function initializes a Router with default configurations, including:
// - Path auto-correction (both trailing slash correction and fixed path correction) enabled.
// - Automatic handling of OPTIONS requests and checks for allowed methods.
// - PanicHandler for recovering from panics in HTTP handlers.
// - GlobalOPTIONS handler to set default CORS headers for OPTIONS requests.
//
// Returns:
// - *Router: A pointer to the newly created Router instance, with default settings.
func NewRouter() *Router {
	r := &Router{
		trees:                  make([]*Tree, 10),
		methodsIndex:           make(map[string]int),
		registeredPaths:        make(map[string][]string),
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      true,
		HandleMethodNotAllowed: true,
		HandleOPTIONS:          true,
		PanicHandler: func(ctx *Ctx, data any) {
			ctx.Status(fasthttp.StatusInternalServerError)

			contentType := strings.ToLower(utils.UnsafeStr(ctx.root.Request.Header.ContentType()))
			if strings.HasPrefix(contentType, MIMEApplicationJSON) {
				_ = ctx.JSON(Data{
					"code":    http.StatusInternalServerError,
					"message": "Internal Server Error",
				})
			} else {
				_ = ctx.HTML("<strong>Internal Server Error</strong>")
			}

			// Logs error details using the specified log levels.
			log.Errorf("PanicHandler:: \n    - Context %s \n    - Error %v", ctx.root.String(), data)
		},
		GlobalOPTIONS: func(ctx *Ctx) error {
			// Set CORs headers
			ctx.root.Response.Header.Set(HeaderAccessControlAllowOrigin, "*")
			ctx.root.Response.Header.Set(HeaderAccessControlAllowMethods, "PUT, POST, GET, DELETE, OPTIONS, PATCH")
			ctx.root.Response.Header.Set(HeaderAccessControlAllowHeaders, "Authorization, Content-Type, x-requested-with, origin, true-client-ip, X-Correlation-ID")

			return nil
		},
	}

	// Initialize standard method indices
	r.methodsIndex[fasthttp.MethodGet] = 0
	r.methodsIndex[fasthttp.MethodHead] = 1
	r.methodsIndex[fasthttp.MethodPost] = 2
	r.methodsIndex[fasthttp.MethodPut] = 3
	r.methodsIndex[fasthttp.MethodPatch] = 4
	r.methodsIndex[fasthttp.MethodDelete] = 5
	r.methodsIndex[fasthttp.MethodConnect] = 6
	r.methodsIndex[fasthttp.MethodOptions] = 7
	r.methodsIndex[fasthttp.MethodTrace] = 8
	r.methodsIndex[MethodWild] = 9

	return r
}

// Group creates a new route group with the specified base path.
//
// Route groups are a convenient way to group related routes under a common
// prefix. All routes registered through the returned group will have their paths
// automatically prefixed with the group's base path.
//
// The group's base path should not end with a trailing slash, except for the root path ("/").
// If the provided path violates this rule, the function will panic.
//
// Parameters:
//   - path: string — The base path for the group.
//
// Returns:
//   - *Group: A pointer to the newly created route group, with the specified base path.
func (r *Router) Group(path string) *Group {
	validatePath(path)

	if path != "/" && strings.HasSuffix(path, "/") {
		panic("group path must not end with a trailing slash")
	}

	return &Group{
		router: r,
		prefix: path,
	}
}

// saveMatchedRoutePath wraps the given handler with a handler that saves the matched route path.
//
// Parameters:
//   - path: The matched route path to save.
//   - handler: The original handler to wrap.
//
// Returns:
//   - IHandler: A new handler that saves the matched route path before delegating to the original handler.
func (r *Router) saveMatchedRoutePath(path string, handler IHandler) IHandler {
	return &saveMatchedRoutePathHandler{
		path:    path,
		handler: handler,
	}
}

// saveMatchedRoutePathHandler is a handler that saves the matched route path to the context.
type saveMatchedRoutePathHandler struct {
	Endpoint
	// path is the matched route path to be saved.
	path string
	// handler is the original handler to delegate to after saving the matched route path.
	handler IHandler
}

// Handle saves the matched route path to the context and invokes the wrapped handler.
//
// Parameters:
//   - c: The current request context.
//
// Returns:
//   - error: An error, if any, returned by the wrapped handler.
func (e *saveMatchedRoutePathHandler) Handle(c *Ctx) error {
	c.root.SetUserValue(matchedRoutePathParam, e.path)
	return e.handler.Handle(c)
}

// methodIndexOf returns the index corresponding to the given HTTP method.
// This index is used to look up the appropriate routing tree.
//
// Parameters:
//   - method: string — The HTTP method (e.g., GET, POST) for which the index is required.
//
// Returns:
//   - int: The index of the given method. Returns -1 if the method is not recognized or registered.
func (r *Router) methodIndexOf(method string) int {
	if i, ok := r.methodsIndex[method]; ok {
		return i
	}

	return -1
}

// Mutable sets whether the routing tree allows updates to the route handlers.
//
// By default, this is disabled, meaning the routing tree is immutable once
// initialized. Enabling it allows modifications to the route handlers.
//
// WARNING: This feature should be used cautiously, as it may lead to
// unexpected behaviors in concurrent scenarios.
//
// Parameters:
//   - v: bool — A boolean value indicating whether the routing tree should be mutable.
//
// Returns:
//   - None
func (r *Router) Mutable(v bool) {
	r.treeMutable = v

	for i := range r.trees {
		tree := r.trees[i]

		if tree != nil {
			tree.Mutable = v
		}
	}
}

// List returns all registered routes grouped by method.
//
// Returns:
//   - map[string][]string: A map where the key is the HTTP method and the value is a slice of registered route paths.
func (r *Router) List() map[string][]string {
	return r.registeredPaths
}

// GET is a shortcut for router.Handle with the GET HTTP method.
//
// Parameters:
//   - path: string — The route path.
//   - handler: IHandler — The handler function for the route.
func (r *Router) GET(path string, handler IHandler) {
	r.Handle(fasthttp.MethodGet, path, handler)
}

// HEAD is a shortcut for router.Handle with the HEAD HTTP method.
//
// Parameters:
//   - path: string — The route path.
//   - handler: IHandler — The handler function for the route.
func (r *Router) HEAD(path string, handler IHandler) {
	r.Handle(fasthttp.MethodHead, path, handler)
}

// POST is a shortcut for router.Handle with the POST HTTP method.
//
// Parameters:
//   - path: string — The route path.
//   - handler: IHandler — The handler function for the route.
func (r *Router) POST(path string, handler IHandler) {
	r.Handle(fasthttp.MethodPost, path, handler)
}

// PUT is a shortcut for router.Handle with the PUT HTTP method.
//
// Parameters:
//   - path: string — The route path.
//   - handler: IHandler — The handler function for the route.
func (r *Router) PUT(path string, handler IHandler) {
	r.Handle(fasthttp.MethodPut, path, handler)
}

// PATCH is a shortcut for router.Handle with the PATCH HTTP method.
//
// Parameters:
//   - path: string — The route path.
//   - handler: IHandler — The handler function for the route.
func (r *Router) PATCH(path string, handler IHandler) {
	r.Handle(fasthttp.MethodPatch, path, handler)
}

// DELETE is a shortcut for router.Handle with the DELETE HTTP method.
//
// Parameters:
//   - path: string — The route path.
//   - handler: IHandler — The handler function for the route.
func (r *Router) DELETE(path string, handler IHandler) {
	r.Handle(fasthttp.MethodDelete, path, handler)
}

// CONNECT is a shortcut for router.Handle with the CONNECT HTTP method.
//
// Parameters:
//   - path: string — The route path.
//   - handler: IHandler — The handler function for the route.
func (r *Router) CONNECT(path string, handler IHandler) {
	r.Handle(fasthttp.MethodConnect, path, handler)
}

// OPTIONS is a shortcut for router.Handle with the OPTIONS HTTP method.
//
// Parameters:
//   - path: string — The route path.
//   - handler: IHandler — The handler function for the route.
func (r *Router) OPTIONS(path string, handler IHandler) {
	r.Handle(fasthttp.MethodOptions, path, handler)
}

// TRACE is a shortcut for router.Handle with the TRACE HTTP method.
//
// Parameters:
//   - path: string — The route path.
//   - handler: IHandler — The handler function for the route.
func (r *Router) TRACE(path string, handler IHandler) {
	r.Handle(fasthttp.MethodTrace, path, handler)
}

// ServeFiles serves files from the specified root directory for the given route path.
//
// The route `path` must end with "/{filepath:*}", allowing dynamic file serving
// based on the requested filepath. Files are served from the local filesystem
// rooted at the given `rootPath`. For example, if `rootPath` is "/etc" and the
// requested {filepath:*} is "passwd", the handler serves the file located at
// "/etc/passwd".
//
// This function sets up a file server using default settings, such as generating
// index pages and accepting byte ranges when serving files.
//
// Parameters:
//   - path: string — The route path where files will be served. Must end with "/{filepath:*}".
//   - rootPath: string — The root directory of the file system to serve files from.
//
// Example:
//
//	router.ServeFiles("/static/{filepath:*}", "./public")
func (r *Router) ServeFiles(path, rootPath string) {
	r.ServeFilesCustom(path, &fasthttp.FS{
		Root:               rootPath,
		IndexNames:         []string{"index.html"},
		GenerateIndexPages: true,
		AcceptByteRange:    true,
	})
}

// ServeFilesCustom serves files from the given file system settings with custom configurations.
//
// This function allows for serving files from a specific file system (`fasthttp.FS`)
// and provides control over indexing, path rewriting, byte range requests,
// and other customizations via the `fasthttp.FS` configuration.
//
// The route `path` must end with "/{filepath:*}", enabling dynamic serving of
// requested file paths. If the route does not conform to this format, it will panic.
//
// If `fs.PathRewrite` is nil and the route `path` contains slashes, the slashes
// are automatically stripped from the path and used for path rewriting.
//
// Parameters:
//   - path: string — The route path where files will be served. Must end with "/{filepath:*}".
//   - fs: *fasthttp.FS — The custom file system configuration used to serve files.
//
// Example:
//
//	router.ServeFilesCustom("/assets/{filepath:*}", &fasthttp.FS{
//		Root:			   "./assets",
//		IndexNames:		 []string{"index.html"},
//		GenerateIndexPages: false,
//		AcceptByteRange:	true,
//	})
//
// The example above will serve files from the "assets" directory for requests
// that match the route path "/assets/{filepath:*}".
//
// Return:
//
//	No return value. This function registers the route with the router.
func (r *Router) ServeFilesCustom(path string, fs *fasthttp.FS) {
	suffix := "/{filepath:*}"

	if !strings.HasSuffix(path, suffix) {
		panic("path must end with " + suffix + " in path '" + path + "'")
	}

	prefix := path[:len(path)-len(suffix)]
	stripSlashes := strings.Count(prefix, "/")

	if fs.PathRewrite == nil && stripSlashes > 0 {
		fs.PathRewrite = fasthttp.NewPathSlashesStripper(stripSlashes)
	}
	handler := &serveFilesCustomHandler{
		fileHandler: fs.NewRequestHandler(),
	}

	r.GET(path, handler)
}

// serveFilesCustomHandler is a handler for serving files with custom configurations.
type serveFilesCustomHandler struct {
	Endpoint
	fileHandler fasthttp.RequestHandler
}

// Handle processes the incoming request and invokes the file handler.
//
// Parameters:
//   - c: *Ctx — The current request context.
//
// Returns:
//   - error: Always returns nil, as the underlying file handler does not propagate errors.
func (e *serveFilesCustomHandler) Handle(c *Ctx) error {
	e.fileHandler(c.root)

	if (c.root.Response.StatusCode() == fasthttp.StatusNotFound ||
		c.root.Response.StatusCode() == fasthttp.StatusForbidden) &&
		c.Path() != "/" {

		c.root.ResetBody()
		log.Warn("Request not found: ", c.Path())

		if strings.HasPrefix(strings.ToLower(utils.UnsafeStr(c.root.Request.Header.ContentType())), MIMEApplicationJSON) {
			return c.JSON(Data{
				"code": 404,
				"msg":  "404 Not Found",
			})
		}
		return c.View("404", Data{
			"title_page": "404 Not Found",
			"msg":        utils.UnsafeStr(c.root.Response.Body()),
		})
	}

	return nil
}

// Handle registers a new request handler with the given path and method.
//
// This function is intended to register a handler for a specified HTTP method and route path.
// It provides validation for the method and handler parameters, as well as handling optional routes.
//
// If the HTTP method is empty or the handler is nil, the function will panic.
// It also handles the creation of internal structures for method-based routing
// and manages optional route paths (e.g., paths with optional parameters).
//
// Parameters:
//   - method: string — The HTTP method (e.g., "GET", "POST", etc.) for which the handler is registered.
//   - path: string — The route path that the handler will handle. Must be a valid path.
//   - handler: IHandler — The handler interface responsible for processing requests.
//
// Return:
//
//   - No return value. This function registers the specified route and handler.
func (r *Router) Handle(method, path string, handler IHandler) {
	switch {
	case method == "":
		panic("method must not be empty")
	case handler == nil:
		panic("handler must not be nil")
	default:
		validatePath(path)
	}

	r.registeredPaths[method] = append(r.registeredPaths[method], path)

	methodIndex := r.methodIndexOf(method)
	if methodIndex == -1 {
		tree := NewTree()
		tree.Mutable = r.treeMutable

		r.trees = append(r.trees, tree)
		methodIndex = len(r.trees) - 1
		r.methodsIndex[method] = methodIndex
	}

	tree := r.trees[methodIndex]
	if tree == nil {
		tree = NewTree()
		tree.Mutable = r.treeMutable

		r.trees[methodIndex] = tree
		r.globalAllowed = r.allowed("*", "")
	}

	if r.SaveMatchedRoutePath {
		handler = r.saveMatchedRoutePath(path, handler)
	}

	optionalPaths := getOptionalPaths(path)

	// if not has optional paths, adds the original
	if len(optionalPaths) == 0 {
		tree.Add(path, handler)
	} else {
		for _, p := range optionalPaths {
			tree.Add(p, handler)
		}
	}
}

// Lookup retrieves a registered handler for a given HTTP method and route path.
//
// This method is useful for manually inspecting the routing tree, which can be
// helpful for building frameworks or debugging purposes. It returns the handler
// function for the provided method and path, if found.
//
// Parameters:
//   - method: string — The HTTP method (e.g., "GET", "POST", etc.) for which to perform the lookup.
//   - path: string — The route path to look up in the router.
//   - ctx: *Ctx — The current request context, potentially providing additional data
//     necessary for handling the request during lookup.
//
// Returns:
//   - IHandler: The handler interface for processing requests if the method and path match.
//   - bool: Indicates whether a trailing slash redirection should be performed in case of a mismatch.
func (r *Router) Lookup(method, path string, ctx *Ctx) (IHandler, bool) {
	methodIndex := r.methodIndexOf(method)
	if methodIndex == -1 {
		return nil, false
	}

	if tree := r.trees[methodIndex]; tree != nil {
		handler, tsr := tree.Get(path, ctx)
		if handler != nil || tsr {
			return handler, tsr
		}
	}

	if tree := r.trees[r.methodIndexOf(MethodWild)]; tree != nil {
		return tree.Get(path, ctx)
	}

	return nil, false
}

// recv handles recovery from panics and invokes the PanicHandler.
//
// Parameters:
//   - ctx: *Ctx — The current request context.
//
// Return:
//   - No return value.
func (r *Router) recv(ctx *Ctx) {
	if rcv := recover(); rcv != nil {
		r.PanicHandler(ctx, rcv)
	}
}

// allowed determines and returns the allowed HTTP methods for a specific path.
//
// Parameters:
//   - path: string — The route path to check allowed methods for.
//   - reqMethod: string — The HTTP method of the incoming request.
//
// Returns:
//   - string: A comma-separated list of allowed HTTP methods for the given path.
func (r *Router) allowed(path, reqMethod string) (allow string) {
	allowed := make([]string, 0, 9)

	if path == "*" || path == "/*" { // server-wide routes
		// empty method is used for internal calls to refresh the cache
		if reqMethod == "" {
			for method := range r.registeredPaths {
				if method == fasthttp.MethodOptions {
					continue
				}
				// Add request method to list of allowed methods
				allowed = append(allowed, method)
			}
		} else {
			return r.globalAllowed
		}
	} else { // specific path
		for method := range r.registeredPaths {
			// Skip the requested method - we already tried this one
			if method == reqMethod || method == fasthttp.MethodOptions {
				continue
			}

			handle, _ := r.trees[r.methodIndexOf(method)].Get(path, nil)
			if handle != nil {
				// Add request method to list of allowed methods
				allowed = append(allowed, method)
			}
		}
	}

	if len(allowed) > 0 {
		// Add request method to list of allowed methods
		allowed = append(allowed, fasthttp.MethodOptions)

		// Sort allowed methods.
		// sort.Strings(allowed) unfortunately causes unnecessary allocations
		// due to allowed being moved to the heap and interface conversion
		for i, l := 1, len(allowed); i < l; i++ {
			for j := i; j > 0 && allowed[j] < allowed[j-1]; j-- {
				allowed[j], allowed[j-1] = allowed[j-1], allowed[j]
			}
		}

		// return as comma-separated list
		return strings.Join(allowed, ", ")
	}
	return
}

// tryRedirect attempts to redirect the request based on the given path and method.
//
// Parameters:
//   - ctx: *Ctx — The current request context.
//   - tree: *Tree — The routing tree used for resolving paths.
//   - tsr: bool — Indicates if a trailing slash redirect should be performed.
//   - method: string — The HTTP method of the incoming request.
//   - path: string — The original request path.
//
// Returns:
//   - bool: Returns true if a redirect was performed, otherwise false.
func (r *Router) tryRedirect(ctx *Ctx, tree *Tree, tsr bool, method, path string) bool {
	// Moved Permanently, request with GET method
	code := fasthttp.StatusMovedPermanently
	if method != fasthttp.MethodGet {
		// Permanent Redirect, request with same method
		code = fasthttp.StatusPermanentRedirect
	}

	// Get a buffer from the pool once and reuse it
	uri := bytebufferpool.Get()
	defer bytebufferpool.Put(uri)

	// Handle trailing slash redirects
	if tsr && r.RedirectTrailingSlash {
		// Reset buffer to ensure it's empty
		uri.Reset()

		if len(path) > 1 && path[len(path)-1] == '/' {
			uri.SetString(path[:len(path)-1])
		} else {
			uri.SetString(path)
			if uri.WriteByte('/') != nil {
				return false
			}
		}

		// Add query string if present
		if queryBuf := ctx.root.URI().QueryString(); len(queryBuf) > 0 {
			if uri.WriteByte(questionMark) != nil {
				return false
			}
			if _, err := uri.Write(queryBuf); err != nil {
				return false
			}
		}

		ctx.root.Redirect(uri.String(), code)
		return true
	}

	// Try to fix the request path
	if r.RedirectFixedPath {
		path2 := utils.UnsafeStr(ctx.root.Request.URI().Path())

		// Reset buffer to ensure it's empty
		uri.Reset()

		found := tree.FindCaseInsensitivePath(
			cleanPath(path2),
			r.RedirectTrailingSlash,
			uri,
		)

		if found {
			// Add query string if present
			if queryBuf := ctx.root.URI().QueryString(); len(queryBuf) > 0 {
				if uri.WriteByte(questionMark) != nil {
					return false
				}
				if _, err := uri.Write(queryBuf); err != nil {
					return false
				}
			}

			ctx.root.Redirect(uri.String(), code)
			return true
		}
	}

	return false
}

// Handler processes incoming HTTP requests and routes them to the appropriate handler.
// It supports panic recovery, request validation, method-based routing, and error handling.
//
// Parameters:
//   - ctx: *Ctx — The context of the current request, containing request and response details.
//
// Returns:
//   - error: Returns any error encountered during request handling or validation. If no error occurs, returns nil.
func (r *Router) Handler(ctx *Ctx) error {
	if r.PanicHandler != nil {
		defer r.recv(ctx)
	}

	path := utils.UnsafeStr(ctx.root.Request.URI().PathOriginal())
	method := utils.UnsafeStr(ctx.root.Request.Header.Method())
	methodIndex := r.methodIndexOf(method)

	if methodIndex > -1 {
		if tree := r.trees[methodIndex]; tree != nil {
			if handler, tsr := tree.Get(path, ctx); handler != nil {
				if err := handler.Validate(ctx); err != nil {
					return errorHandler(ctx, err, StatusBadRequest)
				}

				if err := handler.Handle(ctx); err != nil {
					return errorHandler(ctx, err, StatusInternalServerError)
				}

				return nil
			} else if method != fasthttp.MethodConnect && path != "/" {
				if ok := r.tryRedirect(ctx, tree, tsr, method, path); ok {
					return nil
				}
			}
		}
	}

	// Try to search in the wild method tree
	if tree := r.trees[r.methodIndexOf(MethodWild)]; tree != nil {
		if handler, tsr := tree.Get(path, ctx); handler != nil {
			if err := handler.Validate(ctx); err != nil {
				return errorHandler(ctx, err, StatusBadRequest)
			}

			if err := handler.Handle(ctx); err != nil {
				return errorHandler(ctx, err, StatusInternalServerError)
			}

			return nil
		} else if method != fasthttp.MethodConnect && path != "/" {
			if ok := r.tryRedirect(ctx, tree, tsr, method, path); ok {
				return nil
			}
		}
	}

	if r.HandleOPTIONS && method == fasthttp.MethodOptions {
		// Handle OPTIONS requests

		if allow := r.allowed(path, fasthttp.MethodOptions); allow != "" {
			ctx.root.Response.Header.Set("Allow", allow)
			if r.GlobalOPTIONS != nil {
				return r.GlobalOPTIONS(ctx)
			}
			return nil
		}
	} else if r.HandleMethodNotAllowed {
		// Handle 405

		if allow := r.allowed(path, method); allow != "" {
			ctx.root.Response.Header.Set("Allow", allow)
			if r.MethodNotAllowed != nil {
				return r.MethodNotAllowed(ctx)
			} else {
				ctx.root.SetStatusCode(fasthttp.StatusMethodNotAllowed)
				ctx.root.SetBodyString(fasthttp.StatusMessage(fasthttp.StatusMethodNotAllowed))
			}
			return nil
		}
	}

	// Handle 404
	if r.NotFound != nil {
		return r.NotFound(ctx)
	} else {
		ctx.root.Error(fasthttp.StatusMessage(fasthttp.StatusNotFound), fasthttp.StatusNotFound)
	}

	return nil
}

// errorHandler handles request errors and sets the appropriate response status and body.
//
// Parameters:
//   - ctx: *Ctx — The current request context containing request and response details.
//   - err: error — The error that occurred during request processing.
//   - code: int — The HTTP status code to be set in the response.
//
// Returns:
//   - error: The same error that was passed in after handling.
func errorHandler(ctx *Ctx, err error, code int) error {
	// Force status error
	if ctx.root.Response.StatusCode() == StatusOK {
		ctx.Status(code)
	}

	if !errors.Is(err, errors.UnknownError) {
		log.Errorf("Bad request: %v", err)
	}

	// Set default body
	body := ctx.root.Response.Body()
	if len(body) == 0 {
		contentType := strings.ToLower(utils.UnsafeStr(ctx.root.Request.Header.ContentType()))
		if strings.HasPrefix(contentType, MIMEApplicationJSON) {
			_ = ctx.JSON(Data{
				"error": err.Error(),
			})
		} else {
			_ = ctx.HTML(err.Error())
		}
	}

	return err
}
