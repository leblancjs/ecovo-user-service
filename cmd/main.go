package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"azure.com/ecovo/user-service/cmd/handler"
	"azure.com/ecovo/user-service/cmd/middleware/auth"
	"azure.com/ecovo/user-service/pkg/db"
	"azure.com/ecovo/user-service/pkg/user"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

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

	userRepository, err := user.NewMongoRepository(db.Users)
	if err != nil {
		log.Fatal(err)
	}
	userUseCase := user.NewService(userRepository)

	r := mux.NewRouter()
	r.Handle("/users/me", handler.RequestID(handler.Auth(authValidator, handler.GetUserFromAuth(userUseCase)))).
		Methods("GET")
	r.Handle("/users/{id}", handler.RequestID(handler.Auth(authValidator, handler.GetUserByID(userUseCase)))).
		Methods("GET").
		Headers("Content-Type", "application/json")
	r.Handle("/users/{id}", handler.RequestID(handler.Auth(authValidator, handler.UpdateUser(userUseCase)))).
		Methods("PATCH").
		HeadersRegexp("Content-Type", "application/(json|json; charset=utf8)")
	r.Handle("/users", handler.RequestID(handler.Auth(authValidator, handler.CreateUser(userUseCase)))).
		Methods("POST").
		HeadersRegexp("Content-Type", "application/(json|json; charset=utf8)")

	log.Fatal(http.ListenAndServe(":"+port, handlers.LoggingHandler(os.Stdout, r)))
}
