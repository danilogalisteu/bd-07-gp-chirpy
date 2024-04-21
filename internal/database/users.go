package database

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

// CreateUser creates a new user and saves it to disk
func (db *DB) CreateUser(email string) (User, error) {
	user := User{}

	dbStructure, err := db.loadDB()
	if err != nil {
		return user, err
	}

	id := len(dbStructure.Users) + 1
	user.ID = id
	user.Email = email

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
