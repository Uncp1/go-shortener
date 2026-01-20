package main

import (
	"io"
	"math/rand"
	"net/http"
)

var urlStore = make(map[string]string)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateId(length int) string {
	b := make([]byte, length)

	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)
}

func shorten(w http.ResponseWriter, r *http.Request) {
	originalURL, err := io.ReadAll(r.Body)

	if err != nil || len(originalURL) == 0 {
		http.Error(w, "cannot read body", http.StatusBadRequest)
		return
	}

	shortId := generateId(8)
	urlStore[shortId] = string(originalURL)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("http://localhost:8080/" + shortId))
}

func redirect(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	originalURL, ok := urlStore[id]

	if !ok {
		http.Error(w, "not found", http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
