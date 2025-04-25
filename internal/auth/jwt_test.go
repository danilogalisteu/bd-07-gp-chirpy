package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	userID := uuid.New()
	secret := "aksjf qw e83947 5y3947987t 5(*&*90 7gq9-v8rhu)"
	expiresIn := 5 * time.Minute

	token, err := MakeJWT(userID, secret, expiresIn)
	if err != nil {
		t.Fatalf(`MakeJWT(%q, %q, %v) returned an error: %v`, userID, secret, expiresIn, err)
	}

	tokenID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf(`ValidateJWT(%q, %q) returned an error: %v`, token, secret, err)
	}
	if userID != tokenID {
		t.Errorf(`ValidateJWT(%q, %q) = %q, want %q`, token, secret, userID, "abc123")
	}
}
