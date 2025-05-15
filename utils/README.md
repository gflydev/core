# Utils Package

The utils package provides a comprehensive collection of utility functions for common operations in the gFly framework. These utilities are designed to be efficient, easy to use, and to follow consistent patterns.

## Overview

The package is organized into multiple files, each focusing on a specific type of utility:

- **Bytes**: Utilities for byte slice manipulation
- **Strings**: Utilities for string manipulation
- **File**: Utilities for file system operations
- **Env**: Utilities for environment variable handling
- **Hash**: Utilities for hashing operations
- **Map**: Utilities for map operations
- **Number**: Utilities for number operations
- **Password**: Utilities for password hashing and verification
- **Reflection**: Utilities for reflection operations
- **Request**: Utilities for HTTP request handling
- **Slice**: Utilities for slice operations
- **Time**: Utilities for time operations
- **Token**: Utilities for generating unique tokens

## Installation

```bash
go get github.com/gflydev/core/utils
```

## Usage Examples

### Bytes Utilities

```go
package main

import (
    "fmt"
    "github.com/gflydev/core/utils"
)

func main() {
    // Generate random bytes
    randomBytes := utils.RandByte(make([]byte, 16))
    fmt.Println(randomBytes)

    // Example byte slice
    originalBytes := []byte{1, 2, 3}

    // Extend a byte slice
    extendedBytes := utils.ExtendByte(originalBytes, 5)
    fmt.Println(extendedBytes)

    // Prepend bytes
    newBytes := utils.PrependByte(originalBytes, 0x01, 0x02, 0x03)
    fmt.Println(newBytes)

    // Copy bytes
    bytesCopy := utils.CopyByte(originalBytes)
    fmt.Println(bytesCopy)

    // Compare bytes
    bytes1 := []byte{1, 2, 3}
    bytes2 := []byte{1, 2, 3}
    areEqual := utils.EqualByte(bytes1, bytes2)
    fmt.Println(areEqual)
}
```

### String Utilities

```go
package main

import (
    "fmt"
    "github.com/gflydev/core/utils"
)

func main() {
    // Convert between strings and bytes without allocation
    myString := "Hello, World!"
    bytes := utils.UnsafeBytes(myString)
    fmt.Println(bytes)

    myBytes := []byte("Hello, World!")
    str := utils.UnsafeStr(myBytes)
    fmt.Println(str)

    // Create a copy of a string
    strCopy := utils.CopyStr(myString)
    fmt.Println(strCopy)

    // Check if a string is in a slice
    mySlice := []string{"apple", "banana", "cherry"}
    isIncluded := utils.IncludeStr(mySlice, "banana")
    fmt.Println(isIncluded)

    // Truncate a string
    longString := "This is a very long string that needs to be truncated"
    truncated := utils.Truncate(longString, 20, true)
    fmt.Println(truncated)

    // Check if a string is empty or contains only whitespace
    isEmpty := utils.IsEmpty("")
    fmt.Println(isEmpty)

    isBlank := utils.IsBlank("   ")
    fmt.Println(isBlank)

    // Convert string case
    camel := utils.ToCamelCase("my-variable-name")
    fmt.Println(camel)

    pascal := utils.ToPascalCase("my-class-name")
    fmt.Println(pascal)

    snake := utils.ToSnakeCase("MyMethodName")
    fmt.Println(snake)
}
```

### File Utilities

```go
package main

import (
    "fmt"
    "github.com/gflydev/core/utils"
    "os"
)

func main() {
    // Get a file extension
    ext := utils.FileExt("document.pdf")
    fmt.Println("File extension:", ext)

    // Check if a file exists
    exists := utils.FileExists("example.txt")
    fmt.Println("File exists:", exists)

    // Create a test file
    err := os.WriteFile("example.txt", []byte("Hello, World!"), 0644)
    if err != nil {
        fmt.Println("Error creating test file:", err)
        return
    }

    // Read a file as string
    content, err := utils.ReadFileAsString("example.txt")
    if err != nil {
        fmt.Println("Error reading file:", err)
    } else {
        fmt.Println("File content:", content)
    }

    // Write a string to a file
    err = utils.WriteStringToFile("example2.txt", "Hello, Modified World!", 0644)
    if err != nil {
        fmt.Println("Error writing to file:", err)
    }

    // Copy a file
    err = utils.CopyFile("example.txt", "example_copy.txt", 0644)
    if err != nil {
        fmt.Println("Error copying file:", err)
    }

    // Create a directory if it doesn't exist
    err = utils.MkdirIfNotExists("example_dir", 0755)
    if err != nil {
        fmt.Println("Error creating directory:", err)
    }

    // Clean up
    os.Remove("example.txt")
    os.Remove("example2.txt")
    os.Remove("example_copy.txt")
    os.RemoveAll("example_dir")
}
```

### Environment Utilities

```go
package main

import (
    "fmt"
    "github.com/gflydev/core/utils"
    "os"
)

func main() {
    // Set some environment variables for demonstration
    os.Setenv("PORT", "9090")
    os.Setenv("TIMEOUT", "60")
    os.Setenv("DEBUG", "true")

    // Get an environment variable with a default value
    port := utils.Env("PORT", "8080")
    fmt.Println("Port:", port)

    // Get a non-existent environment variable (will use default)
    defaultPort := utils.Env("NONEXISTENT_PORT", "8080")
    fmt.Println("Default port:", defaultPort)

    // Get an environment variable as an integer
    timeout := utils.EnvInt("TIMEOUT", 30)
    fmt.Println("Timeout:", timeout)

    // Get an environment variable as a boolean
    debug := utils.EnvBool("DEBUG", false)
    fmt.Println("Debug mode:", debug)
}
```

### Token Utilities

```go
package main

import (
    "fmt"
    "github.com/gflydev/core/utils"
)

func main() {
    // Generate a random unique token
    token := utils.Token()
    fmt.Println("Random token:", token)

    // Generate a token with a specific prefix
    prefixedToken := utils.TokenWithPrefix("user_")
    fmt.Println("Prefixed token:", prefixedToken)

    // Generate multiple tokens to demonstrate uniqueness
    fmt.Println("Multiple tokens:")
    for i := 0; i < 3; i++ {
        fmt.Println("  -", utils.Token())
    }
}
```

### Password Utilities

```go
package main

import (
    "fmt"
    "github.com/gflydev/core/utils"
)

func main() {
    password := "my-secure-password"

    // Hash a password
    hashedPassword, err := utils.HashPassword(password)
    if err != nil {
        fmt.Println("Error hashing password:", err)
        return
    }
    fmt.Println("Hashed password:", hashedPassword)

    // Verify a password against a hash (correct password)
    isValid := utils.CheckPasswordHash(password, hashedPassword)
    fmt.Println("Password valid:", isValid)

    // Verify a password against a hash (incorrect password)
    isInvalid := utils.CheckPasswordHash("wrong-password", hashedPassword)
    fmt.Println("Wrong password valid:", isInvalid)
}
```

## Best Practices

1. **Error Handling**: Always check for errors returned by functions that can fail
2. **Performance**: Use the unsafe functions only when performance is critical and you understand the implications
3. **Immutability**: Use copy functions when you need to ensure data is not modified unexpectedly
4. **Consistency**: Follow the patterns established in the package for your own utility functions

## Contributing

Contributions to the utils package are welcome. Please ensure that any new utility functions:

1. Are well-documented with comments explaining parameters, returns, and behavior
2. Include appropriate error handling
3. Are tested thoroughly
4. Follow the existing code style and patterns

For more information on contributing, see the main project's CONTRIBUTING.md file.
