package entity

import (
	"fmt"
	"time"
)

// Vehicule contains a vehicule's information.
type Vehicule struct {
	ID          ID       `json:"id" bson:"_id,omitempty"`
	UserID      ID       `json:"userId" bson:"userId"`
	Year        int      `json:"year" bson:"year"`
	Make        string   `json:"make" bson:"make"`
	Model       string   `json:"model" bson:"model"`
	Color       string   `json:"color" bson:"color"`
	Photo       string   `json:"photo" bson:"photo"`
	Seats       int      `json:"seats" bson:"seats"`
	Accessories []string `json:"accessories" bson:"accessories"`
}

const (
	// YearMinimum represents the minimum year of a car.
	YearMinimum = 1900
)

// Validate validates that the vehicules's required fields are filled out correctly.
func (v *Vehicule) Validate() error {
	if v.UserID.IsZero() {
		return ValidationError{fmt.Sprintf("id must not be nil")}
	}

	if v.Year <= YearMinimum && v.Year > time.Now().Year() {
		return ValidationError{fmt.Sprintf("year must me between %d and %d", YearMinimum, time.Now().Year())}
	}

	if v.Make == "" {
		return ValidationError{"make is missing"}
	}

	if v.Color == "" {
		return ValidationError{"color is missing"}
	}

	if v.Seats < 1 {
		return ValidationError{"minimum number of seats is 1"}
	}

	return nil
}
