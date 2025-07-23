package main

import (
	"net/http"
	"sync/atomic"
	"time"

	"github.com/PharmacyDoc2018/chirpy/internal/database"
	"github.com/google/uuid"
)

const maxChirpLength = 140
const maxTokenLifetime = 3600 //seconds

type apiConfig struct {
	platfrom       string
	secret         string
	db             *database.Queries
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

type chirp struct {
	Body string `json:"body"`
}

type chirpResponse struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

type chirpError struct {
	err        error
	writtenErr []byte
}

type returnErr struct {
	Error string `json:"error"`
}

type loginRequest struct {
	Email     string        `json:"email"`
	Password  string        `json:"password"`
	ExpiresIn time.Duration `json:"expires_in_seconds,omitempty"`
}

type userResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type loginResponse struct {
	Token string `json:"token"`
	userResponse
}
