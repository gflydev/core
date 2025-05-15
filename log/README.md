# Log

Write logs with color-coded output for different log levels and structured logging support.

## Usage

### Basic Logging

```golang
import "github.com/gflydev/core/log"

// Basic logging
log.Info("Server started")
log.Error("Invalid field name")

// Formatted logging
log.Infof("Server started on port %d", 8080)
log.Errorf("Invalid field name: %s", "username")

// Structured logging with key-value pairs
log.Infow("Request processed", 
    "method", "GET",
    "path", "/api/users",
    "status", 200,
    "duration_ms", 42.5,
)
```

### Structured Logging with JSON

The logger supports structured logging in JSON format:

```golang
logger := log.DefaultLogger()

// Enable structured logging with JSON format
logger.EnableStructuredLogging(true)

// Log structured data
logger.Infow("Request completed",
    "request_id", "req-456",
    "status", 200,
    "duration_ms", 42.5,
)

// This will output JSON like:
// {"time":"2023-05-15T14:30:45Z","level":"INFO","message":"Request completed","request_id":"req-456","status":200,"duration_ms":42.5}
```

### Customizing the JSON Formatter

You can customize the JSON formatter:

```golang
// Create a custom JSON formatter
formatter := log.NewJSONFormatter()
formatter.TimeKey = "timestamp"
formatter.LevelKey = "severity"
formatter.MessageKey = "msg"
formatter.TimeFormat = time.RFC3339Nano

// Set the formatter
logger.SetFormatter(formatter)
```

## Color Coding

The logger uses color-coded output for different log levels to improve readability:

- **TRACE**: Gray
- **DEBUG**: Blue
- **INFO**: Green
- **WARN**: Yellow
- **ERROR**: Red
- **FATAL**: Purple
- **PANIC**: Cyan

Colors are automatically applied when logging to a terminal.

## Log Formats

### Text Format (Default)

The logger outputs log messages with the following format:

```
LEVEL message
```

For example:
```
ERROR Invalid field name
```

The log level is color-coded as described above.

### JSON Format (Structured)

When structured logging is enabled, the logger outputs log messages in JSON format:

```json
{
  "time": "2023-05-15T14:30:45Z",
  "level": "INFO",
  "message": "Request completed",
  "request_id": "req-456",
  "status": 200,
  "duration_ms": 42.5
}
```

This format is ideal for log aggregation and analysis tools.
