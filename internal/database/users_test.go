package database

import (
	"os"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestDBUsers(t *testing.T) {
	fname := "database.json"

	err := os.Remove(fname)
	if err != nil {
		t.Logf("Error removing DB file %s:\n%v", fname, err)
	}

	db, err := NewDB(fname)
	if err != nil {
		t.Fatalf("Error creating DB with file %s:\n%v", fname, err)
	}

	users, err := db.GetUsers()
	if err != nil {
		t.Fatalf("Error getting users from DB:\n%v", err)
	}

	if len(users) > 0 {
		t.Errorf("The users array should have zero length instead of %d", len(users))
	}

	user_email := "user@example.com"
	user_password := "123456"
	user, err := db.CreateUser(user_email, user_password)
	if err != nil {
		t.Errorf("Error creating user on DB:\n%v", err)
	}

	t.Logf("New user: ID = %d, body = '%s'", user.ID, user.Email)

	users, err = db.GetUsers()
	if err != nil {
		t.Fatalf("Error getting items from DB:\n%v", err)
	}

	if len(users) != 1 {
		t.Fatalf("The users DB should have length 1 instead of %d", len(users))
	}

	if users[0].Email != user_email {
		t.Errorf("Email in the DB ('%s') doesn't match the input text ('%s')", users[0].Email, user_email)
	}

	if bcrypt.CompareHashAndPassword([]byte(users[0].Password), []byte(user_password)) != nil {
		t.Errorf("Password hash in the DB ('%s') doesn't match the input text ('%s')", users[0].Password, user_password)
	}

	new_user_email := "admin@example.com"
	new_user_password := "abcdef"
	user, err = db.CreateUser(new_user_email, new_user_password)
	if err != nil {
		t.Errorf("Error creating user on DB:\n%v", err)
	}

	t.Logf("New user: ID = %d, body = '%s'", user.ID, user.Email)

	users, err = db.GetUsers()
	if err != nil {
		t.Fatalf("Error getting users from DB:\n%v", err)
	}

	if len(users) != 2 {
		t.Fatalf("The user DB should have length 2 instead of %d", len(users))
	}

	if users[1].Email != new_user_email {
		t.Errorf("Email in the DB ('%s') doesn't match the input text ('%s')", users[1].Email, new_user_email)
	}

	if bcrypt.CompareHashAndPassword([]byte(users[1].Password), []byte(new_user_password)) != nil {
		t.Errorf("Password hash in the DB ('%s') doesn't match the input text ('%s')", users[1].Password, new_user_password)
	}
}