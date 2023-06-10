package http

import (
	"errors"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

func enableCORS(router *mux.Router) {
	router.Use(middlewareCors)
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		if req.URL.Path != "/authentication" && req.URL.Path != "/logout" {
			notAuthorized := validateToken(req, w)
			if notAuthorized {
				return
			}
		}

		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("FRONTEND_URL"))
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if req.Method == http.MethodOptions {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, req)
	})
}

func validateToken(req *http.Request, w http.ResponseWriter) bool {
	session, err := store.Get(req, "session-name")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return true
	}

	token, ok := session.Values["token"].(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return true
	}

	_, errValidatingToken := ValidateToken(token, os.Getenv("SECRET_KEY_TOKEN"))
	if errValidatingToken != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return true
	}
	return false
}

func ValidateToken(tokenString string, secretKey string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if token.Method.Alg() != jwt.SigningMethodHS256.Name {
			return nil, errors.New("unexpected signing method")
		}
		// Return the secret key used for signing the token
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
