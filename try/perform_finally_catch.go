package try

// RethrowPanic is a special marker used to indicate that the original panic should be rethrown.
// This is used internally by the Throw function when called with nil.
const RethrowPanic = "___throw_it___"

// Define types for the Try-Finally-Catch pattern
type (
	// F represents a function with no parameters and no return value.
	// Used for try and finally blocks.
	F func()

	// E represents any type of error or panic value.
	// This is an alias for any to allow catching any type of panic.
	E any

	// EF represents a function that takes an error parameter.
	// Used for catch blocks to handle errors.
	EF func(err E)

	// It is the main structure that chains try, finally, and catch blocks.
	// It holds the state of the error handling process.
	It struct {
		finally     F    // The function to execute in the finally block
		Error       E    // The error/panic value if one occurred
		caught      bool // Whether Catch has already processed the error
		finallyDone bool // Whether the finally block has already executed
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
	} else {
		// Throw the specific error provided
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
	o = &It{}

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

	// Run the finally block now when there is nothing left to wait for:
	//   - no error occurred (Catch would be a no-op), or
	//   - the error was already handled by a preceding Catch (reversed order:
	//     Perform(f).Catch(c).Finally(g)).
	// When an unhandled error is pending, defer execution to Catch so the
	// documented order (catch runs before finally) is preserved.
	if o.Error == nil || o.caught {
		o.runFinally()
	}

	return o
}

// runFinally executes the registered finally block exactly once.
func (o *It) runFinally() {
	if o.finally != nil && !o.finallyDone {
		o.finallyDone = true
		o.finally()
	}
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
//
// Example:
//
//	Perform(func() {
//	    // Code that might panic
//	}).Catch(func(e E) {
//	    // Handle the error
//	    log.Printf("Error: %v", e)
//	})
func (o *It) Catch(funcCaught EF) *It {
	// Mark the error as handled so a Finally chained *after* this Catch still
	// runs (reversed order: Perform(f).Catch(c).Finally(g)).
	o.caught = true

	// Check if an error occurred in the try block
	if o.Error != nil {
		// Set up recovery to catch any panics from the catch block
		defer func() {
			// Always execute the finally block if one is registered
			o.runFinally()

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
	} else {
		// If no error occurred but a finally block is registered, execute it
		o.runFinally()
	}

	return o
}
