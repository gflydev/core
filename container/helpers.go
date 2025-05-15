// Package container provides a dependency injection container for the gFly framework.
package container

import (
	"reflect"
)

// RegisterType is a generic helper for registering services by type.
// It registers a concrete type as an implementation of an interface.
//
// Example:
//
//	container.RegisterType[ILogger, ConsoleLogger](c, Singleton)
func RegisterType[T any, TImpl any](c *Container, lifetime ServiceLifetime) error {
	serviceType := reflect.TypeOf((*T)(nil)).Elem()
	implType := reflect.TypeOf((*TImpl)(nil)).Elem()
	return c.Register(serviceType, implType, lifetime)
}

// RegisterInstance is a generic helper for registering instances by type.
// It registers an existing instance as a singleton service.
//
// Example:
//
//	logger := &ConsoleLogger{}
//	container.RegisterInstance[ILogger](c, logger)
func RegisterInstance[T any](c *Container, instance interface{}) error {
	serviceType := reflect.TypeOf((*T)(nil)).Elem()
	return c.RegisterInstance(serviceType, instance)
}

// RegisterFactory is a generic helper for registering factory functions by type.
// It registers a factory function for creating service instances.
//
// Example:
//
//	container.RegisterFactory[ILogger](c, func(c *Container) (interface{}, error) {
//	    return &ConsoleLogger{}, nil
//	}, Singleton)
func RegisterFactory[T any](c *Container, factory func(c *Container) (interface{}, error), lifetime ServiceLifetime) {
	serviceType := reflect.TypeOf((*T)(nil)).Elem()
	c.RegisterFactory(serviceType, factory, lifetime)
}

// RegisterSingleton is a convenience method for registering a singleton service.
//
// Example:
//
//	container.RegisterSingleton[ILogger, ConsoleLogger](c)
func RegisterSingleton[T any, TImpl any](c *Container) error {
	return RegisterType[T, TImpl](c, Singleton)
}

// RegisterTransient is a convenience method for registering a transient service.
//
// Example:
//
//	container.RegisterTransient[ILogger, ConsoleLogger](c)
func RegisterTransient[T any, TImpl any](c *Container) error {
	return RegisterType[T, TImpl](c, Transient)
}

// RegisterScoped is a convenience method for registering a scoped service.
//
// Example:
//
//	container.RegisterScoped[ILogger, ConsoleLogger](c)
func RegisterScoped[T any, TImpl any](c *Container) error {
	return RegisterType[T, TImpl](c, Scoped)
}

// RegisterSingletonFactory is a convenience method for registering a singleton factory.
//
// Example:
//
//	container.RegisterSingletonFactory[ILogger](c, func(c *Container) (interface{}, error) {
//	    return &ConsoleLogger{}, nil
//	})
func RegisterSingletonFactory[T any](c *Container, factory func(c *Container) (interface{}, error)) {
	RegisterFactory[T](c, factory, Singleton)
}

// RegisterTransientFactory is a convenience method for registering a transient factory.
//
// Example:
//
//	container.RegisterTransientFactory[ILogger](c, func(c *Container) (interface{}, error) {
//	    return &ConsoleLogger{}, nil
//	})
func RegisterTransientFactory[T any](c *Container, factory func(c *Container) (interface{}, error)) {
	RegisterFactory[T](c, factory, Transient)
}

// RegisterScopedFactory is a convenience method for registering a scoped factory.
//
// Example:
//
//	container.RegisterScopedFactory[ILogger](c, func(c *Container) (interface{}, error) {
//	    return &ConsoleLogger{}, nil
//	})
func RegisterScopedFactory[T any](c *Container, factory func(c *Container) (interface{}, error)) {
	RegisterFactory[T](c, factory, Scoped)
}

// Resolve is a generic helper for resolving services by type.
//
// Example:
//
//	logger, err := container.Resolve[ILogger](c)
func Resolve[T any](c *Container) (T, error) {
	return ResolveByType[T](c)
}

// MustResolve is a generic helper for resolving services by type that panics on error.
//
// Example:
//
//	logger := container.MustResolve[ILogger](c)
func MustResolve[T any](c *Container) T {
	return MustResolveByType[T](c)
}

// BuildServiceProvider creates a new container and configures it with the provided
// configuration function. This is a convenience method for setting up a container
// with multiple service registrations.
//
// Example:
//
//	container := container.BuildServiceProvider(func(c *container.Container) {
//	    container.RegisterSingleton[ILogger, ConsoleLogger](c)
//	    container.RegisterTransient[IUserService, UserService](c)
//	})
func BuildServiceProvider(configure func(c *Container)) *Container {
	container := New()
	configure(container)
	return container
}
