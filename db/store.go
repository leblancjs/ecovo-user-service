package db

import (
	"context"
	"fmt"

	"azure.com/ecovo/user-service/models"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

// Store is an interface representing the ability to access data stored in a
// database.
type Store interface {
	FindUserByID(ID string) (*models.User, error)
	FindUserByAuth0ID(ID string) (*models.User, error)
	CreateUser(user *models.User) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(user *models.User) error
}

// FindUserByID looks for a user with the given ID in the database and returns
// it if it is found.
func (db *DB) FindUserByID(ID string) (*models.User, error) {
	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, fmt.Errorf("store: failed to parse object ID (%s)", err)
	}

	filter := bson.D{{"_id", objectID}}
	var user models.User
	err = db.users.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("store: no user found with ID \"%s\" (%s)", ID, err)
	}
	return &user, nil
}

// FindUserByAuth0ID looks for a user with the given Auth0ID in the database
// and returns it if it is found.
func (db *DB) FindUserByAuth0ID(ID string) (*models.User, error) {
	filter := bson.D{{"auth0ID", ID}}
	var user models.User
	err := db.users.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("store: no user found with Auth0ID \"%s\" (%s)", ID, err)
	}
	return &user, nil
}

// CreateUser creates a user in the database, populates the given user's ID
// field and returns a reference to it.
func (db *DB) CreateUser(user *models.User) (*models.User, error) {
	res, err := db.users.InsertOne(context.TODO(), user)
	if err != nil {
		return nil, fmt.Errorf("store: failed to create user (%s)", err)
	}

	ID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("store: failed to get ID of created user (%s)", err)
	}

	user.ID = ID

	return user, nil
}

// UpdateUser updates a user in the database based on the non-zero fields of
// the given user.
func (db *DB) UpdateUser(user *models.User) error {
	filter := bson.D{{"_id", user.ID}}
	update := bson.D{
		bson.E{"$set", user},
	}
	res, err := db.users.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("store: failed to update user (%s)", err)
	}

	if res.MatchedCount <= 0 {
		return fmt.Errorf("store: no matching user was found")
	}

	return nil
}

// DeleteUser removes a user from the database.
func (db *DB) DeleteUser(user *models.User) error {
	filter := bson.D{{"_id", user.ID}}
	_, err := db.users.DeleteOne(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("store: failed to delete user (%s)", err)
	}

	return nil
}
