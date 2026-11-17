package api

import (
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Tunable limits for the public endpoints. Kept as package constants so the
// guards stay dependency-free and behave identically in every environment.
const (
	// /api/secret/get: per-IP request ceiling.
	secretGetIPLimit  = 30
	secretGetIPWindow = time.Minute

	// /api/secret/get: per-extraction-key failed-attempt lockout. This is the
	// anti-brute-force control for the short extraction code: a wrong password
	// returns 404 without consuming the secret, so unlimited guessing must be
	// cut off explicitly.
	secretGetFailLimit  = 20
	secretGetFailWindow = 10 * time.Minute
	secretGetBlockFor   = 10 * time.Minute

	// /api/transfer/signal: per-IP WebSocket handshake ceiling.
	signalIPLimit  = 30
	signalIPWindow = time.Minute
)

// MaxRequestBodyBytes caps the size of a request body accepted by the API.
// Content is limited to 10,000 characters separately; the overhead of JSON
// framing plus UTF-8 expansion fits comfortably here.
const MaxRequestBodyBytes int64 = 128 * 1024

// windowLimiter is a fixed-window request counter keyed by an arbitrary string
// (typically the client IP). Safe for concurrent use.
type windowLimiter struct {
	mu        sync.Mutex
	entries   map[string]*windowEntry
	limit     int
	window    time.Duration
	lastSweep time.Time
}

type windowEntry struct {
	count int
	reset time.Time
}

func newWindowLimiter(limit int, window time.Duration) *windowLimiter {
	return &windowLimiter{
		entries:   make(map[string]*windowEntry),
		limit:     limit,
		window:    window,
		lastSweep: time.Now(),
	}
}

// allow reports whether a request from key may proceed, plus how long to wait
// when it may not.
func (l *windowLimiter) allow(key string) (bool, time.Duration) {
	now := time.Now()
	l.mu.Lock()
	defer l.mu.Unlock()
	l.sweep(now)

	if e, ok := l.entries[key]; ok && now.Before(e.reset) {
		if e.count >= l.limit {
			return false, time.Until(e.reset)
		}
		e.count++
		return true, 0
	}
	l.entries[key] = &windowEntry{count: 1, reset: now.Add(l.window)}
	return true, 0
}

// sweep drops expired entries so the map cannot grow without bound.
func (l *windowLimiter) sweep(now time.Time) {
	if now.Sub(l.lastSweep) < l.window {
		return
	}
	for k, e := range l.entries {
		if !now.Before(e.reset) {
			delete(l.entries, k)
		}
	}
	l.lastSweep = now
}

// failureLock blocks a specific extraction key after too many failed retrievals
// within a window, defeating online brute force of the short code.
type failureLock struct {
	mu        sync.Mutex
	entries   map[string]*failEntry
	limit     int
	window    time.Duration
	blockFor  time.Duration
	lastSweep time.Time
}

type failEntry struct {
	count        int
	windowStart  time.Time
	blockedUntil time.Time
}

func newFailureLock(limit int, window, blockFor time.Duration) *failureLock {
	return &failureLock{
		entries:   make(map[string]*failEntry),
		limit:     limit,
		window:    window,
		blockFor:  blockFor,
		lastSweep: time.Now(),
	}
}

func (f *failureLock) blocked(key string) (bool, time.Duration) {
	now := time.Now()
	f.mu.Lock()
	defer f.mu.Unlock()
	if e, ok := f.entries[key]; ok && now.Before(e.blockedUntil) {
		return true, time.Until(e.blockedUntil)
	}
	return false, 0
}

func (f *failureLock) recordFailure(key string) {
	now := time.Now()
	f.mu.Lock()
	defer f.mu.Unlock()
	f.sweep(now)

	e, ok := f.entries[key]
	if !ok || now.Sub(e.windowStart) > f.window {
		e = &failEntry{windowStart: now}
		f.entries[key] = e
	}
	e.count++
	if e.count >= f.limit {
		e.blockedUntil = now.Add(f.blockFor)
		e.count = 0
		e.windowStart = now
	}
}

func (f *failureLock) reset(key string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.entries, key)
}

func (f *failureLock) sweep(now time.Time) {
	if now.Sub(f.lastSweep) < f.window {
		return
	}
	for k, e := range f.entries {
		if now.After(e.blockedUntil) && now.Sub(e.windowStart) > f.window {
			delete(f.entries, k)
		}
	}
	f.lastSweep = now
}

func retryAfter(retry time.Duration) string {
	secs := int(math.Ceil(retry.Seconds()))
	if secs < 1 {
		secs = 1
	}
	return strconv.Itoa(secs)
}

// NewSecretGetGuard throttles the retrieval endpoint per client IP and locks a
// specific extraction key after repeated password failures. A successful
// retrieval (200) clears that key's failure history; a miss (404) counts.
func NewSecretGetGuard() gin.HandlerFunc {
	ipLimiter := newWindowLimiter(secretGetIPLimit, secretGetIPWindow)
	lock := newFailureLock(secretGetFailLimit, secretGetFailWindow, secretGetBlockFor)
	return func(c *gin.Context) {
		if allowed, retry := ipLimiter.allow(c.ClientIP()); !allowed {
			c.Header("Retry-After", retryAfter(retry))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, ErrorResponse{Error: "Too many requests, please slow down"})
			return
		}

		key := c.Query("key")
		if key != "" {
			if blocked, retry := lock.blocked(key); blocked {
				c.Header("Retry-After", retryAfter(retry))
				c.AbortWithStatusJSON(http.StatusTooManyRequests, ErrorResponse{Error: "Too many failed attempts for this code, please try again later"})
				return
			}
		}

		c.Next()

		if key != "" {
			switch c.Writer.Status() {
			case http.StatusNotFound:
				lock.recordFailure(key)
			case http.StatusOK:
				lock.reset(key)
			}
		}
	}
}

// NewSignalGuard throttles WebSocket signaling handshakes per client IP. It
// does not validate the password (the signal endpoint never reveals whether a
// secret exists), so only a connection-rate ceiling is applied here.
func NewSignalGuard() gin.HandlerFunc {
	ipLimiter := newWindowLimiter(signalIPLimit, signalIPWindow)
	return func(c *gin.Context) {
		if allowed, retry := ipLimiter.allow(c.ClientIP()); !allowed {
			c.Header("Retry-After", retryAfter(retry))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, ErrorResponse{Error: "Too many requests, please slow down"})
			return
		}
		c.Next()
	}
}

// NewSecurityHeaders sets hardening headers on API responses. The SPA document
// is served directly by nginx, which sets its own policy there.
func NewSecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.Writer.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "no-referrer")
		c.Next()
	}
}

// NewBodyLimit caps the request body size to prevent memory-abuse denial of
// service. An over-limit body surfaces as *http.MaxBytesError from the reader.
func NewBodyLimit(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Body != nil {
			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		}
		c.Next()
	}
}
