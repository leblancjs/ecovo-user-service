package user

import (
	"context"
	"fmt"
	"time"

	"azure.com/ecovo/user-service/pkg/entity"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
)

// A MongoRepository is a repository that performs CRUD operations on users in
// a MongoDB collection.
type MongoRepository struct {
	collection *mongo.Collection
}

type document struct {
	ID           primitive.ObjectID  `bson:"_id,omitempty"`
	SubID        string              `bson:"subId"`
	Email        string              `bson:"email"`
	FirstName    string              `bson:"firstName"`
	LastName     string              `bson:"lastName"`
	DateOfBirth  time.Time           `bson:"dateOfBirth"`
	PhoneNumber  string              `bson:"phoneNumber"`
	Gender       string              `bson:"gender"`
	Photo        string              `bson:"photo"`
	Description  string              `bson:"description"`
	Preferences  *entity.Preferences `bson:"preferences"`
	SignUpPhase  string              `bson:"signUpPhase"`
	UserRating   int                 `json:"userRating" bson:"userRating"`
	DriverRating int                 `json:"driverRating" bson:"driverRating"`
}

func newDocumentFromEntity(u *entity.User) (*document, error) {
	if u == nil {
		return nil, fmt.Errorf("user.MongoRepository: entity is nil")
	}

	var id primitive.ObjectID
	if u.ID.IsZero() {
		id = primitive.NilObjectID
	} else {
		objectID, err := primitive.ObjectIDFromHex(u.ID.Hex())
		if err != nil {
			return nil, fmt.Errorf("user.MongoRepository: failed to create object")
		}

		id = objectID
	}

	return &document{
		id,
		u.SubID,
		u.Email,
		u.FirstName,
		u.LastName,
		u.DateOfBirth,
		u.PhoneNumber,
		u.Gender,
		u.Photo,
		u.Description,
		u.Preferences,
		u.SignUpPhase,
		u.UserRating,
		u.DriverRating,
	}, nil
}

func (d document) Entity() *entity.User {
	return &entity.User{
		entity.NewIDFromHex(d.ID.Hex()),
		d.SubID,
		d.Email,
		d.FirstName,
		d.LastName,
		d.DateOfBirth,
		d.PhoneNumber,
		d.Gender,
		d.Photo,
		d.Description,
		d.Preferences,
		d.SignUpPhase,
		d.UserRating,
		d.DriverRating,
	}
}

// NewMongoRepository creates a user repository for a MongoDB collection.
func NewMongoRepository(collection *mongo.Collection) (Repository, error) {
	if collection == nil {
		return nil, fmt.Errorf("user.MongoRepository: collection is nil")
	}

	return &MongoRepository{collection}, nil
}

// FindByID retrieves the user with the given ID, if it exists.
func (r *MongoRepository) FindByID(ID entity.ID) (*entity.User, error) {
	objectID, err := primitive.ObjectIDFromHex(string(ID))
	if err != nil {
		return nil, fmt.Errorf("user.MongoRepository: failed to create object ID")
	}

	filter := bson.D{{"_id", objectID}}
	var d document
	err = r.collection.FindOne(context.TODO(), filter).Decode(&d)
	if err != nil {
		return nil, fmt.Errorf("user.MongoRepository: no user found with ID \"%s\" (%s)", ID, err)
	}
	return d.Entity(), nil
}

// FindBySubID retrieves the user with the given subscription ID, if it exists.
func (r *MongoRepository) FindBySubID(subID string) (*entity.User, error) {
	filter := bson.D{{"subId", subID}}
	var d document
	err := r.collection.FindOne(context.TODO(), filter).Decode(&d)
	if err != nil {
		return nil, fmt.Errorf("user.MongoRepository: no user found with subscription ID \"%s\" (%s)", subID, err)
	}
	return d.Entity(), nil
}

// Create stores the new user in the database and returns the unique
// identifier that was generated for it.
func (r *MongoRepository) Create(u *entity.User) (entity.ID, error) {
	if u == nil {
		return entity.NilID, fmt.Errorf("user.MongoRepository: failed to create user (user is nil)")
	}

	d, err := newDocumentFromEntity(u)
	if err != nil {
		return entity.NilID, fmt.Errorf("user.MongoRepository: failed to create user document from entity (%s)", err)
	}

	res, err := r.collection.InsertOne(context.TODO(), d)
	if err != nil {
		return entity.NilID, fmt.Errorf("user.MongoRepository: failed to create user (%s)", err)
	}

	ID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return entity.NilID, fmt.Errorf("user.MongoRepository: failed to get ID of created user (%s)", err)
	}

	return entity.ID(ID.Hex()), nil
}

// Update updates the user in the database.
func (r *MongoRepository) Update(u *entity.User) error {
	d, err := newDocumentFromEntity(u)
	if err != nil {
		return fmt.Errorf("user.MongoRepository: failed to create user document from entity (%s)", err)
	}

	filter := bson.D{{"_id", d.ID}}
	update := bson.D{
		bson.E{"$set", d},
	}
	res, err := r.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("user.MongoRepository: failed to update user with ID \"%s\" (%s)", u.ID, err)
	}

	if res.MatchedCount <= 0 {
		return fmt.Errorf("user.MongoRepository: no matching user was found")
	}

	return nil
}

// Delete removes the user with the given ID from the database.
func (r *MongoRepository) Delete(ID entity.ID) error {
	objectID, err := primitive.ObjectIDFromHex(string(ID))
	if err != nil {
		return fmt.Errorf("user.MongoRepository: failed to create object ID")
	}

	filter := bson.D{{"_id", objectID}}
	_, err = r.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("user.MongoRepository: failed to delete user with ID \"%s\" (%s)", ID, err)
	}

	return nil
}
