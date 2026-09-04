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

func getIP(r *http.Request) string {
	// Preferir la IP real enviada por Nginx.
	realIP := strings.TrimSpace(r.Header.Get("X-Real-IP"))

	if realIP != "" && net.ParseIP(realIP) != nil {
		return realIP
	}

	// Fallback: IP de la conexión directa.
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}

	return strings.TrimSpace(r.RemoteAddr)
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
