package main

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

	db, err := NewDB(fname)
	if err != nil {
		t.Fatalf("Error creating DB with file %s:\n%v", fname, err)
	}

	chirps, err := db.GetChirps()
	if err != nil {
		t.Fatalf("Error getting items from DB:\n%v", err)
	}

	if len(chirps) > 0 {
		t.Errorf("The chirp DB should have zero length instead of %d", len(chirps))
	}

	chirp_text := "test chirp"
	chirp, err := db.CreateChirp(chirp_text)
	if err != nil {
		t.Errorf("Error creating chirp on DB:\n%v", err)
	}

	t.Logf("New chirp: ID = %d, body = '%s'", chirp.ID, chirp.Body)

	chirps, err = db.GetChirps()
	if err != nil {
		t.Fatalf("Error getting items from DB:\n%v", err)
	}

	if len(chirps) != 1 {
		t.Fatalf("The chirp DB should have length 1 instead of %d", len(chirps))
	}

	if chirps[0].Body != chirp_text {
		t.Errorf("Content in the DB ('%s') doesn't match the input text ('%s')", chirps[0].Body, chirp_text)
	}

	new_chirp_text := "new test chirp"
	chirp, err = db.CreateChirp(new_chirp_text)
	if err != nil {
		t.Errorf("Error creating chirp on DB:\n%v", err)
	}

	t.Logf("New chirp: ID = %d, body = '%s'", chirp.ID, chirp.Body)

	chirps, err = db.GetChirps()
	if err != nil {
		t.Fatalf("Error getting items from DB:\n%v", err)
	}

	if len(chirps) != 2 {
		t.Fatalf("The chirp DB should have length 2 instead of %d", len(chirps))
	}

	if chirps[1].Body != new_chirp_text {
		t.Errorf("Content in the DB ('%s') doesn't match the input text ('%s')", chirps[1].Body, new_chirp_text)
	}
}
