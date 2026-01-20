package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /", shorten)
	mux.HandleFunc("GET /{id}", redirect)

	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
