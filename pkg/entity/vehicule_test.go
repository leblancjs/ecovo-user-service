package entity

import (
	"testing"
)

func TestVehiculeValidation(t *testing.T) {
	var accessories = []string{"A/C", "Heating Seats"}

	var vehicule = Vehicule{
		Photo:       "https://hide-the-pain.meme/harold.png",
		Year:        2018,
		Make:        "Audi",
		Model:       "A4",
		Color:       "Noir",
		Accessories: accessories,
	}

	t.Run("Should fail when user ID is empty", func(t *testing.T) {
		v := vehicule
		v.UserID = ""

		if _, ok := v.Validate().(ValidationError); !ok {
			t.Fail()
		}
	})

	t.Run("Should fail when year is not valid", func(t *testing.T) {
		v := vehicule
		v.Year = 0

		if _, ok := v.Validate().(ValidationError); !ok {
			t.Fail()
		}
	})

	t.Run("Should fail when make is empty", func(t *testing.T) {
		v := vehicule
		v.Make = ""

		if _, ok := v.Validate().(ValidationError); !ok {
			t.Fail()
		}
	})

	t.Run("Should fail when model is empty", func(t *testing.T) {
		v := vehicule
		v.Model = ""

		if _, ok := v.Validate().(ValidationError); !ok {
			t.Fail()
		}
	})

	t.Run("Should fail when color is empty", func(t *testing.T) {
		v := vehicule
		v.Color = ""

		if _, ok := v.Validate().(ValidationError); !ok {
			t.Fail()
		}
	})

}
