package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gflydev/core/utils"
)

// Environment constants
const (
	EnvLocal       = "local"
	EnvDevelopment = "dev"
	EnvStaging     = "stag"
	EnvProduction  = "prod"
)

// Manager is responsible for loading and providing access to configuration values
type Manager struct {
	// Current environment
	environment string

	// Configuration values
	values map[string]interface{}

	// Configuration file paths
	configPaths []string

	// Whether to load environment variables
	loadEnvVars bool
}

// NewManager creates a new configuration manager
func NewManager() *Manager {
	return &Manager{
		environment: utils.Getenv("APP_ENV", EnvLocal),
		values:      make(map[string]interface{}),
		configPaths: []string{},
		loadEnvVars: true,
	}
}

// SetEnvironment sets the current environment
func (m *Manager) SetEnvironment(env string) *Manager {
	m.environment = env
	return m
}

// GetEnvironment returns the current environment
func (m *Manager) GetEnvironment() string {
	return m.environment
}

// AddConfigPath adds a configuration file path
func (m *Manager) AddConfigPath(path string) *Manager {
	m.configPaths = append(m.configPaths, path)
	return m
}

// SetLoadEnvVars sets whether to load environment variables
func (m *Manager) SetLoadEnvVars(load bool) *Manager {
	m.loadEnvVars = load
	return m
}

// Load loads configuration from all sources
func (m *Manager) Load() error {
	// Load configuration files
	for _, path := range m.configPaths {
		if err := m.loadConfigFile(path); err != nil {
			return err
		}
	}

	// Load environment-specific configuration files
	for _, path := range m.configPaths {
		envPath := m.getEnvironmentConfigPath(path)
		if _, err := os.Stat(envPath); err == nil {
			if err := m.loadConfigFile(envPath); err != nil {
				return err
			}
		}
	}

	// Load environment variables if enabled
	if m.loadEnvVars {
		m.loadEnvironmentVariables()
	}

	return nil
}

// Get retrieves a configuration value by key
func (m *Manager) Get(key string, defaultValue interface{}) interface{} {
	// Check if the key exists in the configuration
	if value, ok := m.getNestedValue(key); ok {
		return value
	}

	// Return the default value if the key doesn't exist
	return defaultValue
}

// GetString retrieves a string configuration value by key
func (m *Manager) GetString(key string, defaultValue string) string {
	value := m.Get(key, defaultValue)
	if strValue, ok := value.(string); ok {
		return strValue
	}
	return defaultValue
}

// GetInt retrieves an integer configuration value by key
func (m *Manager) GetInt(key string, defaultValue int) int {
	value := m.Get(key, defaultValue)
	if intValue, ok := value.(int); ok {
		return intValue
	}
	if floatValue, ok := value.(float64); ok {
		return int(floatValue)
	}
	return defaultValue
}

// GetBool retrieves a boolean configuration value by key
func (m *Manager) GetBool(key string, defaultValue bool) bool {
	value := m.Get(key, defaultValue)
	if boolValue, ok := value.(bool); ok {
		return boolValue
	}
	return defaultValue
}

// GetFloat retrieves a float configuration value by key
func (m *Manager) GetFloat(key string, defaultValue float64) float64 {
	value := m.Get(key, defaultValue)
	if floatValue, ok := value.(float64); ok {
		return floatValue
	}
	if intValue, ok := value.(int); ok {
		return float64(intValue)
	}
	return defaultValue
}

// Set sets a configuration value by key
func (m *Manager) Set(key string, value interface{}) *Manager {
	keys := strings.Split(key, ".")
	m.setNestedValue(keys, value)
	return m
}

// Has checks if a configuration key exists
func (m *Manager) Has(key string) bool {
	_, ok := m.getNestedValue(key)
	return ok
}

// All returns all configuration values
func (m *Manager) All() map[string]interface{} {
	return m.values
}

// loadConfigFile loads configuration from a JSON file
func (m *Manager) loadConfigFile(path string) error {
	// Read the file
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	// Parse the JSON
	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse config file %s: %w", path, err)
	}

	// Merge the configuration
	m.mergeConfig(config)

	return nil
}

// getEnvironmentConfigPath returns the environment-specific configuration file path
func (m *Manager) getEnvironmentConfigPath(path string) string {
	ext := filepath.Ext(path)
	base := strings.TrimSuffix(path, ext)
	return fmt.Sprintf("%s.%s%s", base, m.environment, ext)
}

// loadEnvironmentVariables loads configuration from environment variables
func (m *Manager) loadEnvironmentVariables() {
	// Get all environment variables
	for _, env := range os.Environ() {
		// Split the environment variable into key and value
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		value := parts[1]

		// Convert the key to lowercase for case-insensitive matching
		key = strings.ToLower(key)

		// Set the configuration value
		m.Set(key, value)
	}
}

