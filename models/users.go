package models

import "time"

const (
	// Never represents that a user never smokes, tolerates animals, listens to
	// music or talks
	Never = 0

	// Occasionally represents that a user occasionally smokes, tolerates
	// animals, listens to music or talks
	Occasionally = 1

	// Regularly represents that a user regularly smokes, tolerates animals,
	// listens to music or talks
	Regularly = 2
)

// UserPreferences contains a user's preferences when it comes to smoking,
// animals, conversation and music.
type UserPreferences struct {
	Smoking      int `json:"smoking"`
	Animals      int `json:"animals"`
	Conversation int `json:"conversation"`
	Music        int `json:"music"`
}

// User contains a user's profile.
type User struct {
	ID          string           `json:"id"`
	Auth0ID     string           `json:"-"`
	Email       string           `json:"email"`
	FirstName   string           `json:"firstName"`
	LastName    string           `json:"lastName"`
	DateOfBirth time.Time        `json:"dateOfBirth"`
	PhoneNumber string           `json:"phoneNumber"`
	Gender      string           `json:"gender"`
	Photo       string           `json:"photo"`
	Description string           `json:"description"`
	Preferences *UserPreferences `json:"preferences"`
	SignUpPhase *int             `json:"signUpPhase"`
}
