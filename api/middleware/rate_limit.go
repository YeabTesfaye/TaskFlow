package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"
)

// RateLimiter represents a token bucket rate limiter
type RateLimiter struct {
	tokens     float64
	capacity   float64
	rate       float64
	lastUpdate time.Time
	mutex      sync.Mutex
}

// ClientRateLimiter stores rate limiters for each client
type ClientRateLimiter struct {
	limiters map[string]*RateLimiter
	mutex    sync.RWMutex
}

// NewClientRateLimiter creates a new client rate limiter
func NewClientRateLimiter() *ClientRateLimiter {
	return &ClientRateLimiter{
		limiters: make(map[string]*RateLimiter),
	}
}

// GetLimiter gets or creates a rate limiter for a client
func (c *ClientRateLimiter) GetLimiter(clientIP string) *RateLimiter {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	limiter, exists := c.limiters[clientIP]
	if !exists {
		limiter = &RateLimiter{
			tokens:     30,  // Initial tokens
			capacity:   30,  // Maximum tokens
			rate:       0.5, // Tokens per second
			lastUpdate: time.Now(),
		}
		c.limiters[clientIP] = limiter
	}
	return limiter
}

// Allow checks if a request should be allowed and returns remaining tokens
func (rl *RateLimiter) Allow() (bool, float64) {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastUpdate).Seconds()
	rl.lastUpdate = now

	// Add new tokens based on elapsed time
	rl.tokens += elapsed * rl.rate
	if rl.tokens > rl.capacity {
		rl.tokens = rl.capacity
	}

	if rl.tokens < 1 {
		return false, rl.tokens
	}

	rl.tokens--
	return true, rl.tokens
}

// getClientIP gets the real client IP from various headers
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		return ip
	}

	// Check X-Real-IP header
	ip = r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	// Fall back to RemoteAddr
	return r.RemoteAddr
}

// Global rate limiter instance
var globalLimiter = NewClientRateLimiter()

// RateLimit middleware function
func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := getClientIP(r)

		// Get rate limiter for this client
		limiter := globalLimiter.GetLimiter(clientIP)

		// Check if request is allowed
		allowed, remaining := limiter.Allow()

		// Add rate limit headers
		w.Header().Set("X-RateLimit-Limit", strconv.FormatFloat(limiter.capacity, 'f', 0, 64))
		w.Header().Set("X-RateLimit-Remaining", strconv.FormatFloat(remaining, 'f', 2, 64))
		w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(time.Second).Unix(), 10))

		if !allowed {
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("Rate limit exceeded. Please try again later."))
			return
		}

		next.ServeHTTP(w, r)
	})
}
