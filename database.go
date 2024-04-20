package main

import (
	"encoding/json"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	db := DB{path: path, mux: &sync.RWMutex{}}
	err := db.ensureDB()
	return &db, err
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

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	db.mux.Lock()
	defer db.mux.Unlock()

	f, err := os.OpenFile(db.path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbStructure := DBStructure{
		Chirps: make(map[int]Chirp, 0),
	}

	f, err := os.Open(db.path)
	if err != nil {
		return dbStructure, err
	}

	fi, err := f.Stat()
	if err != nil {
		return dbStructure, err
	}

	if fi.Size() > 0 {
		decoder := json.NewDecoder(f)
	
		if err := decoder.Decode(&dbStructure); err != nil {
			return dbStructure, err
		}
	}


	if err := f.Close(); err != nil {
		return dbStructure, err
	}

	return dbStructure, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	db.ensureDB()

	db.mux.Lock()
	defer db.mux.Unlock()

	f, err := os.OpenFile(db.path, os.O_RDWR, 0666)
	if err != nil {
		return err
	}

	data, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}
	return nil
}
