package core

import (
	"fmt"
	"github.com/gflydev/core/log"
	"github.com/gflydev/core/utils"
	"github.com/valyala/fasthttp"
)

// ===========================================================================================================
// 												gFly
// ===========================================================================================================

var (
	// Application

	AppName  = utils.Getenv("APP_NAME", "gFly")
	AppCode  = utils.Getenv("APP_CODE", "gfly")
	AppURL   = utils.Getenv("APP_URL", "http://localhost:7789")
	AppEnv   = utils.Getenv("APP_ENV", "local")
	AppDebug = utils.Getenv("APP_DEBUG", true)

	// Storage directory

	StorageDir = utils.Getenv("STORAGE_DIR", "storage")  // Directory `{APP}/storage`
	TempDir    = utils.Getenv("TEMP_DIR", "storage/tmp") // Directory `{APP}/storage/temp`
	LogDir     = utils.Getenv("LOG_DIR", "storage/logs") // Directory `{APP}/storage/log`
	AppDir     = utils.Getenv("APP_DIR", "storage/app")  // Directory `{APP}/storage/app`

	// Internal variable

	fnHookMiddlewares []FnHookMiddleware // Hook global middlewares
	fnHookRoutes      []FnHookRoute      // Hook routers
)

type FnHookMiddleware func(fly IFlyMiddleware)
type FnHookRoute func(fly IFly)

// IFly Interface to declare all methods for gFly struct.
type IFly interface {
	IFlyRouter
	IFlyMiddleware
	// Run start application
	Run()
	// Router web router
	Router() *Router
	// RegisterMiddleware register middlewares
	RegisterMiddleware(fn ...FnHookMiddleware)
	// RegisterRouter register router
	RegisterRouter(fn ...FnHookRoute)
}

// GFly Struct define main elements in app.
type GFly struct {
	router      *Router             // Keep reference router
	server      *fasthttp.Server    // Keep reference web server
	config      Config              // App configuration
	middleware  IMiddleware         // Keep referenceMiddleware
	middlewares []MiddlewareHandler // Global middleware handler
}

