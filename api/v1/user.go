package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth"

	"github.com/go-chi/chi"

	jwt "github.com/dgrijalva/jwt-go"
)

// CreateUserRequest info needed to create user
type CreateUserRequest struct {
	UserID string `json:"user_id"`
}

// CreateUserResponse sends info back
type CreateUserResponse struct {
	Message string `json:"message"`
	UserID  string `json:"user_id"`
	JWT     string `json:"token"`
}

// CreateUser will create users, currently just makes token
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var request CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	timestamp := time.Now().UnixNano()

	// For debugging/example purposes, we generate and print
	// a sample jwt token with claims `user_id:123` here:
	_, tokenString, _ := tokenAuth.Encode(jwt.MapClaims{
		"user_id":   request.UserID,
		"timestamp": timestamp,
	})

	response := CreateUserResponse{
		Message: "user created successfully",
		UserID:  request.UserID,
		JWT:     tokenString,
	}

	jsonResponse(w, response, http.StatusOK)
}

// GetUserResponse checks JWT token for user info/time created
type GetUserResponse struct {
	UserID    string  `json:"user_id"`
	Timestamp float64 `json:"time_created"`
}

// GetUser requires JWT, and gets user info
func GetUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "user_id")

	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if userID != claims["user_id"] {
		http.Error(w, "Token does not match user", 500)
	}

	fmt.Println(claims)

	response := GetUserResponse{
		UserID:    claims["user_id"].(string),
		Timestamp: claims["timestamp"].(float64),
	}

	jsonResponse(w, response, http.StatusOK)
}
