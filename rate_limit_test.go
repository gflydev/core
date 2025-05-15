package core

import (
	"testing"
	"time"

	"github.com/gflydev/core/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
)

func TestTokenBucketLimiter(t *testing.T) {
	t.Run("Allow requests within limit", func(t *testing.T) {
		// Create a limiter with 10 requests per second and a burst of 5
		limiter := NewTokenBucketLimiter(10, 5, time.Minute, time.Hour)
		defer limiter.Stop()

		// Key for testing
		key := "test-key"

		// First request should be allowed
		assert.True(t, limiter.Allow(key), "First request should be allowed")

		// Remaining tokens should be 4 (5 - 1)
		assert.Equal(t, 4, limiter.GetRemaining(key), "Should have 4 tokens remaining")

		// Make 4 more requests, all should be allowed
		for i := 0; i < 4; i++ {
			assert.True(t, limiter.Allow(key), "Request %d should be allowed", i+2)
		}

		// Remaining tokens should be 0
		assert.Equal(t, 0, limiter.GetRemaining(key), "Should have 0 tokens remaining")

		// Next request should be denied
		assert.False(t, limiter.Allow(key), "Request after limit should be denied")
	})

	t.Run("Refill tokens over time", func(t *testing.T) {
		// Create a limiter with 10 tokens per second and a burst of 1
		limiter := NewTokenBucketLimiter(10, 1, time.Minute, time.Hour)
		defer limiter.Stop()

		// Key for testing
		key := "test-key"

		// First request should be allowed
		assert.True(t, limiter.Allow(key), "First request should be allowed")

		// Remaining tokens should be 0
		assert.Equal(t, 0, limiter.GetRemaining(key), "Should have 0 tokens remaining")

		// Next request should be denied
		assert.False(t, limiter.Allow(key), "Request after limit should be denied")

		// Wait for 100ms, which should add 1 token (10 tokens per second * 0.1 seconds)
		time.Sleep(100 * time.Millisecond)

		// Next request should be allowed
		assert.True(t, limiter.Allow(key), "Request after refill should be allowed")
	})

	t.Run("Reset limiter", func(t *testing.T) {
		// Create a limiter with 10 requests per second and a burst of 1
		limiter := NewTokenBucketLimiter(10, 1, time.Minute, time.Hour)
		defer limiter.Stop()

		// Key for testing
		key := "test-key"

		// First request should be allowed
		assert.True(t, limiter.Allow(key), "First request should be allowed")

		// Remaining tokens should be 0
		assert.Equal(t, 0, limiter.GetRemaining(key), "Should have 0 tokens remaining")

		// Reset the limiter
		limiter.Reset(key)

		// Remaining tokens should be back to 1
		assert.Equal(t, 1, limiter.GetRemaining(key), "Should have 1 token after reset")

		// Next request should be allowed
		assert.True(t, limiter.Allow(key), "Request after reset should be allowed")
	})

	t.Run("Multiple keys", func(t *testing.T) {
		// Create a limiter with 10 requests per second and a burst of 1
		limiter := NewTokenBucketLimiter(10, 1, time.Minute, time.Hour)
		defer limiter.Stop()

		// Keys for testing
		key1 := "test-key-1"
		key2 := "test-key-2"

		// First request for key1 should be allowed
		assert.True(t, limiter.Allow(key1), "First request for key1 should be allowed")

		// First request for key2 should be allowed
		assert.True(t, limiter.Allow(key2), "First request for key2 should be allowed")

		// Next request for key1 should be denied
		assert.False(t, limiter.Allow(key1), "Second request for key1 should be denied")

		// Next request for key2 should be denied
		assert.False(t, limiter.Allow(key2), "Second request for key2 should be denied")
	})

	t.Run("Cleanup expired buckets", func(t *testing.T) {
		// Create a limiter with 10 requests per second, a burst of 1, and a short expiration
		limiter := NewTokenBucketLimiter(10, 1, 100*time.Millisecond, 200*time.Millisecond)
		defer limiter.Stop()

		// Key for testing
		key := "test-key"

		// First request should be allowed
		assert.True(t, limiter.Allow(key), "First request should be allowed")

		// Wait for the bucket to expire
		time.Sleep(300 * time.Millisecond)

		// The bucket should be removed by the cleanup goroutine
		// Next request should be allowed as if it's the first request
		assert.True(t, limiter.Allow(key), "Request after expiration should be allowed")
	})
}

