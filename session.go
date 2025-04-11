package core

// ====================================================================
//                            Session Structure
// ====================================================================

type ISession interface {
	// Set stores a value in the session associated with the given key.
	// Parameters:
	//  - c: The context object used to manage session information.
	//  - key: The key to associate with the value.
	//  - value: The value to store in the session.
	Set(c *Ctx, key string, value interface{})

	// Get retrieves a value from the session associated with the given key.
	// Parameters:
	//  - c: The context object used to manage session information.
	//  - key: The key associated with the value to retrieve.
	// Returns:
	//  The value stored in the session, or nil if not found.
	Get(c *Ctx, key string) interface{}
}

// ====================================================================
//                            Default Session
// ====================================================================

var sessionError = "Session manager is NULL. Please use import session module and run `session.Setup()`"

type DefaultSession struct {
}

// Set stores a value in the session associated with the given key.
// Parameters:
//   - c: The context object used to manage session information.
//   - key: The key to associate with the value.
//   - value: The value to store in the session.
func (v *DefaultSession) Set(c *Ctx, key string, value interface{}) {
	panic(sessionError)
}

// Get retrieves a value from the session associated with the given key.
// Parameters:
//   - c: The context object used to manage session information.
//   - key: The key associated with the value to retrieve.
//
// Returns:
//
//	The value stored in the session, or nil if not found.
func (v *DefaultSession) Get(c *Ctx, key string) interface{} {
	panic(sessionError)
}

var session ISession = &DefaultSession{}

// ====================================================================
//                              Functions
// ====================================================================

// RegisterSession registers a custom implementation of the ISession interface.
// Parameters:
//   - s: The custom implementation of the ISession interface.
func RegisterSession(s ISession) {
	session = s
}
