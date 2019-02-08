package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!\n"))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/users", helloHandler)

	log.Fatal(http.ListenAndServe(":8080", r))
}
