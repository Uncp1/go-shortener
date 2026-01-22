package main

import (
	"io"
	"math/rand"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// TO DO: небезопасно для concurrent access
var urlStore = make(map[string]string)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateID(length int) string {
	b := make([]byte, length)

	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)
}

func handleShorten(w http.ResponseWriter, r *http.Request) {
	originalURL, err := io.ReadAll(r.Body)

	if err != nil || len(originalURL) == 0 {
		http.Error(w, "cannot read body", http.StatusBadRequest)
		return
	}

	shortID := generateID(8)
	urlStore[shortID] = string(originalURL)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("http://localhost:8080/" + shortID))
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	originalURL, ok := urlStore[id]

	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
