package core

import (
	"strconv"
	"time"

	"github.com/gflydev/core/errors"
)

// Default rate limit settings
const (
	// DefaultRateLimit is the default number of requests allowed per minute
	DefaultRateLimit = 60

	// DefaultRateLimitBurst is the default burst size
	DefaultRateLimitBurst = 10

	// DefaultRateLimitCleanupInterval is the default interval for cleaning up expired buckets
	DefaultRateLimitCleanupInterval = 5 * time.Minute

	// DefaultRateLimitExpiration is the default duration after which a bucket is considered expired
	DefaultRateLimitExpiration = 1 * time.Hour
)

// Rate limit header names
const (
	// RateLimitHeader is the header that contains the rate limit
	RateLimitHeader = "X-RateLimit-Limit"

	// RateLimitRemainingHeader is the header that contains the remaining requests
	RateLimitRemainingHeader = "X-RateLimit-Remaining"

	// RateLimitResetHeader is the header that contains the time when the rate limit will reset
	RateLimitResetHeader = "X-RateLimit-Reset"
)

// RateLimitKeyFunc is a function that extracts a key from the request context
type RateLimitKeyFunc func(ctx *Ctx) string

// DefaultRateLimitKeyFunc returns the client's IP address as the rate limit key
func DefaultRateLimitKeyFunc(ctx *Ctx) string {
	return ctx.root.RemoteIP().String()
}

// RateLimitOptions contains options for the rate limit middleware
type RateLimitOptions struct {
	// Limiter is the rate limiter to use
	Limiter RateLimiter

	// KeyFunc is the function used to extract a key from the request
	KeyFunc RateLimitKeyFunc

	// StatusCode is the HTTP status code to return when rate limit is exceeded
	StatusCode int

	// Message is the message to return when rate limit is exceeded
	Message string

	// SkipFailedRequests determines whether failed requests (non-2xx responses) should count towards the rate limit
	SkipFailedRequests bool

	// SkipSuccessfulRequests determines whether successful requests (2xx responses) should count towards the rate limit
	SkipSuccessfulRequests bool

	// Headers determines whether to add rate limit headers to the response
	Headers bool
}

// DefaultRateLimitOptions returns the default options for the rate limit middleware
func DefaultRateLimitOptions() RateLimitOptions {
	limiter := NewTokenBucketLimiter(
		DefaultRateLimit/60, // Convert from requests per minute to requests per second
		DefaultRateLimitBurst,
		DefaultRateLimitCleanupInterval,
		DefaultRateLimitExpiration,
	)

	return RateLimitOptions{
		Limiter:                limiter,
		KeyFunc:                DefaultRateLimitKeyFunc,
		StatusCode:             StatusTooManyRequests,
		Message:                "Rate limit exceeded. Please try again later.",
		SkipFailedRequests:     false,
		SkipSuccessfulRequests: false,
		Headers:                true,
	}
}

// RateLimitMiddleware creates a middleware that limits the number of requests based on a key
func RateLimitMiddleware(options ...func(*RateLimitOptions)) MiddlewareHandler {
	// Use default options
	opts := DefaultRateLimitOptions()

	// Apply custom options
	for _, option := range options {
		option(&opts)
	}

	return func(ctx *Ctx) error {
		// Extract the key from the request
		key := opts.KeyFunc(ctx)

		// Check if the request is allowed
		allowed := opts.Limiter.Allow(key)

		// Add rate limit headers if enabled
		if opts.Headers {
			limit := opts.Limiter.GetLimit(key)
			remaining := opts.Limiter.GetRemaining(key)
			reset := opts.Limiter.GetResetTime(key)

			ctx.SetHeader(RateLimitHeader, strconv.Itoa(limit))
			ctx.SetHeader(RateLimitRemainingHeader, strconv.Itoa(remaining))
			ctx.SetHeader(RateLimitResetHeader, strconv.FormatInt(reset.Unix(), 10))
		}

		// If not allowed, return a 429 Too Many Requests response
		if !allowed {
			return errors.Newf(
				errors.CodeForbidden,
				"%s (key: %s)",
				opts.Message,
				key,
			)
		}

		return nil
	}
}

// WithRateLimiter sets a custom rate limiter
func WithRateLimiter(limiter RateLimiter) func(*RateLimitOptions) {
	return func(o *RateLimitOptions) {
		o.Limiter = limiter
	}
}

// WithRateLimitKeyFunc sets a custom key function
func WithRateLimitKeyFunc(keyFunc RateLimitKeyFunc) func(*RateLimitOptions) {
	return func(o *RateLimitOptions) {
		o.KeyFunc = keyFunc
	}
}

// WithRateLimitStatusCode sets a custom status code
func WithRateLimitStatusCode(statusCode int) func(*RateLimitOptions) {
	return func(o *RateLimitOptions) {
		o.StatusCode = statusCode
	}
}

// WithRateLimitMessage sets a custom message
func WithRateLimitMessage(message string) func(*RateLimitOptions) {
	return func(o *RateLimitOptions) {
		o.Message = message
	}
}

// WithRateLimitSkipFailedRequests sets whether to skip failed requests
func WithRateLimitSkipFailedRequests(skip bool) func(*RateLimitOptions) {
	return func(o *RateLimitOptions) {
		o.SkipFailedRequests = skip
	}
}

// WithRateLimitSkipSuccessfulRequests sets whether to skip successful requests
func WithRateLimitSkipSuccessfulRequests(skip bool) func(*RateLimitOptions) {
	return func(o *RateLimitOptions) {
		o.SkipSuccessfulRequests = skip
	}
}

// WithRateLimitHeaders sets whether to add rate limit headers
func WithRateLimitHeaders(headers bool) func(*RateLimitOptions) {
	return func(o *RateLimitOptions) {
		o.Headers = headers
	}
}

// WithTokenBucketLimiter creates a token bucket limiter with the specified parameters
func WithTokenBucketLimiter(rate, burst int, cleanupInterval, expiration time.Duration) func(*RateLimitOptions) {
	return func(o *RateLimitOptions) {
		o.Limiter = NewTokenBucketLimiter(rate, burst, cleanupInterval, expiration)
	}
}

// WithRequestsPerMinute sets the rate limiter to allow the specified number of requests per minute
func WithRequestsPerMinute(requestsPerMinute int) func(*RateLimitOptions) {
	return func(o *RateLimitOptions) {
		// Convert from requests per minute to requests per second
		rate := requestsPerMinute / 60
		if rate < 1 {
			rate = 1
		}

		// Use the same burst size as the rate by default
		burst := requestsPerMinute / 10
		if burst < 1 {
			burst = 1
		}

		o.Limiter = NewTokenBucketLimiter(
			rate,
			burst,
			DefaultRateLimitCleanupInterval,
			DefaultRateLimitExpiration,
		)
	}
}

// WithRequestsPerSecond sets the rate limiter to allow the specified number of requests per second
func WithRequestsPerSecond(requestsPerSecond int) func(*RateLimitOptions) {
	return func(o *RateLimitOptions) {
		// Use the same burst size as the rate by default
		burst := requestsPerSecond
		if burst < 1 {
			burst = 1
		}

		o.Limiter = NewTokenBucketLimiter(
			requestsPerSecond,
			burst,
			DefaultRateLimitCleanupInterval,
			DefaultRateLimitExpiration,
		)
	}
}
