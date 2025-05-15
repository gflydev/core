package core

import (
	"github.com/gflydev/core/utils"
)

// RequestIDHeader is the header key for the request ID
const RequestIDHeader = "X-Request-ID"

// RequestIDContextKey is the key used to store the request ID in the context
const RequestIDContextKey = "request_id"

// RequestIDMiddleware adds a unique request ID to each request.
// If the request already has a request ID header, it will use that value.
// Otherwise, it will generate a new unique ID.
//
// The request ID is stored in the context with the key "request_id" and
// is also added to the response headers as "X-Request-ID".
//
// Parameters:
//   - ctx (*Ctx): The request context.
//
// Returns:
//   - error: Always returns nil to continue the middleware chain.
func RequestIDMiddleware(ctx *Ctx) error {
	// Check if request already has a request ID
	requestID := string(ctx.root.Request.Header.Peek(RequestIDHeader))

	// If no request ID is present, generate a new one
	if requestID == "" {
		requestID = utils.Token("request")
	}

	// Store the request ID in the context
	ctx.SetData(RequestIDContextKey, requestID)

	// Add the request ID to the response headers
	ctx.SetHeader(RequestIDHeader, requestID)

	return nil
}

// GetRequestID retrieves the request ID from the context.
//
// Parameters:
//   - ctx (*Ctx): The request context.
//
// Returns:
//   - string: The request ID, or an empty string if not found.
func GetRequestID(ctx *Ctx) string {
	requestID := ctx.GetData(RequestIDContextKey)
	if requestID == nil {
		return ""
	}
	return requestID.(string)
}
