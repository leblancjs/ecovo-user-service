package models

import (
	"errors"
	"fmt"
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

func (prefs *UserPreferences) validate() error {
	if prefs.Smoking < Never || prefs.Smoking > Regularly {
		return fmt.Errorf("smoking out of bounds (%d)", prefs.Smoking)
	}

	if prefs.Conversation < Never || prefs.Conversation > Regularly {
		return fmt.Errorf("conversation out of bounds (%d)", prefs.Conversation)
	}

	if prefs.Music < Never || prefs.Music > Regularly {
		return fmt.Errorf("music out of bounds (%d)", prefs.Music)
	}

	return nil
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

// Validate looks at the user's contents to ensure it has all the required
// fields.
func (user *User) Validate() error {
	if user.Auth0ID == "" {
		return errors.New("missing Auth0ID")
	}

	if user.Email == "" {
		return errors.New("missing email")
	}

	if user.FirstName == "" {
		return errors.New("missing first name")
	}

	if user.LastName == "" {
		return errors.New("missing last name")
	}

	if user.DateOfBirth.IsZero() {
		return errors.New("missing date of birth")
	}

	if user.PhoneNumber == "" {
		return errors.New("missing phone number")
	}

	if user.Gender == "" {
		return errors.New("missing gender")
	}

	if user.Preferences == nil {
		return errors.New("missing preferences")
	}

	err := user.Preferences.validate()
	if err != nil {
		return err
	}

	if user.SignUpPhase == nil {
		return errors.New("missing sign up phase")
	}

	return nil
}
