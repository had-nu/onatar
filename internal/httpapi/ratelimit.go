package httpapi

import (
	"net"
	"net/http"
	"sync"
	"time"
)

// tokenBucket is a per-key token bucket (PRD threat model §7: rate limit
// 10 req/min per IP on POST /build).
type tokenBucket struct {
	mu       sync.Mutex
	rate     float64 // tokens per second
	capacity float64
	buckets  map[string]*bucket
	stop     chan struct{}
	stopped  bool
}

type bucket struct {
	tokens float64
	last   time.Time
}

func newTokenBucket(ratePerMin, capacity float64) *tokenBucket {
	t := &tokenBucket{
		rate:     ratePerMin / 60,
		capacity: capacity,
		buckets:  make(map[string]*bucket),
		stop:     make(chan struct{}),
	}
	go t.sweep()
	return t
}

func (t *tokenBucket) allow(key string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	b, ok := t.buckets[key]
	if !ok {
		b = &bucket{tokens: t.capacity, last: now}
		t.buckets[key] = b
	}
	elapsed := now.Sub(b.last).Seconds()
	b.tokens = min(t.capacity, b.tokens+elapsed*t.rate)
	b.last = now
	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

func (t *tokenBucket) sweep() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-t.stop:
			return
		case <-ticker.C:
			t.mu.Lock()
			for k, b := range t.buckets {
				if time.Since(b.last) > 5*time.Minute {
					delete(t.buckets, k)
				}
			}
			t.mu.Unlock()
		}
	}
}

func (t *tokenBucket) close() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.stopped {
		return
	}
	t.stopped = true
	close(t.stop)
}

// rateLimit wraps a route with per-IP token bucket limiting.
func (s *Server) rateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !s.limiter.allow(clientIP(r)) {
			w.Header().Set("Retry-After", "60")
			writeError(w, http.StatusTooManyRequests, "RATE_LIMITED",
				"rate limit exceeded (10 req/min per IP)", nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
