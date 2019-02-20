package db

import (
	"context"
	"errors"

	"azure.com/ecovo/user-service/models"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

type Store interface {
	FindUserByID(ID string) (*models.User, error)
	FindUserByAuth0ID(ID string) (*models.User, error)
	CreateUser(u *models.User) (*models.User, error)
	UpdateUser(u *models.User) error
	DeleteUser(u *models.User) error
}

var nextID int

// FindUserByID looks for a user with the given ID in the database and returns
// it if it is found.
func (db *DB) FindUserByID(ID string) (*models.User, error) {
	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		// TODO: Return more informative error message
		return nil, err
	}

	filter := bson.D{{"_id", objectID}}
	var user models.User
	err = db.users.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		// TODO: Return more informative error message
		return nil, err
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
		// TODO: Return more informative error message
		return nil, err
	}
	return &user, nil
}

// CreateUser creates a user in the database, populates the given user's ID
// field and returns a reference to it.
func (db *DB) CreateUser(u *models.User) (*models.User, error) {
	res, err := db.users.InsertOne(context.TODO(), u)
	if err != nil {
		// TODO: Return more informative error message
		return nil, err
	}

	ID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		// TODO: Return more informative error message
		return nil, err
	}

	u.ID = ID

	return u, nil
}

// UpdateUser updates a user in the database based on the non-zero fields of
// the given user.
func (db *DB) UpdateUser(u *models.User) error {
	filter := bson.D{{"_id", u.ID}}
	update := bson.D{
		bson.E{"$set", u}}
	res, err := db.users.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		// TODO: Return more informative error message
		return err
	}

	if res.MatchedCount <= 0 {
		// TODO: Return more informative error message
		return errors.New("did not find any matches for update")
	}

	if res.ModifiedCount <= 0 {
		// TODO: Return more informative error message
		return errors.New("did not modify any document")
	}

	return nil
}

// DeleteUser removes a user from the database.
func (db *DB) DeleteUser(u *models.User) error {
	filter := bson.D{{"_id", u.ID}}
	res, err := db.users.DeleteOne(context.TODO(), filter)
	if err != nil {
		// TODO: Return more informative error message
		return err
	}

	if res.DeletedCount <= 0 {
		// TODO: Return more informative error message
		return errors.New("did not delete any document")
	}

	return nil
}
