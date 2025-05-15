package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gflydev/core/config"
)

func main() {
	// Create a directory for config files if it doesn't exist
	if err := os.MkdirAll("config", 0755); err != nil {
		log.Fatalf("Failed to create config directory: %v", err)
	}

	// Create base configuration file
	baseConfig := `{
		"app": {
			"name": "gFly Example App",
			"version": "1.0.0",
			"debug": false
		},
		"database": {
			"host": "localhost",
			"port": 3306,
			"username": "root",
			"password": "password",
			"database": "gfly"
		},
		"server": {
			"host": "0.0.0.0",
			"port": 8080
		}
	}`
	if err := os.WriteFile("config/app.json", []byte(baseConfig), 0644); err != nil {
		log.Fatalf("Failed to write base config file: %v", err)
	}

	// Create development environment configuration file
	devConfig := `{
		"app": {
			"debug": true
		},
		"database": {
			"host": "127.0.0.1",
			"password": "dev_password"
		},
		"server": {
			"port": 8081
		}
	}`
	if err := os.WriteFile("config/app.dev.json", []byte(devConfig), 0644); err != nil {
		log.Fatalf("Failed to write dev config file: %v", err)
	}

	// Create production environment configuration file
	prodConfig := `{
		"app": {
			"debug": false
		},
		"database": {
			"host": "db.example.com",
			"password": "prod_password"
		},
		"server": {
			"host": "0.0.0.0",
			"port": 80
		}
	}`
	if err := os.WriteFile("config/app.prod.json", []byte(prodConfig), 0644); err != nil {
		log.Fatalf("Failed to write prod config file: %v", err)
	}

	// Set environment variables for demonstration
	_ = os.Setenv("APP_ENV", "dev")
	_ = os.Setenv("DATABASE_USERNAME", "dev_user")

	// Example 1: Using the default configuration manager
	fmt.Println("Example 1: Using the default configuration manager")
	fmt.Println("--------------------------------------------------")

	// Add configuration file path
	config.AddConfigPath("config/app.json")

	// Load configuration
	if err := config.Load(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Get configuration values
	appName := config.GetString("app.name", "Default App Name")
	appDebug := config.GetBool("app.debug", false)
	dbHost := config.GetString("database.host", "localhost")
	dbPort := config.GetInt("database.port", 3306)
	dbUsername := config.GetString("database.username", "root")
	serverPort := config.GetInt("server.port", 8080)

	fmt.Printf("Environment: %s\n", config.GetEnvironment())
	fmt.Printf("App Name: %s\n", appName)
	fmt.Printf("App Debug: %t\n", appDebug)
	fmt.Printf("DB Host: %s\n", dbHost)
	fmt.Printf("DB Port: %d\n", dbPort)
	fmt.Printf("DB Username: %s\n", dbUsername)
	fmt.Printf("Server Port: %d\n", serverPort)

	// Example 2: Using a custom configuration manager for production
	fmt.Println("\nExample 2: Using a custom configuration manager for production")
	fmt.Println("----------------------------------------------------------")

	// Create a new configuration manager
	prodManager := config.NewManager()

	// Configure the manager
	prodManager.SetEnvironment(config.EnvProduction)
	prodManager.AddConfigPath("config/app.json")
	prodManager.SetLoadEnvVars(true)

	// Load configuration
	if err := prodManager.Load(); err != nil {
		log.Fatalf("Failed to load production configuration: %v", err)
	}

	// Get configuration values
	prodAppName := prodManager.GetString("app.name", "Default App Name")
	prodAppDebug := prodManager.GetBool("app.debug", false)
	prodDbHost := prodManager.GetString("database.host", "localhost")
	prodDbPort := prodManager.GetInt("database.port", 3306)
	prodDbUsername := prodManager.GetString("database.username", "root")
	prodServerPort := prodManager.GetInt("server.port", 8080)

	fmt.Printf("Environment: %s\n", prodManager.GetEnvironment())
	fmt.Printf("App Name: %s\n", prodAppName)
	fmt.Printf("App Debug: %t\n", prodAppDebug)
	fmt.Printf("DB Host: %s\n", prodDbHost)
	fmt.Printf("DB Port: %d\n", prodDbPort)
	fmt.Printf("DB Username: %s\n", prodDbUsername)
	fmt.Printf("Server Port: %d\n", prodServerPort)

	// Example 3: Setting and overriding configuration values
	fmt.Println("\nExample 3: Setting and overriding configuration values")
	fmt.Println("---------------------------------------------------")

	// Create a new configuration manager
	customManager := config.NewManager()

	// Set some configuration values
	customManager.Set("app.name", "Custom App")
	customManager.Set("app.version", "2.0.0")
	customManager.Set("database.host", "custom-db.example.com")

	// Load configuration (will merge with existing values)
	customManager.AddConfigPath("config/app.json")
	if err := customManager.Load(); err != nil {
		log.Fatalf("Failed to load custom configuration: %v", err)
	}

	// Get configuration values
	customAppName := customManager.GetString("app.name", "Default App Name")
	customAppVersion := customManager.GetString("app.version", "1.0.0")
	customDbHost := customManager.GetString("database.host", "localhost")

	fmt.Printf("App Name: %s\n", customAppName)
	fmt.Printf("App Version: %s\n", customAppVersion)
	fmt.Printf("DB Host: %s\n", customDbHost)

	// Clean up
	os.RemoveAll("config")
}
