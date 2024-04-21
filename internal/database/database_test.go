package database

import (
	"os"
	"testing"
)

func TestDB(t *testing.T) {
	fname := "database.json"

	err := os.Remove(fname)
	if err != nil {
		t.Logf("Error removing DB file %s:\n%v", fname, err)
	}

	_, err = NewDB(fname)
	if err != nil {
		t.Fatalf("Error creating DB with file %s:\n%v", fname, err)
	}

	if _, err := os.Stat(fname); os.IsNotExist(err) {
		t.Fatalf("DB file %s should exist:\n%v", fname, err)
	}
}
