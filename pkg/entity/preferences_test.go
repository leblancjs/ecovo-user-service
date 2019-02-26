package entity

import "testing"

func TestPreferencesValidation(t *testing.T) {
	preferences := Preferences{
		Smoking:      PreferenceNever,
		Music:        PreferenceRegularly,
		Conversation: PreferenceOccasionally,
	}

	t.Run("Should fail when smoking preference is below lower bound", func(t *testing.T) {
		p := preferences
		p.Smoking = PreferenceNever - 1

		if _, ok := p.Validate().(ValidationError); !ok {
			t.Fail()
		}
	})

	t.Run("Should fail when smoking preference is above upper bound", func(t *testing.T) {
		p := preferences
		p.Smoking = PreferenceRegularly + 1

		if _, ok := p.Validate().(ValidationError); !ok {
			t.Fail()
		}
	})

	t.Run("Should fail when music preference is below lower bound", func(t *testing.T) {
		p := preferences
		p.Music = PreferenceNever - 1

		if _, ok := p.Validate().(ValidationError); !ok {
			t.Fail()
		}
	})

	t.Run("Should fail when music preference is above upper bound", func(t *testing.T) {
		p := preferences
		p.Music = PreferenceRegularly + 1

		if _, ok := p.Validate().(ValidationError); !ok {
			t.Fail()
		}
	})

	t.Run("Should fail when conversation preference is below lower bound", func(t *testing.T) {
		p := preferences
		p.Conversation = PreferenceNever - 1

		if _, ok := p.Validate().(ValidationError); !ok {
			t.Fail()
		}
	})

	t.Run("Should fail when conversation preference is above upper bound", func(t *testing.T) {
		p := preferences
		p.Conversation = PreferenceRegularly + 1

		if _, ok := p.Validate().(ValidationError); !ok {
			t.Fail()
		}
	})
}
