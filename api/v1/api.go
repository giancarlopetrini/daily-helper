package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"log"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte("changethissecret"), nil)
}

type statusReponse struct {
	RemoteAddr string `json:"remoteAddr"`
	Status     string `json:"status"`
}

func jsonResponse(w http.ResponseWriter, data interface{}, c int) {
	dj, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(c)
	fmt.Fprintf(w, "%s", dj)
}

func status(w http.ResponseWriter, r *http.Request) {
	response := statusReponse{
		RemoteAddr: r.RemoteAddr,
		Status:     "OK",
	}

	jsonResponse(w, response, http.StatusOK)
}

// NewRouter returns an HTTP handler that implements the routes for the API
func NewRouter() http.Handler {
	r := chi.NewRouter()

	// Unprotected routes
	r.Group(func(r chi.Router) {
		r.Post("/user", CreateUser)
		r.Get("/status", status)
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)
		// Register the API routes
		r.Get("/user/{user_id}", GetUser)
	})

	return r
}
