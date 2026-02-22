package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter provides per-user rate limiting using a token bucket algorithm.
type RateLimiter struct {
	mu       sync.Mutex
	buckets  map[int64]*bucket
	rate     int           // tokens added per interval
	capacity int           // max tokens
	interval time.Duration // refill interval
}

type bucket struct {
	tokens   int
	lastFill time.Time
}

// NewRateLimiter creates a rate limiter.
// rate = requests allowed per interval, capacity = burst size.
func NewRateLimiter(rate, capacity int, interval time.Duration) *RateLimiter {
	rl := &RateLimiter{
		buckets:  make(map[int64]*bucket),
		rate:     rate,
		capacity: capacity,
		interval: interval,
	}
	// Periodically clean up stale buckets
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			rl.cleanup()
		}
	}()
	return rl
}

func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	cutoff := time.Now().Add(-10 * time.Minute)
	for k, b := range rl.buckets {
		if b.lastFill.Before(cutoff) {
			delete(rl.buckets, k)
		}
	}
}

func (rl *RateLimiter) allow(userID int64) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	b, ok := rl.buckets[userID]
	if !ok {
		b = &bucket{tokens: rl.capacity, lastFill: now}
		rl.buckets[userID] = b
	}

	// Refill tokens based on elapsed time
	elapsed := now.Sub(b.lastFill)
	refills := int(elapsed / rl.interval)
	if refills > 0 {
		b.tokens += refills * rl.rate
		if b.tokens > rl.capacity {
			b.tokens = rl.capacity
		}
		b.lastFill = b.lastFill.Add(time.Duration(refills) * rl.interval)
	}

	if b.tokens > 0 {
		b.tokens--
		return true
	}
	return false
}

// Middleware returns a Gin middleware that rate-limits authenticated users.
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt64("user_id")
		if userID == 0 {
			// Not authenticated – let other middleware handle it
			c.Next()
			return
		}
		if !rl.allow(userID) {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests. Please slow down."})
			c.Abort()
			return
		}
		c.Next()
	}
}

// IPRateLimiter provides per-IP rate limiting for unauthenticated endpoints.
type IPRateLimiter struct {
	mu       sync.Mutex
	buckets  map[string]*bucket
	rate     int
	capacity int
	interval time.Duration
}

// NewIPRateLimiter creates an IP-based rate limiter.
func NewIPRateLimiter(rate, capacity int, interval time.Duration) *IPRateLimiter {
	rl := &IPRateLimiter{
		buckets:  make(map[string]*bucket),
		rate:     rate,
		capacity: capacity,
		interval: interval,
	}
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			rl.mu.Lock()
			cutoff := time.Now().Add(-10 * time.Minute)
			for k, b := range rl.buckets {
				if b.lastFill.Before(cutoff) {
					delete(rl.buckets, k)
				}
			}
			rl.mu.Unlock()
		}
	}()
	return rl
}

func (rl *IPRateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	b, ok := rl.buckets[ip]
	if !ok {
		b = &bucket{tokens: rl.capacity, lastFill: now}
		rl.buckets[ip] = b
	}

	elapsed := now.Sub(b.lastFill)
	refills := int(elapsed / rl.interval)
	if refills > 0 {
		b.tokens += refills * rl.rate
		if b.tokens > rl.capacity {
			b.tokens = rl.capacity
		}
		b.lastFill = b.lastFill.Add(time.Duration(refills) * rl.interval)
	}

	if b.tokens > 0 {
		b.tokens--
		return true
	}
	return false
}

// Middleware returns a Gin middleware for IP rate limiting.
func (rl *IPRateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !rl.allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests. Please slow down."})
			c.Abort()
			return
		}
		c.Next()
	}
}
