package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"azure.com/ecovo/user-service/auth"
	"azure.com/ecovo/user-service/db"
	"azure.com/ecovo/user-service/models"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Env struct {
	port  string
	store db.Store
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	authConfig := auth.Config{
		Domain: os.Getenv("AUTH_DOMAIN")}
	authValidator, err := auth.NewTokenValidator(&authConfig)
	if err != nil {
		log.Fatal(err)
	}

	dbConnectionTimeout, err := time.ParseDuration(os.Getenv("DB_CONNECTION_TIMEOUT") + "s")
	if err != nil {
		dbConnectionTimeout = db.DefaultConnectionTimeout
	}
	dbConfig := db.Config{
		Host:              os.Getenv("DB_HOST"),
		Username:          os.Getenv("DB_USERNAME"),
		Password:          os.Getenv("DB_PASSWORD"),
		Name:              os.Getenv("DB_NAME"),
		ConnectionTimeout: dbConnectionTimeout}
	db, err := db.NewDB(&dbConfig)
	if err != nil {
		log.Fatal(err)
	}

	env := &Env{
		port:  port,
		store: db}

	r := mux.NewRouter()
	r.HandleFunc("/users/me", auth.ValidationMiddleware(authValidator, env.getUserFromAuthHandler)).
		Methods("GET")
	r.HandleFunc("/users/{id}", auth.ValidationMiddleware(authValidator, env.getUserByIDHandler)).
		Methods("GET").
		Headers("Content-Type", "application/json")
	r.HandleFunc("/users/{id}", auth.ValidationMiddleware(authValidator, env.updateUserByIDHandler)).
		Methods("PATCH").
		HeadersRegexp("Content-Type", "application/(json|json; charset=utf8)")
	r.HandleFunc("/users", auth.ValidationMiddleware(authValidator, env.createUserHandler)).
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
		user, err := env.store.FindUserByAuth0ID(userInfo.ID)
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

	user, err := env.store.FindUserByID(id)
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
	user, err := env.store.FindUserByID(id)
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
			user.Preferences.Conversation = modifiedUser.Preferences.Conversation
			user.Preferences.Music = modifiedUser.Preferences.Music
		}
		if modifiedUser.SignUpPhase != nil {
			if user.SignUpPhase == nil {
				user.SignUpPhase = new(int)
			}
			*user.SignUpPhase = *modifiedUser.SignUpPhase
		}

		env.store.UpdateUser(user)

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
		user, _ := env.store.FindUserByAuth0ID(userInfo.ID)
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
				user.Preferences.Conversation = models.Occasionally
				user.Preferences.Music = models.Occasionally

				user.SignUpPhase = new(int)
				*user.SignUpPhase = 1

				// TODO: Write the user to a database
				user, err = env.store.CreateUser(user)
				if err != nil {
					log.Print(err)

					// TODO: Write more informative error message
					w.WriteHeader(http.StatusInternalServerError)
				} else {
					err = json.NewEncoder(w).Encode(user)
					if err != nil {
						log.Print(err)

						err = env.store.DeleteUser(user)
						if err != nil {
							log.Print(err)
						}

						// TODO: Write more informative error message
						w.WriteHeader(http.StatusInternalServerError)
					} else {
						w.WriteHeader(http.StatusCreated)
					}
				}
			}
		}
	}
}
