package database

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	chirp := Chirp{}

	dbStructure, err := db.loadDB()
	if err != nil {
		return chirp, err
	}

	id := len(dbStructure.Chirps) + 1
	chirp.ID = id
	chirp.Body = body

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
