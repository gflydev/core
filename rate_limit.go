package core

import (
	"sync"
	"time"
)

// RateLimiter defines the interface for rate limiting implementations
type RateLimiter interface {
	// Allow checks if a request is allowed based on the key
	// Returns true if the request is allowed, false otherwise
	Allow(key string) bool

	// Reset resets the rate limiter for the given key
	Reset(key string)

	// GetLimit returns the current limit for the given key
	GetLimit(key string) int

	// GetRemaining returns the remaining tokens for the given key
	GetRemaining(key string) int

	// GetResetTime returns the time when the rate limit will reset for the given key
	GetResetTime(key string) time.Time
}

// TokenBucketLimiter implements the token bucket algorithm for rate limiting
type TokenBucketLimiter struct {
	// rate is the number of tokens added per second
	rate int

	// capacity is the maximum number of tokens in the bucket
	capacity int

	// tokens maps keys to their token bucket state
	tokens map[string]*tokenBucket

	// cleanupInterval is the interval at which expired buckets are cleaned up
	cleanupInterval time.Duration

	// expiration is the duration after which a bucket is considered expired
	expiration time.Duration

	// mutex for thread safety
	mu sync.RWMutex

	// stopCleanup is a channel to signal the cleanup goroutine to stop
	stopCleanup chan struct{}
}

// tokenBucket represents the state of a token bucket for a specific key
type tokenBucket struct {
	// tokens is the current number of tokens in the bucket
	tokens float64

	// lastRefill is the last time the bucket was refilled
	lastRefill time.Time

	// resetTime is the time when the bucket will be fully refilled
	resetTime time.Time
}

// NewTokenBucketLimiter creates a new TokenBucketLimiter
func NewTokenBucketLimiter(rate, capacity int, cleanupInterval, expiration time.Duration) *TokenBucketLimiter {
	limiter := &TokenBucketLimiter{
		rate:            rate,
		capacity:        capacity,
		tokens:          make(map[string]*tokenBucket),
		cleanupInterval: cleanupInterval,
		expiration:      expiration,
		stopCleanup:     make(chan struct{}),
	}

	// Start the cleanup goroutine
	go limiter.cleanup()

	return limiter
}

// Allow checks if a request is allowed based on the key
func (l *TokenBucketLimiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()

	// Get or create the token bucket for the key
	bucket, exists := l.tokens[key]
	if !exists {
		bucket = &tokenBucket{
			tokens:     float64(l.capacity),
			lastRefill: now,
			resetTime:  now.Add(time.Second),
		}
		l.tokens[key] = bucket
		return true
	}

	// Calculate the number of tokens to add based on the time elapsed since the last refill
	elapsed := now.Sub(bucket.lastRefill).Seconds()
	tokensToAdd := elapsed * float64(l.rate)

	// Refill the bucket
	bucket.tokens = bucket.tokens + tokensToAdd
	if bucket.tokens > float64(l.capacity) {
		bucket.tokens = float64(l.capacity)
	}

	// Update the last refill time
	bucket.lastRefill = now

	// Calculate the time until the bucket is fully refilled
	if bucket.tokens < float64(l.capacity) {
		timeToFull := time.Duration(float64(l.capacity-int(bucket.tokens)) / float64(l.rate) * float64(time.Second))
		bucket.resetTime = now.Add(timeToFull)
	} else {
		bucket.resetTime = now
	}

	// Check if there are enough tokens for this request
	if bucket.tokens >= 1 {
		bucket.tokens--
		return true
	}

	return false
}

// Reset resets the rate limiter for the given key
func (l *TokenBucketLimiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	delete(l.tokens, key)
}

// GetLimit returns the current limit for the given key
func (l *TokenBucketLimiter) GetLimit(key string) int {
	return l.capacity
}

// GetRemaining returns the remaining tokens for the given key
func (l *TokenBucketLimiter) GetRemaining(key string) int {
	l.mu.RLock()
	defer l.mu.RUnlock()

	bucket, exists := l.tokens[key]
	if !exists {
		return l.capacity
	}

	return int(bucket.tokens)
}

// GetResetTime returns the time when the rate limit will reset for the given key
func (l *TokenBucketLimiter) GetResetTime(key string) time.Time {
	l.mu.RLock()
	defer l.mu.RUnlock()

	bucket, exists := l.tokens[key]
	if !exists {
		return time.Now()
	}

	return bucket.resetTime
}

// cleanup periodically removes expired buckets
func (l *TokenBucketLimiter) cleanup() {
	ticker := time.NewTicker(l.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			l.removeExpiredBuckets()
		case <-l.stopCleanup:
			return
		}
	}
}

// removeExpiredBuckets removes buckets that haven't been used for a while
func (l *TokenBucketLimiter) removeExpiredBuckets() {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	for key, bucket := range l.tokens {
		if now.Sub(bucket.lastRefill) > l.expiration {
			delete(l.tokens, key)
		}
	}
}

// Stop stops the cleanup goroutine
func (l *TokenBucketLimiter) Stop() {
	close(l.stopCleanup)
}
