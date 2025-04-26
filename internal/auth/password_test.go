package auth

import (
	"testing"
)

func TestPassword(t *testing.T) {
	password := "abc123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf(`HashPassword(%q) returned an error: %v`, password, err)
	}

	err = CheckPasswordHash(hash, password)
	if err != nil {
		t.Fatalf(`CheckPasswordHash(%q, %q) returned an error: %v`, hash, password, err)
	}
}
