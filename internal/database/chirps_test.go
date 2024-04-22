package database

import (
	"fmt"
	"os"
	"testing"
)

func TestDBChirps(t *testing.T) {
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
		t.Fatalf("Error getting messages from DB:\n%v", err)
	}

	if len(chirps) > 0 {
		t.Errorf("The messages array should have zero length instead of %d", len(chirps))
	}

	chirp_author_id := 1
	chirp_text := "test chirp"
	chirp, err := db.CreateChirp(chirp_author_id, chirp_text)
	if err != nil {
		t.Errorf("Error creating chirp on DB:\n%v", err)
	}

	t.Logf("New chirp: ID = %d, body = '%s'", chirp.ID, chirp.Body)

	chirps, err = db.GetChirps()
	if err != nil {
		t.Fatalf("Error getting messages from DB:\n%v", err)
	}

	if len(chirps) != 1 {
		t.Fatalf("The messages DB should have length 1 instead of %d", len(chirps))
	}

	if chirps[0].Body != chirp_text {
		t.Errorf("Content in the DB ('%s') doesn't match the input text ('%s')", chirps[0].Body, chirp_text)
	}

	new_chirp_text := "new test chirp"
	chirp, err = db.CreateChirp(chirp_author_id, new_chirp_text)
	if err != nil {
		t.Errorf("Error creating message on DB:\n%v", err)
	}

	t.Logf("New chirp: ID = %d, body = '%s'", chirp.ID, chirp.Body)

	chirps, err = db.GetChirps()
	if err != nil {
		t.Fatalf("Error getting messages from DB:\n%v", err)
	}

	if len(chirps) != 2 {
		t.Fatalf("The messages DB should have length 2 instead of %d", len(chirps))
	}

	if chirps[1].Body != new_chirp_text {
		t.Errorf("Content in the DB ('%s') doesn't match the input text ('%s')", chirps[1].Body, new_chirp_text)
	}

	deleted_chirp_id := 1
	wrong_chirp_author_id := 2
	err = db.DeleteChirpById(deleted_chirp_id, wrong_chirp_author_id)
	if err != ErrChirpAuthorInvalid {
		t.Fatalf("Error deleting message with ID '%d' and wrong author ID '%d' from DB:\n%v", deleted_chirp_id, wrong_chirp_author_id, err)
	}

	err = db.DeleteChirpById(deleted_chirp_id, chirp_author_id)
	if err != nil {
		t.Fatalf("Error deleting message with ID '%d' from DB:\n%v", deleted_chirp_id, err)
	}

	chirps, err = db.GetChirps()
	if err != nil {
		t.Fatalf("Error getting messages from DB:\n%v", err)
	}

	if len(chirps) != 1 {
		t.Fatalf("The messages DB should have length 1 instead of %d", len(chirps))
	}

	if chirps[0].ID == deleted_chirp_id {
		t.Errorf("ID in the DB ('%d') shouldn't match the deleted ID ('%d')", chirps[0].ID, deleted_chirp_id)
	}

	chirp_new_author_id := 2
	new_author_chirp_text := "new author test chirp"
	chirp, err = db.CreateChirp(chirp_new_author_id, new_author_chirp_text)
	if err != nil {
		t.Errorf("Error creating message on DB:\n%v", err)
	}

	t.Logf("New chirp: ID = %d, body = '%s'", chirp.ID, chirp.Body)

	chirps, err = db.GetChirps()
	if err != nil {
		t.Fatalf("Error getting messages from DB:\n%v", err)
	}
	fmt.Println(chirps)

	if len(chirps) != 2 {
		t.Fatalf("The messages DB should have length 2 instead of %d", len(chirps))
	}

	if chirps[1].AuthorID != chirp_new_author_id {
		t.Errorf("Content in the DB ('%d') doesn't match the input author ('%d')", chirps[1].AuthorID, chirp_new_author_id)
	}

	if chirps[1].Body != new_author_chirp_text {
		t.Errorf("Content in the DB ('%s') doesn't match the input text ('%s')", chirps[1].Body, new_author_chirp_text)
	}

}
