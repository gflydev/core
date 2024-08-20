package core

// ==================================================================================
// 									Session Structure
// ==================================================================================

type ISession interface {
	Set(c *Ctx, key string, value interface{})
	Get(c *Ctx, key string) interface{}
}

// ==================================================================================
//                                   Default Session
// ==================================================================================

var sessionError = "Session manager is NULL. Please use import session module and run `session.Setup()`"

type DefaultSession struct {
}

func (v *DefaultSession) Set(c *Ctx, key string, value interface{}) {
	panic(sessionError)
}

func (v *DefaultSession) Get(c *Ctx, key string) interface{} {
	panic(sessionError)
}

var session ISession = &DefaultSession{}

// ==================================================================================
//                                   Functions
// ==================================================================================

// RegisterSession Get data from session.
func RegisterSession(s ISession) {
	session = s
}
