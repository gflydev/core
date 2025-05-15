package core

import (
	"fmt"
	"github.com/gflydev/core/container"
	"github.com/gflydev/core/log"
	"github.com/gflydev/core/plugin"
	"github.com/gflydev/core/utils"
	"github.com/valyala/fasthttp"
)

// ====================================================================
//                              gFly
// ====================================================================

var (
	// Application

	// AppName holds the application name fetched from `APP_NAME` environment variable or defaults to "gFly".
	AppName = utils.Getenv("APP_NAME", "gFly")
	// AppCode holds the application code fetched from `APP_CODE` environment variable or defaults to "gfly".
	AppCode = utils.Getenv("APP_CODE", "gfly")
	// AppURL holds the application URL fetched from `APP_URL` environment variable or defaults to "http://localhost:7789".
	AppURL = utils.Getenv("APP_URL", "http://localhost:7789")
	// AppEnv holds the application environment fetched from `APP_ENV` environment variable or defaults to "local".
	AppEnv = utils.Getenv("APP_ENV", "local")
	// AppDebug indicates if the application is in debug mode, fetched from `APP_DEBUG` environment variable or defaults to true.
	AppDebug = utils.Getenv("APP_DEBUG", true)

	// Storage directory

	// StorageDir specifies the directory for application storage, fetched from `STORAGE_DIR` environment variable or defaults to "storage".
	StorageDir = utils.Getenv("STORAGE_DIR", "storage") // Directory `{APP}/storage`
	// TempDir specifies the temporary directory within the storage, fetched from `TEMP_DIR` environment variable or defaults to "storage/tmp".
	TempDir = utils.Getenv("TEMP_DIR", "storage/tmp") // Directory `{APP}/storage/temp`
	// LogDir specifies the directory for logs within the storage, fetched from `LOG_DIR` environment variable or defaults to "storage/logs".
	LogDir = utils.Getenv("LOG_DIR", "storage/logs") // Directory `{APP}/storage/log`
	// AppDir specifies the application-specific directory within the storage, fetched from `APP_DIR` environment variable or defaults to "storage/app".
	AppDir = utils.Getenv("APP_DIR", "storage/app") // Directory `{APP}/storage/app`

	// Internal variables

	// fnHookMiddlewares holds the list of middleware functions to be executed globally.
	fnHookMiddlewares []FnHookMiddleware // Hook global middlewares
	// fnHookRoutes holds the list of route hook functions to define application routes.
	fnHookRoutes []FnHookRoute // Hook routers
)

// FnHookMiddleware Function type for registering middleware hooks.
// Parameters:
//   - fly: An instance of IFlyMiddleware used for middleware setup.
type FnHookMiddleware func(fly IFlyMiddleware)

// FnHookRoute Function type for registering route hooks.
// Parameters:
//   - fly: An instance of IFly used for route setup.
type FnHookRoute func(fly IFly)

// IFly Interface to declare all methods for gFly struct.
type IFly interface {
	// Run Starts the application.
	Run()

	// Router Gets the web router instance.
	// Returns:
	//   - *Router: The root router instance.
	Router() *Router

	// Container returns the dependency injection container.
	// Returns:
	//   - *container.Container: The dependency injection container.
	Container() *container.Container

	// RegisterMiddleware Registers middleware hooks.
	// Parameters:
	//   - fn: Variadic parameter of FnHookMiddleware functions to set up global middleware.
	RegisterMiddleware(fn ...FnHookMiddleware)

	// RegisterRouter Registers route hooks.
	// Parameters:
	//   - fn: Variadic parameter of FnHookRoute functions to set up application routes.
	RegisterRouter(fn ...FnHookRoute)

	// Inherits methods from IFlyRouter, IFlyMiddleware, and IFlyPlugin.

	IFlyRouter
	IFlyMiddleware
	plugin.IFlyPlugin
}

// GFly Struct defining main elements in the application.
type GFly struct {
	// router Keep reference to the root router.
	router *Router
	// server Keep reference to the web server.
	server *fasthttp.Server
	// config App configuration settings.
	config Config
	// middleware Keep reference to the middleware manager.
	middleware IMiddleware
	// middlewares List of global middleware handlers.
	middlewares []MiddlewareHandler
	// container Dependency injection container for service management.
	container *container.Container
	// pluginManager Plugin manager for handling plugins.
	pluginManager *plugin.Manager
}

// Router Retrieves the root router in the gFly application.
// Returns:
//   - *Router: The root router instance of the application.
func (fly *GFly) Router() *Router {
	return fly.router
}

