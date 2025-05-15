package plugin

import (
	"testing"

	"github.com/gflydev/core/container"
	"github.com/stretchr/testify/assert"
)

// MockPlugin is a mock implementation of the Plugin interface for testing.
type MockPlugin struct {
	name        string
	description string
	version     string
	initialized bool
}

// Name returns the name of the plugin.
func (p *MockPlugin) Name() string {
	return p.name
}

// Description returns a description of the plugin.
func (p *MockPlugin) Description() string {
	return p.description
}

// Version returns the version of the plugin.
func (p *MockPlugin) Version() string {
	return p.version
}

// Initialize is called when the plugin is loaded.
func (p *MockPlugin) Initialize(manager *Manager) error {
	p.initialized = true
	return nil
}

// NewMockPlugin creates a new mock plugin.
func NewMockPlugin(name, description, version string) Plugin {
	return &MockPlugin{
		name:        name,
		description: description,
		version:     version,
	}
}

// TestPluginManager tests the plugin manager.
func TestPluginManager(t *testing.T) {
	// Create a new container
	c := container.New()

	// Create a new plugin manager
	manager := NewManager(c)

	// Create a mock plugin
	plugin1 := &MockPlugin{
		name:        "plugin1",
		description: "Test plugin 1",
		version:     "1.0.0",
	}

	// Register the plugin
	err := manager.Register(plugin1)
	assert.NoError(t, err)

	// Try to register the same plugin again
	err = manager.Register(plugin1)
	assert.Error(t, err)

	// Create another mock plugin
	plugin2 := &MockPlugin{
		name:        "plugin2",
		description: "Test plugin 2",
		version:     "1.0.0",
	}

	// Register the second plugin
	err = manager.Register(plugin2)
	assert.NoError(t, err)

	// Get a plugin by name
	p := manager.GetPlugin("plugin1")
	assert.NotNil(t, p)
	assert.Equal(t, "plugin1", p.Name())

	// Get a non-existent plugin
	p = manager.GetPlugin("non-existent")
	assert.Nil(t, p)

	// Get all plugins
	plugins := manager.GetPlugins()
	assert.Len(t, plugins, 2)
	assert.Contains(t, plugins, "plugin1")
	assert.Contains(t, plugins, "plugin2")

	// Initialize the plugins
	err = manager.Initialize()
	assert.NoError(t, err)

	// Check that the plugins were initialized
	assert.True(t, plugin1.initialized)
	assert.True(t, plugin2.initialized)

	// Try to register a plugin after initialization
	plugin3 := &MockPlugin{
		name:        "plugin3",
		description: "Test plugin 3",
		version:     "1.0.0",
	}
	err = manager.Register(plugin3)
	assert.Error(t, err)

	// Try to initialize the plugins again
	err = manager.Initialize()
	assert.Error(t, err)
}

// TestRegisterPluginManager tests the RegisterPluginManager function.
func TestRegisterPluginManager(t *testing.T) {
	// Create a new container
	c := container.New()

	// Register the plugin manager
	err := RegisterPluginManager(c)
	assert.NoError(t, err)

	// Get the plugin manager
	manager, err := GetPluginManager(c)
	assert.NoError(t, err)
	assert.NotNil(t, manager)

	// Get the plugin manager using MustGetPluginManager
	manager = MustGetPluginManager(c)
	assert.NotNil(t, manager)
}
