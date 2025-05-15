package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestRequestIDMiddleware(t *testing.T) {
	tests := map[string]struct {
		requestIDHeader string
		expectHeader    bool
	}{
		"No request ID header": {
			requestIDHeader: "",
			expectHeader:    true,
		},
		"With request ID header": {
			requestIDHeader: "test-request-id",
			expectHeader:    true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Create a new fasthttp request context
			ctx := &fasthttp.RequestCtx{}

			// Set request ID header if provided
			if tt.requestIDHeader != "" {
				ctx.Request.Header.Set(RequestIDHeader, tt.requestIDHeader)
			}

			// Create a new gFly context
			gflyCtx := &Ctx{
				root: ctx,
				data: make(Data),
			}

			// Call the middleware
			err := RequestIDMiddleware(gflyCtx)

			// Check that no error was returned
			assert.NoError(t, err)

			// Get the request ID from the context
			requestID := GetRequestID(gflyCtx)

			// Check that a request ID was set
			assert.NotEmpty(t, requestID)

			// If a request ID header was provided, check that it was used
			if tt.requestIDHeader != "" {
				assert.Equal(t, tt.requestIDHeader, requestID)
			}

			// Check that the response header was set
			if tt.expectHeader {
				responseHeader := string(ctx.Response.Header.Peek(RequestIDHeader))
				assert.Equal(t, requestID, responseHeader)
			}
		})
	}
}

func TestRequestIDLogMiddleware(t *testing.T) {
	// Create a new fasthttp request context
	ctx := &fasthttp.RequestCtx{}

	// Create a new gFly context
	gflyCtx := &Ctx{
		root: ctx,
		data: make(Data),
	}

	// Set a request ID in the context
	requestID := "test-request-id"
	gflyCtx.SetData(RequestIDContextKey, requestID)

	// Call the middleware
	err := RequestIDLogMiddleware(gflyCtx)

	// Check that no error was returned
	assert.NoError(t, err)

	// Get the logger from the context
	logger := GetLogger(gflyCtx)

	// Check that a logger was set
	assert.NotNil(t, logger)
}

func TestGetRequestID(t *testing.T) {
	tests := map[string]struct {
		requestID string
		expected  string
	}{
		"With request ID": {
			requestID: "test-request-id",
			expected:  "test-request-id",
		},
		"Without request ID": {
			requestID: "",
			expected:  "",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Create a new gFly context
			gflyCtx := &Ctx{
				data: make(Data),
			}

			// Set a request ID in the context if provided
			if tt.requestID != "" {
				gflyCtx.SetData(RequestIDContextKey, tt.requestID)
			}

			// Get the request ID from the context
			requestID := GetRequestID(gflyCtx)

			// Check that the correct request ID was returned
			assert.Equal(t, tt.expected, requestID)
		})
	}
}

func TestGetLogger(t *testing.T) {
	tests := map[string]struct {
		setLogger bool
	}{
		"With logger": {
			setLogger: true,
		},
		"Without logger": {
			setLogger: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Create a new gFly context
			gflyCtx := &Ctx{
				data: make(Data),
			}

			// Set a request ID and call the middleware if needed
			if tt.setLogger {
				requestID := "test-request-id"
				gflyCtx.SetData(RequestIDContextKey, requestID)
				_ = RequestIDLogMiddleware(gflyCtx)
			}

			// Get the logger from the context
			logger := GetLogger(gflyCtx)

			// Check that a logger was returned
			assert.NotNil(t, logger)
		})
	}
}
