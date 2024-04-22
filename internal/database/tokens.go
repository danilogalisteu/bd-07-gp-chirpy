package database

import (
	"errors"
	"time"
)

type Token struct {
	Revoked   bool      `json:"revoked"`
	RevokedAt time.Time `json:"revoked_at"`
}

var ErrTokenExists = errors.New("token already exists")
var ErrTokenNotFound = errors.New("token not found")

// CreateToken creates a new token and saves it to disk
func (db *DB) CreateToken(tokenString string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	_, ok := dbStructure.Tokens[tokenString]
	if ok {
		return ErrTokenExists
	}

	dbStructure.Tokens[tokenString] = Token{}

	err = db.writeDB(dbStructure)

	return err
}

// ValidateToken checks that the token is not revoked
func (db *DB) ValidateToken(tokenString string) (bool, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return false, err
	}

	token, ok := dbStructure.Tokens[tokenString]
	if !ok {
		return false, ErrTokenNotFound
	}

	return !token.Revoked, nil
}

// RevokeToken revokes the token
func (db *DB) RevokeToken(tokenString string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	token, ok := dbStructure.Tokens[tokenString]
	if !ok {
		return ErrTokenNotFound
	}

	token.Revoked = true
	token.RevokedAt = time.Now()

	dbStructure.Tokens[tokenString] = token

	err = db.writeDB(dbStructure)

	return err
}
