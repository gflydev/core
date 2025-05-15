package main

import (
	"fmt"
	"time"

	"github.com/gflydev/core"
)

// =========================================================================================
//                                     Home API
// =========================================================================================

// HomeAPI handles the home endpoint
type HomeAPI struct {
	core.Api
}

// NewHomeAPI creates a new HomeAPI
func NewHomeAPI() *HomeAPI {
	return &HomeAPI{}
}

// Handle implements the IHandler interface
func (h *HomeAPI) Handle(ctx *core.Ctx) error {
	return ctx.Success(core.Data{
		"message": "Hello, World!",
		"time":    time.Now().Format(time.RFC3339),
	})
}

// =========================================================================================
//                                     Users API
// =========================================================================================

// UsersAPI handles the users endpoint
type UsersAPI struct {
	core.Api
}

// NewUsersAPI creates a new UsersAPI
func NewUsersAPI() *UsersAPI {
	return &UsersAPI{}
}

// Handle implements the IHandler interface
func (h *UsersAPI) Handle(ctx *core.Ctx) error {
	return ctx.Success(core.Data{
		"users": []core.Data{
			{"id": 1, "name": "Alice"},
			{"id": 2, "name": "Bob"},
			{"id": 3, "name": "Charlie"},
		},
	})
}

// =========================================================================================
//                                     User Profile API
// =========================================================================================

// UserProfileAPI handles the user profile endpoint
type UserProfileAPI struct {
	core.Api
}

// NewUserProfileAPI creates a new UserProfileAPI
func NewUserProfileAPI() *UserProfileAPI {
	return &UserProfileAPI{}
}

// Handle implements the IHandler interface
func (h *UserProfileAPI) Handle(ctx *core.Ctx) error {
	userID := ctx.QueryStr("user_id")
	if userID == "" {
		return ctx.Error(core.Data{
			"error": "User ID is required",
		})
	}

	return ctx.Success(core.Data{
		"user": core.Data{
			"id":    userID,
			"name":  "Example User",
			"email": "user@example.com",
		},
	})
}

// =========================================================================================
//                                     Router
// =========================================================================================

// setupRateLimitRoutes defines all routes for the rate limit example
func setupRateLimitRoutes(app core.IFly) {
	// Apply global rate limiting middleware (60 requests per minute)
	app.UseWithOptions(
		core.RateLimitMiddleware(
			core.WithRequestsPerMinute(60),
		),
		core.WithName("global-rate-limit"),
	)

	// Define a route that demonstrates rate limiting
	app.GET("/", NewHomeAPI())

	// Define a route group with a more restrictive rate limit for sensitive operations
	app.Group("/api", func(g *core.Group) {
		// Apply a more restrictive rate limit to the API group (10 requests per minute)
		g.UseWithOptions(
			core.RateLimitMiddleware(
				core.WithRequestsPerMinute(10),
				core.WithRateLimitMessage("API rate limit exceeded. Please slow down your requests."),
			),
			core.WithName("api-rate-limit"),
		)

		// Define API routes
		g.GET("/users", NewUsersAPI())
	})

	// Define a route group with user-specific rate limiting
	app.Group("/user", func(g *core.Group) {
		// Apply user-specific rate limiting (5 requests per minute per user)
		g.UseWithOptions(
			core.RateLimitMiddleware(
				core.WithRequestsPerMinute(5),
				// Use a custom key function that uses the user ID from the query parameter
				core.WithRateLimitKeyFunc(func(ctx *core.Ctx) string {
					userID := ctx.QueryStr("user_id")
					if userID == "" {
						// If no user ID is provided, use the request path as a fallback
						return "anonymous-" + ctx.Path()
					}
					return userID
				}),
				core.WithRateLimitMessage("User-specific rate limit exceeded."),
			),
			core.WithName("user-rate-limit"),
		)

		// Define user routes
		g.GET("/profile", NewUserProfileAPI())
	})
}

// =========================================================================================
//                                     Application
// =========================================================================================

func main() {
	// Create a new gFly application
	app := core.New()

	// Register router
	app.RegisterRouter(setupRateLimitRoutes)

	// Start the server
	fmt.Println("Server is running on http://localhost:8080")
	fmt.Println("Try the following endpoints:")
	fmt.Println("  - http://localhost:8080/ (60 requests per minute)")
	fmt.Println("  - http://localhost:8080/api/users (10 requests per minute)")
	fmt.Println("  - http://localhost:8080/user/profile?user_id=123 (5 requests per minute per user)")
	app.Run()
}
