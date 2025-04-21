# Log

Write logs with color-coded output for different log levels.

### Usage

Quick usage
```golang
import "github.com/gflydev/core/log"

log.Errorf("Invalid field name")
```

### Color Coding

The logger uses color-coded output for different log levels to improve readability:

- **TRACE**: Gray
- **DEBUG**: Blue
- **INFO**: Green
- **WARN**: Yellow
- **ERROR**: Red
- **FATAL**: Purple
- **PANIC**: Cyan

Colors are automatically applied when logging to a terminal.

### Log Format

The logger outputs log messages with the following format:

```
LEVEL message
```

For example:
```
ERROR Invalid field name
```

The log level is color-coded as described above.
