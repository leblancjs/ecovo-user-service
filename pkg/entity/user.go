package entity

import (
	"fmt"
	"strings"
	"time"
)

// User contains a user's profile.
type User struct {
	ID          ID           `json:"id" bson:"_id,omitempty"`
	SubID       string       `json:"-" bson:"subId"`
	Email       string       `json:"email" bson:"email"`
	FirstName   string       `json:"firstName" bson:"firstName"`
	LastName    string       `json:"lastName" bson:"lastName"`
	DateOfBirth time.Time    `json:"dateOfBirth" bson:"dateOfBirth"`
	PhoneNumber string       `json:"phoneNumber" bson:"phoneNumber"`
	Gender      string       `json:"gender" bson:"gender"`
	Photo       string       `json:"photo" bson:"photo"`
	Description string       `json:"description" bson:"description"`
	Preferences *Preferences `json:"preferences" bson:"preferences"`
	SignUpPhase string       `json:"signUpPhase" bson:"signUpPhase"`
	UserRating	int			`json:"userRating" bson:"userRating"`
	DriverRating	int			`json:"driverRating" bson:"driverRating"`
}

const (
	// GenderMale represents that the user is of the male gender.
	GenderMale = "Male"

	// GenderFemale represents that the user is of the female gender.
	GenderFemale = "Female"

	// GenderOther represents that the user is of, or identifies with, a gender
	// other than male or female.
	GenderOther = "Other"

	// AgeMinimum represents the minimum age a user must have (18 years).
	AgeMinimum = 18 * 365.25 * 24 * time.Hour

	// SignUpPhasePersonalInfo means that the user is at the first sign up
	// phase where it needs to enter its personal information.
	SignUpPhasePersonalInfo = "personalInfo"

	// SignUpPhasePreferences means that the user is at the second sign up
	// phase where it needs to enter its preferences.
	SignUpPhasePreferences = "preferences"

	// SignUpPhaseDone means that the user has completed all sign up phases.
	SignUpPhaseDone = "done"
)

// Validate validates that the user's required fields are filled out correctly.
func (u *User) Validate() error {
	if u.SubID == "" {
		return ValidationError{"subscription ID is missing"}
	}

	if u.FirstName == "" {
		return ValidationError{"first name is missing"}
	}

	if u.LastName == "" {
		return ValidationError{"last name is missing"}
	}

	if u.DateOfBirth.IsZero() {
		return ValidationError{"date of birth is missing"}
	}

	if time.Since(u.DateOfBirth) < AgeMinimum {
		return ValidationError{"must be 18 years of age or older"}
	}

	if u.Gender == "" {
		return ValidationError{"gender is missing"}
	}

	if strings.Compare(u.Gender, GenderMale) != 0 &&
		strings.Compare(u.Gender, GenderFemale) != 0 &&
		strings.Compare(u.Gender, GenderOther) != 0 {
		return ValidationError{fmt.Sprintf("gender must be %s, %s or %s", GenderMale, GenderFemale, GenderOther)}
	}

	if u.Preferences != nil {
		err := u.Preferences.Validate()
		if err != nil {
			return err
		}
	}

	if strings.Compare(u.SignUpPhase, SignUpPhasePersonalInfo) != 0 &&
		strings.Compare(u.SignUpPhase, SignUpPhasePreferences) != 0 &&
		strings.Compare(u.SignUpPhase, SignUpPhaseDone) != 0 {
		return ValidationError{fmt.Sprintf("sign up phase must be %s, %s or %s", SignUpPhasePersonalInfo, SignUpPhasePreferences, SignUpPhaseDone)}
	}

	if u.UserRating < 0 || u.UserRating > 5 {
		return ValidationError{"user rating is not between 0 and 5"}
	}

	if u.DriverRating < 0 || u.DriverRating > 5 {
		return ValidationError{"driver rating is not between 0 and 5"}
	}

	return nil
}
