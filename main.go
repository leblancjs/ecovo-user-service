package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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
	ID          int             `json:"id"`
	Email       string          `json:"email"`
	FirstName   string          `json:"firstName"`
	LastName    string          `json:"lastName"`
	DateOfBirth time.Time       `json:"dateOfBirth"`
	PhoneNumber string          `json:"phoneNumber"`
	Gender      string          `json:"gender"`
	Photo       string          `json:"photo"`
	Description string          `json:"description"`
	Preferences userPreferences `json:"preferences"`
}

var mockPrefs = userPreferences{
	Smoking:      regularly,
	Animals:      never,
	Conversation: occasionally,
	Music:        occasionally}

var mockUser = user{
	ID:          12345,
	Email:       "harold@hide-the-pain.com",
	FirstName:   "Harold",
	LastName:    "The Great",
	DateOfBirth: time.Now(),
	PhoneNumber: "5141234567",
	Gender:      "Male",
	Photo:       "https://hungarytoday.hu/wp-content/uploads/2018/02/18ps27.jpg",
	Description: "So much pain.",
	Preferences: mockPrefs}

func getUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, _ := strconv.Atoi(vars["id"])

	// TODO: Remove mock response
	mockUser.ID = id
	json, _ := json.Marshal(mockUser)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func main() {
	r := mux.NewRouter()
	r.Headers("Content-Type", "application/json")
	r.HandleFunc("/users/{id}", getUserByIDHandler).
		Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}
