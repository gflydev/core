package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewManager(t *testing.T) {
	// Set environment variable for testing
	os.Setenv("APP_ENV", EnvDevelopment)
	defer os.Unsetenv("APP_ENV")

	manager := NewManager()
	assert.Equal(t, EnvDevelopment, manager.GetEnvironment())
	assert.Empty(t, manager.configPaths)
	assert.True(t, manager.loadEnvVars)
	assert.Empty(t, manager.values)
}

func TestSetEnvironment(t *testing.T) {
	manager := NewManager()
	manager.SetEnvironment(EnvProduction)
	assert.Equal(t, EnvProduction, manager.GetEnvironment())
}

func TestAddConfigPath(t *testing.T) {
	manager := NewManager()
	manager.AddConfigPath("config/app.json")
	assert.Equal(t, []string{"config/app.json"}, manager.configPaths)

	manager.AddConfigPath("config/database.json")
	assert.Equal(t, []string{"config/app.json", "config/database.json"}, manager.configPaths)
}

func TestSetLoadEnvVars(t *testing.T) {
	manager := NewManager()
	assert.True(t, manager.loadEnvVars)

	manager.SetLoadEnvVars(false)
	assert.False(t, manager.loadEnvVars)
}

func TestSet(t *testing.T) {
	manager := NewManager()

	// Test setting a simple value
	manager.Set("app.name", "Test App")
	assert.Equal(t, "Test App", manager.GetString("app.name", ""))

	// Test setting a nested value
	manager.Set("database.connection.host", "localhost")
	assert.Equal(t, "localhost", manager.GetString("database.connection.host", ""))

	// Test overriding a value
	manager.Set("app.name", "New App Name")
	assert.Equal(t, "New App Name", manager.GetString("app.name", ""))
}

func TestGet(t *testing.T) {
	manager := NewManager()
	manager.Set("app.name", "Test App")
	manager.Set("app.version", 1)
	manager.Set("app.debug", true)
	manager.Set("app.rate", 1.5)

	// Test getting a string value
	assert.Equal(t, "Test App", manager.GetString("app.name", "Default"))
	assert.Equal(t, "Default", manager.GetString("app.unknown", "Default"))

	// Test getting an int value
	assert.Equal(t, 1, manager.GetInt("app.version", 0))
	assert.Equal(t, 0, manager.GetInt("app.unknown", 0))

	// Test getting a bool value
	assert.True(t, manager.GetBool("app.debug", false))
	assert.False(t, manager.GetBool("app.unknown", false))

	// Test getting a float value
	assert.Equal(t, 1.5, manager.GetFloat("app.rate", 0.0))
	assert.Equal(t, 0.0, manager.GetFloat("app.unknown", 0.0))

	// Test getting a value with type conversion
	assert.Equal(t, 1, manager.GetInt("app.version", 0))
	assert.Equal(t, 1.0, manager.GetFloat("app.version", 0.0))
}

func TestHas(t *testing.T) {
	manager := NewManager()
	manager.Set("app.name", "Test App")

	assert.True(t, manager.Has("app.name"))
	assert.False(t, manager.Has("app.unknown"))
}

func TestAll(t *testing.T) {
	manager := NewManager()
	manager.Set("app.name", "Test App")
	manager.Set("app.version", 1)

	all := manager.All()
	assert.NotNil(t, all)
	assert.NotNil(t, all["app"])

	appMap, ok := all["app"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "Test App", appMap["name"])
	assert.Equal(t, 1, appMap["version"])
}

func TestLoadEnvironmentVariables(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("TEST_APP_NAME", "Test App")
	os.Setenv("TEST_APP_VERSION", "2")
	os.Setenv("TEST_APP_DEBUG", "true")
	defer func() {
		os.Unsetenv("TEST_APP_NAME")
		os.Unsetenv("TEST_APP_VERSION")
		os.Unsetenv("TEST_APP_DEBUG")
	}()

	manager := NewManager()
	manager.loadEnvironmentVariables()

	assert.Equal(t, "Test App", manager.GetString("test_app_name", ""))
	assert.Equal(t, "2", manager.GetString("test_app_version", ""))
	assert.Equal(t, "true", manager.GetString("test_app_debug", ""))
}

func TestGetEnvironmentConfigPath(t *testing.T) {
	manager := NewManager()
	manager.SetEnvironment(EnvDevelopment)

	path := manager.getEnvironmentConfigPath("config/app.json")
	assert.Equal(t, "config/app.dev.json", path)

	path = manager.getEnvironmentConfigPath("config/database.yml")
	assert.Equal(t, "config/database.dev.yml", path)
}

func TestDefaultManagerFunctions(t *testing.T) {
	// Test default manager functions
	SetEnvironment(EnvProduction)
	assert.Equal(t, EnvProduction, GetEnvironment())

	AddConfigPath("config/app.json")
	Set("app.name", "Test App")
	assert.Equal(t, "Test App", GetString("app.name", ""))
	assert.True(t, Has("app.name"))
	assert.NotEmpty(t, All())
}

// Helper function to create a temporary JSON config file
func createTempConfigFile(t *testing.T, content string) string {
	tmpfile, err := os.CreateTemp("", "config-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	return tmpfile.Name()
}

func TestLoadConfigFile(t *testing.T) {
	// Create a temporary config file
	configContent := `{
		"app": {
			"name": "Test App",
			"version": 1,
			"debug": true
		},
		"database": {
			"host": "localhost",
			"port": 3306
		}
	}`
	configFile := createTempConfigFile(t, configContent)
	defer os.Remove(configFile)

	manager := NewManager()
	err := manager.loadConfigFile(configFile)
	assert.NoError(t, err)

	assert.Equal(t, "Test App", manager.GetString("app.name", ""))
	assert.Equal(t, 1, manager.GetInt("app.version", 0))
	assert.True(t, manager.GetBool("app.debug", false))
	assert.Equal(t, "localhost", manager.GetString("database.host", ""))
	assert.Equal(t, 3306, manager.GetInt("database.port", 0))
}

func TestLoad(t *testing.T) {
	// Create a base config file
	baseConfigContent := `{
		"app": {
			"name": "Base App",
			"version": 1,
			"debug": false
		},
		"database": {
			"host": "localhost",
			"port": 3306
		}
	}`
	baseConfigFile := createTempConfigFile(t, baseConfigContent)
	defer os.Remove(baseConfigFile)

	// Create an environment-specific config file
	envConfigContent := `{
		"app": {
			"name": "Dev App",
			"debug": true
		},
		"database": {
			"port": 3307
		}
	}`
	envConfigFile := createTempConfigFile(t, envConfigContent)
	defer os.Remove(envConfigFile)

	// Rename the environment config file to match the expected pattern
	envConfigPath := baseConfigFile + ".dev.json"
	os.Rename(envConfigFile, envConfigPath)
	defer os.Remove(envConfigPath)

	// Set environment variables
	os.Setenv("APP_VERSION", "2")
	defer os.Unsetenv("APP_VERSION")

	manager := NewManager()
	manager.SetEnvironment(EnvDevelopment)
	manager.AddConfigPath(baseConfigFile)
	err := manager.Load()
	assert.NoError(t, err)

	// Check that values are correctly loaded and overridden
	assert.Equal(t, "Dev App", manager.GetString("app.name", ""))        // From env config
	assert.Equal(t, "2", manager.GetString("app.version", ""))           // From env var
	assert.True(t, manager.GetBool("app.debug", false))                  // From env config
	assert.Equal(t, "localhost", manager.GetString("database.host", "")) // From base config
	assert.Equal(t, 3307, manager.GetInt("database.port", 0))            // From env config
}
