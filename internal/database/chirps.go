package database

import "errors"

type Chirp struct {
	ID       int    `json:"id"`
	AuthorID int    `json:"author_id"`
	Body     string `json:"body"`
}

var ErrChirpIdNotFound = errors.New("chirp id was not found")
var ErrChirpAuthorInvalid = errors.New("chirp author invalid")

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(author_id int, body string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := 0
	for _, chirp := range dbStructure.Chirps {
		id = max(id, chirp.ID)
	}
	id += 1

	chirp := Chirp{
		ID:       id,
		AuthorID: author_id,
		Body:     body,
	}

	dbStructure.Chirps[id] = chirp

	err = db.writeDB(dbStructure)

	return chirp, err
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	chirps := make([]Chirp, 0)

	dbStructure, err := db.loadDB()
	if err != nil {
		return chirps, err
	}

	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

// DeleteChirpById removes a chirp from the database if the author is verified
func (db *DB) DeleteChirpById(id int, author_id int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	for map_id, chirp := range dbStructure.Chirps {
		if chirp.ID == id {
			if chirp.AuthorID == author_id {
				delete(dbStructure.Chirps, map_id)
				err = db.writeDB(dbStructure)
				return err
			}
			return ErrChirpAuthorInvalid
		}
	}

	return ErrChirpIdNotFound
}
