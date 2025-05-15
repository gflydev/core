package try

import (
	"github.com/gflydev/core/errors"
)

// RethrowPanic is a special marker used to indicate that the original panic should be rethrown.
// This is used internally by the Throw function when called with nil.
const RethrowPanic = "___throw_it___"

// Define types for the Try-Finally-Catch pattern
type (
	// F represents a function with no parameters and no return value.
	// Used for try and finally blocks.
	F func()

	// E represents any type of error or panic value.
	// This is an alias for interface{} to allow catching any type of panic.
	E interface{}

	// EF represents a function that takes an error parameter.
	// Used for catch blocks to handle errors.
	EF func(err E)

	// It is the main structure that chains try, finally, and catch blocks.
	// It holds the state of the error handling process.
	It struct {
		finally F // The function to execute in the finally block
		Error   E // The error/panic value if one occurred
	}
)

// Throw explicitly throws a panic that can be caught by a catch block.
// This is useful for rethrowing errors or creating custom errors.
//
// Parameters:
//   - e: The error to be thrown. If nil, the original panic value will be rethrown.
//
// Behavior:
//   - If e is nil, a special RethrowPanic marker is used to indicate the original error should be used.
//   - If e is a standard error, it will be wrapped with the appropriate gFly error type.
//   - If e is already a gFly error, it will be used directly.
//   - Otherwise, the provided value is directly used as the panic value.
//
// Example:
//
//	Perform(func() {
//	    if err := someOperation(); err != nil {
//	        Throw(err) // Throw a specific error
//	    }
//	}).Catch(func(e E) {
//	    // Handle the error
//	})
func Throw(e E) {
	if e == nil {
		// Use the special marker to indicate we should rethrow the original error
		panic(RethrowPanic)
	} else if stdErr, ok := e.(error); ok {
		// Check if it's already a gFly error
		if _, ok := stdErr.(errors.Error); !ok {
			// Wrap standard errors with the appropriate gFly error type
			panic(errors.Wrap(stdErr, errors.CodeInternal, "Error thrown in try-catch block"))
		} else {
			// It's already a gFly error, use it directly
			panic(stdErr)
		}
	} else {
		// Throw the specific value provided (not an error)
		panic(e)
	}
}

// Perform starts a try-catch-finally block by executing the provided function.
// This is the entry point for the error handling pattern.
//
// Parameters:
//   - funcToTry: The main logic to execute within the try block.
//
// Returns:
//   - *It: An instance of the It structure for chaining Finally and Catch calls.
//
// Example:
//
//	Perform(func() {
//	    // Code that might panic
//	}).Finally(func() {
//	    // Cleanup code that always runs
//	}).Catch(func(e E) {
//	    // Error handling code
//	})
func Perform(funcToTry F) (o *It) {
	// Create a new It instance with no finally function and no error
	o = &It{nil, nil}

	// Set up recovery to catch any panics from the try block
	defer func() {
		o.Error = recover()
	}()

	// Execute the try block
	funcToTry()

	// Return It instance for chaining
	return
}

// Finally registers a function to be executed at the end of the try-catch block,
// regardless of whether an error occurred or not.
//
// Parameters:
//   - finallyFunc: The function containing cleanup or finalization logic.
//
// Returns:
//   - *It: The same instance of It structure for chaining.
//
// Panics:
//   - If Finally is called more than once on the same It instance.
//
// Example:
//
//	Perform(func() {
//	    // Open a file
//	}).Finally(func() {
//	    // Close the file, regardless of errors
//	}).Catch(func(e E) {
//	    // Handle any errors
//	})
func (o *It) Finally(finallyFunc F) *It {
	if o.finally != nil {
		panic("Finally function already registered. Cannot register multiple Finally blocks.")
	}

	o.finally = finallyFunc
	return o
}

// Catch registers an error-handling function that is executed if an error occurs
// in the try block. If no error occurred, the catch block is skipped.
//
// Parameters:
//   - funcCaught: The function to handle the error, which receives the error as a parameter.
//
// Returns:
//   - *It: The same instance of It structures for chaining.
//
// Behavior:
//   - If an error occurred in the try block, the catch function is executed with the error.
//   - If no error occurred, the catch function is skipped.
//   - The final function (if registered) is always executed, even if the catch block panics.
//   - If the catch block panics, the panic is propagated after the finally block executes.
//   - If Throw(nil) is called in the catch block, the original error is rethrown.
//   - Standard errors are automatically wrapped as gFly errors if they aren't already.
//
// Example:
//
//	Perform(func() {
//	    // Code that might panic
//	}).Catch(func(e E) {
//	    // Handle the error
//	    if gflyErr, ok := e.(errors.Error); ok {
//	        // Handle gFly error with additional context
//	        log.Printf("Error code: %s, Message: %s", gflyErr.Code(), gflyErr.Message())
//	    } else {
//	        // Handle standard error
//	        log.Printf("Error: %v", e)
//	    }
//	})
func (o *It) Catch(funcCaught EF) *It {
	// Check if an error occurred in the try block
	if o.Error != nil {
		// Ensure the error is a gFly error if it's a standard error
		if stdErr, ok := o.Error.(error); ok && !isGFlyError(stdErr) {
			o.Error = errors.Wrap(stdErr, errors.CodeInternal, "Error caught in try-catch block")
		}

		// Set up recovery to catch any panics from the catch block
		defer func() {
			// Always execute the finally block if one is registered
			if o.finally != nil {
				o.finally()
			}

			// Check if the catch block panicked
			if err := recover(); err != nil {
				// If the special RethrowPanic marker was used, replace it with the original error
				if err == RethrowPanic {
					err = o.Error
				}
				// Propagate the panic
				panic(err)
			}
		}()

		// Execute the catch block with the error
		funcCaught(o.Error)
	} else if o.finally != nil {
		// If no error occurred but a finally block is registered, execute it
		o.finally()
	}

	return o
}

// isGFlyError checks if an error is a gFly error
func isGFlyError(err error) bool {
	_, ok := err.(errors.Error)
	return ok
}