// mergeConfig merges a configuration map into the current configuration
func (m *Manager) mergeConfig(config map[string]interface{}) {
	for key, value := range config {
		if nestedMap, ok := value.(map[string]interface{}); ok {
			// If the value is a nested map, merge it recursively
			if existingValue, exists := m.values[key]; exists {
				if existingMap, ok := existingValue.(map[string]interface{}); ok {
					// Merge the nested maps
					for nestedKey, nestedValue := range nestedMap {
						existingMap[nestedKey] = nestedValue
					}
					continue
				}
			}
		}
		// Set the value directly
		m.values[key] = value
	}
}

// getNestedValue retrieves a nested configuration value by key
func (m *Manager) getNestedValue(key string) (interface{}, bool) {
	keys := strings.Split(key, ".")
	var current interface{} = m.values

	for _, k := range keys {
		// Convert the key to lowercase for case-insensitive matching
		k = strings.ToLower(k)

		// Check if the current value is a map
		currentMap, ok := current.(map[string]interface{})
		if !ok {
			return nil, false
		}

		// Get the value for the current key
		value, exists := currentMap[k]
		if !exists {
			return nil, false
		}

		// Update the current value
		current = value
	}

	return current, true
}

// setNestedValue sets a nested configuration value by key
func (m *Manager) setNestedValue(keys []string, value interface{}) {
	if len(keys) == 0 {
		return
	}

	// Convert the key to lowercase for case-insensitive matching
	key := strings.ToLower(keys[0])

	if len(keys) == 1 {
		// Set the value directly
		m.values[key] = value
		return
	}

	// Get or create the nested map
	var nestedMap map[string]interface{}
	if existingValue, exists := m.values[key]; exists {
		if existingMap, ok := existingValue.(map[string]interface{}); ok {
			nestedMap = existingMap
		} else {
			nestedMap = make(map[string]interface{})
			m.values[key] = nestedMap
		}
	} else {
		nestedMap = make(map[string]interface{})
		m.values[key] = nestedMap
	}

	// Set the value in the nested map
	var current interface{} = nestedMap
	for i := 1; i < len(keys)-1; i++ {
		// Convert the key to lowercase for case-insensitive matching
		k := strings.ToLower(keys[i])

		// Check if the current value is a map
		currentMap, ok := current.(map[string]interface{})
		if !ok {
			// Create a new map if the current value is not a map
			currentMap = make(map[string]interface{})
			current = currentMap
		}

		// Get or create the nested map
		if nestedValue, exists := currentMap[k]; exists {
			if nestedMap, ok := nestedValue.(map[string]interface{}); ok {
				current = nestedMap
			} else {
				newMap := make(map[string]interface{})
				currentMap[k] = newMap
				current = newMap
			}
		} else {
			newMap := make(map[string]interface{})
			currentMap[k] = newMap
			current = newMap
		}
	}

	// Set the value in the last nested map
	lastKey := strings.ToLower(keys[len(keys)-1])
	if currentMap, ok := current.(map[string]interface{}); ok {
		currentMap[lastKey] = value
	}
}

// Default configuration manager instance
var defaultManager = NewManager()

// SetEnvironment sets the current environment for the default manager
func SetEnvironment(env string) {
	defaultManager.SetEnvironment(env)
}

// GetEnvironment returns the current environment from the default manager
func GetEnvironment() string {
	return defaultManager.GetEnvironment()
}

// AddConfigPath adds a configuration file path to the default manager
func AddConfigPath(path string) {
	defaultManager.AddConfigPath(path)
}

// SetLoadEnvVars sets whether to load environment variables for the default manager
func SetLoadEnvVars(load bool) {
	defaultManager.SetLoadEnvVars(load)
}

// Load loads configuration from all sources for the default manager
func Load() error {
	return defaultManager.Load()
}

// Get retrieves a configuration value by key from the default manager
func Get(key string, defaultValue interface{}) interface{} {
	return defaultManager.Get(key, defaultValue)
}

// GetString retrieves a string configuration value by key from the default manager
func GetString(key string, defaultValue string) string {
	return defaultManager.GetString(key, defaultValue)
}

// GetInt retrieves an integer configuration value by key from the default manager
func GetInt(key string, defaultValue int) int {
	return defaultManager.GetInt(key, defaultValue)
}

// GetBool retrieves a boolean configuration value by key from the default manager
func GetBool(key string, defaultValue bool) bool {
	return defaultManager.GetBool(key, defaultValue)
}

// GetFloat retrieves a float configuration value by key from the default manager
func GetFloat(key string, defaultValue float64) float64 {
	return defaultManager.GetFloat(key, defaultValue)
}

// Set sets a configuration value by key in the default manager
func Set(key string, value interface{}) {
	defaultManager.Set(key, value)
}

// Has checks if a configuration key exists in the default manager
func Has(key string) bool {
	return defaultManager.Has(key)
}

// All returns all configuration values from the default manager
func All() map[string]interface{} {
	return defaultManager.All()
}
