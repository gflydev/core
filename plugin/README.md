# Plugin System

The plugin system allows for extending the gFly framework's functionality through plugins. Plugins can register services, middleware, routes, and other components to enhance the framework.

## Overview

The plugin system consists of the following components:

1. **Plugin Interface**: Defines the contract that all plugins must implement.
2. **Plugin Manager**: Manages the registration and initialization of plugins.
3. **Integration with GFly**: The GFly struct implements the IFlyPlugin interface to provide plugin support.

## Plugin Interface

All plugins must implement the `Plugin` interface:

```go
type Plugin interface {
    // Name returns the name of the plugin.
    Name() string

    // Description returns a description of the plugin.
    Description() string

    // Version returns the version of the plugin.
    Version() string

    // Initialize is called when the plugin is loaded.
    // It should register any services, middleware, or routes that the plugin provides.
    // The manager parameter can be used to access the container and other plugins.
    Initialize(manager *Manager) error
}
```

## Plugin Manager

The `Manager` struct manages the registration and initialization of plugins:

```go
type Manager struct {
    // plugins holds the registered plugins.
    plugins map[string]Plugin

    // container is the dependency injection container.
    container *container.Container

    // initialized indicates whether the plugins have been initialized.
    initialized bool
}
```

The manager provides methods for registering plugins, retrieving plugins, and initializing plugins.

## Using Plugins

### Creating a Plugin

To create a plugin, implement the `Plugin` interface:

```go
type MyPlugin struct{}

func (p *MyPlugin) Name() string {
    return "my-plugin"
}

func (p *MyPlugin) Description() string {
    return "My awesome plugin"
}

func (p *MyPlugin) Version() string {
    return "1.0.0"
}

func (p *MyPlugin) Initialize(manager *plugin.Manager) error {
    // Register services, middleware, routes, etc.
    app, err := plugin.ResolvePlugin[core.IFly](manager.Container())
    if err != nil {
        return err
    }

    // Add middleware
    app.Use(func(c *core.Ctx) error {
        // Middleware logic
        return nil
    })

    // Register routes
    app.GET("/my-plugin", &MyHandler{})

    return nil
}

func NewMyPlugin() plugin.Plugin {
    return &MyPlugin{}
}
```

### Registering a Plugin

To register a plugin with the gFly application:

```go
app := core.New()
err := app.RegisterPlugin(NewMyPlugin())
if err != nil {
    log.Fatalf("Failed to register plugin: %v", err)
}
```

### Plugin Initialization

Plugins are automatically initialized during application startup, after the global middlewares and routes are set up, but before the server starts listening for requests.

## Example Plugins

### Logger Plugin

The logger plugin logs HTTP requests:

```go
type LoggerPlugin struct{}

func (p *LoggerPlugin) Name() string {
    return "logger"
}

func (p *LoggerPlugin) Description() string {
    return "A plugin that logs requests"
}

func (p *LoggerPlugin) Version() string {
    return "1.0.0"
}

func (p *LoggerPlugin) Initialize(manager *plugin.Manager) error {
    app, err := plugin.ResolvePlugin[core.IFly](manager.Container())
    if err != nil {
        return err
    }

    app.Use(func(c *core.Ctx) error {
        method := string(c.Root().Request.Header.Method())
        path := c.Path()
        log.Infof("Request: %s %s", method, path)
        return nil
    })

    return nil
}

func NewLoggerPlugin() plugin.Plugin {
    return &LoggerPlugin{}
}
```

## Best Practices

1. **Plugin Naming**: Use unique, descriptive names for plugins to avoid conflicts.
2. **Error Handling**: Handle errors properly in the Initialize method to ensure the application can start correctly.
3. **Dependency Injection**: Use the container to resolve dependencies and register services.
4. **Documentation**: Document your plugin's functionality, configuration options, and any middleware or routes it adds.
5. **Versioning**: Use semantic versioning for your plugins to indicate compatibility.