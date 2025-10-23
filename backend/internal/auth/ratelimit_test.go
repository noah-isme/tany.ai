package auth

import "testing"

func TestRateLimiterAllow(t *testing.T) {
	limiter := NewRateLimiter(1, 1, defaultLimiterTTL)
	if !limiter.Allow("user:1") {
		t.Fatal("expected first request to pass")
	}
	if limiter.Allow("user:1") {
		t.Fatal("expected second request to be rate limited")
	}
	if !limiter.Allow("user:2") {
		t.Fatal("expected different key to pass")
	}
}