// Run Start gFly app.
func (fly *GFly) Run() {
	// --------------- Setup Server ---------------
	fly.server = &fasthttp.Server{
		Handler:                       fasthttp.CompressHandler(fly.serveFastHTTP),
		ErrorHandler:                  fly.errorHandler,
		Name:                          fly.config.Name,
		Concurrency:                   fly.config.Concurrency,
		ReadTimeout:                   fly.config.ReadTimeout,
		WriteTimeout:                  fly.config.WriteTimeout,
		IdleTimeout:                   fly.config.IdleTimeout,
		ReadBufferSize:                fly.config.ReadBufferSize,
		WriteBufferSize:               fly.config.WriteBufferSize,
		NoDefaultDate:                 fly.config.NoDefaultDate,
		NoDefaultContentType:          fly.config.NoDefaultContentType,
		DisableHeaderNamesNormalizing: fly.config.DisableHeaderNamesNormalizing,
		DisableKeepalive:              fly.config.DisableKeepalive,
		MaxRequestBodySize:            fly.config.MaxRequestBodySize,
		NoDefaultServerHeader:         fly.config.NoDefaultServerHeader, // True when `Name` Empty
		GetOnly:                       fly.config.GetOnly,
		ReduceMemoryUsage:             fly.config.ReduceMemoryUsage,
		StreamRequestBody:             fly.config.StreamRequestBody,
		DisablePreParseMultipartForm:  fly.config.DisablePreParseMultipartForm,
	}

	url := fmt.Sprintf(
		"%s:%v",
		utils.Getenv("SERVER_HOST", "0.0.0.0"),
		utils.Getenv("SERVER_PORT", 7789),
	)

	// --------------- Print startup message ---------------
	if !fly.config.DisableStartupMessage {
		startupMessage(url, AppName, AppEnv)
	}

	// --------------- Setup Logs ---------------
	setupLog()

	// --------------- Checking service  ---------------
	// TODO: Need to add more checking

	// --------------- Initialize plugins ---------------
	if err := plugin.InitializePlugins(fly.container); err != nil {
		log.Fatalf("Error initializing plugins: %v", err)
	}

	// --------------- Global middlewares  ---------------
	for _, fn := range fnHookMiddlewares {
		fn(fly)
	}

	// --------------- Router  ---------------
	for _, fn := range fnHookRoutes {
		fn(fly)
	}

	// --------------- Serve static file ---------------
	serveFiles(fly)

	certFile := utils.Getenv("SERVER_TLS_CERT", "")
	keyFile := utils.Getenv("SERVER_TLS_KEY", "")

	switch {
	case certFile != "" && keyFile != "":
		if err := fly.server.ListenAndServeTLS(url, certFile, keyFile); err != nil {
			log.Fatalf("Error start server %v", err)
		}
	default:
		log.Fatal(fly.server.ListenAndServe(url))
	}
}

// serveFastHTTP Serve FastHTTP via HTTP function
// The linking between fasthttp.RequestHandler to gFly's Ctx.
// Parameters:
//   - ctx (*fasthttp.RequestCtx): The context of the current HTTP request.
func (fly *GFly) serveFastHTTP(ctx *fasthttp.RequestCtx) {
	// handlerCtx encapsulates gFly application context, the root HTTP context,
	// and custom data storage for the request.
	handlerCtx := &Ctx{
		app:  fly,    // Reference to the gFly application instance.
		root: ctx,    // Reference to the current HTTP request context.
		data: Data{}, // Custom data storage initialized as an empty map.
	}

	// Pass the request context to the router for further processing.
	_ = fly.router.Handler(handlerCtx)
}

// errorHandler Server error handler.
// Logs errors with Debug and Error log levels.
// Parameters:
//   - ctx (*fasthttp.RequestCtx): The context of the current HTTP request.
//   - err (error): The error that has occurred.
func (fly *GFly) errorHandler(ctx *fasthttp.RequestCtx, err error) {
	// Logs error details using the specified log levels.
	log.Debugf("Error %s", ctx.String()) // Debug log with context details.
	log.Errorf("Error happens %v", err)  // Error log with error details.
}

