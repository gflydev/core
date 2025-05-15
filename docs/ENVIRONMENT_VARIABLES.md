# gFly Core Environment Variables

This document provides a comprehensive list of all environment variables used by the gFly Core framework. These environment variables can be used to configure various aspects of the framework without modifying code.

## Application Configuration

| Variable | Description | Default Value | Type | Component |
|----------|-------------|---------------|------|-----------|
| `APP_NAME` | The name of the application | `"gFly"` | string | Core |
| `APP_CODE` | The code identifier for the application | `"gfly"` | string | Core |
| `APP_URL` | The base URL for the application | `"http://localhost:7789"` | string | Core |
| `APP_ENV` | The environment the application is running in (local, dev, stag, prod) | `"local"` | string | Core, Config |
| `APP_DEBUG` | Whether the application is in debug mode | `true` | boolean | Core |

## Server Configuration

| Variable | Description | Default Value | Type | Component |
|----------|-------------|---------------|------|-----------|
| `SERVER_HOST` | The host address to bind the server to | `"0.0.0.0"` | string | Core |
| `SERVER_PORT` | The port to run the server on | `7789` | integer | Core |
| `SERVER_TLS_CERT` | Path to the TLS certificate file | `""` | string | Core |
| `SERVER_TLS_KEY` | Path to the TLS key file | `""` | string | Core |

## Storage Configuration

| Variable | Description | Default Value | Type | Component |
|----------|-------------|---------------|------|-----------|
| `STORAGE_DIR` | The main storage directory for the application | `"storage"` | string | Core |
| `TEMP_DIR` | The directory for temporary files | `"storage/tmp"` | string | Core |
| `LOG_DIR` | The directory for log files | `"storage/logs"` | string | Core, Log |
| `APP_DIR` | The directory for application-specific storage | `"storage/app"` | string | Core |
| `STATIC_PATH` | The directory for static files to be served | `"public"` | string | Core |

## Logging Configuration

| Variable | Description | Default Value | Type | Component |
|----------|-------------|---------------|------|-----------|
| `LOG_CHANNEL` | The output channel for logs (file, stdout) | `"file"` | string | Log |
| `LOG_FILE` | The name of the log file | `"gfly.log"` | string | Log |
| `LOG_LEVEL` | The minimum log level to output (trace, debug, info, warn, error, fatal, panic) | `"trace"` | string | Log |

## Configuration System

The configuration system in gFly Core supports loading configuration from environment variables. All environment variables are automatically loaded into the configuration system and can be accessed using the `config` package.

### Environment-Specific Configuration

The framework supports different environments through the `APP_ENV` environment variable. The supported environments are:

- `local`: Local development environment
- `dev`: Development environment
- `stag`: Staging environment
- `prod`: Production environment

Environment-specific configuration files are loaded based on the current environment. For example, if `APP_ENV` is set to `dev`, the framework will load `app.dev.json` after loading `app.json`.

## Using Environment Variables

### In Application Code

You can access environment variables in your application code using the `utils.Getenv` function:

```go
import "github.com/gflydev/core/utils"

// Get a string environment variable with a default value
appName := utils.Getenv("APP_NAME", "Default App Name")

// Get an integer environment variable with a default value
serverPort := utils.Getenv("SERVER_PORT", 8080)

// Get a boolean environment variable with a default value
debugMode := utils.Getenv("APP_DEBUG", false)
```

### Using .env Files

The framework supports loading environment variables from `.env` files using the `github.com/joho/godotenv` package. To use this feature, include the following import in your application:

```go
import _ "github.com/joho/godotenv/autoload"
```

This will automatically load environment variables from a `.env` file in the root directory of your application.

## Best Practices

1. **Use environment variables for configuration that changes between environments**
   - Database connection strings
   - API keys and secrets
   - Feature flags

2. **Provide sensible defaults for all environment variables**
   - This ensures your application can run without extensive configuration

3. **Document all environment variables used by your application**
   - Include them in your project's README or a dedicated documentation file

4. **Use a `.env.example` file to document required environment variables**
   - This helps other developers understand what environment variables they need to set

5. **Never commit sensitive values in `.env` files to version control**
   - Add `.env` to your `.gitignore` file
   - Use `.env.example` with placeholder values instead