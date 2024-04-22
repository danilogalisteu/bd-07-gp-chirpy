package database

type Chirp struct {
	ID       int    `json:"id"`
	AuthorID int    `json:"author_id"`
	Body     string `json:"body"`
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(author_id int, body string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := len(dbStructure.Chirps) + 1

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

// DeleteChirp removes a chirp from the database
func (db *DB) DeleteChirp(id int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	for map_id, chirp := range dbStructure.Chirps {
		if chirp.ID == id {
			delete(dbStructure.Chirps, map_id)
			break
		}
	}

	err = db.writeDB(dbStructure)

	return err
}