// New Create a new gFly application instance.
// Parameters:
//   - config (variadic Config): Optional application configuration(s) for overriding defaults.
//
// Returns:
//   - IFly: Instance of gFly application.
func New(config ...Config) IFly {
	// Create a new dependency injection container
	diContainer := container.New()

	// Create a new GFly instance with initialized router and middleware.
	app := &GFly{
		router:     NewRouter(),     // Instantiate a new Router for the app.
		middleware: NewMiddleware(), // Instantiate a new Middleware manager for the app.
		container:  diContainer,     // Instantiate a new dependency injection container.
	}

	// Override default configuration if additional config is provided.
	if len(config) > 0 {
		app.config = config[0] // Use the first provided configuration.
	} else {
		app.config = DefaultConfig // Otherwise, use the default configuration.
	}

	// Register core services in the container
	app.registerCoreServices()

	// Initialize the plugin manager
	app.pluginManager = plugin.NewManager(diContainer)

	return app // Return the initialized gFly application instance.
}

// RegisterMiddleware Register middleware hooks to process global middleware.
// Parameters:
//   - fn (variadic FnHookMiddleware): List of middleware hook functions to register globally.
func (fly *GFly) RegisterMiddleware(fn ...FnHookMiddleware) {
	fnHookMiddlewares = fn // Assign provided middleware hooks to the global list.
}

// RegisterRouter Register route hooks to define application routes.
// Parameters:
//   - fn (variadic FnHookRoute): List of route hook functions to register globally.
func (fly *GFly) RegisterRouter(fn ...FnHookRoute) {
	fnHookRoutes = fn // Assign provided route hooks to the global list.
}

// registerCoreServices registers the core services of the framework in the container.
func (fly *GFly) registerCoreServices() {
	// Register the router
	_ = container.RegisterInstance[*Router](fly.container, fly.router)

	// Register the middleware manager
	_ = container.RegisterInstance[IMiddleware](fly.container, fly.middleware)

	// Register the application itself
	_ = container.RegisterInstance[IFly](fly.container, fly)

	// Register the configuration
	_ = container.RegisterInstance[Config](fly.container, fly.config)

	// Register the logger (using the default logger)
	_ = container.RegisterInstance[log.AllLogger](fly.container, log.DefaultLogger())

	// Register the plugin manager
	_ = plugin.RegisterPluginManager(fly.container)
}

// Container returns the dependency injection container.
//
// Returns:
//   - *container.Container: The dependency injection container.
func (fly *GFly) Container() *container.Container {
	return fly.container
}

// ====================================================================
//                        gFly - Middleware methods
// ====================================================================

// IFlyMiddleware Interface to declare all Middleware methods for gFly struct.
type IFlyMiddleware interface {
	// Use adds middleware for global (all requests).
	//
	// Parameters:
	//   - middlewares ([]MiddlewareHandler): One or more middleware handlers to be applied globally.
	Use(middlewares ...MiddlewareHandler)

	// UseWithOptions adds middleware with configuration options for global (all requests).
	//
	// Parameters:
	//   - handler (MiddlewareHandler): The middleware handler to be applied globally.
	//   - options ([]MiddlewareOption): Configuration options for the middleware.
	UseWithOptions(handler MiddlewareHandler, options ...MiddlewareOption)

	// Middleware creates a middleware chain handler for grouping.
	//
	// Parameters:
	//   - middleware ([]MiddlewareHandler): Middleware handlers to be grouped.
	//
	// Returns:
	//   - func(IHandler) IHandler: A function that applies the middleware handlers to an IHandler.
	Middleware(middleware ...MiddlewareHandler) func(IHandler) IHandler

	// MiddlewareWithOptions creates a middleware chain handler with configuration options.
	//
	// Parameters:
	//   - configs ([]MiddlewareConfig): Middleware configurations to be grouped.
	//
	// Returns:
	//   - func(IHandler) IHandler: A function that applies the middleware handlers to an IHandler.
	MiddlewareWithOptions(configs ...MiddlewareConfig) func(IHandler) IHandler
}

// Use adds middleware for global (all requests).
//
// Example usage:
//
//	group.Use(middleware.RuleMiddlewareFunc, middleware.AuthMiddlewareFunc)
//
// Parameters:
//   - middlewares ([]MiddlewareHandler): One or more middleware handlers to be added globally.
func (fly *GFly) Use(middlewares ...MiddlewareHandler) {
	// Append the provided middleware handlers to the global middleware list.
	fly.middlewares = append(fly.middlewares, middlewares...)
}

// UseWithOptions adds middleware with configuration options for global (all requests).
//
// Example usage:
//
//	fly.UseWithOptions(authMiddleware, WithName("auth"), WithPriority(1))
//
// Parameters:
//   - handler (MiddlewareHandler): The middleware handler to be applied globally.
//   - options ([]MiddlewareOption): Configuration options for the middleware.
func (fly *GFly) UseWithOptions(handler MiddlewareHandler, options ...MiddlewareOption) {
	// Create a middleware configuration with the provided options
	config := fly.middleware.Use(handler, options...)

	// Convert the middleware configuration to a handler and append it to the global middleware list
	fly.middlewares = append(fly.middlewares, config.Handler)
}

