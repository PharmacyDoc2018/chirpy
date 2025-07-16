package main

import (
	"net/http"
	"sync/atomic"
	"time"

	"github.com/PharmacyDoc2018/chirpy/internal/database"
	"github.com/google/uuid"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
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

type returnVal struct {
	Valid       bool   `json:"valid"`
	CleanedBody string `json:"cleaned_body"`
}

type returnErr struct {
	Error string `json:"error"`
}

type newUserRequest struct {
	Email string `json:"email"`
}

type newUserResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}
