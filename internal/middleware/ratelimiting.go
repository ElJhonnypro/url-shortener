package middleware

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type client struct {
	requests  int
	windowEnd time.Time
}

type RateLimiter struct {
	mu sync.Mutex

	clients map[string]client

	limit  int
	window time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		clients: make(map[string]client),
		limit:   limit,
		window:  window,
	}
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ip := getIP(r)

		if !rl.Allow(ip) {
			w.Header().Set("Retry-After", "60")
			http.Error(
				w,
				"rate limit exceeded",
				http.StatusTooManyRequests,
			)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) Allow(ip string) bool {
	now := time.Now()

	rl.mu.Lock()
	defer rl.mu.Unlock()

	c, exists := rl.clients[ip]

	if !exists || now.After(c.windowEnd) {
		rl.clients[ip] = client{
			requests:  1,
			windowEnd: now.Add(rl.window),
		}

		return true
	}

	if c.requests >= rl.limit {
		return false
	}

	c.requests++

	rl.clients[ip] = c

	return true
}

func getIP(r *http.Request) string {
	// Para desarrollo/local.
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}

	// Fallback por si RemoteAddr no contiene puerto.
	return strings.TrimSpace(r.RemoteAddr)
}
