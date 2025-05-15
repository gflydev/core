// Example demonstrating the use of the dependency injection container in a gFly application.
package main

import (
	"fmt"
	"github.com/gflydev/core"
	"github.com/gflydev/core/container"
	"github.com/gflydev/core/log"
)

// Define interfaces and implementations for the example

// UserService defines methods for user management
type UserService interface {
	GetUserName(id int) string
}

// DefaultUserService is a simple implementation of UserService
type DefaultUserService struct {
	logger log.AllLogger
}

// NewDefaultUserService creates a new DefaultUserService with the given logger
func NewDefaultUserService(logger log.AllLogger) *DefaultUserService {
	return &DefaultUserService{
		logger: logger,
	}
}

// GetUserName returns a user name for the given ID
func (s *DefaultUserService) GetUserName(id int) string {
	s.logger.Infof("Getting user name for ID: %d", id)
	return fmt.Sprintf("User %d", id)
}

// UserAPI is a handler for user-related API endpoints
type UserAPI struct {
	core.Api
	userService UserService
}

// NewUserAPI creates a new UserAPI with the given UserService
func NewUserAPI(userService UserService) *UserAPI {
	return &UserAPI{
		userService: userService,
	}
}

// Handle handles the API request
func (h *UserAPI) Handle(c *core.Ctx) error {
	// Get user ID from query parameter, default to 1 if not provided
	userID, err := c.QueryInt("id")
	if err != nil {
		userID = 1 // Default to 1 if there's an error
	}

	// Get user name from service
	userName := h.userService.GetUserName(userID)

	// Return JSON response
	return c.JSON(core.Data{
		"id":   userID,
		"name": userName,
	})
}

// setupDependencyInjection configures the dependency injection container
func setupDependencyInjection(app core.IFly) {
	// Get the container from the application
	c := app.Container()

	// Register the UserService as a singleton
	container.RegisterSingletonFactory[UserService](c, func(c *container.Container) (interface{}, error) {
		// Resolve the logger from the container
		logger, err := container.Resolve[log.AllLogger](c)
		if err != nil {
			return nil, err
		}

		// Create and return the UserService
		return NewDefaultUserService(logger), nil
	})

	// Log that services have been registered
	logger, _ := container.Resolve[log.AllLogger](c)
	logger.Info("Dependency injection container configured")
}

// setupRoutes configures the application routes
func setupRoutes(app core.IFly) {
	// Get the container from the application
	c := app.Container()

	// Define API routes
	app.Group("/api", func(apiRouter *core.Group) {
		// Resolve the UserService from the container
		userService, err := container.Resolve[UserService](c)
		if err != nil {
			panic(err)
		}

		// Register the UserAPI route
		apiRouter.GET("/user", NewUserAPI(userService))
	})
}

// Example usage in main function:
/*
func main() {
	// Create a new gFly application
	app := core.New()

	// Set up dependency injection
	setupDependencyInjection(app)

	// Register routes
	app.RegisterRouter(setupRoutes)

	// Run the application
	app.Run()
}
*/

// Note: This example is not meant to be run directly.
// It demonstrates how to use the dependency injection container in a gFly application.
// To run this example, incorporate the setupDependencyInjection and setupRoutes functions
// into your main.go file.
