# gFly Core Development Guidelines

This document provides essential information for developers working on the gFly Core project.

## Build/Configuration Instructions

### Prerequisites

- Go 1.24.0 or higher
- Make (for running build commands)
- Code quality tools:
  - gocritic
  - gosec
  - govulncheck
  - golangci-lint

### Setting Up the Development Environment

1. Clone the repository:
   ```bash
   git clone https://github.com/gflydev/core.git
   cd core
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Install required code quality tools:
   ```bash
   go install github.com/go-critic/go-critic/cmd/gocritic@latest
   go install github.com/securego/gosec/v2/cmd/gosec@latest
   go install golang.org/x/vuln/cmd/govulncheck@latest
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   ```

## Testing Information

### Running Tests

The project uses Go's standard testing package along with the `github.com/stretchr/testify` package for assertions.

#### Running All Tests

```bash
make test
```

#### Running Tests for Specific Packages

```bash
make test.try    # Run tests for the try package
make test.utils  # Run tests for the utils package
make test.log    # Run tests for the log package
make test.errors # Run tests for the errors package
```

#### Running Tests with Coverage

```bash
make test.cover
```
This command generates a coverage report and opens it in your browser.

### Writing Tests

1. Create a test file with the naming convention `<filename>_test.go` in the same package as the code being tested.

2. Use table-driven tests for comprehensive test coverage:

```go
func TestFunction(t *testing.T) {
    tests := map[string]struct {
        input    string
        expected string
    }{
        "Test case 1": {
            input:    "input1",
            expected: "expected1",
        },
        "Test case 2": {
            input:    "input2",
            expected: "expected2",
        },
    }

    for name, tt := range tests {
        t.Run(name, func(t *testing.T) {
            result := Function(tt.input)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

3. Use the `testify/assert` or `testify/require` packages for assertions:
   - Use `assert` for non-fatal assertions that allow the test to continue
   - Use `require` for fatal assertions that stop the test on failure

### Example Test

Here's an example test for the `IncludeStr` function in the `utils` package:

```go
func TestIncludeStr(t *testing.T) {
    tests := map[string]struct {
        slice    []string
        search   string
        expected bool
    }{
        "Empty slice": {
            slice:    []string{},
            search:   "hello",
            expected: false,
        },
        "String found": {
            slice:    []string{"hello", "world", "test"},
            search:   "world",
            expected: true,
        },
        "String not found": {
            slice:    []string{"hello", "world", "test"},
            search:   "golang",
            expected: false,
        },
    }

    for name, tt := range tests {
        t.Run(name, func(t *testing.T) {
            result := IncludeStr(tt.slice, tt.search)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

## Code Quality

The project uses several tools to ensure code quality:

### Static Analysis

```bash
make critic      # Run gocritic for static code analysis
make security    # Run gosec for security scanning
make vulncheck   # Run govulncheck for vulnerability checking
make lint        # Run golangci-lint for linting
```

### Running All Quality Checks

```bash
make all
```

This command runs all code quality checks and tests.

## Code Style Guidelines

1. **Documentation**: All exported functions, types, and constants should have proper documentation comments.
   ```go
   // FunctionName does something specific.
   //
   // Parameters:
   //   - param1 type: Description of param1.
   //
   // Returns:
   //   - type: Description of return value.
   func FunctionName(param1 type) type {
       // Implementation
   }
   ```

2. **Error Handling**: Use the custom error types in the `errors` package for consistent error handling.

3. **Testing**: Write comprehensive tests for all new functionality using table-driven tests.

4. **Naming Conventions**:
   - Use camelCase for variable names
   - Use PascalCase for exported functions, types, and constants
   - Use snake_case for test function names

5. **Package Structure**: Keep related functionality in the same package. Create new packages for distinct functionality.

## Project Structure

- `errors/`: Custom error types and error handling utilities
- `log/`: Logging functionality
- `try/`: Error handling utilities
- `utils/`: General utility functions
- `examples/`: Example applications using the framework

## Additional Resources

- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Effective Go](https://golang.org/doc/effective_go)