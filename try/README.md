# Try-Finally-Catch

A Go package that provides structured error handling similar to try-catch-finally blocks in other languages, but adapted for Go's panic/recover mechanism.

## Overview

The `try` package offers a clean and structured way to handle errors and panics in Go code. It allows you to:

- Execute code that might panic in a controlled environment
- Ensure cleanup code always runs, even if an error occurs
- Handle errors in a structured way
- Rethrow or transform errors as needed

## Installation

```bash
go get github.com/gflydev/core
```

## Usage

### Basic Example

```
package main

import (
    "github.com/gflydev/core/try"
    "log"
)

func main() {
    try.Perform(func() {
        // Code that might panic
        result := riskyOperation()
        processResult(result)
    }).Finally(func() {
        // Cleanup code that always runs
        cleanup()
    }).Catch(func(e try.E) {
        // Error handling code
        log.Printf("Error occurred: %v", e)
    })
}
```

### Error Handling Patterns

#### Simple Error Handling

```
package example

import (
    "github.com/gflydev/core/try"
    "log"
)

func example() {
    try.Perform(func() {
        if err := someOperation(); err != nil {
            panic(err)
        }
    }).Catch(func(e try.E) {
        log.Printf("Operation failed: %v", e)
    })
}
```

#### Using Throw for Custom Errors

```
package example

import (
    "github.com/gflydev/core/try"
    "log"
)

func example() {
    try.Perform(func() {
        if value < 0 {
            try.Throw("Value cannot be negative")
        }
    }).Catch(func(e try.E) {
        log.Printf("Validation error: %v", e)
    })
}
```

#### Rethrowing Errors

```
package example

import (
    "github.com/gflydev/core/try"
    "log"
)

func example() {
    try.Perform(func() {
        // Some code that might panic
    }).Catch(func(e try.E) {
        log.Printf("Error occurred: %v", e)

        // Rethrow the original error
        try.Throw(nil)
    })
}
```

#### Resource Management

```
package example

import (
    "github.com/gflydev/core/try"
    "log"
    "os"
)

func example() {
    var file *os.File

    try.Perform(func() {
        var err error
        file, err = os.Open("example.txt")
        if err != nil {
            panic(err)
        }

        // Process the file
        processFile(file)
    }).Finally(func() {
        // Close the file if it was opened
        if file != nil {
            file.Close()
        }
    }).Catch(func(e try.E) {
        log.Printf("File operation failed: %v", e)
    })
}
```

## Best Practices

1. **Keep Try Blocks Focused**: Each try block should focus on a specific operation that might fail.

2. **Always Use Finally for Cleanup**: Use Finally for any cleanup operations that should always run, regardless of whether an error occurred.

3. **Be Specific in Catch Blocks**: Handle specific error types when possible, rather than catching all errors.

4. **Avoid Nested Try Blocks**: Instead of nesting try blocks, chain them sequentially for better readability.

5. **Don't Overuse**: While this package provides a convenient way to handle errors, it's not a replacement for Go's standard error handling. Use it for complex error handling scenarios where the structured approach adds clarity.

## API Reference

### Functions

- `Perform(func())`: Starts a try-catch-finally block by executing the provided function.
- `Throw(error)`: Explicitly throws a panic that can be caught by a catch block.

### Methods

- `Finally(func())`: Registers a function to be executed at the end of the try-catch block.
- `Catch(func(error))`: Registers an error-handling function that is executed if an error occurs.

## License

This package is part of the gFlyDev Core library and is licensed under the same terms.
