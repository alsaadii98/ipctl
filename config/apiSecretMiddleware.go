package config

import (
	"encoding/json"
	"net/http"
	"os"
)

type APIResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func sendError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(APIResponse{
		Status:  statusCode,
		Message: message,
	})
}

func APISecretMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiSecret := r.Header.Get("x-api-secret")
		expectedSecret := os.Getenv("API_SECRET")

		if apiSecret != expectedSecret {
			sendError(w, http.StatusForbidden, "Forbidden: Invalid API secret")
			return
		}

		next.ServeHTTP(w, r)
	})
}
