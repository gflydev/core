// Package plugin provides a plugin system for the gFly framework.
// It allows for extending the framework's functionality through plugins.
package plugin

import (
	"github.com/gflydev/core/container"
	"github.com/gflydev/core/errors"
	"github.com/gflydev/core/log"
)

// Plugin is the interface that all plugins must implement.
// Plugins can extend the framework's functionality by registering
// services, middleware, routes, and other components.
type Plugin interface {
	// Name returns the name of the plugin.
	Name() string

	// Description returns a description of the plugin.
	Description() string

	// Version returns the version of the plugin.
	Version() string

	// Initialize is called when the plugin is loaded.
	// It should register any services, middleware, or routes that the plugin provides.
	// The container parameter can be used to register services.
	Initialize(manager *Manager) error
}

// Manager manages the registration and initialization of plugins.
// It provides methods for registering plugins and accessing the container.
type Manager struct {
	// plugins holds the registered plugins.
	plugins map[string]Plugin

	// container is the dependency injection container.
	container *container.Container

	// initialized indicates whether the plugins have been initialized.
	initialized bool
}

// NewManager creates a new plugin manager with the given container.
func NewManager(container *container.Container) *Manager {
	return &Manager{
		plugins:     make(map[string]Plugin),
		container:   container,
		initialized: false,
	}
}

// Register registers a plugin with the manager.
// It returns an error if a plugin with the same name is already registered.
func (m *Manager) Register(plugin Plugin) error {
	if m.initialized {
		return errors.Internalf("cannot register plugin after initialization")
	}

	name := plugin.Name()
	if _, exists := m.plugins[name]; exists {
		return errors.Conflictf("plugin with name '%s' already registered", name)
	}

	m.plugins[name] = plugin
	log.Infof("Registered plugin: %s (v%s)", name, plugin.Version())
	return nil
}

// Initialize initializes all registered plugins.
// It calls the Initialize method on each plugin.
func (m *Manager) Initialize() error {
	if m.initialized {
		return errors.Internalf("plugins already initialized")
	}

	for name, plugin := range m.plugins {
		log.Infof("Initializing plugin: %s", name)
		if err := plugin.Initialize(m); err != nil {
			return errors.Wrapf(err, errors.CodeInternal, "failed to initialize plugin '%s'", name)
		}
	}

	m.initialized = true
	return nil
}

// GetPlugin returns a plugin by name.
// It returns nil if the plugin is not found.
func (m *Manager) GetPlugin(name string) Plugin {
	return m.plugins[name]
}

// GetPlugins returns all registered plugins.
func (m *Manager) GetPlugins() map[string]Plugin {
	return m.plugins
}

// Container returns the dependency injection container.
func (m *Manager) Container() *container.Container {
	return m.container
}
