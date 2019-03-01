package handler

import (
	"encoding/json"
	"net/http"

	"azure.com/ecovo/user-service/cmd/middleware/auth"
	"azure.com/ecovo/user-service/pkg/entity"
	"azure.com/ecovo/user-service/pkg/user"
	"azure.com/ecovo/user-service/pkg/vehicule"
	"github.com/gorilla/mux"
)

// CreateVehicule handles a request to create a vehicule. We first verify the User related to the
// vehicules exists.
func CreateVehicule(uService user.UseCase, vService vehicule.UseCase) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)

		var v *entity.Vehicule
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			return err
		}

		userInfo, err := auth.FromContext(r.Context())
		if err != nil {
			return err
		}

		userID := entity.NewIDFromHex(vars["userId"])
		v.UserID = userID

		v, err = vService.Register(v, userInfo.SubID)
		if err != nil {
			return err
		}

		w.WriteHeader(http.StatusCreated)

		err = json.NewEncoder(w).Encode(v)
		if err != nil {
			_ = vService.Delete(entity.ID(v.ID), userID, userInfo.SubID)

			return err
		}

		return nil
	}
}

// DeleteVehicule handles a request to delete a vehicule.
func DeleteVehicule(uService user.UseCase, vService vehicule.UseCase) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		vars := mux.Vars(r)

		userInfo, err := auth.FromContext(r.Context())
		if err != nil {
			return err
		}

		id := entity.NewIDFromHex(vars["id"])
		userID := entity.NewIDFromHex(vars["userId"])

		err = vService.Delete(id, userID, userInfo.SubID)
		if err != nil {
			return err
		}

		w.WriteHeader(http.StatusOK)

		return nil
	}
}

// GetVehiculeByID handles a request to retrieve a vehicule by its unique identifier.
func GetVehiculeByID(uService user.UseCase, vService vehicule.UseCase) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)

		id := entity.NewIDFromHex(vars["id"])
		v, err := vService.FindByID(id)
		if err != nil {
			return err
		}

		err = json.NewEncoder(w).Encode(v)
		if err != nil {
			return err
		}

		w.WriteHeader(http.StatusOK)

		return nil
	}
}

// GetVehiculesByUserID handles a request to retrieve the authenticated user's vehicule.
func GetVehiculesByUserID(uService user.UseCase, vService vehicule.UseCase) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)

		userID := entity.NewIDFromHex(vars["userId"])
		v, err := vService.FindByUserID(userID)
		if err != nil {
			return err
		}

		err = json.NewEncoder(w).Encode(v)
		if err != nil {
			return err
		}

		return nil
	}
}
