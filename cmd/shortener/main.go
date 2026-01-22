package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	r.Post("/", handleShorten)
	r.Get("/{id}", handleRedirect)

	err := http.ListenAndServe(":8080", r)
	log.Fatal(err)
}
