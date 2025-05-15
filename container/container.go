// Package container provides a dependency injection container for the gFly framework.
// It allows for registering, resolving, and managing the lifecycle of services.
package container

import (
	"reflect"
	"sync"

	"github.com/gflydev/core/errors"
)

// ServiceLifetime defines how the container manages a service instance.
type ServiceLifetime int

const (
	// Singleton services are created once and shared across all resolutions.
	Singleton ServiceLifetime = iota
	// Transient services are created each time they are resolved.
	Transient
	// Scoped services are created once per scope.
	Scoped
)

// ServiceDescriptor contains the information needed to create and manage a service.
type ServiceDescriptor struct {
	// ServiceType is the type of the service being registered.
	ServiceType reflect.Type
	// ImplementationType is the concrete type that implements the service.
	ImplementationType reflect.Type
	// Factory is a function that creates an instance of the service.
	Factory func(c *Container) (interface{}, error)
	// Lifetime defines how the service instance is managed.
	Lifetime ServiceLifetime
	// Instance holds the singleton instance if Lifetime is Singleton.
	Instance interface{}
}

// Container is a dependency injection container that manages service registration and resolution.
type Container struct {
	// services maps service types to their descriptors.
	services map[reflect.Type]*ServiceDescriptor
	// mutex protects concurrent access to the container.
	mutex sync.RWMutex
	// parent is the parent container for scoped containers.
	parent *Container
	// scopes tracks all child scopes created from this container.
	scopes []*Container
}

// New creates a new dependency injection container.
func New() *Container {
	return &Container{
		services: make(map[reflect.Type]*ServiceDescriptor),
		scopes:   make([]*Container, 0),
	}
}

// Register registers a service with the container.
// It takes the service interface type, the implementation type, and the service lifetime.
func (c *Container) Register(serviceType, implementationType reflect.Type, lifetime ServiceLifetime) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !implementationType.Implements(serviceType) && serviceType.Kind() != reflect.Interface {
		return errors.InvalidInputf("type %v does not implement %v", implementationType, serviceType)
	}

	c.services[serviceType] = &ServiceDescriptor{
		ServiceType:        serviceType,
		ImplementationType: implementationType,
		Lifetime:           lifetime,
		Factory: func(c *Container) (interface{}, error) {
			return reflect.New(implementationType).Interface(), nil
		},
	}

	return nil
}

// RegisterInstance registers an existing instance as a singleton service.
func (c *Container) RegisterInstance(serviceType reflect.Type, instance interface{}) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	instanceType := reflect.TypeOf(instance)
	if !instanceType.Implements(serviceType) && serviceType.Kind() != reflect.Interface {
		return errors.InvalidInputf("instance of type %v does not implement %v", instanceType, serviceType)
	}

	c.services[serviceType] = &ServiceDescriptor{
		ServiceType:        serviceType,
		ImplementationType: instanceType,
		Lifetime:           Singleton,
		Instance:           instance,
	}

	return nil
}

// RegisterFactory registers a factory function for creating service instances.
func (c *Container) RegisterFactory(serviceType reflect.Type, factory func(c *Container) (interface{}, error), lifetime ServiceLifetime) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.services[serviceType] = &ServiceDescriptor{
		ServiceType: serviceType,
		Factory:     factory,
		Lifetime:    lifetime,
	}
}

// Resolve resolves a service from the container.
func (c *Container) Resolve(serviceType reflect.Type) (interface{}, error) {
	c.mutex.RLock()
	descriptor, exists := c.services[serviceType]
	c.mutex.RUnlock()

	if !exists {
		if c.parent != nil {
			// Try to resolve from parent container if this is a scoped container
			return c.parent.Resolve(serviceType)
		}
		return nil, errors.NotFoundf("service of type %v not registered", serviceType)
	}

	switch descriptor.Lifetime {
	case Singleton:
		return c.resolveSingleton(descriptor)
	case Transient:
		return c.resolveTransient(descriptor)
	case Scoped:
		return c.resolveScoped(descriptor)
	default:
		return nil, errors.Internalf("unknown service lifetime: %v", descriptor.Lifetime)
	}
}

// resolveSingleton resolves a singleton service.
func (c *Container) resolveSingleton(descriptor *ServiceDescriptor) (interface{}, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if descriptor.Instance != nil {
		return descriptor.Instance, nil
	}

	if descriptor.Factory != nil {
		instance, err := descriptor.Factory(c)
		if err != nil {
			return nil, errors.Wrap(err, errors.CodeInternal, "failed to create service instance using factory")
		}
		descriptor.Instance = instance
		return instance, nil
	}

	instance := reflect.New(descriptor.ImplementationType).Interface()
	descriptor.Instance = instance
	return instance, nil
}

// resolveTransient resolves a transient service.
func (c *Container) resolveTransient(descriptor *ServiceDescriptor) (interface{}, error) {
	if descriptor.Factory != nil {
		return descriptor.Factory(c)
	}
	return reflect.New(descriptor.ImplementationType).Interface(), nil
}

// resolveScoped resolves a scoped service.
func (c *Container) resolveScoped(descriptor *ServiceDescriptor) (interface{}, error) {
	// If this is not a scoped container, create a new scope
	if c.parent == nil {
		scope := c.CreateScope()
		return scope.Resolve(descriptor.ServiceType)
	}

	// This is a scoped container, resolve as singleton within this scope
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if descriptor.Instance != nil {
		return descriptor.Instance, nil
	}

	if descriptor.Factory != nil {
		instance, err := descriptor.Factory(c)
		if err != nil {
			return nil, errors.Wrap(err, errors.CodeInternal, "failed to create scoped service instance using factory")
		}
		descriptor.Instance = instance
		return instance, nil
	}

	instance := reflect.New(descriptor.ImplementationType).Interface()
	descriptor.Instance = instance
	return instance, nil
}

// CreateScope creates a new scoped container.
func (c *Container) CreateScope() *Container {
	scope := &Container{
		services: make(map[reflect.Type]*ServiceDescriptor),
		parent:   c,
		scopes:   make([]*Container, 0),
	}

	c.mutex.Lock()
	c.scopes = append(c.scopes, scope)
	c.mutex.Unlock()

	return scope
}

// MustResolve resolves a service from the container and panics if it cannot be resolved.
func (c *Container) MustResolve(serviceType reflect.Type) interface{} {
	instance, err := c.Resolve(serviceType)
	if err != nil {
		panic(err)
	}
	return instance
}

// ResolveByType is a generic helper for resolving services by type.
func ResolveByType[T any](c *Container) (T, error) {
	var zero T
	serviceType := reflect.TypeOf((*T)(nil)).Elem()
	instance, err := c.Resolve(serviceType)
	if err != nil {
		return zero, err
	}
	return instance.(T), nil
}

// MustResolveByType is a generic helper for resolving services by type that panics on error.
func MustResolveByType[T any](c *Container) T {
	instance, err := ResolveByType[T](c)
	if err != nil {
		panic(err)
	}
	return instance
}

// Dispose cleans up all resources held by the container and its scopes.
func (c *Container) Dispose() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Dispose all child scopes first
	for _, scope := range c.scopes {
		scope.Dispose()
	}
	c.scopes = nil

	// Clear all services
	c.services = make(map[reflect.Type]*ServiceDescriptor)
}
