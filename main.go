package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"azure.com/ecovo/user-service/auth"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const (
	never        = 0
	occasionally = 1
	regularly    = 2
)

type UserPreferences struct {
	Smoking      int `json:"smoking"`
	Animals      int `json:"animals"`
	Conversation int `json:"conversation"`
	Music        int `json:"music"`
}

type User struct {
	auth0ID     string
	ID          string           `json:"id"`
	Email       string           `json:"email"`
	FirstName   string           `json:"firstName"`
	LastName    string           `json:"lastName"`
	DateOfBirth time.Time        `json:"dateOfBirth"`
	PhoneNumber string           `json:"phoneNumber"`
	Gender      string           `json:"gender"`
	Photo       string           `json:"photo"`
	Description string           `json:"description"`
	Preferences *UserPreferences `json:"preferences"`
	SignUpPhase *int             `json:"signUpPhase"`
}

// TODO: Store users in a database
var usersByID map[string]*User
var nextID int

func getUserFromAuthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo, err := auth.UserInfoFromContext(r.Context())
	if err != nil {
		// TODO: Write more informative error message
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		// TODO: Get user from a database
		var user *User
		for _, u := range usersByID {
			if u.auth0ID == userInfo.ID {
				user = u
				break
			}
		}

		if user == nil {
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

func getUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	// TODO: Get user from a database
	user, present := usersByID[id]
	if !present {
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

func updateUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	// TODO: Get user from a database
	user, present := usersByID[id]
	if !present {
		// TODO: Write more informative error message
		w.WriteHeader(http.StatusNotFound)
	} else {
		var modifiedUser User
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

		// TODO: Update user in a database
		// ...

		w.WriteHeader(http.StatusOK)
	}
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo, err := auth.UserInfoFromContext(r.Context())
	if err != nil {
		// TODO: Write more informative error message
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		// TODO: Check if user exists in a database
		userAlreadyExists := false
		for _, u := range usersByID {
			if u.auth0ID == userInfo.ID {
				userAlreadyExists = true
				break
			}
		}
		if userAlreadyExists {
			// TODO: Write more informative error message
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			var user User
			err := json.NewDecoder(r.Body).Decode(&user)
			if err != nil {
				log.Print(err)

				// TODO: Write more informative error message
				w.WriteHeader(http.StatusBadRequest)
			} else {
				// TODO: Generate ID through database
				user.auth0ID = userInfo.ID
				user.ID = strconv.Itoa(nextID)
				nextID++

				// TODO: Validate presence of required fields
				// ...

				if user.Preferences == nil {
					user.Preferences = &UserPreferences{}
				}
				user.Preferences.Smoking = occasionally
				user.Preferences.Animals = occasionally
				user.Preferences.Conversation = occasionally
				user.Preferences.Music = occasionally

				user.SignUpPhase = new(int)
				*user.SignUpPhase = 1

				// TODO: Write the user to a database
				usersByID[user.ID] = &user

				err = json.NewEncoder(w).Encode(user)
				if err != nil {
					log.Print(err)

					delete(usersByID, user.ID)
					nextID--

					// TODO: Write more informative error message
					w.WriteHeader(http.StatusInternalServerError)
				} else {
					w.WriteHeader(http.StatusCreated)
				}
			}
		}
	}
}

func main() {
	// TODO: Initialize connection with a database
	usersByID = make(map[string]*User)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := mux.NewRouter()
	r.HandleFunc("/users/me", auth.TokenValidationMiddleware(getUserFromAuthHandler)).
		Methods("GET")
	r.HandleFunc("/users/{id}", auth.TokenValidationMiddleware(getUserByIDHandler)).
		Methods("GET").
		Headers("Content-Type", "application/json")
	r.HandleFunc("/users/{id}", auth.TokenValidationMiddleware(updateUserByIDHandler)).
		Methods("PATCH").
		Headers("Content-Type", "application/json")
	r.HandleFunc("/users", auth.TokenValidationMiddleware(createUserHandler)).
		Methods("POST").
		Headers("Content-Type", "application/json")

	log.Fatal(http.ListenAndServe(":"+port, handlers.LoggingHandler(os.Stdout, r)))
}
