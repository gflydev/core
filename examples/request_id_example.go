package main

import (
	"github.com/gflydev/core"
	"github.com/gflydev/core/log"
)

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

	// Register routes
	app.RegisterRouter(func(fly core.IFly) {
		// Define a route that uses request ID
		fly.GET("/example", &ExampleHandler{})
	})

	app.Run()
}

// ExampleHandler demonstrates how to use request ID in a handler
type ExampleHandler struct {
	core.Endpoint
}

// Handle implements the IHandler interface
func (h *ExampleHandler) Handle(ctx *core.Ctx) error {
	// Get the request ID
	requestID := core.GetRequestID(ctx)

	// Get the request-specific logger
	logger := core.GetLogger(ctx)

	// Log with request ID
	logger.Infof("Processing request %s", requestID)

	// Do something...

	// Log the result
	logger.Infow("Request processed",
		"status", "success",
		"request_id", requestID)

	// Include request ID in the response
	return ctx.Success(map[string]interface{}{
		"status":     "success",
		"message":    "Request processed successfully",
		"request_id": requestID,
	})
}

// ExternalServiceCall demonstrates how to propagate request ID to external services
func ExternalServiceCall(requestID string) {
	// Create a logger with request ID
	logger := log.DefaultLogger()
	requestLogger := log.NewRequestIDLogger(logger, requestID)

	// Log with request ID
	requestLogger.Info("Making external service call")

	// Include request ID in the external service call headers
	// headers := map[string]string{
	//     "X-Request-ID": requestID,
	// }
	// Make the external service call with the headers...

	requestLogger.Info("External service call completed")
}
