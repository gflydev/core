// Package plugin provides a plugin system for the gFly framework.
package plugin

import (
	"github.com/gflydev/core/container"
	"github.com/gflydev/core/errors"
	"github.com/gflydev/core/log"
)

// IFlyPlugin is an interface that extends the core IFly interface with plugin functionality.
// This interface will be implemented by the GFly struct to provide plugin support.
type IFlyPlugin interface {
	// RegisterPlugin registers a plugin with the application.
	// It returns an error if the plugin cannot be registered.
	RegisterPlugin(plugin Plugin) error

	// GetPlugin returns a plugin by name.
	// It returns nil if the plugin is not found.
	GetPlugin(name string) Plugin

	// GetPlugins returns all registered plugins.
	GetPlugins() map[string]Plugin

	// PluginManager returns the plugin manager.
	PluginManager() *Manager
}

// RegisterPluginInstance registers a plugin instance with the container.
// This is a helper function for registering plugins as services in the container.
func RegisterPluginInstance[T any](c *container.Container, plugin Plugin) error {
	return container.RegisterInstance[T](c, plugin)
}

// ResolvePlugin resolves a plugin from the container by type.
// This is a helper function for resolving plugins from the container.
func ResolvePlugin[T any](c *container.Container) (T, error) {
	return container.Resolve[T](c)
}

// MustResolvePlugin resolves a plugin from the container by type and panics on error.
// This is a helper function for resolving plugins from the container.
func MustResolvePlugin[T any](c *container.Container) T {
	return container.MustResolve[T](c)
}

// RegisterPluginManager registers the plugin manager with the container.
// This is called by the GFly struct during initialization.
func RegisterPluginManager(c *container.Container) error {
	manager := NewManager(c)
	return container.RegisterInstance[*Manager](c, manager)
}

// GetPluginManager retrieves the plugin manager from the container.
// This is a helper function for getting the plugin manager from the container.
func GetPluginManager(c *container.Container) (*Manager, error) {
	return container.Resolve[*Manager](c)
}

// MustGetPluginManager retrieves the plugin manager from the container and panics on error.
// This is a helper function for getting the plugin manager from the container.
func MustGetPluginManager(c *container.Container) *Manager {
	return container.MustResolve[*Manager](c)
}

// InitializePlugins initializes all registered plugins.
// This is called by the GFly struct during application startup.
func InitializePlugins(c *container.Container) error {
	manager, err := GetPluginManager(c)
	if err != nil {
		return errors.Wrapf(err, errors.CodeInternal, "failed to get plugin manager")
	}

	log.Info("Initializing plugins...")
	if err := manager.Initialize(); err != nil {
		return errors.Wrapf(err, errors.CodeInternal, "failed to initialize plugins")
	}
	log.Info("Plugins initialized successfully")

	return nil
}