func TestRateLimitMiddleware(t *testing.T) {
	t.Run("Basic rate limiting", func(t *testing.T) {
		// Create a mock context
		ctx := &Ctx{
			data: make(map[string]interface{}),
		}

		// Create a rate limit middleware with 2 requests per second
		middleware := RateLimitMiddleware(WithRequestsPerSecond(2))

		// First request should be allowed
		err := middleware(ctx)
		assert.NoError(t, err, "First request should be allowed")

		// Second request should be allowed
		err = middleware(ctx)
		assert.NoError(t, err, "Second request should be allowed")

		// Third request should be denied
		err = middleware(ctx)
		require.Error(t, err, "Third request should be denied")
		assert.Equal(t, errors.CodeForbidden, err.(errors.Error).Code(), "Error should have the correct code")
	})

	t.Run("Custom key function", func(t *testing.T) {
		// Create a mock context
		ctx := &Ctx{
			data: make(map[string]interface{}),
		}

		// Set user ID in context
		ctx.SetData("user_id", "123")

		// Create a custom key function that uses the user ID
		keyFunc := func(ctx *Ctx) string {
			userID := ctx.GetData("user_id")
			if userID == nil {
				return "anonymous"
			}
			return userID.(string)
		}

		// Create a rate limit middleware with the custom key function
		middleware := RateLimitMiddleware(
			WithRequestsPerSecond(1),
			WithRateLimitKeyFunc(keyFunc),
		)

		// First request should be allowed
		err := middleware(ctx)
		assert.NoError(t, err, "First request should be allowed")

		// Second request should be denied
		err = middleware(ctx)
		require.Error(t, err, "Second request should be denied")
		assert.Equal(t, errors.CodeForbidden, err.(errors.Error).Code(), "Error should have the correct code")

		// Change the user ID
		ctx.SetData("user_id", "456")

		// Request with different user ID should be allowed
		err = middleware(ctx)
		assert.NoError(t, err, "Request with different user ID should be allowed")
	})

	t.Run("Custom error message", func(t *testing.T) {
		// Create a mock context
		ctx := &Ctx{
			data: make(map[string]interface{}),
		}

		// Create a rate limit middleware with a custom error message
		customMessage := "Custom rate limit exceeded message"
		middleware := RateLimitMiddleware(
			WithRequestsPerSecond(1),
			WithRateLimitMessage(customMessage),
		)

		// First request should be allowed
		err := middleware(ctx)
		assert.NoError(t, err, "First request should be allowed")

		// Second request should be denied with custom message
		err = middleware(ctx)
		require.Error(t, err, "Second request should be denied")
		assert.Contains(t, err.Error(), customMessage, "Error should contain the custom message")
	})

	t.Run("Rate limit headers", func(t *testing.T) {
		// Create a mock context with response headers
		ctx := &Ctx{
			data: make(map[string]interface{}),
			root: &fasthttp.RequestCtx{},
		}

		// Create a rate limit middleware
		middleware := RateLimitMiddleware(WithRequestsPerSecond(2))

		// First request should be allowed and set headers
		err := middleware(ctx)
		assert.NoError(t, err, "First request should be allowed")

		// Check that headers were set
		limitHeader := string(ctx.root.Response.Header.Peek(RateLimitHeader))
		remainingHeader := string(ctx.root.Response.Header.Peek(RateLimitRemainingHeader))
		resetHeader := string(ctx.root.Response.Header.Peek(RateLimitResetHeader))

		assert.Equal(t, "2", limitHeader, "Limit header should be set correctly")
		assert.Equal(t, "1", remainingHeader, "Remaining header should be set correctly")
		assert.NotEmpty(t, resetHeader, "Reset header should be set")
	})
}
