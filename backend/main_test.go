package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestAttachmentMetadataRouteIsRemoved guards against reintroducing the E69
// non-consuming password oracle: GET /api/secret/attachment must not exist.
func TestAttachmentMetadataRouteIsRemoved(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupRouter()

	for _, route := range router.Routes() {
		if route.Path == "/api/secret/attachment" {
			t.Fatalf("insecure route %q must not be registered", route.Path)
		}
	}

	request := httptest.NewRequest(http.MethodGet, "/api/secret/attachment?key=k&password=p", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("GET /api/secret/attachment status = %d, want %d", response.Code, http.StatusNotFound)
	}
}

// TestSecretRetrievalRouteStillRegistered ensures the legitimate one-time
// retrieval endpoint is unaffected by removing the attachment-metadata route.
func TestSecretRetrievalRouteStillRegistered(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupRouter()

	found := false
	for _, route := range router.Routes() {
		if route.Method == http.MethodPost && route.Path == "/api/secret/get" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("POST /api/secret/get must remain registered")
	}
}
