package middleware

import (
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

// We use a map to store a rate limiter for each IP address.
// The Mutex prevents race conditions when multiple requests arrive simultaneously.
var (
	limiters = make(map[string]*rate.Limiter)
	mu       sync.Mutex
)

func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := limiters[ip]
	if !exists {
		// STRICTER LIMIT FOR TESTING: 1 token per second, max burst of 2
		limiter = rate.NewLimiter(1, 2)
		limiters[ip] = limiter
	}
	return limiter
}

func RateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Note: If you deploy behind an AWS load balancer or Cloudflare, 
		// you'd want to check the "X-Forwarded-For" header instead of RemoteAddr.
		ip := r.RemoteAddr 

		limiter := getLimiter(ip)

		// Check if the IP has exceeded its token bucket allowance
		if !limiter.Allow() {
			http.Error(w, "429 Too Many Requests - Slow down!", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}