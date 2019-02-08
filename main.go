package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const (
	never        = 0
	occasionally = 1
	regularly    = 2
)

type userPreferences struct {
	Smoking      int `json:"smoking"`
	Animals      int `json:"animals"`
	Conversation int `json:"conversation"`
	Music        int `json:"music"`
}

type user struct {
	ID          string           `json:"id"`
	Email       string           `json:"email"`
	FirstName   string           `json:"firstName"`
	LastName    string           `json:"lastName"`
	DateOfBirth time.Time        `json:"dateOfBirth"`
	PhoneNumber string           `json:"phoneNumber"`
	Gender      string           `json:"gender"`
	Photo       string           `json:"photo"`
	Description string           `json:"description"`
	Preferences *userPreferences `json:"preferences"`
	SignUpPhase *int             `json:"signUpPhase"`
}

var mockPrefs = userPreferences{
	Smoking:      regularly,
	Animals:      never,
	Conversation: occasionally,
	Music:        occasionally}

var mockUser = user{
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

func getUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	// TODO: Remove mock response
	mockUser.ID = id

	err := json.NewEncoder(w).Encode(mockUser)
	if err != nil {
		log.Print(err)
	}
}

func updateUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	// TODO: Remove mock response
	mockUser.ID = id

	var user user
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Print(err)
	}

	if user.FirstName != "" {
		mockUser.FirstName = user.FirstName
	}
	if user.LastName != "" {
		mockUser.LastName = user.LastName
	}
	if !user.DateOfBirth.IsZero() {
		mockUser.DateOfBirth = user.DateOfBirth
	}
	if user.PhoneNumber != "" {
		mockUser.PhoneNumber = user.PhoneNumber
	}
	if user.Gender != "" {
		mockUser.Gender = user.Gender
	}
	if user.Photo != "" {
		mockUser.Photo = user.Photo
	}
	if user.Description != "" {
		mockUser.Description = user.Description
	}
	if user.Preferences != nil {
		mockUser.Preferences.Smoking = user.Preferences.Smoking
		mockUser.Preferences.Animals = user.Preferences.Animals
		mockUser.Preferences.Conversation = user.Preferences.Conversation
		mockUser.Preferences.Music = user.Preferences.Music
	}
	if user.SignUpPhase != nil {
		mockUser.SignUpPhase = user.SignUpPhase
	}

	w.WriteHeader(http.StatusOK)
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user user
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Print(err)
	}

	// TODO: Remove mock response
	if user.Preferences == nil {
		user.Preferences = &userPreferences{}
	}
	user.Preferences.Smoking = occasionally
	user.Preferences.Animals = occasionally
	user.Preferences.Conversation = occasionally
	user.Preferences.Music = occasionally

	user.SignUpPhase = new(int)
	*user.SignUpPhase = 0

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		log.Print(err)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/users/{id}", getUserByIDHandler).
		Methods("GET").
		Headers("Content-Type", "application/json")
	r.HandleFunc("/users/{id}", updateUserByIDHandler).
		Methods("PATCH").
		Headers("Content-Type", "application/json")
	r.HandleFunc("/users", createUserHandler).
		Methods("POST").
		Headers("Content-Type", "application/json")

	log.Fatal(http.ListenAndServe(":8080", r))
}
