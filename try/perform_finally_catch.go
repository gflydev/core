package try

// RethrowPanic Special re-throw panic
const RethrowPanic = "___throw_it___"

// Define types of Try-Finally-Catch
type (
	// F Function type
	F func()
	// E Error type
	E interface{}
	// EF Error function type
	EF func(err E)
	// It structure
	It struct {
		finally F
		Error   E
	}
)

// Throw function (return or rethrow an exception)
// Parameters:
//   - e: The error to be thrown. If nil, a default panic with RethrowPanic value is thrown.
func Throw(e E) {
	// Throw default error
	if e == nil {
		panic(RethrowPanic)
	} else {
		// Throw a specific exception
		panic(e)
	}
}

// Perform registers the main-logic function and executes it.
// Parameters:
//   - funcToTry: The main logic to execute within the try block.
//
// Returns:
//   - *It: An instance of the It structure containing the result and error (if any).
func Perform(funcToTry F) (o *It) {
	// Initial exception object with null values
	o = &It{nil, nil}

	// Catch throw in from main logic
	defer func() {
		o.Error = recover()
	}()

	// Perform main logic
	funcToTry()

	// Response instance of It instance
	return
}

// Finally registers the finally-logic function that is executed at the end of the try-catch block.
// Parameters:
//   - finallyFunc: The function containing cleanup or finalization logic to be executed.
//
// Returns:
//   - *It: The same instance of the It structure for chaining.
func (o *It) Finally(finallyFunc F) *It {
	if o.finally != nil {
		panic("Finally Function by default !!")
	} else {
		o.finally = finallyFunc
	}

	return o
}

// Catch registers the error-handling function that is executed if an error occurs.
// Parameters:
//   - funcCaught: The function to handle the error, which receives the error as a parameter.
//
// Returns:
//   - *It: The same instance of the It structure for chaining.
func (o *It) Catch(funcCaught EF) *It {
	// Check if it has Error
	if o.Error != nil {
		// Catch error in from catching logic
		defer func() {
			// Call finally (Before receive error from recovering process)
			if o.finally != nil {
				o.finally()
			}

			// Receive error from recovering process
			if err := recover(); err != nil {
				// If it is just re-throw panic Exception
				if err == RethrowPanic {
					err = o.Error
				}
				panic(err)
			}
		}()

		// Perform catching logic
		funcCaught(o.Error)
	} else if o.finally != nil {
		// Perform finally logic
		o.finally()
	}

	return o
}
