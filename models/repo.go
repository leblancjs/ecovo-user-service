package models

import (
	"errors"
	"strconv"
)

// Repository represents a data access object that can perform CRUD operations.
type Repository interface {
	FindByID(ID string) (*User, error)
	FindByAuth0ID(ID string) (*User, error)
	Create(User) (*User, error)
	Update(User) error
	Delete(User) error
}

var nextID int

// FindByID looks for a user with the given ID in the database and returns it
// if it is found.
func (db *DB) FindByID(ID string) (*User, error) {
	// TODO: Get user from a database
	for _, user := range db.Db {
		if ID == user.ID {
			return user, nil
		}
	}

	return nil, errors.New("no user found with ID \"" + ID + "\"")
}

// FindByAuth0ID looks for a user with the given Auth0ID in the database and
// returns it if it is found.
func (db *DB) FindByAuth0ID(ID string) (*User, error) {
	// TODO: Get user from a database
	for _, user := range db.Db {
		if ID == user.Auth0ID {
			return user, nil
		}
	}

	return nil, errors.New("no user found with Auth0ID \"" + ID + "\"")
}

// Create creates a user in the database, populates the given user's ID field
// and returns a reference to it.
func (db *DB) Create(user *User) (*User, error) {
	// TODO: Save user to a database
	user.ID = strconv.Itoa(nextID)
	nextID++

	db.Db[user.ID] = user

	return user, nil
}

// Update updates a user in the database based on the non-zero fields of the
// given user.
func (db *DB) Update(user *User) error {
	// TODO: Update user in a database
	return nil
}

// Delete removes a user from the database.
func (db *DB) Delete(user *User) error {
	// TODO: Delete user in a database
	delete(db.Db, user.ID)
	return nil
}
