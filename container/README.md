# Dependency Injection Container

This package provides a dependency injection container for the gFly framework. It allows for registering, resolving, and managing the lifecycle of services.

## Features

- **Service Registration**: Register services with different lifetimes (singleton, transient, scoped)
- **Service Resolution**: Resolve services by type
- **Factory Functions**: Create services using factory functions
- **Lifecycle Management**: Manage the lifecycle of services
- **Scoped Containers**: Create scoped containers for request-scoped services
- **Type-Safe API**: Use generics for type-safe service registration and resolution

## Usage

### Creating a Container

```go
// Create a new container
container := container.New()

// Or use BuildServiceProvider for configuring multiple services
container := container.BuildServiceProvider(func(c *container.Container) {
    // Register services here
})
```

### Registering Services

```go
// Register a singleton service
container.RegisterSingleton[ILogger, ConsoleLogger](container)

// Register a transient service
container.RegisterTransient[IUserService, UserService](container)

// Register a scoped service
container.RegisterScoped[IRequestContext, RequestContext](container)

// Register an existing instance
logger := &ConsoleLogger{}
container.RegisterInstance[ILogger](container, logger)

// Register a factory function
container.RegisterSingletonFactory[IDatabase](container, func(c *container.Container) (interface{}, error) {
    return NewDatabase("connection-string"), nil
})
```

### Resolving Services

```go
// Resolve a service (returns error if not found)
logger, err := container.Resolve[ILogger](container)
if err != nil {
    // Handle error
}

// Resolve a service (panics if not found)
logger := container.MustResolve[ILogger](container)
```

### Working with Scopes

```go
// Create a scoped container
scope := container.CreateScope()

// Resolve scoped services
context, err := container.Resolve[IRequestContext](scope)

// Dispose the scope when done
scope.Dispose()
```

### Service Lifetimes

- **Singleton**: Services are created once and shared across all resolutions
- **Transient**: Services are created each time they are resolved
- **Scoped**: Services are created once per scope

## Example

```go
package main

import (
    "fmt"
    "github.com/gflydev/core/container"
)

// Define interfaces and implementations
type IGreeter interface {
    Greet(name string) string
}

type Greeter struct{}

func (g *Greeter) Greet(name string) string {
    return fmt.Sprintf("Hello, %s!", name)
}

type ITimeProvider interface {
    GetCurrentTime() string
}

type TimeProvider struct{}

func (t *TimeProvider) GetCurrentTime() string {
    return time.Now().Format(time.RFC3339)
}

type GreetingService struct {
    greeter      IGreeter
    timeProvider ITimeProvider
}

func NewGreetingService(greeter IGreeter, timeProvider ITimeProvider) *GreetingService {
    return &GreetingService{
        greeter:      greeter,
        timeProvider: timeProvider,
    }
}

func (s *GreetingService) CreateGreeting(name string) string {
    return fmt.Sprintf("%s - %s", s.greeter.Greet(name), s.timeProvider.GetCurrentTime())
}

func main() {
    // Create and configure the container
    c := container.BuildServiceProvider(func(c *container.Container) {
        // Register services
        container.RegisterSingleton[IGreeter, Greeter](c)
        container.RegisterSingleton[ITimeProvider, TimeProvider](c)
        
        // Register a factory that resolves dependencies
        container.RegisterSingletonFactory[GreetingService](c, func(c *container.Container) (interface{}, error) {
            greeter, err := container.Resolve[IGreeter](c)
            if err != nil {
                return nil, err
            }
            
            timeProvider, err := container.Resolve[ITimeProvider](c)
            if err != nil {
                return nil, err
            }
            
            return NewGreetingService(greeter, timeProvider), nil
        })
    })
    
    // Resolve the service
    service, err := container.Resolve[GreetingService](c)
    if err != nil {
        panic(err)
    }
    
    // Use the service
    greeting := service.CreateGreeting("World")
    fmt.Println(greeting)
}
```

## Best Practices

1. **Register services at startup**: Register all services during application startup to catch any registration errors early.
2. **Use interfaces**: Register services by their interfaces rather than concrete types to promote loose coupling.
3. **Prefer singleton lifetime**: Use singleton lifetime for stateless services to improve performance.
4. **Use scoped lifetime for request-scoped services**: Use scoped lifetime for services that should be shared within a request but not across requests.
5. **Dispose scopes**: Always dispose scopes when they are no longer needed to prevent memory leaks.