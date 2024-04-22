package database

import (
	"os"
	"testing"
)

func TestDBTokens(t *testing.T) {
	fname := "database.json"

	err := os.Remove(fname)
	if err != nil {
		t.Logf("Error removing DB file %s:\n%v", fname, err)
	}

	db, err := NewDB(fname)
	if err != nil {
		t.Fatalf("Error creating DB with file %s:\n%v", fname, err)
	}

	tokenString := "1234567890"
	valid, err := db.ValidateToken(tokenString)
	if err != ErrTokenNotFound {
		t.Fatalf("Token should not exist in the DB:\n%v", err)
		if err != nil {
			t.Fatalf("Error validating token:\n%v", err)
		}
	}
	if (err == nil) && !valid {
		t.Fatalf("The token was revoked")
	}

	err = db.CreateToken(tokenString)
	if err != nil {
		t.Fatalf("Error creating token on DB:\n%v", err)
	}

	valid, err = db.ValidateToken(tokenString)
	if err == ErrTokenNotFound {
		t.Fatalf("Token should exist in the DB:\n%v", err)
	}
	if err != nil {
		t.Fatalf("Error validating token:\n%v", err)
	}
	if (err == nil) && !valid {
		t.Fatalf("The token was revoked")
	}

	err = db.RevokeToken(tokenString)
	if err != nil {
		t.Fatalf("Error revoking token:\n%v", err)
	}

	valid, err = db.ValidateToken(tokenString)
	if err == ErrTokenNotFound {
		t.Fatalf("Token should exist in the DB:\n%v", err)
	}
	if err != nil {
		t.Fatalf("Error validating token:\n%v", err)
	}
	if (err == nil) && valid {
		t.Fatalf("The token was not revoked")
	}
}
