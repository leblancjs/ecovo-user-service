package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"azure.com/ecovo/user-service/pkg/auth"
	"azure.com/ecovo/user-service/pkg/db"
	"azure.com/ecovo/user-service/pkg/httperror"
	"azure.com/ecovo/user-service/pkg/models"
	"azure.com/ecovo/user-service/pkg/requestid"
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
	db, err := db.New(&dbConfig)
	if err != nil {
		log.Fatal(err)
	}

	env := &Env{
		port:  port,
		store: db}

	r := mux.NewRouter()
	r.HandleFunc("/users/me", requestid.Middleware(auth.ValidationMiddleware(authValidator, env.getUserFromAuthHandler))).
		Methods("GET")
	r.HandleFunc("/users/{id}", requestid.Middleware(auth.ValidationMiddleware(authValidator, env.getUserByIDHandler))).
		Methods("GET").
		Headers("Content-Type", "application/json")
	r.HandleFunc("/users/{id}", requestid.Middleware(auth.ValidationMiddleware(authValidator, env.updateUserByIDHandler))).
		Methods("PATCH").
		HeadersRegexp("Content-Type", "application/(json|json; charset=utf8)")
	r.HandleFunc("/users", requestid.Middleware(auth.ValidationMiddleware(authValidator, env.createUserHandler))).
		Methods("POST").
		HeadersRegexp("Content-Type", "application/(json|json; charset=utf8)")

	log.Fatal(http.ListenAndServe(":"+env.port, handlers.LoggingHandler(os.Stdout, r)))
}

func (env *Env) getUserFromAuthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo, err := auth.UserInfoFromContext(r.Context())
	if err != nil {
		httperror.Handler(w, r,
			httperror.NewInternalServerError(httperror.ErrInternalServerError, err))
	} else {
		user, _ := env.store.FindUserByAuth0ID(userInfo.ID)
		if user == nil {
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
				SignUpPhase: models.SignUpPhasePersonalInfo,
			})
			if err != nil {
				httperror.Handler(w, r,
					httperror.NewInternalServerError(httperror.ErrInternalServerError, err))
			}
		} else {
			err := json.NewEncoder(w).Encode(user)
			if err != nil {
				httperror.Handler(w, r,
					httperror.NewInternalServerError(httperror.ErrInternalServerError, err))
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
		httperror.Handler(w, r,
			httperror.NewNotFoundError(httperror.ErrUserNotFound, err))
	} else {
		err := json.NewEncoder(w).Encode(user)
		if err != nil {
			httperror.Handler(w, r,
				httperror.NewInternalServerError(httperror.ErrInternalServerError, err))
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}

func (env *Env) updateUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	user, err := env.store.FindUserByID(id)
	if err != nil {
		httperror.Handler(w, r,
			httperror.NewNotFoundError(httperror.ErrUserNotFound, err))
	} else {
		var modifiedUser models.User
		err := json.NewDecoder(r.Body).Decode(&modifiedUser)
		if err != nil {
			httperror.Handler(w, r,
				httperror.NewBadRequestError(httperror.ErrBadRequest, err))
		} else {
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
				user.Preferences.Conversation = modifiedUser.Preferences.Conversation
				user.Preferences.Music = modifiedUser.Preferences.Music
			}

			if modifiedUser.SignUpPhase != "" {
				user.SignUpPhase = modifiedUser.SignUpPhase
			}

			err = user.Validate()
			if err != nil {
				httperror.Handler(w, r,
					httperror.NewBadRequestError(err.Error(), err))
			} else {
				err := env.store.UpdateUser(user)
				if err != nil {
					httperror.Handler(w, r,
						httperror.NewInternalServerError(httperror.ErrInternalServerError, err))
				} else {
					w.WriteHeader(http.StatusOK)
				}
			}
		}
	}
}

func (env *Env) createUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo, err := auth.UserInfoFromContext(r.Context())
	if err != nil {
		httperror.Handler(w, r,
			httperror.NewInternalServerError(httperror.ErrInternalServerError, err))
	} else {
		user, _ := env.store.FindUserByAuth0ID(userInfo.ID)
		if user != nil {
			httperror.Handler(w, r,
				httperror.NewInternalServerError(httperror.ErrInternalServerError, errors.New("user already exists with given Auth0 ID")))
		} else {
			err := json.NewDecoder(r.Body).Decode(&user)
			if err != nil {
				httperror.Handler(w, r,
					httperror.NewBadRequestError(httperror.ErrBadRequest, err))
			} else {
				user.Auth0ID = userInfo.ID
				user.Email = userInfo.Email

				if user.Preferences == nil {
					user.Preferences = &models.UserPreferences{}
				}
				user.Preferences.Smoking = models.Occasionally
				user.Preferences.Conversation = models.Occasionally
				user.Preferences.Music = models.Occasionally

				user.SignUpPhase = models.SignUpPhasePreferences

				err = user.Validate()
				if err != nil {
					httperror.Handler(w, r,
						httperror.NewBadRequestError(err.Error(), err))
				} else {
					user, err = env.store.CreateUser(user)
					if err != nil {
						httperror.Handler(w, r,
							httperror.NewInternalServerError(httperror.ErrInternalServerError, err))
					} else {
						err = json.NewEncoder(w).Encode(user)
						if err != nil {
							env.store.DeleteUser(user)
							httperror.Handler(w, r,
								httperror.NewInternalServerError(httperror.ErrInternalServerError, err))
						} else {
							w.WriteHeader(http.StatusCreated)
						}
					}
				}
			}
		}
	}
}