// Middleware creates a middleware chain handler for grouping.
//
// Example usage:
//
//	group.POST("/one", gfly.IFly.Middleware(middleware.RuleMiddlewareFunc)(api.NewDefaultApi()))
//
// Parameters:
//   - middlewares ([]MiddlewareHandler): Middleware handlers to be grouped.
//
// Returns:
//   - func(IHandler) IHandler: A function that applies the grouped middleware handlers to an IHandler.
func (fly *GFly) Middleware(middlewares ...MiddlewareHandler) func(IHandler) IHandler {
	// Group the provided middleware handlers and return a function for applying them to an IHandler.
	return fly.middleware.Group(middlewares...)
}

// MiddlewareWithOptions creates a middleware chain handler with configuration options.
//
// Example usage:
//
//	authConfig := fly.middleware.Use(authMiddleware, WithName("auth"), WithPriority(1))
//	logConfig := fly.middleware.Use(logMiddleware, WithName("log"), WithPhase(PhasePostRequest))
//	group.POST("/one", fly.MiddlewareWithOptions(authConfig, logConfig)(api.NewDefaultApi()))
//
// Parameters:
//   - configs ([]MiddlewareConfig): Middleware configurations to be grouped.
//
// Returns:
//   - func(IHandler) IHandler: A function that applies the middleware handlers to an IHandler.
func (fly *GFly) MiddlewareWithOptions(configs ...MiddlewareConfig) func(IHandler) IHandler {
	// Group the provided middleware configurations and return a function for applying them to an IHandler.
	return fly.middleware.GroupWithOptions(configs...)
}

// ====================================================================
//                        gFly - HTTP methods
// ====================================================================

// IFlyRouter Interface to declare all HTTP methods for gFly struct.
type IFlyRouter interface {
	// GET Http GET method
	// Parameters:
	//   - path (string): The URL path to route the GET request.
	//   - handler (IHandler): The handler to process the GET request.
	GET(path string, handler IHandler)

	// HEAD Http HEAD method
	// Parameters:
	//   - path (string): The URL path to route the HEAD request.
	//   - handler (IHandler): The handler to process the HEAD request.
	HEAD(path string, handler IHandler)

	// POST Http POST method
	// Parameters:
	//   - path (string): The URL path to route the POST request.
	//   - handler (IHandler): The handler to process the POST request.
	POST(path string, handler IHandler)

	// PUT Http PUT method
	// Parameters:
	//   - path (string): The URL path to route the PUT request.
	//   - handler (IHandler): The handler to process the PUT request.
	PUT(path string, handler IHandler)

	// PATCH Http PATCH method
	// Parameters:
	//   - path (string): The URL path to route the PATCH request.
	//   - handler (IHandler): The handler to process the PATCH request.
	PATCH(path string, handler IHandler)

	// DELETE Http DELETE method
	// Parameters:
	//   - path (string): The URL path to route the DELETE request.
	//   - handler (IHandler): The handler to process the DELETE request.
	DELETE(path string, handler IHandler)

	// CONNECT Http CONNECT method
	// Parameters:
	//   - path (string): The URL path to route the CONNECT request.
	//   - handler (IHandler): The handler to process the CONNECT request.
	CONNECT(path string, handler IHandler)

	// OPTIONS Http OPTIONS method
	// Parameters:
	//   - path (string): The URL path to route the OPTIONS request.
	//   - handler (IHandler): The handler to process the OPTIONS request.
	OPTIONS(path string, handler IHandler)

	// TRACE Http TRACE method
	// Parameters:
	//   - path (string): The URL path to route the TRACE request.
	//   - handler (IHandler): The handler to process the TRACE request.
	TRACE(path string, handler IHandler)

	// Group multi routers
	// Parameters:
	//   - path (string): The URL path prefix for the group.
	//   - groupFunc (func(*Group)): A function defining routes and middleware for the group.
	Group(path string, groupFunc func(*Group))
}

// GET is a shortcut for Router.GET(path, handler).
// Parameters:
//   - path (string): The URL path to route the GET request.
//   - handler (IHandler): The handler to process the GET request.
func (fly *GFly) GET(path string, handler IHandler) {
	fly.router.GET(path, fly.wrapMiddlewares(handler))
}

// HEAD is a shortcut for Router.HEAD(path, handler).
// Parameters:
//   - path (string): The URL path to route the HEAD request.
//   - handler (IHandler): The handler to process the HEAD request.
func (fly *GFly) HEAD(path string, handler IHandler) {
	fly.router.HEAD(path, fly.wrapMiddlewares(handler))
}

