package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestWindowLimiterAllowsUpToLimit(t *testing.T) {
	limiter := newWindowLimiter(3, time.Minute)

	for i := 0; i < 3; i++ {
		if allowed, _ := limiter.allow("ip"); !allowed {
			t.Fatalf("request %d should be allowed", i+1)
		}
	}
	if allowed, retry := limiter.allow("ip"); allowed {
		t.Fatal("request over the limit should be blocked")
	} else if retry <= 0 {
		t.Fatalf("blocked request should report a positive retry duration, got %v", retry)
	}
	if allowed, _ := limiter.allow("other"); !allowed {
		t.Fatal("a different key should have an independent bucket")
	}
}

func TestFailureLockBlocksThenResetClears(t *testing.T) {
	lock := newFailureLock(3, time.Minute, time.Minute)
	const key = "abc12345"

	for i := 0; i < 3; i++ {
		lock.recordFailure(key)
	}
	if blocked, retry := lock.blocked(key); !blocked {
		t.Fatal("key should be blocked after reaching the failure limit")
	} else if retry <= 0 {
		t.Fatalf("blocked key should report a positive retry duration, got %v", retry)
	}

	lock.reset(key)
	if blocked, _ := lock.blocked(key); blocked {
		t.Fatal("reset should clear the failure lock")
	}
}

func TestSecretGetGuardLocksKeyAfterFailures(t *testing.T) {
	router := gin.New()
	router.GET("/api/secret/get", NewSecretGetGuard(), func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Secret not found or already accessed"})
	})

	var last *httptest.ResponseRecorder
	for i := 0; i < secretGetFailLimit+1; i++ {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/secret/get?key=bruteforce-key", nil)
		router.ServeHTTP(rec, req)
		last = rec
	}

	if last.Code != http.StatusTooManyRequests {
		t.Fatalf("request after %d failures = %d, want %d", secretGetFailLimit+1, last.Code, http.StatusTooManyRequests)
	}
	if body := last.Body.String(); !strings.Contains(body, "failed attempts") {
		t.Fatalf("expected lockout error message, got %q", body)
	}
}

func TestSecretGetGuardRateLimitsPerIP(t *testing.T) {
	router := gin.New()
	router.GET("/api/secret/get", NewSecretGetGuard(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})

	var last *httptest.ResponseRecorder
	for i := 0; i < secretGetIPLimit+1; i++ {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/secret/get", nil)
		router.ServeHTTP(rec, req)
		last = rec
	}

	if last.Code != http.StatusTooManyRequests {
		t.Fatalf("request after %d requests from one IP = %d, want %d", secretGetIPLimit+1, last.Code, http.StatusTooManyRequests)
	}
	if last.Header().Get("Retry-After") == "" {
		t.Fatal("throttled response should carry a Retry-After header")
	}
}

func TestSecurityHeaders(t *testing.T) {
	router := gin.New()
	router.Use(NewSecurityHeaders())
	router.GET("/x", func(c *gin.Context) { c.Status(http.StatusOK) })

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/x", nil))

	for header, want := range map[string]string{
		"X-Content-Type-Options": "nosniff",
		"X-Frame-Options":        "DENY",
		"Referrer-Policy":        "no-referrer",
	} {
		if got := rec.Header().Get(header); got != want {
			t.Errorf("%s = %q, want %q", header, got, want)
		}
	}
}

func TestBodyLimitRejectsOversizedBody(t *testing.T) {
	router := gin.New()
	router.POST("/x", NewBodyLimit(16), func(c *gin.Context) {
		if _, err := io.ReadAll(c.Request.Body); err != nil {
			c.Status(http.StatusRequestEntityTooLarge)
			return
		}
		c.Status(http.StatusOK)
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/x", strings.NewReader(strings.Repeat("a", 64)))
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("oversized body status = %d, want %d", rec.Code, http.StatusRequestEntityTooLarge)
	}
}

func TestCreateSecretHandlerRejectsTooLongContent(t *testing.T) {
	router := gin.New()
	router.POST("/api/secret", NewBodyLimit(MaxRequestBodyBytes), CreateSecretHandler)

	body := `{"content":"` + strings.Repeat("a", maxContentLength+1) + `"}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/secret", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("too-long content status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if !strings.Contains(rec.Body.String(), "Content too long") {
		t.Fatalf("expected content-length error, got %q", rec.Body.String())
	}
}

func TestCreateSecretHandlerRejectsOversizedBody(t *testing.T) {
	router := gin.New()
	router.POST("/api/secret", NewBodyLimit(64), CreateSecretHandler)

	body := `{"content":"` + strings.Repeat("a", 200) + `"}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/secret", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("oversized body status = %d, want %d", rec.Code, http.StatusRequestEntityTooLarge)
	}
}
