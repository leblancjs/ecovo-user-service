package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"azure.com/ecovo/user-service/auth"
	"azure.com/ecovo/user-service/models"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Env struct {
	port string
	repo *models.DB
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// TODO: Initialize connection with a database
	db, _ := models.NewDB("")

	env := &Env{
		port: port,
		repo: db}

	r := mux.NewRouter()
	r.HandleFunc("/users/me", auth.TokenValidationMiddleware(env.getUserFromAuthHandler)).
		Methods("GET")
	r.HandleFunc("/users/{id}", auth.TokenValidationMiddleware(env.getUserByIDHandler)).
		Methods("GET").
		Headers("Content-Type", "application/json")
	r.HandleFunc("/users/{id}", auth.TokenValidationMiddleware(env.updateUserByIDHandler)).
		Methods("PATCH").
		HeadersRegexp("Content-Type", "application/(json|json; charset=utf8)")
	r.HandleFunc("/users", auth.TokenValidationMiddleware(env.createUserHandler)).
		Methods("POST").
		HeadersRegexp("Content-Type", "application/(json|json; charset=utf8)")

	log.Fatal(http.ListenAndServe(":"+env.port, handlers.LoggingHandler(os.Stdout, r)))
}

func (env *Env) getUserFromAuthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo, err := auth.UserInfoFromContext(r.Context())
	if err != nil {
		// TODO: Write more informative error message
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		user, err := env.repo.FindByAuth0ID(userInfo.ID)
		if err != nil {
			// TODO: Write more informative error message
			w.WriteHeader(http.StatusNotFound)
		} else {
			err := json.NewEncoder(w).Encode(user)
			if err != nil {
				log.Print(err)
			}
		}
	}
}

func (env *Env) getUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	user, err := env.repo.FindByID(id)
	if err != nil {
		// TODO: Write more informative error message
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)

		err := json.NewEncoder(w).Encode(user)
		if err != nil {
			log.Print(err)
		}
	}
}

func (env *Env) updateUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	// TODO: Get user from a database
	user, err := env.repo.FindByID(id)
	if err != nil {
		// TODO: Write more informative error message
		w.WriteHeader(http.StatusNotFound)
	} else {
		var modifiedUser models.User
		err := json.NewDecoder(r.Body).Decode(&modifiedUser)
		if err != nil {
			log.Print(err)
		}

		if modifiedUser.FirstName != "" {
			user.FirstName = modifiedUser.FirstName
		}
		if modifiedUser.LastName != "" {
			user.LastName = modifiedUser.LastName
		}
		if !modifiedUser.DateOfBirth.IsZero() {
			user.DateOfBirth = modifiedUser.DateOfBirth
		}
		if modifiedUser.PhoneNumber != "" {
			user.PhoneNumber = modifiedUser.PhoneNumber
		}
		if modifiedUser.Gender != "" {
			user.Gender = modifiedUser.Gender
		}
		if modifiedUser.Photo != "" {
			user.Photo = modifiedUser.Photo
		}
		if modifiedUser.Description != "" {
			user.Description = modifiedUser.Description
		}
		// TODO: Fix bug where preferences with null value crush existing values
		if modifiedUser.Preferences != nil {
			user.Preferences.Smoking = modifiedUser.Preferences.Smoking
			user.Preferences.Animals = modifiedUser.Preferences.Animals
			user.Preferences.Conversation = modifiedUser.Preferences.Conversation
			user.Preferences.Music = modifiedUser.Preferences.Music
		}
		if modifiedUser.SignUpPhase != nil {
			if user.SignUpPhase == nil {
				user.SignUpPhase = new(int)
			}
			*user.SignUpPhase = *modifiedUser.SignUpPhase
		}

		env.repo.Update(user)

		w.WriteHeader(http.StatusOK)
	}
}

func (env *Env) createUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo, err := auth.UserInfoFromContext(r.Context())
	if err != nil {
		// TODO: Write more informative error message
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		// TODO: Check if user exists in a database
		user, _ := env.repo.FindByAuth0ID(userInfo.ID)
		if user != nil {
			// TODO: Write more informative error message
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			err := json.NewDecoder(r.Body).Decode(&user)
			if err != nil {
				log.Print(err)

				// TODO: Write more informative error message
				w.WriteHeader(http.StatusBadRequest)
			} else {
				user.Auth0ID = userInfo.ID

				// TODO: Validate presence of required fields
				// ...

				if user.Preferences == nil {
					user.Preferences = &models.UserPreferences{}
				}
				user.Preferences.Smoking = models.Occasionally
				user.Preferences.Animals = models.Occasionally
				user.Preferences.Conversation = models.Occasionally
				user.Preferences.Music = models.Occasionally

				user.SignUpPhase = new(int)
				*user.SignUpPhase = 1

				// TODO: Write the user to a database
				user, err = env.repo.Create(user)
				if err != nil {
					log.Print(err)

					// TODO: Write more informative error message
					w.WriteHeader(http.StatusInternalServerError)
				} else {
					w.WriteHeader(http.StatusCreated)

					err = json.NewEncoder(w).Encode(user)
					if err != nil {
						log.Print(err)

						env.repo.Delete(user)

						// TODO: Write more informative error message
						w.WriteHeader(http.StatusInternalServerError)
					}
				}
			}
		}
	}
}
