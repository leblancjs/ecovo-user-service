package entity

import "fmt"

// Preferences contains a user's preferences when it comes to smoking,
// conversation and music.
type Preferences struct {
	Smoking      int `json:"smoking" bson:"smoking"`
	Conversation int `json:"conversation" bson:"conversation"`
	Music        int `json:"music" bson:"music"`
}

const (
	// PreferenceNever represents that a user never smokes, listens to music or
	// talks.
	PreferenceNever = 0

	// PreferenceOccasionally represents that a user occasionally smokes,
	// listens to music or talks.
	PreferenceOccasionally = 1

	// PreferenceRegularly represents that a user regularly smokes, listens to
	// music or talks.
	PreferenceRegularly = 2
)

// Validate validates that the preferences' required fields are filled out
// correctly.
func (p *Preferences) Validate() error {
	if p.Smoking < PreferenceNever || p.Smoking > PreferenceRegularly {
		return ValidationError{fmt.Sprintf("smoking preference is out of bounds \"%d\"", p.Smoking)}
	}

	if p.Conversation < PreferenceNever || p.Conversation > PreferenceRegularly {
		return ValidationError{fmt.Sprintf("conversation preference is out of bounds \"%d\"", p.Conversation)}
	}

	if p.Music < PreferenceNever || p.Music > PreferenceRegularly {
		return ValidationError{fmt.Sprintf("music preference is out of bounds \"%d\"", p.Music)}
	}

	return nil
}
