package middlewares

import (
	"errors"
	"net/http"

	"github.com/arvinpaundra/ecommerce-api/api/auth"
	"github.com/arvinpaundra/ecommerce-api/api/responses"
)

// Middleware to format all responses into json
func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

// Middleware to check validity of authentication token
func SetMiddlewareAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := auth.TokenValid(r)

		if err != nil {
			responses.ERROR(w, http.StatusUnauthorized, errors.New("Unathorized. Your token is expired."))
			return
		}

		next(w, r)
	}
}
