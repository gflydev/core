// Package plugins provides example plugins for the gFly framework.
package plugins

import (
	"github.com/gflydev/core"
	"github.com/gflydev/core/log"
	"github.com/gflydev/core/plugin"
)

// LoggerPlugin is an example plugin that logs requests and responses.
type LoggerPlugin struct{}

// Name returns the name of the plugin.
func (p *LoggerPlugin) Name() string {
	return "logger"
}

// Description returns a description of the plugin.
func (p *LoggerPlugin) Description() string {
	return "A plugin that logs requests and responses"
}

// Version returns the version of the plugin.
func (p *LoggerPlugin) Version() string {
	return "1.0.0"
}

// Initialize is called when the plugin is loaded.
// It registers a middleware that logs requests and responses.
func (p *LoggerPlugin) Initialize(manager *plugin.Manager) error {
	log.Info("Initializing logger plugin")

	// Register a middleware that logs requests and responses
	app, err := plugin.ResolvePlugin[core.IFly](manager.Container())
	if err != nil {
		return err
	}

	// Add a middleware that logs requests
	app.Use(func(c *core.Ctx) error {
		// Get the HTTP method and path
		method := string(c.Root().Request.Header.Method())
		path := c.Path()

		log.Infof("Request: %s %s", method, path)

		// Call the next middleware/handler by returning nil
		// This will allow the request to continue to the next middleware or handler

		// Note: In this middleware pattern, we can't log the response status code
		// after the handler has executed because middleware functions are executed
		// sequentially before the handler. We would need to use a different approach
		// to log the response status code.

		return nil
	})

	log.Info("Logger plugin initialized")
	return nil
}

// NewLoggerPlugin creates a new logger plugin.
func NewLoggerPlugin() plugin.Plugin {
	return &LoggerPlugin{}
}
