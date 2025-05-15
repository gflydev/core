# Configuration Management System

The `config` package provides a flexible configuration management system with support for different environments (local, development, staging, production). It allows loading configuration from various sources, including JSON files and environment variables.

## Features

- Support for different environments (local, development, staging, production)
- Load configuration from JSON files
- Load environment-specific configuration files
- Load configuration from environment variables
- Access configuration values with type conversion
- Nested configuration values using dot notation
- Default values for missing configuration keys

## Usage

### Basic Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/gflydev/core/config"
)

func main() {
	// Add configuration file paths
	config.AddConfigPath("config/app.json")
	config.AddConfigPath("config/database.json")

	// Load configuration
	if err := config.Load(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Get configuration values
	appName := config.GetString("app.name", "Default App Name")
	appDebug := config.GetBool("app.debug", false)
	dbHost := config.GetString("database.host", "localhost")
	dbPort := config.GetInt("database.port", 3306)

	fmt.Printf("App Name: %s\n", appName)
	fmt.Printf("App Debug: %t\n", appDebug)
	fmt.Printf("DB Host: %s\n", dbHost)
	fmt.Printf("DB Port: %d\n", dbPort)
}
```

### Environment-Specific Configuration

The configuration system supports environment-specific configuration files. For example, if you have a configuration file `config/app.json` and the current environment is `dev`, the system will also load `config/app.dev.json` if it exists.

```go
// Set the environment (default is determined by APP_ENV environment variable)
config.SetEnvironment(config.EnvDevelopment)

// Add configuration file paths
config.AddConfigPath("config/app.json")

// Load configuration
if err := config.Load(); err != nil {
	log.Fatalf("Failed to load configuration: %v", err)
}
```

### Using a Custom Manager

You can create a custom configuration manager if you need multiple independent configurations:

```go
// Create a new configuration manager
manager := config.NewManager()

// Configure the manager
manager.SetEnvironment(config.EnvProduction)
manager.AddConfigPath("config/app.json")
manager.SetLoadEnvVars(false)

// Load configuration
if err := manager.Load(); err != nil {
	log.Fatalf("Failed to load configuration: %v", err)
}

// Get configuration values
appName := manager.GetString("app.name", "Default App Name")
```

## Configuration File Format

Configuration files should be in JSON format:

```json
{
  "app": {
    "name": "My App",
    "debug": true,
    "url": "http://localhost:8080"
  },
  "database": {
    "host": "localhost",
    "port": 3306,
    "username": "root",
    "password": "secret",
    "database": "myapp"
  }
}
```

## Environment-Specific Configuration Files

Environment-specific configuration files should have the environment name as a suffix before the file extension:

- `config/app.json` - Base configuration
- `config/app.local.json` - Local environment configuration
- `config/app.dev.json` - Development environment configuration
- `config/app.stag.json` - Staging environment configuration
- `config/app.prod.json` - Production environment configuration

Environment-specific configuration files override values from the base configuration file.

## Environment Variables

Environment variables are automatically loaded into the configuration. The keys are converted to lowercase and can be accessed using dot notation:

```
APP_NAME=My App
APP_DEBUG=true
DATABASE_HOST=localhost
DATABASE_PORT=3306
```

These environment variables can be accessed as:

```go
appName := config.GetString("app.name", "Default App Name")
appDebug := config.GetBool("app.debug", false)
dbHost := config.GetString("database.host", "localhost")
dbPort := config.GetInt("database.port", 3306)
```

## API Reference

### Constants

- `EnvLocal` - Local environment
- `EnvDevelopment` - Development environment
- `EnvStaging` - Staging environment
- `EnvProduction` - Production environment

### Functions

- `SetEnvironment(env string)` - Sets the current environment
- `GetEnvironment() string` - Returns the current environment
- `AddConfigPath(path string)` - Adds a configuration file path
- `SetLoadEnvVars(load bool)` - Sets whether to load environment variables
- `Load() error` - Loads configuration from all sources
- `Get(key string, defaultValue interface{}) interface{}` - Retrieves a configuration value by key
- `GetString(key string, defaultValue string) string` - Retrieves a string configuration value by key
- `GetInt(key string, defaultValue int) int` - Retrieves an integer configuration value by key
- `GetBool(key string, defaultValue bool) bool` - Retrieves a boolean configuration value by key
- `GetFloat(key string, defaultValue float64) float64` - Retrieves a float configuration value by key
- `Set(key string, value interface{})` - Sets a configuration value by key
- `Has(key string) bool` - Checks if a configuration key exists
- `All() map[string]interface{}` - Returns all configuration values

### Manager Methods

- `NewManager() *Manager` - Creates a new configuration manager
- `SetEnvironment(env string) *Manager` - Sets the current environment
- `GetEnvironment() string` - Returns the current environment
- `AddConfigPath(path string) *Manager` - Adds a configuration file path
- `SetLoadEnvVars(load bool) *Manager` - Sets whether to load environment variables
- `Load() error` - Loads configuration from all sources
- `Get(key string, defaultValue interface{}) interface{}` - Retrieves a configuration value by key
- `GetString(key string, defaultValue string) string` - Retrieves a string configuration value by key
- `GetInt(key string, defaultValue int) int` - Retrieves an integer configuration value by key
- `GetBool(key string, defaultValue bool) bool` - Retrieves a boolean configuration value by key
- `GetFloat(key string, defaultValue float64) float64` - Retrieves a float configuration value by key
- `Set(key string, value interface{}) *Manager` - Sets a configuration value by key
- `Has(key string) bool` - Checks if a configuration key exists
- `All() map[string]interface{}` - Returns all configuration values