// POST is a shortcut for Router.POST(path, handler).
// Parameters:
//   - path (string): The URL path to route the POST request.
//   - handler (IHandler): The handler to process the POST request.
func (fly *GFly) POST(path string, handler IHandler) {
	fly.router.POST(path, fly.wrapMiddlewares(handler))
}

// PUT is a shortcut for Router.PUT(path, handler).
// Parameters:
//   - path (string): The URL path to route the PUT request.
//   - handler (IHandler): The handler to process the PUT request.
func (fly *GFly) PUT(path string, handler IHandler) {
	fly.router.PUT(path, fly.wrapMiddlewares(handler))
}

// PATCH is a shortcut for Router.PATCH(path, handler).
// Parameters:
//   - path (string): The URL path to route the PATCH request.
//   - handler (IHandler): The handler to process the PATCH request.
func (fly *GFly) PATCH(path string, handler IHandler) {
	fly.router.PATCH(path, fly.wrapMiddlewares(handler))
}

// DELETE is a shortcut for Router.DELETE(path, handler).
// Parameters:
//   - path (string): The URL path to route the DELETE request.
//   - handler (IHandler): The handler to process the DELETE request.
func (fly *GFly) DELETE(path string, handler IHandler) {
	fly.router.DELETE(path, fly.wrapMiddlewares(handler))
}

// CONNECT is a shortcut for Router.CONNECT(path, handler).
// Parameters:
//   - path (string): The URL path to route the CONNECT request.
//   - handler (IHandler): The handler to process the CONNECT request.
func (fly *GFly) CONNECT(path string, handler IHandler) {
	fly.router.CONNECT(path, fly.wrapMiddlewares(handler))
}

// OPTIONS is a shortcut for Router.OPTIONS(path, handler).
// Parameters:
//   - path (string): The URL path to route the OPTIONS request.
//   - handler (IHandler): The handler to process the OPTIONS request.
func (fly *GFly) OPTIONS(path string, handler IHandler) {
	fly.router.OPTIONS(path, fly.wrapMiddlewares(handler))
}

// TRACE is a shortcut for Router.TRACE(path, handler).
// Parameters:
//   - path (string): The URL path to route the TRACE request.
//   - handler (IHandler): The handler to process the TRACE request.
func (fly *GFly) TRACE(path string, handler IHandler) {
	fly.router.TRACE(path, fly.wrapMiddlewares(handler))
}

// Group creates a group Handler functions.
// Parameters:
//   - path (string): The URL path prefix for the group.
//   - groupFunc (func(*Group)): A function defining routes and middleware for the group.
func (fly *GFly) Group(path string, groupFunc func(*Group)) {
	group := fly.router.Group(path)

	// Auto-append middleware from gFly to the group.
	group.middlewares = fly.middlewares

	groupFunc(group)
}

// wrapMiddlewares applies global middlewares to the handler.
// Parameters:
//   - handler (IHandler): The handler to wrap with middleware.
//
// Returns:
//   - IHandler: The handler wrapped with all global middlewares applied.
func (fly *GFly) wrapMiddlewares(handler IHandler) IHandler {
	if len(fly.middlewares) > 0 {
		middlewareGroup := NewMiddleware()

		// Convert simple middleware handlers to middleware configs
		configs := make([]MiddlewareConfig, len(fly.middlewares))
		for i, handler := range fly.middlewares {
			configs[i] = MiddlewareConfig{
				Handler:  handler,
				Phase:    PhasePreRequest,
				Priority: i,
			}
		}

		return middlewareGroup.GroupWithOptions(configs...)(handler)
	}

	return handler
}

// ====================================================================
//                        gFly - Plugin methods
// ====================================================================

// RegisterPlugin registers a plugin with the application.
// It returns an error if the plugin cannot be registered.
func (fly *GFly) RegisterPlugin(plugin plugin.Plugin) error {
	return fly.pluginManager.Register(plugin)
}

// GetPlugin returns a plugin by name.
// It returns nil if the plugin is not found.
func (fly *GFly) GetPlugin(name string) plugin.Plugin {
	return fly.pluginManager.GetPlugin(name)
}

// GetPlugins returns all registered plugins.
func (fly *GFly) GetPlugins() map[string]plugin.Plugin {
	return fly.pluginManager.GetPlugins()
}

// PluginManager returns the plugin manager.
func (fly *GFly) PluginManager() *plugin.Manager {
	return fly.pluginManager
}
