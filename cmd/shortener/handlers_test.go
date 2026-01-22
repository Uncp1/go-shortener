package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleShorten(t *testing.T) {
	tests := []struct {
		name             string
		requestBody      string
		expectedStatus   int
		responseContains string
	}{
		{
			name:             "success",
			requestBody:      "https://youtu.be/XHKm_dGFTzM?si=u-j5hqDr34FlUVsr",
			expectedStatus:   http.StatusCreated,
			responseContains: "http://localhost:8080/",
		},
		{
			name:           "empty body",
			requestBody:    "",
			expectedStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.requestBody))
			w := httptest.NewRecorder()

			handleShorten(w, req)

			result := w.Result()
			defer result.Body.Close()
			body, err := io.ReadAll(result.Body)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, result.StatusCode)

			if tt.responseContains != "" {
				assert.Contains(t, string(body), tt.responseContains)
			}
		})
	}
}

func TestHandleRedirect(t *testing.T) {
	urlStore["abc12345"] = "https://youtu.be/XHKm_dGFTzM?si=u-j5hqDr34FlUVsr"

	tests := []struct {
		name             string
		id               string
		expectedStatus   int
		expectedLocation string
	}{
		{
			name:             "success",
			id:               "abc12345",
			expectedStatus:   http.StatusTemporaryRedirect,
			expectedLocation: "https://youtu.be/XHKm_dGFTzM?si=u-j5hqDr34FlUVsr",
		},
		{
			name:           "not found",
			id:             "unknown",
			expectedStatus: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/"+tt.id, nil)
			req.SetPathValue("id", tt.id)
			w := httptest.NewRecorder()

			handleRedirect(w, req)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, tt.expectedStatus, result.StatusCode)

			if tt.expectedLocation != "" {
				assert.Equal(t, tt.expectedLocation, result.Header.Get("Location"))
			}
		})
	}
}
