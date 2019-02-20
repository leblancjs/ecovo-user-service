package models

import (
	"time"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

const (
	// Never represents that a user never smokes, listens to music or talks.
	Never = 0

	// Occasionally represents that a user occasionally smokes, listens to
	// music or talks.
	Occasionally = 1

	// Regularly represents that a user regularly smokes, listens to music or
	// talks.
	Regularly = 2
)

// UserPreferences contains a user's preferences when it comes to smoking,
// conversation and music.
type UserPreferences struct {
	Smoking      int `json:"smoking" bson:"smoking"`
	Conversation int `json:"conversation" bson:"conversation"`
	Music        int `json:"music" bson:"music"`
}

// User contains a user's profile.
type User struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Auth0ID     string             `json:"-" bson:"auth0ID"`
	Email       string             `json:"email" bson:"email"`
	FirstName   string             `json:"firstName" bson:"firstName"`
	LastName    string             `json:"lastName" bson:"lastName"`
	DateOfBirth time.Time          `json:"dateOfBirth" bson:"dateOfBirth"`
	PhoneNumber string             `json:"phoneNumber" bson:"phoneNumber"`
	Gender      string             `json:"gender" bson:"gender"`
	Photo       string             `json:"photo" bson:"photo"`
	Description string             `json:"description" bson:"description"`
	Preferences *UserPreferences   `json:"preferences" bson:"preferences"`
	SignUpPhase *int               `json:"signUpPhase" bson:"signUpPhase"`
}
