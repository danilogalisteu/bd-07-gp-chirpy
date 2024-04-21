package database

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var ErrUserExists = errors.New("user already exists")
var ErrUserIdNotFound = errors.New("user id was not found")
var ErrUserEmailNotFound = errors.New("user email was not found")
var ErrUserInfoNotValid = errors.New("user info is not valid")

// CreateUser creates a new user and saves it to disk
func (db *DB) CreateUser(email string, password string) (User, error) {
	user := User{}

	dbStructure, err := db.loadDB()
	if err != nil {
		return user, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return user, err
	}

	id := len(dbStructure.Users) + 1
	user.ID = id
	user.Email = email
	user.Password = string(passwordHash)

	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)

	return user, err
}

// GetUsers returns all users in the database
func (db *DB) GetUsers() ([]User, error) {
	users := make([]User, 0)

	dbStructure, err := db.loadDB()
	if err != nil {
		return users, err
	}

	for _, user := range dbStructure.Users {
		users = append(users, user)
	}

	return users, nil
}

// ValidateUser checks that the combination of email and password hash is valid
func (db *DB) ValidateUser(email string, password string) (User, error) {
	user := User{}

	dbStructure, err := db.loadDB()
	if err != nil {
		return user, err
	}

	for _, user := range dbStructure.Users {
		if user.Email == email {
			err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
			if err == nil {
				return user, nil
			}
		}
	}

	return user, ErrUserInfoNotValid
}
