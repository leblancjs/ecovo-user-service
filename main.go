package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

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

// TODO: Store Auth0 ID in user
type User struct {
	ID          string           `json:"id"`
	Auth0ID     string           `json:"auth0Id,omitempty"`
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

var mockPrefs = UserPreferences{
	Smoking:      regularly,
	Animals:      never,
	Conversation: occasionally,
	Music:        occasionally}

var mockUser = User{
	ID:          "12345",
	Email:       "harold@hide-the-pain.com",
	FirstName:   "Harold",
	LastName:    "The Great",
	DateOfBirth: time.Now(),
	PhoneNumber: "5141234567",
	Gender:      "Male",
	Photo:       "https://hungarytoday.hu/wp-content/uploads/2018/02/18ps27.jpg",
	Description: "So much pain.",
	Preferences: &mockPrefs,
	SignUpPhase: nil}

// TODO: Store users in a database
var usersByID map[string]User
var nextID int

func getUserInfo(auth string) map[string]interface{} {
	// TODO: Put domain in configuration file
	req, err := http.NewRequest("GET", "https://ecovo.auth0.com/userinfo", nil)
	if err != nil {
		log.Print("failed to create GET /userinfo request")
		return nil
	}
	req.Header.Set("Authorization", auth)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Print("failed to do request")
		return nil
	}

	var userInfo map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	if err != nil {
		log.Print("failed to decode user info")
		return nil
	}
	return userInfo
}

func getUserFromAuthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// TODO: Use helper for context key
	userInfo := r.Context().Value("userInfo")
	if userInfo == nil {
		// TODO: Write more informative error message
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		// TODO: Get user based on Auth0 ID
		// ...

		err := json.NewEncoder(w).Encode(mockUser)
		if err != nil {
			log.Print(err)
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
			*user.SignUpPhase = *modifiedUser.SignUpPhase
		}

		// TODO: Update user in a database
		// ...

		w.WriteHeader(http.StatusOK)
	}
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Print(err)
	}

	// TODO: Check presence in a database based on Auth0 ID
	_, present := usersByID[user.ID]
	if present {
		// TODO: Write more informative error message
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		// TODO: Store Auth0 ID in user and generate ID through database
		user.ID = strconv.Itoa(nextID)
		nextID++

		if user.Preferences == nil {
			user.Preferences = &UserPreferences{}
		}
		user.Preferences.Smoking = occasionally
		user.Preferences.Animals = occasionally
		user.Preferences.Conversation = occasionally
		user.Preferences.Music = occasionally

		user.SignUpPhase = new(int)
		*user.SignUpPhase = 0

		// TODO: Write the user to a database
		usersByID[user.ID] = user

		w.WriteHeader(http.StatusCreated)

		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			log.Print(err)
		}
	}
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userInfo := getUserInfo(r.Header.Get("Authorization"))
		if userInfo == nil {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			// TODO: Create helper for context keys
			ctx := context.WithValue(r.Context(), "userInfo", userInfo)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}

func main() {
	// TODO: Initialize connection with a database
	usersByID = make(map[string]User)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := mux.NewRouter()
	r.HandleFunc("/users/me", authMiddleware(getUserFromAuthHandler)).
		Methods("GET")
	r.HandleFunc("/users/{id}", authMiddleware(getUserByIDHandler)).
		Methods("GET").
		Headers("Content-Type", "application/json")
	r.HandleFunc("/users/{id}", authMiddleware(updateUserByIDHandler)).
		Methods("PATCH").
		Headers("Content-Type", "application/json")
	r.HandleFunc("/users", authMiddleware(createUserHandler)).
		Methods("POST").
		Headers("Content-Type", "application/json")

	log.Fatal(http.ListenAndServe(":"+port, handlers.LoggingHandler(os.Stdout, r)))
}
