package main

import (
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
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

type returnVal struct {
	Valid       bool   `json:"valid"`
	CleanedBody string `json:"cleaned_body"`
}

type returnErr struct {
	Error string `json:"error"`
}
