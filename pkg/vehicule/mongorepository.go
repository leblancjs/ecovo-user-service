package vehicule

import (
	"context"
	"fmt"

	"azure.com/ecovo/user-service/pkg/entity"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
)

// A MongoRepository is a repository that performs CRUD operations on vehicules in
// a MongoDB collection.
type MongoRepository struct {
	collection *mongo.Collection
}

type document struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	UserID      primitive.ObjectID `bson:"userId"`
	Year        int                `bson:"year"`
	Make        string             `bson:"make"`
	Model       string             `bson:"model"`
	Color       string             `bson:"color"`
	Photo       string             `bson:"photo"`
	Accessories []string           `bson:"accessories"`
}

func newDocumentFromEntity(v *entity.Vehicule) (*document, error) {
	if v == nil {
		return nil, fmt.Errorf("vehicule.MongoRepository: entity is nil")
	}

	var id primitive.ObjectID
	if v.ID.IsZero() {
		id = primitive.NilObjectID
	} else {
		objectID, err := primitive.ObjectIDFromHex(v.ID.Hex())
		if err != nil {
			return nil, fmt.Errorf("vehicule.MongoRepository: failed to create object")
		}

		id = objectID
	}

	var userID primitive.ObjectID
	if v.UserID.IsZero() {
		userID = primitive.NilObjectID
	} else {
		objectID, err := primitive.ObjectIDFromHex(v.UserID.Hex())
		if err != nil {
			return nil, fmt.Errorf("vehicule.MongoRepository: failed to create object")
		}

		userID = objectID
	}

	return &document{
		id,
		userID,
		v.Year,
		v.Make,
		v.Model,
		v.Color,
		v.Photo,
		v.Accessories,
	}, nil
}

func (d document) Entity() *entity.Vehicule {
	return &entity.Vehicule{
		entity.NewIDFromHex(d.ID.Hex()),
		entity.NewIDFromHex(d.UserID.Hex()),
		d.Year,
		d.Make,
		d.Model,
		d.Color,
		d.Photo,
		d.Accessories,
	}
}

// NewMongoRepository creates a vehicule repository for a MongoDB collection.
func NewMongoRepository(collection *mongo.Collection) (Repository, error) {
	if collection == nil {
		return nil, fmt.Errorf("vehicule.MongoRepository: collection is nil")
	}

	return &MongoRepository{collection}, nil
}

// FindByID retrieves the vehicule with the given ID, if it exists.
func (r *MongoRepository) FindByID(ID entity.ID) (*entity.Vehicule, error) {
	objectID, err := primitive.ObjectIDFromHex(string(ID))
	if err != nil {
		return nil, fmt.Errorf("vehicule.MongoRepository: failed to create object ID")
	}

	filter := bson.D{{"_id", objectID}}
	var d document
	err = r.collection.FindOne(context.TODO(), filter).Decode(&d)
	if err != nil {
		return nil, fmt.Errorf("vehicule.MongoRepository: no vehicule found with ID \"%s\" (%s)", ID, err)
	}
	return d.Entity(), nil
}

// FindBySubID retrieves the vehicule with the given subscription ID, if it exists.
func (r *MongoRepository) FindByUserID(userID entity.ID) ([]*entity.Vehicule, error) {
	objectID, err := primitive.ObjectIDFromHex(string(userID))
	findOptions := options.Find()
	filter := bson.D{{"userId", objectID}}
	cur, err := r.collection.Find(context.TODO(), filter, findOptions)

	if err != nil {
		return nil, fmt.Errorf("vehicule.MongoRepository: no vehicules found with user ID \"%s\" (%s)", userID, err)
	}

	var vehicules []*entity.Vehicule
	for cur.Next(context.TODO()) {
		var d document
		err := cur.Decode(&d)
		if err != nil {
			return nil, err
		}
		vehicules = append(vehicules, d.Entity())
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(context.TODO())

	return vehicules, nil
}

// Create stores the new vehicule in the database and returns the unique
// identifier that was generated for it.
func (r *MongoRepository) Create(v *entity.Vehicule) (entity.ID, error) {
	if v == nil {
		return entity.NilID, fmt.Errorf("vehicule.MongoRepository: failed to create vehicule (vehicule is nil)")
	}

	d, err := newDocumentFromEntity(v)
	if err != nil {
		return entity.NilID, fmt.Errorf("vehicule.MongoRepository: failed to create vehicule document from entity (%s)", err)
	}

	res, err := r.collection.InsertOne(context.TODO(), d)
	if err != nil {
		return entity.NilID, fmt.Errorf("vehicule.MongoRepository: failed to create vehicule (%s)", err)
	}

	ID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return entity.NilID, fmt.Errorf("vehicule.MongoRepository: failed to get ID of created vehicule (%s)", err)
	}

	return entity.ID(ID.Hex()), nil
}

// Delete removes the vehicule with the given ID from the database.
func (r *MongoRepository) Delete(ID entity.ID) error {
	objectID, err := primitive.ObjectIDFromHex(string(ID))
	if err != nil {
		return fmt.Errorf("vehicule.MongoRepository: failed to create object ID")
	}

	filter := bson.D{{"_id", objectID}}
	_, err = r.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("vehicule.MongoRepository: failed to delete vehicule with ID \"%s\" (%s)", ID, err)
	}

	return nil
}
