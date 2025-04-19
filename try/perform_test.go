package try

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPerform_NormalFlow(t *testing.T) {
	// Test cases for normal flow (no panics)
	tests := map[string]struct {
		withFinally bool
	}{
		"Without finally": {
			withFinally: false,
		},
		"With finally": {
			withFinally: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			calledTry := false
			calledFinally := false
			calledCatch := false

			it := Perform(func() {
				calledTry = true
			})

			if tt.withFinally {
				it.Finally(func() {
					calledFinally = true
				})
			}

			it.Catch(func(_ E) {
				calledCatch = true
			})

			assert.True(t, calledTry, "Try block should be called")
			assert.Equal(t, tt.withFinally, calledFinally, "Finally block should be called if specified")
			assert.False(t, calledCatch, "Catch block should not be called in normal flow")
		})
	}
}

func TestPerform_PanicInTry(t *testing.T) {
	// Test cases for panics in the try block
	tests := map[string]struct {
		withFinally bool
		panicValue  interface{}
	}{
		"String panic without finally": {
			withFinally: false,
			panicValue:  "testing panic",
		},
		"String panic with finally": {
			withFinally: true,
			panicValue:  "testing panic",
		},
		"Error panic with finally": {
			withFinally: true,
			panicValue:  errors.New("error panic"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			calledFinally := false
			calledCatch := false
			var caughtError interface{}

			it := Perform(func() {
				panic(tt.panicValue)
			})

			if tt.withFinally {
				it.Finally(func() {
					calledFinally = true
				})
			}

			it.Catch(func(e E) {
				calledCatch = true
				caughtError = e
			})

			assert.True(t, calledCatch, "Catch block should be called when try panics")
			assert.Equal(t, tt.withFinally, calledFinally, "Finally block should be called if specified")
			assert.Equal(t, tt.panicValue, caughtError, "Caught error should match the panic value")
		})
	}
}

func TestPerform_PanicInCatch(t *testing.T) {
	// Test cases for panics in the catch block
	tests := map[string]struct {
		withFinally     bool
		tryPanicValue   interface{}
		catchPanicValue interface{}
	}{
		"Without finally": {
			withFinally:     false,
			tryPanicValue:   "testing panic",
			catchPanicValue: "another panic",
		},
		"With finally": {
			withFinally:     true,
			tryPanicValue:   "testing panic",
			catchPanicValue: "another panic",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			calledFinally := false

			defer func() {
				err := recover()
				assert.Equal(t, tt.catchPanicValue, err, "Recovered panic should match catch panic value")

				if tt.withFinally {
					assert.True(t, calledFinally, "Finally block should be called before propagating panic")
				}
			}()

			it := Perform(func() {
				panic(tt.tryPanicValue)
			})

			if tt.withFinally {
				it.Finally(func() {
					calledFinally = true
				})
			}

			it.Catch(func(e E) {
				assert.Equal(t, tt.tryPanicValue, e, "Caught error should match try panic value")
				panic(tt.catchPanicValue)
			})

			t.Fatal("Test should not reach this point")
		})
	}
}

func TestPerform_Throw(t *testing.T) {
	// Test cases for using Throw
	tests := map[string]struct {
		withFinally bool
		throwValue  interface{}
		expected    interface{}
	}{
		"Throw nil (rethrow original)": {
			withFinally: true,
			throwValue:  nil,
			expected:    "testing panic",
		},
		"Throw specific error": {
			withFinally: false,
			throwValue:  "custom error",
			expected:    "custom error",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			calledFinally := false

			defer func() {
				err := recover()
				assert.Equal(t, tt.expected, err, "Recovered panic should match expected value")

				if tt.withFinally {
					assert.True(t, calledFinally, "Finally block should be called before propagating panic")
				}
			}()

			it := Perform(func() {
				panic("testing panic")
			})

			if tt.withFinally {
				it.Finally(func() {
					calledFinally = true
				})
			}

			it.Catch(func(e E) {
				assert.Equal(t, "testing panic", e, "Caught error should match try panic value")
				Throw(tt.throwValue)
			})

			t.Fatal("Test should not reach this point")
		})
	}
}

func TestPerform_PanicInFinally(t *testing.T) {
	// Test cases for panics in the finally block
	tests := map[string]struct {
		tryPanics       bool
		tryPanicValue   interface{}
		finallyPanicMsg string
	}{
		"Normal flow with finally panic": {
			tryPanics:       false,
			finallyPanicMsg: "finally panic",
		},
		"Try panics and finally panics": {
			tryPanics:       true,
			tryPanicValue:   "testing panic",
			finallyPanicMsg: "finally panic",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			calledTry := false
			calledCatch := false

			defer func() {
				err := recover()
				assert.Equal(t, tt.finallyPanicMsg, err, "Recovered panic should match finally panic message")
				assert.True(t, calledTry, "Try block should be called")

				if tt.tryPanics {
					assert.True(t, calledCatch, "Catch block should be called if try panics")
				} else {
					assert.False(t, calledCatch, "Catch block should not be called in normal flow")
				}
			}()

			it := Perform(func() {
				calledTry = true
				if tt.tryPanics {
					panic(tt.tryPanicValue)
				}
			})

			it.Finally(func() {
				panic(tt.finallyPanicMsg)
			})

			it.Catch(func(e E) {
				calledCatch = true
				if tt.tryPanics {
					assert.Equal(t, tt.tryPanicValue, e, "Caught error should match try panic value")
				}
			})

			t.Fatal("Test should not reach this point")
		})
	}
}

func TestFinally_CalledTwice(t *testing.T) {
	// Test that calling Finally twice causes a panic
	assert.Panics(t, func() {
		Perform(func() {
			// Do nothing
		}).Finally(func() {
			// First finally
		}).Finally(func() {
			// Second finally - should panic
		})
	}, "Calling Finally twice should panic")
}

func TestThrow_DirectCall(t *testing.T) {
	// Test direct calls to Throw
	tests := map[string]struct {
		throwValue interface{}
		expected   interface{}
	}{
		"Throw nil": {
			throwValue: nil,
			expected:   RethrowPanic,
		},
		"Throw string": {
			throwValue: "custom error",
			expected:   "custom error",
		},
		"Throw error": {
			throwValue: errors.New("error object"),
			expected:   errors.New("error object"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			defer func() {
				err := recover()
				if tt.expected == nil {
					assert.Nil(t, err, "Recovered panic should be nil")
				} else if e, ok := tt.expected.(error); ok && e.Error() == "error object" {
					// For error objects, compare error messages
					require.NotNil(t, err, "Recovered panic should not be nil")
					assert.Equal(t, "error object", err.(error).Error(), "Error message should match")
				} else {
					assert.Equal(t, tt.expected, err, "Recovered panic should match expected value")
				}
			}()

			Throw(tt.throwValue)
			t.Fatal("Test should not reach this point")
		})
	}
}
