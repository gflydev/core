# Utilities

Common utils

### Usage

Quick usage
```go
import "github.com/gflydev/core/utils"

newHello := utils.CopyStr("Hello world")
```

### Bytes

```go
// Random byte
myBytes := utils.RandByte(make([]byte, 15))

// Extend bytes
myBytes := utils.ExtendByte(myBytes, 5)
```

### Token

```go
// Random UNIQUE token
myToken := utils.Token()
```