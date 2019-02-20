package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"
)

// DB represents a database session.
type DB struct {
	client *mongo.Client
	users  *mongo.Collection
}

// const connectionString = "mongodb+srv://ecovo_admin:<PASSWORD>@cluster0-tgosy.mongodb.net/test?retryWrites=true"

// UserCollectionName represents the name of the collection in the database
// that contains the user records.
const UserCollectionName = "users"

// NewDB establishes a connection to a database server and returns the database
// with the given name.
func NewDB(host string, username string, password string, name string) (*DB, error) {
	url := fmt.Sprintf("mongodb://%s:%s@%s", username, password, host)

	client, err := mongo.NewClient(url)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database(name)
	if db == nil {
		return nil, errors.New("could not find database with name \"" + name + "\"")
	}

	users := db.Collection("users")

	return &DB{client, users}, nil
}
