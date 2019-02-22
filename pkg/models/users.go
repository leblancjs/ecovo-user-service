package models

import (
	"fmt"
	"strings"
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

const (
	// SignUpPhasePersonalInfo means that the user is at the first sign up
	// phase where it needs to enter its personal information.
	SignUpPhasePersonalInfo string = "personalInfo"

	// SignUpPhasePreferences means that the user is at the second sign up
	// phase where it needs to enter its preferences.
	SignUpPhasePreferences string = "preferences"

	// SignUpPhaseDone means that the user has completed all sign up phases.
	SignUpPhaseDone string = "done"
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
	SignUpPhase string             `json:"signUpPhase" bson:"signUpPhase"`
}

// Validate looks at the user's contents to ensure it has all the required
// fields.
func (user *User) Validate() error {
	if user.Auth0ID == "" {
		return fmt.Errorf("missing Auth0ID")
	}

	if user.Email == "" {
		return fmt.Errorf("missing email")
	}

	if user.FirstName == "" {
		return fmt.Errorf("missing first name")
	}

	if user.LastName == "" {
		return fmt.Errorf("missing last name")
	}

	if user.DateOfBirth.IsZero() {
		return fmt.Errorf("missing date of birth")
	}

	if user.PhoneNumber == "" {
		return fmt.Errorf("missing phone number")
	}

	if user.Gender == "" {
		return fmt.Errorf("missing gender")
	}

	if user.Preferences == nil {
		return fmt.Errorf("missing preferences")
	}

	err := user.Preferences.validate()
	if err != nil {
		return err
	}

	if user.SignUpPhase == "" {
		return fmt.Errorf("missing sign up phase")
	}

	if strings.Compare(user.SignUpPhase, SignUpPhasePersonalInfo) != 0 && strings.Compare(user.SignUpPhase, SignUpPhasePreferences) != 0 && strings.Compare(user.SignUpPhase, SignUpPhaseDone) != 0 {
		return fmt.Errorf("invalid signup phase \"%s\"", user.SignUpPhase)
	}

	return nil
}
