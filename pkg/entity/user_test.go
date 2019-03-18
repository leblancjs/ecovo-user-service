package entity

import (
	"testing"
	"time"
)

func TestUserValidation(t *testing.T) {
	var preferences = Preferences{
		Smoking:      PreferenceNever,
		Music:        PreferenceRegularly,
		Conversation: PreferenceOccasionally,
	}

	var location, _ = time.LoadLocation("")

	var user = User{
		SubID:        "harold|hide.the.pain",
		Email:        "harold@hide-the-pain.meme",
		FirstName:    "Harold",
		LastName:     "The Great",
		DateOfBirth:  time.Date(1950, time.February, 12, 0, 0, 0, 0, location),
		PhoneNumber:  "(450) 123-4567",
		Gender:       GenderMale,
		Photo:        "https://hide-the-pain.meme/harold.png",
		Description:  "So much pain.",
		Preferences:  &preferences,
		SignUpPhase:  SignUpPhasePersonalInfo,
		UserRating:   4,
		DriverRating: 2,
	}

	t.Run("Should fail when subscription ID is empty", func(t *testing.T) {
		u := user
		u.SubID = ""

		if _, ok := u.Validate().(ValidationError); !ok {
			t.Fail()
		}
	})

	t.Run("Should fail when first name is empty", func(t *testing.T) {
		u := user
		u.FirstName = ""

		if _, ok := u.Validate().(ValidationError); !ok {
			t.Fail()
		}
	})

	t.Run("Should fail when last name is empty", func(t *testing.T) {
		u := user
		u.LastName = ""

		if _, ok := u.Validate().(ValidationError); !ok {
			t.Fail()
		}
	})

	t.Run("Should fail when date of birth is empty", func(t *testing.T) {
		u := user
		u.DateOfBirth = time.Time{}

		if _, ok := u.Validate().(ValidationError); !ok {
			t.Fail()
		}
	})

	t.Run("Should fail when date of birth reveals user is not 18 years of age", func(t *testing.T) {
		u := user
		u.DateOfBirth = time.Now().AddDate(-17, 0, -1)

		if _, ok := u.Validate().(ValidationError); !ok {
			t.Fail()
		}
	})

	t.Run("Should fail when gender is empty", func(t *testing.T) {
		u := user
		u.Gender = ""

		if _, ok := u.Validate().(ValidationError); !ok {
			t.Fail()
		}
	})

	t.Run("Should fail when gender is not valid", func(t *testing.T) {
		u := user
		u.Gender = "Harold"

		if _, ok := u.Validate().(ValidationError); !ok {
			t.Fail()
		}
	})

	t.Run("Should succeed when gender is "+GenderMale, func(t *testing.T) {
		u := user
		u.Gender = "Male"

		err := u.Validate()
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Should succeed when gender is "+GenderFemale, func(t *testing.T) {
		u := user
		u.Gender = "Female"

		err := u.Validate()
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Should succeed when gender is "+GenderOther, func(t *testing.T) {
		u := user
		u.Gender = "Other"

		err := u.Validate()
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Should succeed when email is empty", func(t *testing.T) {
		u := user
		u.Email = ""

		err := u.Validate()
		if err != nil {
			t.Fail()
		}
	})

	t.Run("Should succeed when phone number is empty", func(t *testing.T) {
		u := user
		u.PhoneNumber = ""

		err := u.Validate()
		if err != nil {
			t.Fail()
		}
	})

	t.Run("Should succeed when preferences are nil", func(t *testing.T) {
		u := user
		u.Preferences = nil

		err := u.Validate()
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Should fail when sign up phase is not valid", func(t *testing.T) {
		u := user
		u.SignUpPhase = "Harold"

		if _, ok := u.Validate().(ValidationError); !ok {
			t.Fail()
		}
	})

	t.Run("Should succeed when sign up phase is "+SignUpPhasePersonalInfo, func(t *testing.T) {
		u := user
		u.SignUpPhase = SignUpPhasePersonalInfo

		err := u.Validate()
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Should succeed when sign up phase is "+SignUpPhasePreferences, func(t *testing.T) {
		u := user
		u.SignUpPhase = SignUpPhasePreferences

		err := u.Validate()
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Should succeed when sign up phase is "+SignUpPhaseDone, func(t *testing.T) {
		u := user
		u.SignUpPhase = SignUpPhaseDone

		err := u.Validate()
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Should fail when user rating is over then 5", func(t *testing.T) {
		u := user
		u.UserRating = 6

		if _, ok := u.Validate().(ValidationError); !ok {
			t.Fail()
		}
	})

	t.Run("Should fail when user rating is under then 0", func(t *testing.T) {
		u := user
		u.UserRating = -1

		if _, ok := u.Validate().(ValidationError); !ok {
			t.Fail()
		}
	})

	t.Run("Should fail when driver rating is over then 5", func(t *testing.T) {
		u := user
		u.DriverRating = 6

		if _, ok := u.Validate().(ValidationError); !ok {
			t.Fail()
		}
	})

	t.Run("Should fail when driver rating is under then 0", func(t *testing.T) {
		u := user
		u.DriverRating = -1

		if _, ok := u.Validate().(ValidationError); !ok {
			t.Fail()
		}
	})
}
