package auth

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const defaultLimiterTTL = 15 * time.Minute

// RateLimiter coordinates request limits for authentication endpoints.
type RateLimiter struct {
	mu      sync.Mutex
	limit   rate.Limit
	burst   int
	ttl     time.Duration
	clients map[string]*clientLimiter
}

type clientLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// NewRateLimiter constructs a rate limiter with the provided limits.
func NewRateLimiter(perMinute, burst int, ttl time.Duration) *RateLimiter {
	if perMinute <= 0 {
		panic("perMinute must be positive")
	}
	if burst <= 0 {
		panic("burst must be positive")
	}
	if ttl <= 0 {
		ttl = defaultLimiterTTL
	}
	rl := &RateLimiter{
		limit:   rate.Limit(float64(perMinute) / 60.0),
		burst:   burst,
		ttl:     ttl,
		clients: make(map[string]*clientLimiter),
	}
	return rl
}

// Allow reports whether a request associated with the given key may proceed.
func (r *RateLimiter) Allow(key string) bool {
	now := time.Now()

	r.mu.Lock()
	defer r.mu.Unlock()

	entry, ok := r.clients[key]
	if !ok {
		entry = &clientLimiter{
			limiter:  rate.NewLimiter(r.limit, r.burst),
			lastSeen: now,
		}
		r.clients[key] = entry
	}

	entry.lastSeen = now
	allowed := entry.limiter.Allow()
	r.cleanupLocked(now)
	return allowed
}

func (r *RateLimiter) cleanupLocked(now time.Time) {
	for key, entry := range r.clients {
		if now.Sub(entry.lastSeen) > r.ttl {
			delete(r.clients, key)
		}
	}
}
