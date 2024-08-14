# Try-Finally-Catch

### Usage

Quick usage
```go
import "github.com/gflydev/core/try"

try.Perform(func() {
    calledTry()
}).Finally(func() {
    calledFinally()
}).Catch(func(e try.E) {
    log.Errorf("Catch error %v", e)
})
```