// Router Get root router in gFly app.
func (fly *GFly) Router() *Router { // Get root Router
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

	// --------------- Setup Logs ---------------
	setupLog()

	// Print startup message
	if !fly.config.DisableStartupMessage {
		startupMessage(url, AppName, AppEnv)
	}

	// --------------- Checking service  ---------------
	// TODO: Need to add more checking

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
// The linking between fasthttp.RequestHandler to gFly's Ctx
func (fly *GFly) serveFastHTTP(ctx *fasthttp.RequestCtx) {
	handlerCtx := &Ctx{
		app:  fly,
		root: ctx,
		data: Data{},
	}

	_ = fly.router.Handler(handlerCtx)
}

// errorHandler Server error handler.
func (fly *GFly) errorHandler(ctx *fasthttp.RequestCtx, err error) {
	log.Debugf("Error %s", ctx.String())
	log.Errorf("Error happens %v", err)
}

// New Create new gFly app.
func New(config ...Config) IFly {
	app := &GFly{
		router:     NewRouter(),
		middleware: NewMiddleware(),
	}

	// Override config if provided
	if len(config) > 0 {
		app.config = config[0]
	} else {
		app.config = DefaultConfig
	}

	return app
}

// RegisterMiddleware Register Middleware
func (fly *GFly) RegisterMiddleware(fn ...FnHookMiddleware) {
	fnHookMiddlewares = fn
}

// RegisterRouter Register Router
func (fly *GFly) RegisterRouter(fn ...FnHookRoute) {
	fnHookRoutes = fn
}

// ===========================================================================================================
// 										gFly - Middleware methods
// ===========================================================================================================

// IFlyMiddleware Interface to declare all Middleware methods for gFly struct.
type IFlyMiddleware interface {
	// Use middleware to global (All requests)
	Use(middlewares ...MiddlewareHandler)
	// Middleware is a shortcut for Middleware.Group(middlewares ...MiddlewareHandler)
	Middleware(middleware ...MiddlewareHandler) func(IHandler) IHandler
}

// Use Middleware for global (All requests)
// Example
//
//	group.Use(middleware.RuleMiddlewareFunc, middleware.AuthMiddlewareFunc)
func (fly *GFly) Use(middlewares ...MiddlewareHandler) {
	fly.middlewares = append(fly.middlewares, middlewares...)
}

// Middleware is a shortcut for Middleware.Group(middlewares ...MiddlewareHandler)
// Example
//
//	group.POST("/one", gfly.IFly.Middleware(middleware.RuleMiddlewareFunc)(api.NewDefaultApi()))
func (fly *GFly) Middleware(middlewares ...MiddlewareHandler) func(IHandler) IHandler {
	return fly.middleware.Group(middlewares...)
}

// ===========================================================================================================
// 										gFly - HTTP methods
// ===========================================================================================================

// IFlyRouter Interface to declare all HTTP methods for gFly struct.
type IFlyRouter interface {
	// GET Http GET method
	GET(path string, handler IHandler)
	// HEAD Http HEAD method
	HEAD(path string, handler IHandler)
	// POST Http POST method
	POST(path string, handler IHandler)
	// PUT Http PUT method
	PUT(path string, handler IHandler)
	// PATCH Http PATCH method
	PATCH(path string, handler IHandler)
	// DELETE Http DELETE method
	DELETE(path string, handler IHandler)
	// CONNECT Http CONNECT method
	CONNECT(path string, handler IHandler)
	// OPTIONS Http OPTIONS method
	OPTIONS(path string, handler IHandler)
	// TRACE Http TRACE method
	TRACE(path string, handler IHandler)
	// Group multi routers
	Group(path string, groupFunc func(*Group))
}

// GET is a shortcut for Router.GET(path, handler)
func (fly *GFly) GET(path string, handler IHandler) {
	fly.router.GET(path, fly.wrapMiddlewares(handler))
}

// HEAD is a shortcut for Router.HEAD(path, handler)
func (fly *GFly) HEAD(path string, handler IHandler) {
	fly.router.HEAD(path, fly.wrapMiddlewares(handler))
}

// POST is a shortcut for Router.POST(path, handler)
func (fly *GFly) POST(path string, handler IHandler) {
	fly.router.POST(path, fly.wrapMiddlewares(handler))
}

// PUT is a shortcut for Router.PUT(path, handler)
func (fly *GFly) PUT(path string, handler IHandler) {
	fly.router.PUT(path, fly.wrapMiddlewares(handler))
}

// PATCH is a shortcut for Router.PATCH(path, handler)
func (fly *GFly) PATCH(path string, handler IHandler) {
	fly.router.PATCH(path, fly.wrapMiddlewares(handler))
}

// DELETE is a shortcut for Router.DELETE(path, handler)
func (fly *GFly) DELETE(path string, handler IHandler) {
	fly.router.DELETE(path, fly.wrapMiddlewares(handler))
}

// CONNECT is a shortcut for Router.CONNECT(path, handler)
func (fly *GFly) CONNECT(path string, handler IHandler) {
	fly.router.CONNECT(path, fly.wrapMiddlewares(handler))
}

// OPTIONS is a shortcut for Router.OPTIONS(path, handler)
func (fly *GFly) OPTIONS(path string, handler IHandler) {
	fly.router.OPTIONS(path, fly.wrapMiddlewares(handler))
}

// TRACE is a shortcut for Router.TRACE(path, handler)
func (fly *GFly) TRACE(path string, handler IHandler) {
	fly.router.TRACE(path, fly.wrapMiddlewares(handler))
}

// Group Create a group Handler functions.
func (fly *GFly) Group(path string, groupFunc func(*Group)) {
	group := fly.router.Group(path)

	// Auto append middleware from gfly to the group
	group.middlewares = fly.middlewares

	groupFunc(group)
}

func (fly *GFly) wrapMiddlewares(handler IHandler) IHandler {
	if len(fly.middlewares) > 0 {
		middlewareGroup := NewMiddleware()

		return middlewareGroup.Group(fly.middlewares...)(handler)
	}

	return handler
}
