package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/PharmacyDoc2018/chirpy/internal/auth"
	"github.com/PharmacyDoc2018/chirpy/internal/database"
	"github.com/google/uuid"
)

func handleResourseChirps(mux *http.ServeMux, cfg *apiConfig) {

	mux.HandleFunc("POST /api/chirps", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		token, err := auth.GetBearerToken(req.Header)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(401)
			data, _ := json.Marshal(returnErr{
				Error: fmt.Sprint(err),
			})
			w.Write(data)
		}

		userID, err := auth.ValidateJWT(token, cfg.secret)
		if err != nil {
			returnError := returnErr{
				Error: fmt.Sprintf("401 Unauthorized. %s", err),
			}
			w.WriteHeader(401)
			data, _ := json.Marshal(returnError)
			w.Write(data)
		}

		newChirp, decodeErr := decodeChirp(req)
		if decodeErr != nil {
			fmt.Println(decodeErr.err)
			w.WriteHeader(400)
			w.Write(decodeErr.writtenErr)
			return
		}

		if ok, err := validateChirpLength(newChirp); !ok {
			fmt.Println(err.err)
			w.WriteHeader(400)
			w.Write(err.writtenErr)
			return
		}

		filterProfanity(newChirp)

		chirpParams := database.CreateChirpParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			Body:      newChirp.Body,
			UserID:    userID,
		}

		storedChirp, err := cfg.db.CreateChirp(req.Context(), chirpParams)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		returnedChirp := chirpResponse{
			Id:        storedChirp.ID,
			CreatedAt: storedChirp.CreatedAt,
			UpdatedAt: storedChirp.UpdatedAt,
			Body:      storedChirp.Body,
			UserId:    storedChirp.UserID,
		}

		data, err := json.Marshal(returnedChirp)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		w.WriteHeader(201)
		w.Write(data)

	})

	mux.HandleFunc("GET /api/chirps", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		storedChirps, err := cfg.db.GetChirps(req.Context())
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
			return
		}

		var returnedChirps []chirpResponse

		for _, storedChirp := range storedChirps {
			returnedChirps = append(returnedChirps, chirpResponse{
				Id:        storedChirp.ID,
				CreatedAt: storedChirp.CreatedAt,
				UpdatedAt: storedChirp.UpdatedAt,
				Body:      storedChirp.Body,
				UserId:    storedChirp.UserID,
			})
		}

		data, err := json.Marshal(returnedChirps)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
			return
		}

		w.WriteHeader(200)
		w.Write(data)
	})

	mux.HandleFunc("GET /api/chirps/{id}", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		id, err := uuid.Parse(req.PathValue("id"))
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(400)
			return
		}
		storedChirp, err := cfg.db.GetChirp(req.Context(), id)
		if err != nil {
			if fmt.Sprint(err) == "sql: no rows in result set" {
				w.WriteHeader(404)
				return
			} else {
				fmt.Println(err)
				w.WriteHeader(500)
				return
			}
		}

		data, err := json.Marshal(chirpResponse{
			Id:        storedChirp.ID,
			CreatedAt: storedChirp.CreatedAt,
			UpdatedAt: storedChirp.UpdatedAt,
			Body:      storedChirp.Body,
			UserId:    storedChirp.UserID,
		})
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
			return
		}

		w.WriteHeader(200)
		w.Write(data)
	})

	mux.HandleFunc("DELETE /api/chirps/{id}", func(w http.ResponseWriter, req *http.Request) {
		token, err := auth.GetBearerToken(req.Header)
		if err != nil {
			fmt.Println(err)
			data, err := json.Marshal(returnErr{
				Error: fmt.Sprint(err),
			})
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(500)
				return
			}
			w.WriteHeader(401)
			w.Write(data)
			return
		}

		userID, err := auth.ValidateJWT(token, cfg.secret)
		if err != nil {
			fmt.Println(err)
			data, err := json.Marshal(returnErr{
				Error: fmt.Sprint(err),
			})
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(500)
				return
			}
			w.WriteHeader(401)
			w.Write(data)
			return
		}

		id, err := uuid.Parse(req.PathValue("id"))
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(400)
			return
		}
		storedChirp, err := cfg.db.GetChirp(req.Context(), id)
		if err != nil {
			if fmt.Sprint(err) == "sql: no rows in result set" {
				w.WriteHeader(404)
				return
			} else {
				fmt.Println(err)
				w.WriteHeader(500)
				return
			}
		}

		if userID != storedChirp.UserID {
			data, err := json.Marshal(returnErr{
				Error: "403 Unauthorized",
			})
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(500)
				return
			}
			w.WriteHeader(403)
			w.Write(data)
			return
		}

		_, err = cfg.db.DeleteChirpByID(req.Context(), storedChirp.ID)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
			return
		}

		w.WriteHeader(204)
	})

}
