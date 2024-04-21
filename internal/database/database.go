package database

import (
	"encoding/json"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"userss"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	db := &DB{path: path, mux: &sync.RWMutex{}}
	err := db.ensureDB()
	return db, err
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	db.mux.Lock()
	defer db.mux.Unlock()

	f, err := os.OpenFile(db.path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbStructure := DBStructure{
		Chirps: make(map[int]Chirp, 0),
		Users:  make(map[int]User, 0),
	}

	f, err := os.Open(db.path)
	if err != nil {
		return dbStructure, err
	}
	defer f.Close()

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
	defer f.Close()

	data, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return nil
}
