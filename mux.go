package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/PharmacyDoc2018/chirpy/internal/database"
	"github.com/google/uuid"
)

const filepathRoot = "."
const port = "8080"

func initMux(cfg *apiConfig) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot)))))
	mux.Handle("/app/assets/", cfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot)))))

	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))

	})

	mux.HandleFunc("GET /admin/metrics/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl := template.Must(template.ParseFiles("./admin/metrics/index.html"))
		err := tmpl.Execute(w, struct {
			Hits int32
		}{
			Hits: cfg.fileserverHits.Load(),
		})
		if err != nil {
			fmt.Println(err)
		}
	})

	mux.HandleFunc("POST /admin/reset", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		cfg.fileserverHits.Store(0)
		w.Write([]byte("hit count reset!"))
	})

	mux.HandleFunc("POST /api/validate_chirp", func(w http.ResponseWriter, req *http.Request) {
		genericErrorReturn := returnErr{
			Error: "Something went wrong",
		}

		tooLongErrorReturn := returnErr{
			Error: "Chirp is too long",
		}

		goodReturn := returnVal{
			Valid: true,
		}

		w.Header().Set("Content-Type", "application/json")

		decoder := json.NewDecoder(req.Body) // Decode Request
		defer req.Body.Close()
		newChirp := chirp{}
		err := decoder.Decode(&newChirp)
		if err != nil {
			fmt.Printf("error decoding chirp: %s\n", err)
			w.WriteHeader(500)
			data, err := json.Marshal(genericErrorReturn)
			if err != nil {
				fmt.Println(err)
			}
			w.Write(data)
			return
		}

		lenChirp := len(newChirp.Body) // Encode Response
		if lenChirp > 140 {            // If chirp is too long
			data, err := json.Marshal(tooLongErrorReturn)
			if err != nil {
				w.WriteHeader(500)
				fmt.Println(err)
				return
			}
			w.WriteHeader(400)
			w.Write(data)

		} else { // If chirp is not too long
			goodReturn.CleanedBody = filterProfanity(newChirp.Body) // Clean up the profanity
			data, err := json.Marshal(goodReturn)
			if err != nil {
				w.WriteHeader(500)
				fmt.Println(err)
				return
			}
			w.WriteHeader(200)
			w.Write(data)
		}

	})

	mux.HandleFunc("POST /api/users", func(w http.ResponseWriter, req *http.Request) {
		genericErrorReturn := returnErr{
			Error: "Something went wrong",
		}

		w.Header().Set("Content-Type", "application/json")

		newUserReq := newUserRequest{}
		decoder := json.NewDecoder(req.Body)
		defer req.Body.Close()
		err := decoder.Decode(&newUserReq)
		if err != nil {
			fmt.Printf("error decoding email: %s\n", err)
			w.WriteHeader(500)
			data, err := json.Marshal(genericErrorReturn)
			if err != nil {
				fmt.Println(err)
			}
			w.Write(data)
			return
		}

		params := database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			Email:     newUserReq.Email,
		}

		newUser, err := cfg.db.CreateUser(req.Context(), params)
		if err != nil {
			fmt.Println(err)
			errorReturn := returnErr{
				Error: fmt.Sprint(err),
			}
			data, err := json.Marshal(errorReturn)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(500)
				return
			}
			w.WriteHeader(400)
			w.Write(data)
			return
		}

		returnedNewUser := newUserResponse{
			ID:        newUser.ID,
			CreatedAt: newUser.CreatedAt,
			UpdatedAt: newUser.UpdatedAt,
			Email:     newUser.Email,
		}

		data, err := json.Marshal(returnedNewUser)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(200)
		w.Write(data)
	})

	return mux
}
