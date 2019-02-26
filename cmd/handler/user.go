package handler

import (
	"encoding/json"
	"net/http"

	"azure.com/ecovo/user-service/cmd/middleware/auth"
	"azure.com/ecovo/user-service/pkg/entity"
	"azure.com/ecovo/user-service/pkg/user"
	"github.com/gorilla/mux"
)

// CreateUser handles a request to create a user.
func CreateUser(service user.UseCase) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-Type", "application/json")

		var u *entity.User
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			return err
		}

		userInfo, err := auth.FromContext(r.Context())
		if err != nil {
			return err
		}

		u.SubID = userInfo.SubID
		u.Email = userInfo.Email

		u, err = service.Register(u)
		if err != nil {
			return err
		}

		w.WriteHeader(http.StatusCreated)

		err = json.NewEncoder(w).Encode(u)
		if err != nil {
			_ = service.Delete(entity.ID(u.ID))

			return err
		}

		return nil
	}
}

// UpdateUser handles a request to update a user.
func UpdateUser(service user.UseCase) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)

		var u *entity.User
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			return err
		}

		u.ID = entity.NewIDFromHex(vars["id"])

		err = service.Update(u)
		if err != nil {
			return err
		}

		w.WriteHeader(http.StatusOK)

		return nil
	}
}

// GetUserByID handles a request to retrieve a user by its unique identifier.
func GetUserByID(service user.UseCase) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)

		id := entity.NewIDFromHex(vars["id"])
		u, err := service.FindByID(id)
		if err != nil {
			return err
		}

		err = json.NewEncoder(w).Encode(u)
		if err != nil {
			return err
		}

		w.WriteHeader(http.StatusOK)

		return nil
	}
}

// GetUserFromAuth handles a request to retrieve the authenticated user.
func GetUserFromAuth(service user.UseCase) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-Type", "application/json")

		userInfo, err := auth.FromContext(r.Context())
		if err != nil {
			return err
		}

		u, err := service.FindBySubID(userInfo.SubID)
		if err != nil {
			type tmpUser struct {
				Email       string `json:"email"`
				FirstName   string `json:"firstName"`
				LastName    string `json:"lastName"`
				Photo       string `json:"photo"`
				SignUpPhase string `json:"signUpPhase"`
			}

			err := json.NewEncoder(w).Encode(&tmpUser{
				Email:       userInfo.Email,
				FirstName:   userInfo.FirstName,
				LastName:    userInfo.LastName,
				Photo:       userInfo.Picture,
				SignUpPhase: entity.SignUpPhasePersonalInfo,
			})
			if err != nil {
				return err
			}
		} else {
			err = json.NewEncoder(w).Encode(u)
			if err != nil {
				return err
			}
		}

		return nil
	}
}
