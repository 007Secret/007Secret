package api

import (
	"regexp"
	"testing"
)

func TestGenerateAlphanumericString(t *testing.T) {
	validCode := regexp.MustCompile(`^[A-Za-z0-9]{4}$`)

	for range 1_000 {
		code, err := generateAlphanumericString(4)
		if err != nil {
			t.Fatalf("generateAlphanumericString returned an error: %v", err)
		}
		if !validCode.MatchString(code) {
			t.Fatalf("generated extraction code contains invalid characters: %q", code)
		}
	}
}
